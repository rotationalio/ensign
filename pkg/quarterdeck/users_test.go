package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/emails/mock"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *quarterdeckTestSuite) TestUserDetail() {
	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test passing empty ULID results in StatusNotFound error
	user, err := s.client.UserDetail(ctx, "")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	require.Nil(user, "expected no data returned after an error")

	// Retrieving a user requires authentication
	user, err = s.client.UserDetail(ctx, "01GQFQ4475V3BZDMSXFV5DK6YY")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	require.Nil(user, "expected no data returned after an error")

	// Retrieving a user's detail requires the collaborators:read permission
	claims := &tokens.Claims{
		Name:  "Invalid User",
		Email: "invalid@user.com",
		OrgID: ulids.New().String(),
	}
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.UserDetail(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z")
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")
	require.Nil(user, "expected no data returned after an error")

	// invalid permissions results in a StatusUnauthorized error
	claims.Permissions = []string{perms.ReadAPIKeys}
	ctx = s.AuthContext(ctx, claims)

	user, err = s.client.UserDetail(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z")
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")
	require.Nil(user, "expected no data returned after an error")

	// Invalid requester with correct permissions but in an organization that does not exist cannot retrieve detail of a user
	claims.Permissions = []string{perms.ReadCollaborators}
	ctx = s.AuthContext(ctx, claims)

	// Test that the requester does not have permission to access the user because the orgID does not exist in the database
	user, err = s.client.UserDetail(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z")
	s.CheckError(err, http.StatusForbidden, "requester is not authorized to access this user")
	require.Nil(user, "expected no data returned after an error")

	// set up valid requester with collaborators:read permission but requesting
	// detail for user in a different organization results in StatusForbidden error
	claims = &tokens.Claims{
		Name:        "Edison Edgar Franklin",
		Email:       "eefrank@checkers.io",
		OrgID:       "01GQFQ14HXF2VC7C1HJECS60XX",
		Permissions: []string{perms.ReadCollaborators},
	}
	ctx = s.AuthContext(ctx, claims)

	user, err = s.client.UserDetail(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z")
	s.CheckError(err, http.StatusForbidden, "requester is not authorized to access this user")
	require.Nil(user, "expected no data returned after an error")

	// happy path test
	user, err = s.client.UserDetail(ctx, "01GQYYKY0ECGWT5VJRVR32MFHM")
	require.NoError(err, "could not fetch valid user detail")
	require.NotNil(user, "expected user to be retrieved")
}

func (s *quarterdeckTestSuite) TestUserUpdate() {
	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Need to be authorized to update a user
	in := &api.User{}
	user, err := s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	require.Nil(user, "expected no data returned after an error")

	// Updating a user requires the collaborators:edit permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
		OrgID: ulids.New().String(),
	}
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")
	require.Nil(user, "expected no data returned after an error")

	// validate incomplete user information returns error
	claims.Permissions = []string{perms.EditCollaborators}
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "missing required field: user_id")
	require.Nil(user, "expected no data returned after an error")

	// missing claims subject results in error
	in.UserID = ulids.New()
	in.Name = "Johnny Miller"
	user, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "user claims invalid or unavailable")
	require.Nil(user, "expected no data returned after an error")

	// invalid user_id results in error
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(user, "expected no data returned after an error")

	// passing in user from a different organization results in error
	in.UserID = ulid.MustParse("01GQFQ4475V3BZDMSXFV5DK6XX")
	user, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(user, "expected no data returned after an error")

	// invalid requester orgID results in error
	in.UserID = ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	user, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(user, "expected no data returned after an error")

	// happy path test
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.UserUpdate(ctx, in)
	require.NoError(err, "should have been able to update the user")
	require.NotSame(in, user, "expected a different object to be returned")
}

func (s *quarterdeckTestSuite) TestListUser() {
	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Listing users requires authentication
	req := &api.UserPageQuery{}
	_, err := s.client.UserList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Listing users requires the collaborators:read permission
	claims := &tokens.Claims{
		Name:  "Edison Edgar Franklin",
		Email: "eefrank@checkers.io",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.UserList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GQFQ4475V3BZDMSXFV5DK6XX"
	claims.OrgID = "01GQFQ14HXF2VC7C1HJECS60XX"
	claims.Permissions = []string{perms.ReadCollaborators}
	ctx = s.AuthContext(ctx, claims)

	// Should be able to list all users for the specified organization
	page, err := s.client.UserList(ctx, req)
	require.NoError(err, "could not fetch users")
	require.Len(page.Users, 4, "expected 4 results back from the fixtures")
	require.Empty(page.NextPageToken, "expected no next page token in response")

	// Should be able to paginate the request for the specified organization
	req.PageSize = 1
	page, err = s.client.UserList(ctx, req)
	require.NoError(err, "could not fetch paginated users")
	require.Len(page.Users, 1, "expected 1 result back from the fixtures")
	require.NotEmpty(page.NextPageToken, "expected next page token in response")

	// Test fetching the next page with the next page token
	req.NextPageToken = page.NextPageToken
	page2, err := s.client.UserList(ctx, req)
	require.NoError(err, "could not fetch paginated api keys")
	require.Len(page2.Users, 1, "expected 1 result back from the fixtures")
	require.NotEmpty(page2.NextPageToken, "expected next page token in response")
	require.NotEqual(page.Users[0].Name, page2.Users[0].Name, "expected a new page of results")

	// Limit maximum number of requests to 4, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for {
		page, err = s.client.UserList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Users)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 4, "expected 4 pages")
	require.Equal(nResults, 4, "expected 4 results")
}

