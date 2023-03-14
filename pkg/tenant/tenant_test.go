package tenant_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	emock "github.com/rotationalio/go-ensign/mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type tenantTestSuite struct {
	suite.Suite
	srv         *tenant.Server
	auth        *authtest.Server
	client      api.TenantClient
	quarterdeck *mock.Server
	ensign      *emock.Ensign
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
		Ensign: config.EnsignConfig{
			Insecure: true,
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
	ensignClient, err := suite.ensign.Client(context.Background())
	assert.NoError(err, "could not initialize the Ensign mock client")
	suite.srv.SetEnsignClient(ensignClient)
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

func TestTenant(t *testing.T) {
	suite.Run(t, &tenantTestSuite{})
}

func (s *tenantTestSuite) requireError(err error, status int, message string, msgAndArgs ...interface{}) {
	require := s.Require()
	require.EqualError(err, fmt.Sprintf("[%d] %s", status, message), msgAndArgs...)
}
