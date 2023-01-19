package tenant_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/tenant/mock"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type tenantTestSuite struct {
	suite.Suite
	srv         *tenant.Server
	auth        *authtest.Server
	client      api.TenantClient
	quarterdeck *mock.Server
	stop        chan bool
}

// Runs once before all tests are executed
func (suite *tenantTestSuite) SetupSuite() {
	var err error
	require := suite.Require()
	suite.stop = make(chan bool, 1)

	// Discards logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Start the authtest server for authentication verification
	suite.auth, err = authtest.NewServer()
	require.NoError(err, "could not start the authtest server")

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
		Database: config.DatabaseConfig{
			Testing: true,
		},
	}.Mark()
	require.NoError(err, "test configuration is invalid")

	suite.srv, err = tenant.New(conf)
	require.NoError(err, "could not create the tenant api server from the test configuration")

	// Start an httptest server to handle mock requests to Quarterdeck
	suite.quarterdeck, err = mock.NewServer()

	// Starts the Tenant server. Server will run for the duration of all tests.
	// Implements reset methods to ensure the server state doesn't change
	// between tests in Before/After.
	go func() {
		suite.srv.Serve()
		suite.stop <- true
	}()

	// Wait for 500ms to ensure the API server starts.
	time.Sleep(500 * time.Millisecond)

	// Creates a Tenant client to make requests to the server.
	require.NotEmpty(suite.srv.URL(), "no url to connect the client on")
	suite.client, err = api.New(suite.srv.URL())
	require.NoError(err, "could not initialize the Tenant client")
}

func (suite *tenantTestSuite) TearDownSuite() {
	require := suite.Require()

	// Shutdown the quarterdeck mock server
	suite.quarterdeck.Close()

	// Shutdown the authtest server
	suite.auth.Close()

	// Shuts down the tenant API server.
	err := suite.srv.Shutdown()
	require.NoError(err, "could not gracefully shut down the tenant test server")

	// Waits for server to stop in order to prevent race conditions.
	<-suite.stop

	// Cleanup logger
	logger.ResetLogger()
}

func (suite *tenantTestSuite) AfterTest(suiteName, testName string) {
	// Ensure any credentials set on the client are reset
	suite.client.(*api.APIv1).SetCredentials("")
	suite.client.(*api.APIv1).SetCSRFProtect(false)
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
