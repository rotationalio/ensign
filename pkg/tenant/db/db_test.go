package db_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
)

func TestDBTestingMode(t *testing.T) {
	// Should be able to connect and close when db is in testing mode
	conf := config.DatabaseConfig{Testing: true}
	require.NoError(t, db.Connect(conf), "could not connect to database in testing mode")

	// TODO: this will no longer be true when we add a trtl mock
	require.False(t, db.IsConnected(), "expected database to not be connected in testing mode")

	require.NoError(t, db.Close(), "could not close database in testing mode")
	require.False(t, db.IsConnected(), "expected database to be not connected after close")
}
