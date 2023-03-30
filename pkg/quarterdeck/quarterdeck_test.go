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
	"github.com/rotationalio/ensign/pkg/utils/emails"
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
	// Note use assert instead of require so that go routines are properly handled in
	// tests; assert uses t.Error while require uses t.FailNow and multiple go routines
	// might lead to incorrect testing behavior.
	assert := s.Assert()
	s.stop = make(chan bool)

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Create a temporary test database for the tests
	var err error
	s.dbPath, err = os.MkdirTemp("", "quarterdeck-*")
	assert.NoError(err, "could not create temporary directory for database")

	// Create a test configuration to run the Quarterdeck API server as a fully
	// functional server on an open port using the local-loopback for networking.
	s.conf, err = config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost:3000"},
		EmailURL: config.URLConfig{
			Base:   "http://localhost:3000",
			Invite: "/invite",
			Verify: "/verify",
		},
		SendGrid: emails.Config{
			FromEmail:  "quarterdeck@rotational.io",
			AdminEmail: "admins@rotationa.io",
			Testing:    true,
		},
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
			RefreshOverlap:  -10 * time.Minute,
		},
		RateLimit: config.RateLimitConfig{
			PerSecond: 20.00,
			Burst:     120,
			TTL:       5 * time.Minute,
		},
	}.Mark()
	assert.NoError(err, "test configuration is invalid")

	s.srv, err = quarterdeck.New(s.conf)
	assert.NoError(err, "could not create the quarterdeck api server from the test configuration")

	// Start the Quarterdeck server - the goal of the tests is to have the server run
	// for the entire duration of the tests. Implement reset methods to ensure the
	// server state doesn't change between tests in Before/After.
	go func() {
		if err := s.srv.Serve(); err != nil {
			s.T().Logf("error occurred during service: %s", err)
		}
		s.stop <- true
	}()

	// Wait for 500ms to ensure the API server starts up
	time.Sleep(500 * time.Millisecond)

	// Load database fixtures
	assert.NoError(s.LoadDatabaseFixtures(), "could not load database fixtures")

	// Create a Quarterdeck client for making requests to the server
	assert.NotEmpty(s.srv.URL(), "no url to connect the client on")
	s.client, err = api.New(s.srv.URL())
	assert.NoError(err, "could not initialize the Quarterdeck client")
}

func (s *quarterdeckTestSuite) TearDownSuite() {
	assert := s.Assert()

	// Shutdown the quarterdeck API server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.srv.Shutdown(ctx)
	assert.NoError(err, "could not gracefully shutdown the quarterdeck test server")

	// Wait for server to stop to prevent race conditions
	<-s.stop

	// Cleanup logger
	logger.ResetLogger()

	// Cleanup temporary test directory
	err = os.RemoveAll(s.dbPath)
	assert.NoError(err, "could not cleanup temporary database")
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
	s.Assert().NoError(s.resetDatabase())
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

// Stop the task manager, waiting for all the tasks to finish. Tests should defer
// ResetTasks() to ensure that the task manager is available to the other tests.
func (s *quarterdeckTestSuite) StopTasks() {
	tasks := s.srv.GetTaskManager()
	tasks.Stop()
}

// Reset the task manager to ensure that other tests have access to it.
func (s *quarterdeckTestSuite) ResetTasks() {
	s.srv.ResetTaskManager()
}

func (s *quarterdeckTestSuite) AuthContext(ctx context.Context, claims *tokens.Claims) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return api.ContextWithToken(ctx, s.srv.AccessToken(claims))
}

func (s *quarterdeckTestSuite) CheckError(err error, status int, msg string) {
	require := s.Require()
	require.Error(err, "expected an error but didn't get one")

	var serr *api.StatusError
	require.True(errors.As(err, &serr), "error is not a status error: %v", err)
	require.Equal(status, serr.StatusCode, "status code does not match expected status: %s", serr.Error())

	if msg != "" {
		require.Equal(msg, serr.Reply.Error, "error message does not match expected error")
	}
}

func TestQuarterdeck(t *testing.T) {
	suite.Run(t, &quarterdeckTestSuite{})
}
