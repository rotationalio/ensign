package models_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestGetAPIKey() {
	require := m.Require()

	// Test get by client ID
	apikey, err := models.GetAPIKey(context.Background(), "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa")
	require.NoError(err, "could not fetch api key by client ID")
	require.NotNil(apikey)

	// Ensure the model is fully populated
	require.Equal("01GME02TJP2RRP39MKR525YDQ6", apikey.ID.String())
	require.Equal("DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa", apikey.KeyID)
	require.Equal("$argon2id$v=19$m=65536,t=1,p=2$5tE7XLSdqM36DUmzeSppvA==$eTfRYSCuBssAcuxxFv/eh92CyL1NuNqBPkhlLoIAVAw=", apikey.Secret)
	require.Equal("Eagle Publishers", apikey.Name)
	require.Equal(ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"), apikey.OrgID)
	require.Equal(ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"), apikey.ProjectID)
	require.Equal(ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"), apikey.CreatedBy)
	require.Equal("Beacon UI", apikey.Source.String)
	require.Equal("Quarterdeck API/v1", apikey.UserAgent.String)
	require.Equal("2023-01-22T13:26:25.394129Z", apikey.LastUsed.String)

	permissions, err := apikey.Permissions(context.Background(), false)
	require.NoError(err)
	require.Len(permissions, 5)

	// Ensure GetAPIKey returns not found
	apikey, err = models.GetAPIKey(context.Background(), keygen.KeyID())
	require.ErrorIs(err, models.ErrNotFound)
	require.Nil(apikey)
}

func (m *modelTestSuite) TestRetrieveAPIKey() {
	require := m.Require()

	// Test get by ulid
	apikey, err := models.RetrieveAPIKey(context.Background(), ulid.MustParse("01GME02TJP2RRP39MKR525YDQ6"))
	require.NoError(err, "could not fetch api key by id")
	require.NotNil(apikey)

	// Ensure the model is fully populated
	require.Equal("01GME02TJP2RRP39MKR525YDQ6", apikey.ID.String())
	require.Equal("DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa", apikey.KeyID)
	require.Empty(apikey.Secret, "client secret should not be returned on retrieve")
	require.Equal("Eagle Publishers", apikey.Name)
	require.Equal(ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"), apikey.OrgID)
	require.Equal(ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"), apikey.ProjectID)
	require.Equal(ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"), apikey.CreatedBy)
	require.Equal("Beacon UI", apikey.Source.String)
	require.Equal("Quarterdeck API/v1", apikey.UserAgent.String)
	require.Equal("2023-01-22T13:26:25.394129Z", apikey.LastUsed.String)

	permissions, err := apikey.Permissions(context.Background(), false)
	require.NoError(err)
	require.Len(permissions, 5)

	// Ensure RetrieveAPIKey returns not found
	apikey, err = models.RetrieveAPIKey(context.Background(), ulids.New())
	require.ErrorIs(err, models.ErrNotFound)
	require.Nil(apikey)
}

func (m *modelTestSuite) TestDeleteAPIKey() {
	defer m.ResetDB()
	require := m.Require()

	keyID := ulid.MustParse("01GME02TJP2RRP39MKR525YDQ6")
	orgID := ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")

	// Should not be able to delete a key with the wrong organization
	err := models.DeleteAPIKey(context.Background(), keyID, ulids.New())
	require.ErrorIs(err, models.ErrNotFound)

	// Should not be able to delete a key that is not found
	err = models.DeleteAPIKey(context.Background(), ulids.New(), orgID)
	require.ErrorIs(err, models.ErrNotFound)

	// Should be able to delete a key
	err = models.DeleteAPIKey(context.Background(), keyID, orgID)
	require.NoError(err)

	// Should not be able to retrieve a key once its deleted
	_, err = models.RetrieveAPIKey(context.Background(), keyID)
	require.ErrorIs(err, models.ErrNotFound)

	// Key should be in revoked keys database
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	require.NoError(err, "could not create transaction")
	defer tx.Rollback()

	var permissions string
	err = tx.QueryRow("SELECT permissions FROM revoked_api_keys WHERE id=$1 AND organization_id=$2", keyID, orgID).Scan(&permissions)
	require.NoError(err, "could not fetched revoked key")
	require.Equal(`["topics:create","topics:read","metrics:read","publisher","subscriber"]`, permissions, "permissions not serialized correctly")
}

func (m *modelTestSuite) TestCreateAPIKey() {
	defer m.ResetDB()
	require := m.Require()

	// Create an API key with minimal information
	apikey := &models.APIKey{
		Name:      "Testing API Key",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulids.New(),
	}
	apikey.SetPermissions("publisher", "subscriber")

	err := apikey.Create(context.Background())
	require.NoError(err, "could not create a valid apikey")

	// Ensure that the model was populated correctly
	require.False(ulids.IsZero(apikey.ID), "no API key was set")
	require.NotZero(apikey.KeyID)
	require.NotZero(apikey.Secret)
	require.NotZero(apikey.Created)
	require.NotZero(apikey.Modified)

	// Fetch the apikey key from the database
	cmpt, err := models.GetAPIKey(context.Background(), apikey.KeyID)
	require.NoError(err, "no model was created in the database")
	require.NotSame(apikey, cmpt, "something went wrong with the tests")
	require.Equal(apikey, cmpt, "fetched model not identical to saved model")

	expectedPermissions, _ := apikey.Permissions(context.Background(), false)
	actualPermissions, _ := apikey.Permissions(context.Background(), false)
	require.Equal(expectedPermissions, actualPermissions, "permissions not saved to database")
}

func (m *modelTestSuite) TestUpdateAPIKey() {
	defer m.ResetDB()
	require := m.Require()

	// Test workflow of updating a key from scratch (e.g. as the API handler does it)
	// without retrieving the key from the database.
	// Cannot update a key without an ID
	key := &models.APIKey{}
	err := key.Update(context.Background())
	require.ErrorIs(err, models.ErrMissingModelID)

	// Cannot update a key without an orgID
	key.ID = ulid.MustParse("01GME02TJP2RRP39MKR525YDQ6")
	err = key.Update(context.Background())
	require.ErrorIs(err, models.ErrMissingOrgID)

	// Cannot update a key without a name
	key.OrgID = ulids.New()
	err = key.Update(context.Background())
	require.ErrorIs(err, models.ErrMissingKeyName)

	// Cannot update a key witout the correct orgID (important for security)
	key.Name = "not the original name"
	err = key.Update(context.Background())
	require.ErrorIs(err, models.ErrNotFound)

	// Should be able to update the key with the correct orgID
	key.OrgID = ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	err = key.Update(context.Background())
	require.NoError(err, "could not update a valid key")

	// Ensure the model was populated on update
	require.Equal("01GME02TJP2RRP39MKR525YDQ6", key.ID.String())
	require.Equal("DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa", key.KeyID)
	require.Empty(key.Secret, "client secret should not be returned on retrieve")
	require.Equal("not the original name", key.Name)
	require.Equal(ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"), key.OrgID)
	require.Equal(ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"), key.ProjectID)
	require.Equal(ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"), key.CreatedBy)
	require.Equal("Beacon UI", key.Source.String)
	require.Equal("Quarterdeck API/v1", key.UserAgent.String)
	require.Equal("2023-01-22T13:26:25.394129Z", key.LastUsed.String)
	require.NotEmpty(key.Created)
	require.NotEmpty(key.Modified)

	permissions, err := key.Permissions(context.Background(), false)
	require.NoError(err)
	require.Len(permissions, 5)

	// Ensure the modified timestamp was set
	modified, err := key.GetModified()
	require.NoError(err)
	require.LessOrEqual(time.Since(modified), 1*time.Second)

	// Retrieve the key from the database and make sure it matches
	cmpt, err := models.RetrieveAPIKey(context.Background(), key.ID)
	require.NoError(err, "could not retrieve row from database")
	require.Equal(key, cmpt)

	// Ensure we can update the cmpt model after retrieved from the db
	cmpt.Name = "changed yet again"
	err = cmpt.Update(context.Background())
	require.NoError(err, "could not update a fully populated model")

	// Ensure Update returns not found if the key is not in the database
	key.ID = ulids.New()
	err = key.Update(context.Background())
	require.ErrorIs(err, models.ErrNotFound)
}

func (m *modelTestSuite) TestAPIKeyValidation() {
	require := m.Require()

	// Empty model is not valid
	apikey := &models.APIKey{}
	require.ErrorIs(apikey.Validate(), models.ErrMissingModelID)

	// KeyID and Secret is required
	apikey.ID = ulids.New()
	require.ErrorIs(apikey.Validate(), models.ErrMissingKeyMaterial)

	// Secret must be a derived key
	apikey.KeyID = keygen.KeyID()
	apikey.Secret = keygen.Secret()
	require.ErrorIs(apikey.Validate(), models.ErrInvalidSecret)

	// Name is required
	apikey.Secret, _ = passwd.CreateDerivedKey("supersecret")
	require.ErrorIs(apikey.Validate(), models.ErrMissingKeyName)

	// OrganizationID is required
	apikey.Name = "testing123"
	require.ErrorIs(apikey.Validate(), models.ErrMissingOrgID)

	// ProjectID is required
	apikey.OrgID = ulids.New()
	require.ErrorIs(apikey.Validate(), models.ErrMissingProjectID)

	// Permissions are required
	apikey.ProjectID = ulids.New()
	require.ErrorIs(apikey.Validate(), models.ErrNoPermissions)

	// Valid API Key
	apikey.ID = ulids.Null
	apikey.AddPermissions("foo:read", "foo:write")
	apikey.ID = ulids.New()
	require.NoError(apikey.Validate())
}

func (m *modelTestSuite) TestAPIKeyUpdateLastSeen() {
	defer m.ResetDB()

	require := m.Require()
	apikey, err := models.GetAPIKey(context.Background(), "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa")
	require.NoError(err, "could not fetch api key by client ID")

	// The apikey pointer will be modified so get a second copy for comparison
	prev, err := models.GetAPIKey(context.Background(), "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa")
	require.NoError(err, "could not fetch api key by client ID")

	err = apikey.UpdateLastUsed(context.Background())
	require.NoError(err, "could not update last used: %+v", err)

	// Fetch the record from the database for comparison purposes.
	cmpr, err := models.GetAPIKey(context.Background(), "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa")
	require.NoError(err, "could not fetch api key by client ID")

	// Nothing but last used and modified should have changed.
	require.Equal(prev.ID, cmpr.ID)
	require.Equal(prev.KeyID, cmpr.KeyID)
	require.Equal(prev.Secret, cmpr.Secret)
	require.Equal(prev.Name, cmpr.Name)
	require.Equal(prev.ProjectID, cmpr.ProjectID)
	require.Equal(prev.CreatedBy, cmpr.CreatedBy)
	require.Equal(prev.Created, cmpr.Created)

	// Last Used and Modified should have changed to the same timestamp
	require.Equal(cmpr.LastUsed.String, cmpr.Modified, "expected modified and last used to be equal")
	require.NotEqual(prev.LastUsed.String, cmpr.LastUsed.String)
	require.NotEqual(prev.Modified, cmpr.Modified)

	// The pointer should have been updated to match what's in the database
	require.Equal(apikey.LastUsed.String, cmpr.LastUsed.String)
	require.Equal(apikey.Modified, cmpr.Modified)

	// Last Used and Modified should be after the previous Last Used and Modified
	ll, err := cmpr.GetLastUsed()
	require.NoError(err, "could not parse last used")
	require.False(ll.IsZero())

	pll, err := prev.GetLastUsed()
	require.NoError(err, "could not parse last used fixture")
	require.True(ll.After(pll), "cmpr last used %q is not after prev last used %q", cmpr.LastUsed.String, prev.LastUsed.String)

	mod, err := cmpr.GetModified()
	require.NoError(err, "could not parse modified")
	require.False(mod.IsZero())

	pmod, err := prev.GetModified()
	require.NoError(err, "could not parse modified fixture")
	require.True(mod.After(pmod), "cmpr modified %q is not after prev modified %q", cmpr.Modified, prev.Modified)
}

func TestAPIKeyLastSeen(t *testing.T) {
	apikey := &models.APIKey{}

	ts, err := apikey.GetLastUsed()
	require.NoError(t, err, "could not get null last used")
	require.Zero(t, ts, "expected zero-valued timestamp")

	now := time.Now()
	apikey.SetLastUsed(now)

	ts, err = apikey.GetLastUsed()
	require.NoError(t, err, "could not get non-null last used")
	require.True(t, now.Equal(ts))
}

func (m *modelTestSuite) TestAPIKeyPermissions() {
	require := m.Require()

	// Create a user with only a user ID
	apikey := &models.APIKey{ID: ulid.MustParse("01GME02TJP2RRP39MKR525YDQ6")}

	// Fetch the permissions for the user
	permissions, err := apikey.Permissions(context.Background(), false)
	require.NoError(err, "could not fetch permissions for api key")
	require.Len(permissions, 5)
}

func (m *modelTestSuite) TestAPIKeyAddSetPermissions() {
	require := m.Require()

	// Should not be able to add or set permissions to an existing APIKey.
	apikey, err := models.GetAPIKey(context.Background(), "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa")
	require.NoError(err, "could not fetch api key by client ID")

	err = apikey.AddPermissions("read:foo", "write:foo", "delete:foo")
	require.ErrorIs(err, models.ErrModifyPermissions)

	err = apikey.SetPermissions("read:foo", "write:foo", "delete:foo")
	require.ErrorIs(err, models.ErrModifyPermissions)

	// Should be able to add permissions to a new APIKey
	apikey = &models.APIKey{}
	require.NoError(apikey.AddPermissions("read:foo", "write:foo", "delete:foo"))
	perms, _ := apikey.Permissions(context.Background(), false)
	require.Len(perms, 3)

	require.NoError(apikey.AddPermissions("read:bar", "write:bar"))
	perms, _ = apikey.Permissions(context.Background(), false)
	require.Len(perms, 5)

	// SetPermissions should overwrite the old permissions
	require.NoError(apikey.SetPermissions("topics", "publisher"))
	perms, _ = apikey.Permissions(context.Background(), false)
	require.Len(perms, 2)

	// should be able to set permissions on a new APIKey
	apikey = &models.APIKey{}
	require.NoError(apikey.SetPermissions("read:foo", "write:foo", "delete:foo"))
	perms, _ = apikey.Permissions(context.Background(), false)
	require.Len(perms, 3)

	// add permissions should not have duplicates even when already set on the key
	require.NoError(apikey.AddPermissions("read:foo", "write:foo", "delete:foo", "delete:foo", "write:foo", "read:foo", "write:foo"))
	perms, _ = apikey.Permissions(context.Background(), false)
	require.Len(perms, 3)

	// set permissions should not have duplicates
	require.NoError(apikey.SetPermissions("read:foo", "write:foo", "delete:foo", "write:foo", "read:foo", "read:foo", "delete:foo"))
	perms, _ = apikey.Permissions(context.Background(), false)
	require.Len(perms, 3)
}
