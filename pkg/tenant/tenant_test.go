package tenant_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	sdk "github.com/rotationalio/go-ensign"
	emock "github.com/rotationalio/go-ensign/mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tenantTestSuite struct {
	suite.Suite
	srv         *tenant.Server
	auth        *authtest.Server
	client      api.TenantClient
	quarterdeck *mock.Server
	ensign      *emock.Ensign
	subscriber  *tenant.TopicSubscriber
	stop        chan bool
}

// Runs once before all tests are executed
func (suite *tenantTestSuite) SetupSuite() {
	// Note use assert instead of require so that go routines are properly handled in
	// tests; assert uses t.Error while require uses t.FailNow and multiple go routines
	// might lead to incorrect testing behavior.
	var err error
	assert := suite.Assert()
	suite.stop = make(chan bool, 1)

	// Discards logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Start the authtest server for authentication verification
	suite.auth, err = authtest.NewServer()
	assert.NoError(err, "could not start the authtest server")

	// Start an httptest server to handle mock requests to Quarterdeck
	suite.quarterdeck, err = mock.NewServer()
	assert.NoError(err, "could not start the quarterdeck mock server")

	// Ensure Quarterdeck returns a 200 on status so Tenant knows it's ready
	suite.quarterdeck.OnStatus(mock.UseStatus(http.StatusOK))

	// Start a server to handle mock requests to Ensign
	suite.ensign = emock.New(nil)

	// Creates a test configuration to run the Tenant API server as a fully
	// functional server on an open port using the local-loopback for networking.
	conf, err := config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost:3000"},
		Auth: config.AuthConfig{
			Audience:     authtest.Audience,
			Issuer:       authtest.Issuer,
			KeysURL:      suite.auth.KeysURL(),
			CookieDomain: "localhost",
		},
		SendGrid: emails.Config{
			FromEmail:  "ensign@rotational.io",
			AdminEmail: "admins@rotational.io",
			Testing:    true,
		},
		Quarterdeck: config.QuarterdeckConfig{
			URL:          suite.quarterdeck.URL(),
			WaitForReady: 1 * time.Second,
		},
		Database: config.DatabaseConfig{
			Testing: true,
		},
		Ensign: config.SDKConfig{
			Enabled:          true,
			ClientID:         "testing",
			ClientSecret:     "testing",
			Endpoint:         "bufconn",
			Insecure:         true,
			NoAuthentication: true,
		},
	}.Mark()
	assert.NoError(err, "test configuration is invalid")

	suite.srv, err = tenant.New(conf)
	assert.NoError(err, "could not create the tenant api server from the test configuration")

	// Starts the Tenant server. Server will run for the duration of all tests.
	// Implements reset methods to ensure the server state doesn't change
	// between tests in Before/After.
	go func() {
		if err := suite.srv.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// This is a bad enough error that we should panic, otherwise tests will
			// fail and it will be hard to debug why the tests failed.
			panic(err)
		}
		suite.stop <- true
	}()

	// Wait for 500ms to ensure the API server starts.
	time.Sleep(500 * time.Millisecond)

	// Creates a Tenant client to make requests to the server.
	assert.NotEmpty(suite.srv.URL(), "no url to connect the client on")
	suite.client, err = api.New(suite.srv.URL())
	assert.NoError(err, "could not initialize the Tenant client")

	// Set the Ensign client on the server
	sdkClient, err := sdk.New(sdk.WithMock(suite.ensign, grpc.WithTransportCredentials(insecure.NewCredentials())))
	assert.NoError(err, "could not connect an sdk client to the mock ensign server")

	ensignClient := &tenant.EnsignClient{}
	ensignClient.SetClient(sdkClient)
	ensignClient.SetOpts(conf.Ensign)

	suite.srv.SetEnsignClient(ensignClient)
	suite.subscriber = tenant.NewTopicSubscriber(ensignClient)
}

func (suite *tenantTestSuite) TearDownSuite() {
	assert := suite.Assert()

	// Shutdown the quarterdeck mock server
	suite.quarterdeck.Close()

	// Shutdown the ensign mock server
	suite.ensign.Shutdown()

	// Shutdown the authtest server
	suite.auth.Close()

	// Shuts down the tenant API server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := suite.srv.Shutdown(ctx)
	assert.NoError(err, "could not gracefully shut down the tenant test server")

	// Waits for server to stop in order to prevent race conditions.
	<-suite.stop

	// Cleanup logger
	logger.ResetLogger()
}

func (suite *tenantTestSuite) AfterTest(suiteName, testName string) {
	// Ensure any credentials set on the client are reset
	suite.client.(*api.APIv1).SetCredentials("")
	suite.client.(*api.APIv1).SetCSRFProtect(false)

	// Reset the quarterdeck mock server
	suite.quarterdeck.Reset()
}

// Stop the task manager, waiting for all the tasks to finish. Tests should defer
// ResetTasks() to ensure that the task manager is available to the other tests.
func (suite *tenantTestSuite) StopTasks() {
	tasks := suite.srv.GetTaskManager()
	tasks.Stop()
}

// Reset the task manager to ensure that other tests have access to it.
func (suite *tenantTestSuite) ResetTasks() {
	suite.srv.ResetTaskManager()
}

// Helper function to set cookies for CSRF protection on the tenant client
func (s *tenantTestSuite) SetClientCSRFProtection() error {
	s.client.(*api.APIv1).SetCSRFProtect(true)
	return nil
}

// Helper function to set the credentials on the test client from claims, reducing 3 or
// 4 lines of code into a single helper function call to make tests more readable.
func (s *tenantTestSuite) SetClientCredentials(claims *tokens.Claims) error {
	token, err := s.auth.CreateAccessToken(claims)
	if err != nil {
		return err
	}

	s.client.(*api.APIv1).SetCredentials(token)
	return nil
}

// Helper function to add the user claims to the current context.
func (s *tenantTestSuite) ContextWithClaims(ctx context.Context, claims *tokens.Claims) (c context.Context, err error) {
	token, err := s.auth.CreateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	return qd.ContextWithToken(ctx, token), nil
}

func TestTenant(t *testing.T) {
	suite.Run(t, &tenantTestSuite{})
}

func statusMessage(status int, message string) string {
	return fmt.Sprintf("[%d] %s", status, message)
}

func (s *tenantTestSuite) requireError(err error, status int, message string, msgAndArgs ...interface{}) {
	require := s.Require()
	require.EqualError(err, statusMessage(status, message), msgAndArgs...)
}

func (s *tenantTestSuite) requireMultiError(err error, messages ...string) {
	require := s.Require()
	require.IsType(&multierror.Error{}, err)

	var actual []string
	for _, e := range err.(*multierror.Error).Errors {
		actual = append(actual, e.Error())
	}

	require.ElementsMatch(messages, actual)
}
