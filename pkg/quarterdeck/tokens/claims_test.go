package tokens_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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

func TestClaimsProjectID(t *testing.T) {
	claims := &tokens.Claims{}
	require.False(t, claims.ValidateProject(ulids.New()), "empty project ID should not validate")

	claims.ProjectID = "foo"
	require.False(t, claims.ValidateProject(ulids.New()), "invalid project ID should not validate")

	claims.ProjectID = ulids.MustParse("01GTW1R9MH8723JQDRMFE16CZ7").String()
	require.False(t, claims.ValidateProject(ulids.New()), "incorrect project ID should not validate")

	require.True(t, claims.ValidateProject(ulids.MustParse("01GTW1R9MH8723JQDRMFE16CZ7")), "correct project ID should be valid")
}

func TestClaimsParseOrgID(t *testing.T) {
	claims := &tokens.Claims{}
	require.Equal(t, ulids.Null, claims.ParseOrgID())

	claims.OrgID = "notvalid"
	require.Equal(t, ulids.Null, claims.ParseOrgID())

	orgID := ulids.New()
	claims.OrgID = orgID.String()
	require.Equal(t, orgID, claims.ParseOrgID())
}

func TestClaimsParseUserID(t *testing.T) {
	claims := &tokens.Claims{}
	require.Equal(t, ulids.Null, claims.ParseUserID())

	claims.Subject = "notvalid"
	require.Equal(t, ulids.Null, claims.ParseUserID())

	userID := ulids.New()
	claims.Subject = userID.String()
	require.Equal(t, userID, claims.ParseUserID())
}
