package permissions_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/stretchr/testify/require"
)

func TestRolePermissions(t *testing.T) {
	// This test ensures that each role is associated with the correct permissions to
	// make it easier to change roles and permissions in the future since the migrations
	// generally only contain private keys.
	connectDB(t)
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	require.NoError(t, err, "could not begin database transaction")
	defer tx.Rollback()

	// Ensure the expected number of roles exist
	var nRoles int
	err = tx.QueryRow("SELECT count(id) FROM roles").Scan(&nRoles)
	require.NoError(t, err, "could not count number of roles in database")
	require.Equal(t, 4, nRoles, "expected 4 roles in the database")

	testCases := []struct {
		Role        string
		Permissions []string
	}{
		{
			Role: perms.RoleOwner,
			Permissions: []string{
				perms.EditOrganizations,
				perms.DetailOrganizations,
				perms.AddCollaborators,
				perms.RemoveCollaborators,
				perms.EditCollaborators,
				perms.ReadCollaborators,
				perms.EditProjects,
				perms.DeleteProjects,
				perms.ReadProjects,
				perms.EditAPIKeys,
				perms.DeleteAPIKeys,
				perms.ReadAPIKeys,
				perms.CreateTopics,
				perms.EditTopics,
				perms.DestroyTopics,
				perms.ReadTopics,
				perms.ReadMetrics,
			},
		},
		{
			Role: perms.RoleAdmin,
			Permissions: []string{
				perms.DetailOrganizations,
				perms.AddCollaborators,
				perms.RemoveCollaborators,
				perms.EditCollaborators,
				perms.ReadCollaborators,
				perms.EditProjects,
				perms.DeleteProjects,
				perms.ReadProjects,
				perms.EditAPIKeys,
				perms.DeleteAPIKeys,
				perms.ReadAPIKeys,
				perms.CreateTopics,
				perms.EditTopics,
				perms.DestroyTopics,
				perms.ReadTopics,
				perms.ReadMetrics,
			},
		},
		{
			Role: perms.RoleMember,
			Permissions: []string{
				perms.DetailOrganizations,
				perms.ReadCollaborators,
				perms.EditProjects,
				perms.DeleteProjects,
				perms.ReadProjects,
				perms.DeleteAPIKeys,
				perms.ReadAPIKeys,
				perms.CreateTopics,
				perms.EditTopics,
				perms.DestroyTopics,
				perms.ReadTopics,
				perms.ReadMetrics,
			},
		},
		{
			Role: perms.RoleObserver,
			Permissions: []string{
				perms.ReadCollaborators,
				perms.ReadProjects,
				perms.ReadAPIKeys,
				perms.ReadTopics,
				perms.ReadMetrics,
			},
		},
	}

	for _, tc := range testCases {
		permissions := make([]string, 0)
		rows, err := tx.Query("SELECT p.name FROM role_permissions rp JOIN permissions p on rp.permission_id=p.id JOIN roles r on rp.role_id=r.id WHERE r.name=$1", tc.Role)
		require.NoError(t, err, "could not fetch permissions for role")

		for rows.Next() {
			var permission string
			err = rows.Scan(&permission)
			require.NoError(t, err, "could not fetch row from cursor")
			permissions = append(permissions, permission)
		}

		require.NoError(t, rows.Close())
		require.Len(t, permissions, len(tc.Permissions), "incorrect permissions for role %s", tc.Role)
		require.Equal(t, tc.Permissions, permissions, "incorrect permissions for role %s", tc.Role)
	}
}

func connectDB(t *testing.T) {
	dbpath := filepath.Join(t.TempDir(), "test.db")
	dsn := "sqlite:///" + dbpath

	err := db.Connect(dsn, false)
	require.NoError(t, err, "could not connect to the database")

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err, "could not close database connection")
	})
}
