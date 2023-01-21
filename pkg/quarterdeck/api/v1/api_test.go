package api_test

import (
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
)

func TestValidateCreate(t *testing.T) {
	// Create a key with all fields restricted
	key := &api.APIKey{
		ID:           ulids.New(),
		ClientID:     "foo",
		ClientSecret: "bar",
		OrgID:        ulids.New(),
		CreatedBy:    ulids.New(),
		UserAgent:    "zap",
		LastUsed:     time.Now(),
	}

	// Remove restrictions one at a time
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: id")

	key.ID = ulids.Null
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: client_id")

	key.ClientID = ""
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: client_secret")

	key.ClientSecret = ""
	require.ErrorIs(t, key.ValidateCreate(), api.ErrMissingField)
	require.EqualError(t, key.ValidateCreate(), "missing required field: name")

	key.Name = "bob"
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: org_id")

	key.OrgID = ulids.Null
	require.ErrorIs(t, key.ValidateCreate(), api.ErrMissingField)
	require.EqualError(t, key.ValidateCreate(), "missing required field: project_id")

	key.ProjectID = ulids.New()
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: created_by")

	key.CreatedBy = ulids.Null
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: user_agent")

	key.UserAgent = ""
	require.ErrorIs(t, key.ValidateCreate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateCreate(), "field restricted for request: last_used")

	key.LastUsed = time.Time{}
	require.ErrorIs(t, key.ValidateCreate(), api.ErrMissingField)
	require.EqualError(t, key.ValidateCreate(), "missing required field: permissions")

	// Remove last restriction and test valid key
	key.Permissions = []string{"foo", "barr"}
	require.NoError(t, key.ValidateCreate())
}
