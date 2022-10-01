package db_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/stretchr/testify/require"
)

func TestConnectClose(t *testing.T) {
	dsn := "sqlite3:///" + filepath.Join(t.TempDir(), "test.db")

	_, err := db.BeginTx(context.Background(), nil)
	require.ErrorIs(t, err, db.ErrNotconnected, "should not be able to open a transaction without connecting")

	err = db.Close()
	require.NoError(t, err, "should be able to close the db without error when not connected")

	// Connect to the DB
	err = db.Connect(dsn, false)
	require.NoError(t, err, "could not connect to the db")

	err = db.Connect(dsn, false)
	require.NoError(t, err, "multiple connects should not cause an error")

	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err, "could not create transaction")

	require.NoError(t, tx.Rollback(), "could not abort transaction")
	require.NoError(t, db.Close(), "could not close db")

	require.NoError(t, db.Connect(dsn, false), "could not reconnect to the db")
	require.NoError(t, db.Close(), "could not close db")
}

func TestReadOnly(t *testing.T) {
	// Ensure the DB is closed so it opens in readonly mode.
	require.NoError(t, db.Close(), "could not close database")

	// Connect to the DB in readonly mode
	dsn := "sqlite3:///" + filepath.Join(t.TempDir(), "test.db")
	require.NoError(t, db.Connect(dsn, true), "could not connect to db")

	_, err := db.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: false})
	require.ErrorIs(t, err, db.ErrReadOnly, "should not be able to open a write tx in read only mode")

	_, err = db.BeginTx(context.Background(), nil)
	require.NoError(t, err, "could not create transaction from nil tx options")

	require.NoError(t, db.Close(), "could not close database")
}

func TestParseDSN(t *testing.T) {
	testCases := []struct {
		uri    string
		scheme string
		path   string
		err    error
	}{
		{"sqlite3:///path/to/test.db", "sqlite3", "path/to/test.db", nil},
		{"sqlite3:////absolute/path/test.db", "sqlite3", "/absolute/path/test.db", nil},
		{"", "", "", db.ErrCannotParseDSN},
		{"sqlite3://", "", "", db.ErrCannotParseDSN},
	}

	for _, tc := range testCases {
		dsn, err := db.ParseDSN(tc.uri)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected dsn parsing error")
			continue
		}

		require.Equal(t, tc.scheme, dsn.Scheme)
		require.Equal(t, tc.path, dsn.Path)
	}
}
