package api_test

import (
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestRegisterRequestValidate(t *testing.T) {
	req := &api.RegisterRequest{}

	// Remove restrictions one at a time
	require.ErrorIs(t, req.Validate(), api.ErrMissingField)
	require.EqualError(t, req.Validate(), "missing required field: email")

	req.Email = "jb@example.com"
	require.ErrorIs(t, req.Validate(), api.ErrMissingField)
	require.EqualError(t, req.Validate(), "missing required field: password")

	req.Password = "password"
	require.ErrorIs(t, req.Validate(), api.ErrPasswordMismatch)

	req.PwCheck = "password"
	require.ErrorIs(t, req.Validate(), api.ErrPasswordTooWeak)

	req.Password = "super4secret"
	req.PwCheck = "super4secret"
	require.ErrorIs(t, req.Validate(), api.ErrMissingField)
	require.EqualError(t, req.Validate(), "missing required field: terms_agreement")

	req.AgreeToS = true
	require.ErrorIs(t, req.Validate(), api.ErrMissingField)
	require.EqualError(t, req.Validate(), "missing required field: privacy_agreement")

	req.AgreePrivacy = true
	require.NoError(t, req.Validate())
}

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

func TestValidateUpdate(t *testing.T) {
	// Create a key with all fields restricted
	key := &api.APIKey{
		ID:           ulids.Null,
		ClientID:     "foo",
		ClientSecret: "bar",
		OrgID:        ulids.New(),
		CreatedBy:    ulids.New(),
		ProjectID:    ulids.New(),
		Source:       "ding",
		UserAgent:    "zap",
		LastUsed:     time.Now(),
		Permissions:  []string{"foo", "bar", "baz"},
	}

	// Remove restrictions one at a time
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrMissingField)
	require.EqualError(t, key.ValidateUpdate(), "missing required field: id")

	key.ID = ulids.New()
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: client_id")

	key.ClientID = ""
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: client_secret")

	key.ClientSecret = ""
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrMissingField)
	require.EqualError(t, key.ValidateUpdate(), "missing required field: name")

	key.Name = "bob"
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: org_id")

	key.OrgID = ulids.Null
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: project_id")

	key.ProjectID = ulids.Null
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: created_by")

	key.CreatedBy = ulids.Null
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: source")

	key.Source = ""
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: user_agent")

	key.UserAgent = ""
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: last_used")

	key.LastUsed = time.Time{}
	require.ErrorIs(t, key.ValidateUpdate(), api.ErrRestrictedField)
	require.EqualError(t, key.ValidateUpdate(), "field restricted for request: permissions")

	// Remove last restriction and test valid key
	key.Permissions = nil
	require.NoError(t, key.ValidateUpdate())
}

func TestValidateUserUpdate(t *testing.T) {
	// create empty User object
	user := &api.User{}

	// Check for missing fields
	require.ErrorIs(t, user.ValidateUpdate(), api.ErrMissingField)
	require.EqualError(t, user.ValidateUpdate(), "missing required field: user_id")

	userID := ulids.New()
	user.UserID = userID
	require.ErrorIs(t, user.ValidateUpdate(), api.ErrMissingField)
	require.EqualError(t, user.ValidateUpdate(), "missing required field: name")

	name := "Sonali Mehra"
	user.Name = name
	require.NoError(t, user.ValidateUpdate())
}
