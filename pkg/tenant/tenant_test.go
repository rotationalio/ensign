package tenant_test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type tenantTestSuite struct {
	suite.Suite
	srv    *tenant.Server
	client api.TenantClient
	stop   chan bool
}

// Runs once before all tests are executed
func (suite *tenantTestSuite) SetupSuite() {
	require := suite.Require()
	suite.stop = make(chan bool, 1)

	// Discards logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Creates a test configuration to run the Tenant API server as a fully
	// functional server on an open port using the local-loopback for networking.
	conf, err := config.Config{
		Maintenance: false,
		BindAddr:    "127.0.0.1:0",
		Mode:        gin.TestMode,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
	}.Mark()
	require.NoError(err, "test configuration is invalid")

	suite.srv, err = tenant.New(conf)
	require.NoError(err, "could not create the tenant api server from the test configuration")

	// Starts the BFF server. Server will run for the duration of all tests.
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

	// Shuts down the tenant API server.
	err := suite.srv.Shutdown()
	require.NoError(err, "could not gracefully shut down the tenant test server")

	// Waits for servr to stop in order to prevent race conditions.
	<-suite.stop

	// Cleanup logger
	logger.ResetLogger()
}

func TestTenant(t *testing.T) {
	suite.Run(t, &tenantTestSuite{})
}
