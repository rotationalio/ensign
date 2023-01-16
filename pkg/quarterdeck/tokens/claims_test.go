package tokens_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/stretchr/testify/require"
)

func TestClaims(t *testing.T) {
	claims := &tokens.Claims{
		Permissions: []string{"read:foo", "write:foo", "delete:foo", "read:bar"},
	}

	// Test Permissions
	require.False(t, claims.HasPermission("write:bar"), "unexpected permission returned")
	require.True(t, claims.HasPermission("write:foo"), "expected permission to be true")
	require.False(t, claims.HasAllPermissions("write:foo", "write:bar"), "only has one permission")
	require.False(t, claims.HasAllPermissions("delete:bar", "write:bar"), "has no permissions")
	require.True(t, claims.HasAllPermissions("delete:foo", "write:foo", "read:foo"), "has all permissions")
}
