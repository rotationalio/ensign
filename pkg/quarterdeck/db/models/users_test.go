package models_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"

	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestGetUser() {
	require := m.Require()

	// Test get by ID string
	user, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.NoError(err, "could not fetch user by string ID")
	require.NotNil(user)
	require.Equal("01GKHJSK7CZW0W282ZN3E9W86Z", user.ID.String())
	require.Equal("Jannel P. Hudson", user.Name)

	// Test get by ULID
	user2, err := models.GetUser(context.Background(), ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"))
	require.NoError(err, "could not fetch user by ulid")
	require.Equal("01GKHJSK7CZW0W282ZN3E9W86Z", user2.ID.String())
	require.Equal(user, user2)

	// Ensure we cannot fetch a user by integer
	_, err = models.GetUser(context.Background(), 1)
	require.Error(err, "should not be able to pass a number in as an ID")

	// Test get by email
	user3, err := models.GetUserEmail(context.Background(), "jannel@example.com")
	require.NoError(err, "could not fetch user by email")
	require.Equal("01GKHJSK7CZW0W282ZN3E9W86Z", user3.ID.String())
	require.Equal(user, user3)

	// Test Not Found by ID
	_, err = models.GetUser(context.Background(), "01GKHKS95XD0J25GHR14KT3WX1")
	require.ErrorIs(err, models.ErrNotFound, "should return not found error")

	_, err = models.GetUserEmail(context.Background(), "notvalid@testing.io")
	require.ErrorIs(err, models.ErrNotFound, "should return not found error")

	// Test cannot parse ULID
	_, err = models.GetUser(context.Background(), "zedy")
	require.EqualError(err, "ulid: bad data size when unmarshaling")
}

func (m *modelTestSuite) TestUserSave() {
	defer m.ResetDB()

	require := m.Require()
	user, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.NoError(err, "could not fetch user by string ID")
	require.Equal("Jannel P. Hudson", user.Name)

	prevModified := user.Modified
	user.Name = "New Name"
	user.Email = "new@example.com"
	user.Password = "Invalid Password"
	user.SetLastLogin(time.Now())

	err = user.Save(context.Background())
	require.NoError(err, "could not update user")

	cmpr, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.NoError(err, "could not fetch user by string ID")

	// Everything but modified should be the same on compare
	require.Equal(user.Name, cmpr.Name)
	require.Equal(user.Email, cmpr.Email)
	require.Equal(user.Password, cmpr.Password)
	require.Equal(user.LastLogin, cmpr.LastLogin)
	require.Equal(user.Created, cmpr.Created)
	require.Equal(user.Modified, cmpr.Modified)
	require.NotEqual(prevModified, cmpr.Modified)
}

func (m *modelTestSuite) TestUserUpdateLastLogin() {
	defer m.ResetDB()

	require := m.Require()
	user, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.NoError(err, "could not fetch user by string ID")

	// The user pointer will be modified so get a second copy for comparison
	prev, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.NoError(err, "could not fetch user by string ID")

	err = user.UpdateLastLogin(context.Background())
	require.NoError(err, "could not update last login: %+v", err)

	// Fetch the record from the database for comparison purposes.
	cmpr, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.NoError(err, "could not fetch user by string ID")

	// Nothing but last login and modified should have changed.
	require.Equal(prev.Name, cmpr.Name)
	require.Equal(prev.Email, cmpr.Email)
	require.Equal(prev.Password, cmpr.Password)
	require.Equal(prev.Created, cmpr.Created)

	// Last Login and Modified should have changed to the same timestamp
	require.Equal(cmpr.LastLogin.String, cmpr.Modified, "expected modified and last login to be equal")
	require.NotEqual(prev.LastLogin.String, cmpr.LastLogin.String)
	require.NotEqual(prev.Modified, cmpr.Modified)

	// The pointer should have been updated to match what's in the database
	require.Equal(user.LastLogin.String, cmpr.LastLogin.String)
	require.Equal(user.Modified, cmpr.Modified)

	// Last Login and Modified should be after the previous Last Login and Modified
	ll, err := cmpr.GetLastLogin()
	require.NoError(err, "could not parse last login")
	require.False(ll.IsZero())

	pll, err := prev.GetLastLogin()
	require.NoError(err, "could not parse last login fixture")
	require.True(ll.After(pll), "cmpr last login %q is not after prev last login %q", cmpr.LastLogin.String, prev.LastLogin.String)

	mod, err := cmpr.GetModified()
	require.NoError(err, "could not parse modified")
	require.False(mod.IsZero())

	pmod, err := prev.GetModified()
	require.NoError(err, "could not parse modified fixture")
	require.True(mod.After(pmod), "cmpr modified %q is not after prev modified %q", cmpr.Modified, prev.Modified)
}

func TestUserLastLogin(t *testing.T) {
	user := &models.User{}

	ts, err := user.GetLastLogin()
	require.NoError(t, err, "could not get null last login")
	require.Zero(t, ts, "expected zero-valued timestamp")

	now := time.Now()
	user.SetLastLogin(now)

	ts, err = user.GetLastLogin()
	require.NoError(t, err, "could not get non-null last login")
	require.True(t, now.Equal(ts))
}
