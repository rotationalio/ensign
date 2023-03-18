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
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestListAPIKey() {
	require := m.Require()

	ctx := context.Background()
	orgID := ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	projectID := ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT")

	keys, cursor, err := models.ListAPIKeys(ctx, ulids.Null, ulids.Null, nil)
	require.ErrorIs(err, models.ErrMissingOrgID, "orgID is required for list queries")
	require.Nil(cursor)
	require.Nil(keys)

	_, _, err = models.ListAPIKeys(ctx, orgID, projectID, &pagination.Cursor{})
	require.ErrorIs(err, models.ErrMissingPageSize, "pagination is required for list queries")

	// Should return all example apikeys in both projects (page cursor not required)
	keys, cursor, err = models.ListAPIKeys(ctx, orgID, ulids.Null, nil)
	require.NoError(err, "could not fetch all apikeys for example org")
	require.Nil(cursor, "should be no next page so no cursor")
	require.Len(keys, 11, "expected 11 keys returned 2 from the birds project and 9 from the test project")

	// Keys should be returned in descending order by ID
	var prevID ulid.ULID
	for _, k := range keys {
		if !ulids.IsZero(prevID) {
			require.True(k.ID.String() < prevID.String(), "expected keys to be sorted by ID descending")
		}
		prevID = k.ID
	}

	// Should return example apikeys in the specified project (page cursor not required)
	keys, cursor, err = models.ListAPIKeys(ctx, orgID, projectID, nil)
	require.NoError(err, "could not fetch project apikeys for example org")
	require.Nil(cursor, "should be no next page so no cursor")
	require.Len(keys, 2, "expected 2 keys returned from the birds project")
}

func (m *modelTestSuite) TestListAPIKeyPagination() {
	require := m.Require()
	ctx := context.Background()
	orgID := ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	projectID := ulid.MustParse("01GQFR0KM5S2SSJ8G5E086VQ9K")

	// Test pagination without a project
	pages := 0
	nRows := 0
	cursor := pagination.New("", "", 3)
	for cursor != nil && pages < 100 {
		keys, nextPage, err := models.ListAPIKeys(ctx, orgID, ulids.Null, cursor)
		require.NoError(err, "could not fetch page from server")

		// Ensure that all keys in this page are sorted by ID descending
		var prevID ulid.ULID
		for _, k := range keys {
			if !ulids.IsZero(prevID) {
				require.True(k.ID.String() < prevID.String(), "expected page keys to be sorted by ID descending")
			}
			prevID = k.ID
		}

		if nextPage != nil {
			require.NotEqual(cursor.StartIndex, nextPage.StartIndex)
			require.NotEqual(cursor.EndIndex, nextPage.EndIndex)
			require.Equal(cursor.PageSize, nextPage.PageSize)
		}

		pages++
		nRows += len(keys)
		cursor = nextPage
	}

	require.Equal(4, pages, "expected 11 results in 4 pages")
	require.Equal(11, nRows, "expected 11 results in 4 pages")

	// Test pagination with a project
	pages = 0
	nRows = 0
	cursor = pagination.New("", "", 3)
	for cursor != nil && pages < 100 {
		keys, nextPage, err := models.ListAPIKeys(ctx, orgID, projectID, cursor)
		require.NoError(err, "could not fetch page from server")

		// Ensure that all keys in this page are sorted by ID descending
		var prevID ulid.ULID
		for _, k := range keys {
			if !ulids.IsZero(prevID) {
				require.True(k.ID.String() < prevID.String(), "expected page keys to be sorted by ID descending")
			}
			prevID = k.ID
		}

		pages++
		nRows += len(keys)
		cursor = nextPage
	}

	require.Equal(3, pages, "expected 9 results in 3 pages")
	require.Equal(9, nRows, "expected 9 results in 3 pages")
}

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
	require.Len(permissions, 7)

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
	require.Len(permissions, 7)

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
	require.Equal(`["topics:create","topics:edit","topics:destroy","topics:read","metrics:read","publisher","subscriber"]`, permissions, "permissions not serialized correctly")
}

