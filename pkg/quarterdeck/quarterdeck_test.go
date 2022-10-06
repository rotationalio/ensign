package quarterdeck_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type quarterdeckTestSuite struct {
	suite.Suite
	srv    *quarterdeck.Server
	client api.QuarterdeckClient
	dbPath string
	stop   chan bool
}

// Run once before all the tests are executed
func (suite *quarterdeckTestSuite) SetupSuite() {
	require := suite.Require()
	suite.stop = make(chan bool, 1)

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Create a temporary test database for the tests
	var err error
	suite.dbPath, err = os.MkdirTemp("", "quarterdeck-*")
	require.NoError(err, "could not create temporary directory for database")

	// Create a test configuration to run the Quarterdeck API server as a fully
	// functional server on an open port using the local-loopback for networking.
	conf, err := config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost:3000"},
		Database: config.DatabaseConfig{
			URL:      "sqlite3:///" + filepath.Join(suite.dbPath, "test.db"),
			ReadOnly: false,
		},
		Token: config.TokenConfig{
			Keys: map[string]string{
				"01GE6191AQTGMCJ9BN0QC3CCVG": "testdata/01GE6191AQTGMCJ9BN0QC3CCVG.pem",
				"01GE62EXXR0X0561XD53RDFBQJ": "testdata/01GE62EXXR0X0561XD53RDFBQJ.pem",
			},
			Audience: "http://localhost:3000",
			Issuer:   "http://quarterdeck.test",
		},
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

	// Cleanup temporary test directory
	err = os.RemoveAll(suite.dbPath)
	require.NoError(err, "could not cleanup temporary database")
}

func (suite *quarterdeckTestSuite) ResetDatabase() (err error) {
	// Truncate all database tables except roles, permissions, and role_permissions
	stmts := []string{
		"DELETE FROM organizations",
		"DELETE FROM users",
		"DELETE FROM organization_users",
		"DELETE FROM projects",
		"DELETE FROM api_keys",
		"DELETE FROM user_roles",
		"DELETE FROM api_key_permissions",
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(context.Background(), nil); err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		if _, err = tx.Exec(stmt); err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func TestQuarterdeck(t *testing.T) {
	suite.Run(t, &quarterdeckTestSuite{})
}
