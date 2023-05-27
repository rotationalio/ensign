package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/sqlite"
	"github.com/stretchr/testify/require"
)

func TestDriver(t *testing.T) {
	db, err := sql.Open("ensign_sqlite3", "testdata/test.db")
	require.NoError(t, err, "could not open connection to testdb")

	conn, err := db.Conn(context.Background())
	require.NoError(t, err, "could not create connection with custom driver")
	require.Equal(t, 1, sqlite.NumConns())

	// Get the underlying sqlite3 connection
	sqlc, ok := sqlite.GetLastConn()
	require.True(t, ok, "connection was not in connection map?")
	require.IsType(t, &sqlite.Conn{}, sqlc, "connection of wrong type returned")

	err = conn.Close()
	require.NoError(t, err, "could not close connection")

	err = db.Close()
	require.NoError(t, err, "could not close database")
	require.Equal(t, 0, sqlite.NumConns())
}
