package quarterdeck_test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type quarterdeckTestSuite struct {
	suite.Suite
	srv    *quarterdeck.Server
	client api.QuarterdeckClient
	stop   chan bool
}

// Run once before all the tests are executed
func (suite *quarterdeckTestSuite) SetupSuite() {
	require := suite.Require()
	suite.stop = make(chan bool, 1)

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Create a test configuration to run the Quarterdeck API server as a fully
	// functional server on an open port using the local-loopback for networking.
	conf, err := config.Config{
		Maintenance: false,
		BindAddr:    "127.0.0.1:0",
		Mode:        gin.TestMode,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
	}.Mark()
	require.NoError(err, "test configuration is invalid")

	suite.srv, err = quarterdeck.New(conf)
	require.NoError(err, "could not create the quarterdeck api server from the test configuration")

	// Start the BFF server - the goal of the tests is to have the server run for the
	// entire duration of the tests. Implement reset methods to ensure the server state
	// doesn't change between tests in Before/After.
	go func() {
		suite.srv.Serve()
		suite.stop <- true
	}()

	// Wait for 500ms to ensure the API server starts up
	time.Sleep(500 * time.Millisecond)

	// Create a Quarterdeck client for making requests to the server
	require.NotEmpty(suite.srv.URL(), "no url to connect the client on")
	suite.client, err = api.New(suite.srv.URL())
	require.NoError(err, "could not initialize the Quarterdeck client")
}

func (suite *quarterdeckTestSuite) TearDownSuite() {
	require := suite.Require()

	// Shutdown the quarterdeck API server
	err := suite.srv.Shutdown()
	require.NoError(err, "could not gracefully shutdown the quarterdeck test server")

	// Wait for server to stop to prevent race conditions
	<-suite.stop

	// Cleanup logger
	logger.ResetLogger()
}

func TestQuarterdeck(t *testing.T) {
	suite.Run(t, &quarterdeckTestSuite{})
}
