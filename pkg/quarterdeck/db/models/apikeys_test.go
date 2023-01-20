package models_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestGetAPIKey() {
	require := m.Require()

	// Test get by client ID
	apikey, err := models.GetAPIKey(context.Background(), "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa")
	require.NoError(err, "could not fetch api key by client ID")
	require.NotNil(apikey)
	require.Equal("01GME02TJP2RRP39MKR525YDQ6", apikey.ID.String())
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
