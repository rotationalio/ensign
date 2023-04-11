package models_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"

	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestGetUser() {
	require := m.Require()

	testCases := []struct {
		userID         any
		orgID          any
		err            error
		validateFields bool
	}{
		// Test GetUser by userID string and default org
		{"01GKHJSK7CZW0W282ZN3E9W86Z", ulids.Null, nil, true},

		// Test GetUser by userID ULID and default org
		{ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z"), ulids.Null, nil, true},

		// Test GetUser by string with specified OrgID
		{"01GQYYKY0ECGWT5VJRVR32MFHM", "01GQFQ14HXF2VC7C1HJECS60XX", nil, true},

		// Test GetUser by ULIDs with specified OrgID
		{ulid.MustParse("01GQYYKY0ECGWT5VJRVR32MFHM"), ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"), nil, true},

		// Should not be able to pass an integer in as the userID
		{42, ulids.Null, ulids.ErrUnknownType, false},

		// Test cannot parse ID
		{"zedy", ulids.Null, ulid.ErrDataSize, false},
		{"01GQYYKY0ECGWT5VJRVR32MFHM", "zedy", ulid.ErrDataSize, false},

		// Test Not Found by userID
		{"01GKHKS95XD0J25GHR14KT3WX1", ulids.Null, models.ErrNotFound, false},

		// Test Not Found by null ID
		{ulids.Null, ulids.Null, models.ErrNotFound, false},

		// Test User not in organization
		{"01GQYYKY0ECGWT5VJRVR32MFHM", "01GKHKS95XD0J25GHR14KT3WX1", models.ErrUserOrganization, false},
	}

	for _, tc := range testCases {
		user, err := models.GetUser(context.Background(), tc.userID, tc.orgID)
		require.ErrorIs(err, tc.err)

		if tc.validateFields {
			// Ensure all fields are returned and not zero valued
			require.False(ulids.IsZero(user.ID))
			require.NotEmpty(user.Name)
			require.NotEmpty(user.Email)
			require.NotEmpty(user.Password)
			require.True(user.AgreeToS.Valid && user.AgreeToS.Bool)
			require.True(user.AgreePrivacy.Valid && user.AgreePrivacy.Bool)
			require.True(user.EmailVerified)
			require.True(user.EmailVerificationExpires.Valid && user.EmailVerificationExpires.String != "")
			require.True(user.EmailVerificationToken.Valid && user.EmailVerificationToken.String != "")
			require.NotEmpty(user.EmailVerificationSecret)
			require.True(user.LastLogin.Valid && user.LastLogin.String != "")
			require.NotEmpty(user.Created)
			require.NotEmpty(user.Modified)

			orgID, err := user.OrgID()
			require.NoError(err, "could not fetch orgID from user")
			require.False(ulids.IsZero(orgID))

			role, err := user.Role()
			require.NoError(err, "could not fetch role from user")
			require.NotEmpty(role)

			perms, err := user.Permissions(context.Background(), false)
			require.NoError(err, "could not fetch permissions for user")
			require.NotEmpty(perms)
		}
	}
}

func (m *modelTestSuite) TestGetUserEmail() {
	require := m.Require()

	testCases := []struct {
		email          string
		orgID          any
		err            error
		validateFields bool
	}{
		// Test GetUser by email and default org
		{"jannel@example.com", ulids.Null, nil, true},

		// Test GetUser by string with specified OrgID
		{"jannel@example.com", "01GKHJRF01YXHZ51YMMKV3RCMK", nil, true},

		// Test GetUser by ULIDs with specified OrgID
		{"jannel@example.com", ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"), nil, true},

		// Test cannot parse org ID
		{"jannel@example.com", "zedy", ulid.ErrDataSize, false},

		// Test Not Found by email address
		{"notvalid@esting.io", ulids.Null, models.ErrNotFound, false},

		// Test Not Found by empty email address
		{"", ulids.Null, models.ErrNotFound, false},

		// Test User not in organization
		{"jannel@example.com", "01GKHKS95XD0J25GHR14KT3WX1", models.ErrUserOrganization, false},
	}

	for i, tc := range testCases {
		user, err := models.GetUserEmail(context.Background(), tc.email, tc.orgID)
		require.ErrorIs(err, tc.err, "could not get user for test %d", i)

		if tc.validateFields {
			// Ensure all fields are returned and not zero valued
			require.False(ulids.IsZero(user.ID))
			require.NotEmpty(user.Name)
			require.NotEmpty(user.Email)
			require.NotEmpty(user.Password)
			require.True(user.AgreeToS.Valid && user.AgreeToS.Bool)
			require.True(user.AgreePrivacy.Valid && user.AgreePrivacy.Bool)
			require.True(user.EmailVerified)
			require.True(user.EmailVerificationExpires.Valid && user.EmailVerificationExpires.String != "")
			require.True(user.EmailVerificationToken.Valid && user.EmailVerificationToken.String != "")
			require.NotEmpty(user.EmailVerificationSecret)
			require.True(user.LastLogin.Valid && user.LastLogin.String != "")
			require.NotEmpty(user.Created)
			require.NotEmpty(user.Modified)

			orgID, err := user.OrgID()
			require.NoError(err, "could not fetch orgID from user")
			require.False(ulids.IsZero(orgID))

			role, err := user.Role()
			require.NoError(err, "could not fetch role from user")
			require.NotEmpty(role)

			perms, err := user.Permissions(context.Background(), false)
			require.NoError(err, "could not fetch permissions for user")
			require.NotEmpty(perms)
		}
	}
}

func (m *modelTestSuite) TestGetUserMultiOrg() {
	require := m.Require()
	testCases := []struct {
		userID any
		orgID  string
		email  string
		role   string
	}{
		{"01GQYYKY0ECGWT5VJRVR32MFHM", "01GKHJRF01YXHZ51YMMKV3RCMK", "zendaya@testing.io", "Observer"},
		{"01GQYYKY0ECGWT5VJRVR32MFHM", "01GQFQ14HXF2VC7C1HJECS60XX", "zendaya@testing.io", "Member"},
	}

	for _, tc := range testCases {
		// Test GetUser by ID
		user, err := models.GetUser(context.Background(), tc.userID, tc.orgID)
		require.NoError(err)

		orgID, _ := user.OrgID()
		require.Equal(tc.orgID, orgID.String())

		role, _ := user.Role()
		require.Equal(tc.role, role)

		// Test GetUser by email
		user, err = models.GetUserEmail(context.Background(), tc.email, tc.orgID)
		require.NoError(err)

		orgID, _ = user.OrgID()
		require.Equal(tc.orgID, orgID.String())

		role, _ = user.Role()
		require.Equal(tc.role, role)
	}
}

func (m *modelTestSuite) TestUserCreateNewOrg() {
	defer m.ResetDB()
	require := m.Require()

	// Ensure the original user and organization count is as expected
	nUsers, err := models.CountUsers(context.Background())
	require.NoError(err, "could not count users")
	require.Equal(nUserFixtures, nUsers, "unexpected user fixtures count")

	nOrgs, err := models.CountOrganizations(context.Background())
	require.NoError(err, "could not count orgs")
	require.Equal(nOrganizationFixtures, nOrgs, "unexpected organization fixtures count")

	// Create a user with as minimal information as possible.
	user := &models.User{
		Name:     "Angelica Hudson",
		Email:    "hudson@example.com",
		Password: "$argon2id$v=19$m=65536,t=1,p=2$xto5+nlVR9oyc6CpJR1MtQ==$KToxSO2i3H6KmD8th1FiP1jh/JvDUOfdtMtj5g1Ilnk=",
	}

	user.SetAgreement(true, true)

	// This organization should not exist in the database
	org := &models.Organization{
		Name:   "Testing Organization",
		Domain: "testing",
	}

	// Create the user, the organization, and associate them with the role "Admin"
	require.NoError(user.Create(context.Background(), org, "Admin"), "could not create user")

	// Ensure that an ID, created, and modified timestamps on the user were created
	require.False(ulids.IsZero(user.ID))
	require.NotZero(user.Created)
	require.NotZero(user.Modified)

	// Ensure that an ID, created, and modified timestamps on the org were created
	require.False(ulids.IsZero(org.ID))
	require.NotZero(org.Created)
	require.NotZero(org.Modified)

	// Ensure that the number of users in the database has increased
	nUsers, err = models.CountUsers(context.Background())
	require.NoError(err, "could not count users")
	require.Equal(nUserFixtures+1, nUsers, "user count not increased after create")

	// Ensure the number of organizations in the database have been increased
	nOrgs, err = models.CountOrganizations(context.Background())
	require.NoError(err, "could not count organizations")
	require.Equal(nOrganizationFixtures+1, nOrgs, "organization count not increased after create")

	// Check that the user has been assigned the organization that was created
	userOrg, _ := user.OrgID()
	require.Equal(org.ID, userOrg)

	// Check that the organization and user are linked with a role
	our, err := models.GetOrgUser(context.Background(), user.ID, org.ID)
	require.NoError(err, "could not fetch organization user mapping with role")

	cmpuser, err := our.User(context.Background(), false)
	require.NoError(err, "could not get user to compare")
	require.Equal(user, cmpuser)

	cmporg, err := our.Organization(context.Background(), false)
	require.NoError(err, "could not get organization to compare")
	require.Equal(org, cmporg)

	role, err := our.Role(context.Background(), false)
	require.NoError(err, "could not get user role fom database")
	require.Equal("Admin", role.Name)

	userPerms, err := user.Permissions(context.Background(), false)
	require.NoError(err, "could not get user permissions")
	rolePerms, err := role.Permissions(context.Background(), false)
	require.NoError(err, "could not get role permissions")

	require.Equal(len(userPerms), len(rolePerms), "user and role permissions do not match")
	for _, perm := range rolePerms {
		require.Contains(userPerms, perm.Name)
	}
}

func (m *modelTestSuite) TestUserCreateExistingOrg() {
	defer m.ResetDB()
	require := m.Require()

	// Ensure the original user and organization count is as expected
	nUsers, err := models.CountUsers(context.Background())
	require.NoError(err, "could not count users")
	require.Equal(nUserFixtures, nUsers, "unexpected user fixtures count")

	nOrgs, err := models.CountOrganizations(context.Background())
	require.NoError(err, "could not count orgs")
	require.Equal(nOrganizationFixtures, nOrgs, "unexpected organization fixtures count")

	// Create a user with as minimal information as possible.
	user := &models.User{
		Name:     "Angelica Hudson",
		Email:    "hudson@example.com",
		Password: "$argon2id$v=19$m=65536,t=1,p=2$xto5+nlVR9oyc6CpJR1MtQ==$KToxSO2i3H6KmD8th1FiP1jh/JvDUOfdtMtj5g1Ilnk=",
	}

	user.SetAgreement(true, true)

	// This organization should not exist in the database
	org := &models.Organization{
		ID: ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"),
	}

	// Create the user, the organization, and associate them with the role "Member"
	require.NoError(user.Create(context.Background(), org, "Member"), "could not create user")

	// Ensure that an ID, created, and modified timestamps on the user were created
	require.False(ulids.IsZero(user.ID))
	require.NotZero(user.Created)
	require.NotZero(user.Modified)

	// Ensure that an ID, created, and modified timestamps on the org were created
	require.False(ulids.IsZero(org.ID))
	require.NotZero(org.Created)
	require.NotZero(org.Modified)

	// Ensure that the number of users in the database has increased
	nUsers, err = models.CountUsers(context.Background())
	require.NoError(err, "could not count users")
	require.Equal(nUserFixtures+1, nUsers, "user count not increased after create")

	// Ensure the number of organizations in the database have been increased
	nOrgs, err = models.CountOrganizations(context.Background())
	require.NoError(err, "could not count organizations")
	require.Equal(nOrganizationFixtures, nOrgs, "organization count not increased after create")

	// Check that the user has been assigned the organization that was created
	userOrg, _ := user.OrgID()
	require.Equal(org.ID, userOrg)

	// Check that the organization and user are linked with a role
	our, err := models.GetOrgUser(context.Background(), user.ID, org.ID)
	require.NoError(err, "could not fetch organization user mapping with role")

	cmpuser, err := our.User(context.Background(), false)
	require.NoError(err, "could not get user to compare")
	require.Equal(user, cmpuser)

	cmporg, err := our.Organization(context.Background(), false)
	require.NoError(err, "could not get organization to compare")
	require.Equal(org, cmporg)

	role, err := our.Role(context.Background(), false)
	require.NoError(err, "could not get user role fom database")
	require.Equal("Member", role.Name)

	userPerms, err := user.Permissions(context.Background(), false)
	require.NoError(err, "could not get user permissions")
	rolePerms, err := role.Permissions(context.Background(), false)
	require.NoError(err, "could not get role permissions")

	require.Equal(len(userPerms), len(rolePerms), "user and role permissions do not match")
	for _, perm := range rolePerms {
		require.Contains(userPerms, perm.Name)
	}
}

func (m *modelTestSuite) TestUserSave() {
	defer m.ResetDB()

	require := m.Require()
	user, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z", ulid.ULID{})
	require.NoError(err, "could not fetch user by string ID")
	require.Equal("Jannel P. Hudson", user.Name)

	prevModified := user.Modified
	user.Name = "New Name"
	user.Email = "new@example.com"
	user.Password = "Invalid Password"
	user.SetLastLogin(time.Now())

	err = user.Save(context.Background())
	require.ErrorIs(err, models.ErrInvalidPassword, "passwords should be argon2 derived keys")

	// Create a correct password
	user.Password, _ = passwd.CreateDerivedKey(user.Password)
	err = user.Save(context.Background())
	require.NoError(err, "could not update user")

	cmpr, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z", ulid.ULID{})
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
	user, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z", ulid.ULID{})
	require.NoError(err, "could not fetch user by string ID")

	// The user pointer will be modified so get a second copy for comparison
	prev, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z", ulid.ULID{})
	require.NoError(err, "could not fetch user by string ID")

	err = user.UpdateLastLogin(context.Background())
	require.NoError(err, "could not update last login: %+v", err)

	// Fetch the record from the database for comparison purposes.
	cmpr, err := models.GetUser(context.Background(), "01GKHJSK7CZW0W282ZN3E9W86Z", ulid.ULID{})
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

func (m *modelTestSuite) TestUserSwitchOrganization() {
	require := m.Require()

	// A zero-valued user cannot switch organizations
	user := &models.User{}
	require.ErrorIs(user.SwitchOrganization(context.Background(), "01GKHJRF01YXHZ51YMMKV3RCMK"), models.ErrUserOrganization)

	// Get the user in their first organization
	user, err := models.GetUser(context.Background(), "01GQYYKY0ECGWT5VJRVR32MFHM", "01GKHJRF01YXHZ51YMMKV3RCMK")
	require.NoError(err, "could not fetch multi-org user from database")

	orgID, err := user.OrgID()
	require.NoError(err, "could not fetch orgID from user")
	require.Equal(ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"), orgID)

	role, err := user.Role()
	require.NoError(err, "Could not fetch role from user")
	require.Equal("Observer", role)

	// Should not be able to switch into an organization that does not exist
	err = user.SwitchOrganization(context.Background(), "01GQZAE0GQRGB37RA1R3SR5XVH")
	require.ErrorIs(err, models.ErrUserOrganization)

	// Should not be able to switch into an organization the user doesn't belong to
	err = user.SwitchOrganization(context.Background(), "01GQZAC80RAZ1XQJKRZJ2R4KNJ")
	require.ErrorIs(err, models.ErrUserOrganization)

	// Should not be able to switch organizations if the orgId doesn't parse
	err = user.SwitchOrganization(context.Background(), "zeddy")
	require.ErrorIs(err, ulid.ErrDataSize)

	// Switch user to a valid other organization
	err = user.SwitchOrganization(context.Background(), "01GQFQ14HXF2VC7C1HJECS60XX")
	require.NoError(err, "could not switch the user's organization")

	orgID, err = user.OrgID()
	require.NoError(err, "could not fetch orgID from user")
	require.Equal(ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"), orgID)

	role, err = user.Role()
	require.NoError(err, "Could not fetch role from user")
	require.Equal("Member", role)
}

func (m *modelTestSuite) TestRemoveOrganization() {
	defer m.ResetDB()

	require := m.Require()
	userID := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	user, err := models.GetUser(context.Background(), userID, ulid.ULID{})
	require.NoError(err, "could not fetch user from database")
	require.Equal("Jannel P. Hudson", user.Name)

	// Test passing in an empty orgID returns an error
	err = user.RemoveOrganization(context.Background(), ulid.ULID{})
	require.ErrorIs(err, models.ErrMissingOrgID, "empty orgID should return an error")

	// Test successful organization removal
	orgID := ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")
	err = user.RemoveOrganization(context.Background(), orgID)
	require.NoError(err, "could not remove user from organization")

	// Ensure all organization API keys for the user were revoked
	keys, _, err := models.ListAPIKeys(context.Background(), orgID, ulids.Null, userID, nil)
	require.NoError(err, "could not list user API keys")
	require.Empty(keys, "expected user keys to be revoked")

	// Ensure the organization mapping was removed
	_, err = models.GetOrgUser(context.Background(), user.ID, orgID)
	require.ErrorIs(err, models.ErrNotFound, "organization user mapping should not exist")

	// Ensure user invitations were removed
	_, err = models.GetUserInvite(context.Background(), "3s855zxQxp-GEk_tgZkAzBxJUgzsWyUTlxIAee_dOJg")
	require.ErrorIs(err, models.ErrNotFound, "user invite should not exist")

	// Remove the user from their only organization
	orgID = ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	err = user.RemoveOrganization(context.Background(), orgID)
	require.NoError(err, "could not remove user from organization")

	// Ensure all organization API keys for the user were revoked
	keys, _, err = models.ListAPIKeys(context.Background(), orgID, ulids.Null, userID, nil)
	require.NoError(err, "could not list user keys")
	require.Empty(keys, "expected user keys to be revoked")

	// Ensure the orgnization mapping was removed
	_, err = models.GetOrgUser(context.Background(), user.ID, orgID)
	require.ErrorIs(err, models.ErrNotFound, "organization user mapping should not exist")

	// Ensure the user was deleted
	_, err = models.GetUser(context.Background(), orgID, ulid.ULID{})
	require.ErrorIs(err, models.ErrNotFound, "user should not exist")
}

func (m *modelTestSuite) TestUserRole() {
	require := m.Require()

	// Create a user with only a user ID
	userID := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	user := &models.User{ID: userID}

	// Fetch the organization roles for the user
	role, err := user.UserRole(context.Background(), ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"), false)
	require.NoError(err, "could not fetch user role for organization")
	require.Equal(role, "Owner")

	// Should not be able to fetch role that doesn't exist
	_, err = user.UserRole(context.Background(), ulid.MustParse("01GQZ77GJ4700TP8N6QXHQEBVF"), false)
	require.ErrorIs(err, models.ErrUserOrganization)
}

func (m *modelTestSuite) TestUserPermissions() {
	require := m.Require()

	// Create a user with only a user ID
	userID := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	user := &models.User{ID: userID}

	// An organization ID is required to fetch permission
	_, err := user.Permissions(context.Background(), false)
	require.ErrorIs(err, models.ErrMissingOrgID)

	// Add the organization to the user
	user.SwitchOrganization(context.Background(), "01GKHJRF01YXHZ51YMMKV3RCMK")

	// Fetch the permissions for the user
	permissions, err := user.Permissions(context.Background(), false)
	require.NoError(err, "could not fetch permissions for user")
	require.Len(permissions, 18, "wrong number of permissions, have the owner role permissions changed?")
}

func (m *modelTestSuite) TestGetUserInvite() {
	require := m.Require()

	ctx := context.Background()
	token := "3s855zxQxp-GEk_tgZkAzBxJUgzsWyUTlxIAee_dOJg"
	orgID := ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")
	userID := ulid.MustParse("01GQFQ4475V3BZDMSXFV5DK6XX")

	// Test passing empty string returns not found error
	_, err := models.GetUserInvite(ctx, "")
	require.ErrorIs(err, models.ErrNotFound)

	// Test passing non-existent token returns not found error
	_, err = models.GetUserInvite(ctx, "not-a-token")
	require.ErrorIs(err, models.ErrNotFound)

	// Test passing valid token returns the invite
	invite, err := models.GetUserInvite(ctx, token)
	require.NoError(err)
	require.Equal(orgID, invite.OrgID)
	require.Equal("jannel@example.com", invite.Email)
	require.NotEmpty(invite.Token, 64)
	require.NotEmpty(invite.Secret)
	require.Equal(userID, invite.CreatedBy)

	expires, err := time.Parse(time.RFC3339Nano, invite.Expires)
	require.NoError(err, "could not parse invite expiration")
	require.NotZero(expires)
	created, err := invite.GetCreated()
	require.NoError(err, "could not parse invite creation time")
	require.NotZero(created)
	modified, err := invite.GetModified()
	require.NoError(err, "could not parse invite modification time")
	require.NotZero(modified)
}

func (m *modelTestSuite) TestCreateUserInvite() {
	require := m.Require()
	defer m.ResetDB()

	ctx := context.Background()
	userID := ulid.MustParse("01GQFQ4475V3BZDMSXFV5DK6XX")
	orgID := ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")

	// Retrieve the user from the database
	user, err := models.GetUser(ctx, userID, orgID)
	require.NoError(err, "could not retrieve user from database")

	// Creating an invite without an email should return an error
	_, err = user.CreateInvite(ctx, "", "Member")
	require.EqualError(err, "email address is required", "did not return error for missing email")

	// Creating an invite without a role should return an error
	_, err = user.CreateInvite(ctx, "gon@hunters.com", "")
	require.EqualError(err, "missing role for user invite", "did not return error for missing role")

	// Should return an error if the user is already in the organization
	_, err = user.CreateInvite(ctx, "eefrank@checkers.io", "Member")
	require.ErrorIs(err, models.ErrUserOrgExists)

	// Create an invite for a new user
	invite, err := user.CreateInvite(ctx, "gon@hunters.com", "Member")
	require.NoError(err, "could not create user invite")
	require.NotEmpty(invite.Token, "did not return invite token")
	expires, err := time.Parse(time.RFC3339Nano, invite.Expires)
	require.NoError(err, "could not parse invite expiration")
	require.True(expires.After(time.Now()), "invite expiration is not in the future")

	// Make sure the invite was created
	invite, err = models.GetUserInvite(ctx, invite.Token)
	require.NoError(err, "could not retrieve invite from database")
	require.NotZero(invite.UserID)
	require.Equal("Member", invite.Role)
	require.Equal(orgID, invite.OrgID)
	require.Equal("gon@hunters.com", invite.Email)
	require.NotEmpty(invite.Token, 64)
	require.NotEmpty(invite.Secret)
	require.Equal(userID, invite.CreatedBy)

	expiresAt, err := time.Parse(time.RFC3339Nano, invite.Expires)
	require.NoError(err, "could not parse invite expiration")
	require.True(expiresAt.After(time.Now()), "invite expiration is not in the future")
	created, err := invite.GetCreated()
	require.NoError(err, "could not parse invite creation time")
	require.NotZero(created)
	modified, err := invite.GetModified()
	require.NoError(err, "could not parse invite modification time")
	require.NotZero(modified)

	// Test creating an invite for an existing user
	invite, err = user.CreateInvite(ctx, "jannel@example.com", "Admin")
	require.NoError(err, "could not create user invite")
	require.NotEmpty(invite.Token, "did not return invite token")
	require.True(expires.After(time.Now()), "invite expiration is not in the future")

	// Make sure the invite was created
	invite, err = models.GetUserInvite(ctx, invite.Token)
	require.NoError(err, "could not retrieve invite from database")
	require.NotZero(invite.UserID)
	require.Equal("Admin", invite.Role)
	require.Equal(orgID, invite.OrgID)
	require.Equal("jannel@example.com", invite.Email)
	require.NotEmpty(invite.Token, 64)
	require.NotEmpty(invite.Secret)
	require.Equal(userID, invite.CreatedBy)
}

func (m *modelTestSuite) TestDeleteInvite() {
	require := m.Require()
	defer m.ResetDB()

	ctx := context.Background()
	token := "3s855zxQxp-GEk_tgZkAzBxJUgzsWyUTlxIAee_dOJg"

	// Test passing valid token does not error
	err := models.DeleteInvite(ctx, token)
	require.NoError(err, "could not delete invite")

	// Invite should be deleted
	_, err = models.GetUserInvite(ctx, token)
	require.ErrorIs(err, models.ErrNotFound)
}

func (m *modelTestSuite) TestListUsers() {
	require := m.Require()

	ctx := context.Background()
	orgID := ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")

	//test passing null orgID results in error
	users, cursor, err := models.ListUsers(ctx, ulids.Null, nil)
	require.ErrorIs(err, models.ErrMissingOrgID, "orgID is required for list queries")
	require.NotNil(err)
	require.Nil(cursor)
	require.Nil(users)

	//test passing invalid orgID results in error
	users, cursor, err = models.ListUsers(ctx, 1, nil)
	require.Contains("cannot parse input: unknown type", err.Error())
	require.NotNil(err)
	require.Nil(cursor)
	require.Nil(users)

	// test passing in pagination.Cursor without page size results in error
	_, _, err = models.ListUsers(ctx, orgID, &pagination.Cursor{})
	require.ErrorIs(err, models.ErrMissingPageSize, "page size is required for list users queries with pagination")

	// Should return all checkers org users (page cursor not required)
	// there are 4 users associated with this org in the fixtures
	users, cursor, err = models.ListUsers(ctx, orgID, nil)
	require.NoError(err, "could not fetch all users for checkers org")
	require.Nil(cursor, "should be no next page so no cursor")
	require.Len(users, 4, "expected 4 users from checkers org")
	user := users[0]
	//verify password is not returned
	require.Empty(user.Password)
	//verify all other values are returned
	require.NotNil(user.ID)
	require.NotNil(user.Name)
	require.NotNil(user.Email)
	require.NotNil(user.AgreeToS)
	require.NotNil(user.AgreePrivacy)
	require.NotNil(user.LastLogin)
	require.NotNil(user.OrgID)
	role, err := user.Role()
	require.Nil(err)
	require.Equal("Owner", role)
	permissions, err := user.Permissions(ctx, false)
	require.Nil(err)
	require.Len(permissions, 18, "expected 18 permissions for user")
}

func (m *modelTestSuite) TestListUsersPagination() {
	require := m.Require()
	ctx := context.Background()
	orgID := ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")

	pages := 0
	nRows := 0
	cursor := pagination.New("", "", 2)
	for cursor != nil && pages < 100 {
		users, nextPage, err := models.ListUsers(ctx, orgID, cursor)
		require.NoError(err, "could not fetch page from server")
		if nextPage != nil {
			require.NotEqual(cursor.StartIndex, nextPage.StartIndex)
			require.NotEqual(cursor.EndIndex, nextPage.EndIndex)
			require.Equal(cursor.PageSize, nextPage.PageSize)
		}

		pages++
		nRows += len(users)
		cursor = nextPage
	}

	require.Equal(2, pages, "expected 2 pages")
	require.Equal(4, nRows, "expected 4 results")
}

func (m *modelTestSuite) TestUpdate() {
	defer m.ResetDB()

	require := m.Require()

	user := &models.User{}
	ctx := context.Background()
	// passing in a zero-valued userID returns error
	err := user.Update(ctx, 0)
	require.ErrorIs(err, models.ErrMissingModelID)

	userID := ulid.MustParse("01GQYYKY0ECGWT5VJRVR32MFHM")
	user = &models.User{ID: userID}
	// passing in a zero-valued orgID returns error
	err = user.Update(ctx, 0)
	require.ErrorIs(err, models.ErrMissingOrgID)

	// passing in a nil orgID returns error
	err = user.Update(ctx, nil)
	require.ErrorIs(err, models.ErrMissingOrgID)

	// passing in a user object without a name returns error
	orgID := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Y")
	err = user.Update(ctx, orgID)
	require.ErrorIs(err, models.ErrInvalidUser)

	// failure to pass in valid orgID returns error
	user.Name = "Sarah Fisher"
	err = user.Update(ctx, orgID)
	require.Equal(models.ErrNotFound, err)

	// passing an orgID that's different from the user's organization results in an error
	orgID = ulid.MustParse("01GQZAC80RAZ1XQJKRZJ2R4KNJ")
	err = user.Update(ctx, orgID)
	require.Equal("object not found in the database", err.Error())

	// happy path test
	orgID = ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	err = user.Update(ctx, orgID)
	require.NoError(err)
}

func TestVerificationToken(t *testing.T) {
	// Should return an error if the user does not have an email
	user := &models.User{}
	require.EqualError(t, user.CreateVerificationToken(), "email address is required", "should return an error if the user does not have an email")

	// Test that the fields are set correctly
	user.Email = "leopold.wentzel@gmail.com"
	require.NoError(t, user.CreateVerificationToken(), "could not create email token")
	require.NotEmpty(t, user.GetVerificationToken(), "email verification token should be set")
	require.True(t, user.EmailVerificationExpires.Valid, "email verification token expiration should be set")
	expiresAt, err := time.Parse(time.RFC3339Nano, user.EmailVerificationExpires.String)
	require.NoError(t, err, "could not parse email verification token expiration")
	require.True(t, expiresAt.After(time.Now()), "email verification token expiration should be in the future")
	require.Len(t, user.EmailVerificationSecret, 128, "wrong length for email verification secret")
}
