package models_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// The models test suite sets up a SQLite database in a temporary directory and connects
// to it with the db package in single node mode (e.g. no replication). The database
// is loaded with fixtures specified by SQL files in the testdata directory and can be
// reset if changes are made to it.
type modelTestSuite struct {
	suite.Suite
	dbpath string
}

func (m *modelTestSuite) SetupSuite() {
	m.CreateDB()
}

// Creates a SQLite database with the current schema and the db module connected in a
// temporary directory that will be cleaned up when the tests are complete.
func (m *modelTestSuite) CreateDB() {
	require := m.Require()

	// Only create the database path on the first call to CreateDB. Otherwise the call
	// to TempDir() will be prefixed with the name of the subtest, which will cause an
	// "attempt to write a read-only database" for subsequent tests because the directory
	// will be deleted when the subtest is complete.
	if m.dbpath == "" {
		m.dbpath = filepath.Join(m.T().TempDir(), "testdb.sqlite")
	}

	err := db.Connect("sqlite:///"+m.dbpath, false)
	require.NoError(err, "could not connect to temporary database")

	// Execute any SQL files in the testdata directory
	paths, err := filepath.Glob("testdata/*.sql")
	require.NoError(err, "could not list testdata directory")

	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(err, "could not begin tx")

	for _, path := range paths {
		stmt, err := os.ReadFile(path)
		require.NoError(err, "could not read query from file")
		_, err = tx.Exec(string(stmt))
		require.NoError(err, "could not execute sql query from testdata")
	}

	tx.Commit()
}

// Closes the connection to the current database and connects to a new database.
func (m *modelTestSuite) ResetDB() {
	require := m.Require()
	require.NoError(db.Close(), "could not close connection to db")
	require.NoError(os.Remove(m.dbpath), "could not delete old db")
	m.CreateDB()
}

func TestModels(t *testing.T) {
	suite.Run(t, new(modelTestSuite))
}

type TestModel struct {
	models.Base
}

func TestBaseEmbedding(t *testing.T) {
	m := &TestModel{
		Base: models.Base{
			Created:  "2022-12-05T14:02:32Z",
			Modified: "2022-12-05T16:27:18Z",
		},
	}

	ts, err := m.GetCreated()
	require.NoError(t, err, "could not get the created timestamp")
	require.True(t, ts.Equal(time.Date(2022, 12, 5, 14, 2, 32, 0, time.UTC)))

	ts, err = m.GetModified()
	require.NoError(t, err, "could not get the modified timestamp")
	require.True(t, ts.Equal(time.Date(2022, 12, 5, 16, 27, 18, 0, time.UTC)))

	now := time.Now()
	m.SetCreated(now)
	m.SetModified(now)

	ts, err = m.GetCreated()
	require.NoError(t, err, "could not get the created timestamp")
	require.True(t, ts.Equal(now))

	ts, err = m.GetModified()
	require.NoError(t, err, "could not get the modified timestamp")
	require.True(t, ts.Equal(now))
}