func (s *quarterdeckTestSuite) TestUserDelete() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: implement the delete user test
	err := s.client.UserDelete(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z")
	require.Error(err, "expected unimplemented error")
}

func (s *quarterdeckTestSuite) TestCreateUserNotAllowed() {
	// Ensure that a user cannot be created via a POST to the /v1/users endpoint.
	require := s.Require()

	apiv1, ok := s.client.(*api.APIv1)
	require.True(ok)

	userID := ulids.New()

	req, err := apiv1.NewRequest(context.TODO(), http.MethodPost, "/v1/users", userID, nil)
	require.NoError(err)

	_, err = apiv1.Do(req, nil, true)
	s.CheckError(err, http.StatusMethodNotAllowed, "method not allowed")
}

func (s *quarterdeckTestSuite) TestUserInvite() {
	require := s.Require()
	defer s.ResetDatabase()
	defer s.ResetTasks()
	defer mock.Reset()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Inviting users requires authentication
	req := &api.UserInviteRequest{}
	_, err := s.client.UserInvite(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Inviting users requires the collaborators:add permission
	claims := &tokens.Claims{
		Name:  "Edison Edgar Franklin",
		Email: "eefrank@checkers.io",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.UserInvite(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GQFQ4475V3BZDMSXFV5DK6XX"
	claims.OrgID = "01GQFQ14HXF2VC7C1HJECS60XX"
	orgID := ulid.MustParse(claims.OrgID)
	subjectID := ulid.MustParse(claims.Subject)
	claims.Permissions = []string{perms.AddCollaborators}
	ctx = s.AuthContext(ctx, claims)

	// Inviting a user requires an email address
	req.Role = perms.RoleMember
	_, err = s.client.UserInvite(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: email")

	// Inviting a user requires a role
	userID := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	req.Email = "jannel@example.com"
	req.Role = ""
	_, err = s.client.UserInvite(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: role")

	// Should return an error if the role is invalid
	req.Role = "not-a-valid-role"
	_, err = s.client.UserInvite(ctx, req)
	s.CheckError(err, http.StatusBadRequest, api.ErrUnknownUserRole.Error())

	// Valid request - invited user already has an account
	req.Role = perms.RoleMember
	sent := time.Now()
	rep, err := s.client.UserInvite(ctx, req)
	require.NoError(err, "could not invite user")
	require.Equal(userID, rep.UserID, "expected user ID to match")
	require.Equal(orgID, rep.OrgID, "expected org ID to match")
	require.Equal(req.Email, rep.Email, "expected email to match")
	require.Equal(req.Role, rep.Role, "expected role to match")
	require.Equal(subjectID, rep.CreatedBy, "expected created by to match")
	require.NotEmpty(rep.Created, "expected created at to be set")
	require.NotEmpty(rep.ExpiresAt, "expected expires at to be set")

	// Valid request - invited user does not have an account
	req.Email = "gon@hunters.com"
	rep, err = s.client.UserInvite(ctx, req)
	require.NoError(err, "could not invite user")
	require.NotEqual(ulid.ULID{}, rep.UserID, "expected user ID to be set")
	require.Equal(orgID, rep.OrgID, "expected org ID to match")
	require.Equal(req.Email, rep.Email, "expected email to match")
	require.Equal(req.Role, rep.Role, "expected role to match")
	require.Equal(subjectID, rep.CreatedBy, "expected created by to match")
	require.NotEmpty(rep.Created, "expected created at to be set")
	require.NotEmpty(rep.ExpiresAt, "expected expires at to be set")

	// Check that both invite emails were sent
	s.StopTasks()
	messages := []*mock.EmailMeta{
		{
			To:        "jannel@example.com",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   emails.InviteRE,
			Timestamp: sent,
		},
		{
			To:        "gon@hunters.com",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   emails.InviteRE,
			Timestamp: sent,
		},
	}
	mock.CheckEmails(s.T(), messages)
}
