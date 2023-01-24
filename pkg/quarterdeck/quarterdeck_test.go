package quarterdeck_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type quarterdeckTestSuite struct {
	suite.Suite
	srv    *quarterdeck.Server
	conf   config.Config
	client api.QuarterdeckClient
	dbPath string
	stop   chan bool
}

// Run once before all the tests are executed
func (s *quarterdeckTestSuite) SetupSuite() {
	require := s.Require()
	s.stop = make(chan bool, 1)

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Create a temporary test database for the tests
	var err error
	s.dbPath, err = os.MkdirTemp("", "quarterdeck-*")
	require.NoError(err, "could not create temporary directory for database")

	// Create a test configuration to run the Quarterdeck API server as a fully
	// functional server on an open port using the local-loopback for networking.
	s.conf, err = config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost:3000"},
		Database: config.DatabaseConfig{
			URL:      "sqlite3:///" + filepath.Join(s.dbPath, "test.db"),
			ReadOnly: false,
		},
		Token: config.TokenConfig{
			Keys: map[string]string{
				"01GE6191AQTGMCJ9BN0QC3CCVG": "testdata/01GE6191AQTGMCJ9BN0QC3CCVG.pem",
				"01GE62EXXR0X0561XD53RDFBQJ": "testdata/01GE62EXXR0X0561XD53RDFBQJ.pem",
			},
			Audience:        "http://localhost:3000",
			Issuer:          "http://quarterdeck.test/",
			AccessDuration:  10 * time.Minute,
			RefreshDuration: 20 * time.Minute,
			RefreshOverlap:  -5 * time.Minute,
		},
	}.Mark()
	require.NoError(err, "test configuration is invalid")

	s.srv, err = quarterdeck.New(s.conf)
	require.NoError(err, "could not create the quarterdeck api server from the test configuration")

	// Start the BFF server - the goal of the tests is to have the server run for the
	// entire duration of the tests. Implement reset methods to ensure the server state
	// doesn't change between tests in Before/After.
	go func() {
		s.srv.Serve()
		s.stop <- true
	}()

	// Wait for 500ms to ensure the API server starts up
	time.Sleep(500 * time.Millisecond)

	// Load database fixtures
	require.NoError(s.LoadDatabaseFixtures(), "could not load database fixtures")

	// Create a Quarterdeck client for making requests to the server
	require.NotEmpty(s.srv.URL(), "no url to connect the client on")
	s.client, err = api.New(s.srv.URL())
	require.NoError(err, "could not initialize the Quarterdeck client")
}

func (s *quarterdeckTestSuite) TearDownSuite() {
	require := s.Require()

	// Shutdown the quarterdeck API server
	err := s.srv.Shutdown()
	require.NoError(err, "could not gracefully shutdown the quarterdeck test server")

	// Wait for server to stop to prevent race conditions
	<-s.stop

	// Cleanup logger
	logger.ResetLogger()

	// Cleanup temporary test directory
	err = os.RemoveAll(s.dbPath)
	require.NoError(err, "could not cleanup temporary database")
}

func (s *quarterdeckTestSuite) LoadDatabaseFixtures() error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("could not begin tx: %w", err)
	}

	if err := s.loadDatabaseFixtures(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *quarterdeckTestSuite) loadDatabaseFixtures(tx *sql.Tx) error {
	// Execute any SQL files in the testdata directory
	paths, err := filepath.Glob("testdata/*.sql")
	if err != nil {
		return fmt.Errorf("could not list testdata directory: %w", err)
	}

	for _, path := range paths {
		stmt, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read query from file: %w", err)
		}

		if _, err = tx.Exec(string(stmt)); err != nil {
			return fmt.Errorf("could not execute query from %s: %w", path, err)
		}
	}
	return nil
}

func (s *quarterdeckTestSuite) ResetDatabase() {
	require := s.Require()
	require.NoError(s.resetDatabase())
}

func (s *quarterdeckTestSuite) resetDatabase() (err error) {
	// Truncate all database tables except roles, permissions, and role_permissions
	// NOTE: Ensure that we delete tables in the order of foreign key relationships
	stmts := []string{
		"DELETE FROM revoked_api_keys",
		"DELETE FROM api_key_permissions",
		"DELETE FROM api_keys",
		"DELETE FROM organization_users",
		"DELETE FROM organization_projects",
		"DELETE FROM organizations",
		"DELETE FROM user_roles",
		"DELETE FROM users",
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

	// Load the fixtures back into the database
	if err = s.loadDatabaseFixtures(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *quarterdeckTestSuite) AuthContext(ctx context.Context, claims *tokens.Claims) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return api.ContextWithToken(ctx, s.srv.AccessToken(claims))
}

func (s *quarterdeckTestSuite) CheckError(err error, status int, msg string) {
	require := s.Require()

	var serr *api.StatusError
	require.True(errors.As(err, &serr), "error is not a status error")
	require.Equal(status, serr.StatusCode, "status code does not match expected status: %s", serr.Error())

	if msg != "" {
		require.Equal(msg, serr.Reply.Error, "error message does not match expected error")
	}
}

func TestQuarterdeck(t *testing.T) {
	suite.Run(t, &quarterdeckTestSuite{})
}
