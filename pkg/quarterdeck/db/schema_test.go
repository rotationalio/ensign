package db_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	migrations, err := db.Migrations()
	require.NoError(t, err, "should have been able to load migrations")
	require.GreaterOrEqual(t, len(migrations), 3, "expected at least three migrations")

	// The first three migrations should match our fixtures
	expected := []*db.Migration{
		{
			ID:   0,
			Name: "Migrations",
			Path: "0000_migrations.sql",
		},
		{
			ID:   1,
			Name: "Initial Schema",
			Path: "0001_initial_schema.sql",
		},
		{
			ID:   2,
			Name: "Default Data",
			Path: "0002_default_data.sql",
		},
	}

	for i, migration := range migrations {
		if i > len(expected) {
			break
		}

		require.Equal(t, expected[i].ID, migration.ID)
		require.Equal(t, expected[i].Name, migration.Name)
		require.Equal(t, expected[i].Path, migration.Path)

		query, err := migration.SQL()
		require.NoError(t, err, "could not load SQL from the migration")
		require.NotEmpty(t, query, "no SQL was returned for the migration")
	}
}