func (m *modelTestSuite) TestCreateAPIKey() {
	defer m.ResetDB()
	require := m.Require()

	// Create an API key with minimal information
	apikey := &models.APIKey{
		Name:      "Testing API Key",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
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
	require.True(apikey.Partial)

	// Fetch the apikey key from the database
	cmpt, err := models.GetAPIKey(context.Background(), apikey.KeyID)
	require.NoError(err, "no model was created in the database")
	require.NotSame(apikey, cmpt, "something went wrong with the tests")
	require.Equal(apikey, cmpt, "fetched model not identical to saved model")

	expectedPermissions, _ := apikey.Permissions(context.Background(), false)
	actualPermissions, _ := apikey.Permissions(context.Background(), false)
	require.Equal(expectedPermissions, actualPermissions, "permissions not saved to database")

	// Test that partial flag is not set on an APIKey with all permissions
	apikey = &models.APIKey{
		Name:      "Full Permissions",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
	}
	apikey.SetPermissions("topics:create", "topics:edit", "topics:destroy", "topics:read", "metrics:read", "publisher", "subscriber")
	err = apikey.Create(context.Background())
	require.NoError(err, "could not create a valid apikey")
	require.False(apikey.Partial, "partial flag should not be set on an APIKey with all permissions")

	// Should not be able to create an APIKey for a project not associated with the orgID
	apikey = &models.APIKey{
		Name:      "Invalid Project",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
	}
	apikey.SetPermissions("publisher", "subscriber")
	apikey.OrgID = ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")
	err = apikey.Create(context.Background())
	require.ErrorIs(err, models.ErrInvalidProjectID)
}

func (m *modelTestSuite) TestCreateAPIKeyUnknownPermission() {
	defer m.ResetDB()
	require := m.Require()

	// Should not be able to create an API key with an unknown permission
	apikey := &models.APIKey{
		Name:      "Unknown Permissions",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
	}
	apikey.SetPermissions("publisher", "subscriber", "notapermission")
	err := apikey.Create(context.Background())
	require.ErrorIs(err, models.ErrInvalidPermission, "expected error when creating an APIKey with a permission not in the database")
}

func (m *modelTestSuite) TestCreateAPIKeyInvalidPermission() {
	defer m.ResetDB()
	require := m.Require()

	// Assert that the organizations:delete permission is not an api key permission
	permission, err := models.GetPermission(context.Background(), "organizations:delete")
	require.NoError(err, "could not fetch permission from the database")
	require.False(permission.AllowAPIKeys, "expected the organizations:delete permission to be not allowed for api keys")

	// Should not be able to create an API key with a permission that is not allowed for api keys
	apikey := &models.APIKey{
		Name:      "Invalid Permissions",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
	}
	apikey.SetPermissions("publisher", "subscriber", "organizations:delete", "topics:create")
	err = apikey.Create(context.Background())
	require.ErrorIs(err, models.ErrInvalidPermission, "expected error when creating an APIKey with a permission is not allowed")
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
	require.Len(permissions, 7)

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
	require.ErrorIs(apikey.Validate(), models.ErrMissingCreatedBy)

	apikey.CreatedBy = ulids.New()
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
	require.Len(permissions, 7)
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

func (m *modelTestSuite) TestStatus() {
	require := m.Require()

	// If LastUsed is not set then the key is unused
	apikey := &models.APIKey{}
	require.Equal(models.APIKeyStatusUnused, apikey.Status())

	apikey = &models.APIKey{}
	apikey.SetLastUsed(time.Time{})
	require.Equal(models.APIKeyStatusUnused, apikey.Status())

	// If not used recently then the key is stale
	apikey = &models.APIKey{}
	apikey.SetLastUsed(time.Now().Add(-time.Hour * 24 * 30 * 4))
	require.Equal(models.APIKeyStatusStale, apikey.Status())

	// If used recently then the key is active
	apikey = &models.APIKey{}
	apikey.SetLastUsed(time.Now().Add(-time.Hour * 24))
	require.Equal(models.APIKeyStatusActive, apikey.Status())
}

func (m *modelTestSuite) TestAPIKeyPArtialPermissionsQuery() {
	defer m.ResetDB()
	require := m.Require()

	query := `SELECT EXISTS (SELECT p.id FROM permissions p WHERE p.allow_api_keys=true EXCEPT SELECT kp.permission_id FROM api_key_permissions kp WHERE kp.api_key_id=:id)`

	// Create an API Key with partial permissions
	partialKey := &models.APIKey{
		Name:      "Partial Permissions",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
	}
	partialKey.SetPermissions("publisher", "subscriber", "topics:create")
	require.NoError(partialKey.Create(context.Background()))

	// Create an API Key with full permissions
	fullKey := &models.APIKey{
		Name:      "Partial Permissions",
		OrgID:     ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
		ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		CreatedBy: ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"),
	}
	fullKey.SetPermissions("publisher", "subscriber", "topics:create", "topics:destroy", "topics:edit", "topics:read", "metrics:read")
	require.NoError(fullKey.Create(context.Background()))

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	require.NoError(err, "could not start transaction")

	err = tx.QueryRow(query, sql.Named("id", partialKey.ID)).Scan(&partialKey.Partial)
	require.NoError(err, "could not execute query for partial key")

	err = tx.QueryRow(query, sql.Named("id", fullKey.ID)).Scan(&fullKey.Partial)
	require.NoError(err, "could not execute query for full key")

	require.True(partialKey.Partial, "expected partial key partial to be true")
	require.False(fullKey.Partial, "expected full key partial to be false")

	tx.Rollback()
}
