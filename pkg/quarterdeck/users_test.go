package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
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
	/*
		require := s.Require()
		data, err := hex.DecodeString("018627cf1e8d2c5f4fdc91ffe01353de")
		require.NoError(err, "could not decode user id")
		id, err := ulids.Parse(data)
		require.NoError(err, "could not parse user id")*/

	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Deleting a user requires authentication
	err := s.client.UserDelete(ctx, "invalid")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Listing users requires the collaborators:read permission
	claims := &tokens.Claims{
		Name:  "Edison Edgar Franklin",
		Email: "eefrank@checkers.io",
	}
	ctx = s.AuthContext(ctx, claims)

	err = s.client.UserDelete(ctx, "invalid")
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GQFQ4475V3BZDMSXFV5DK6XX"
	claims.OrgID = "01GQFQ14HXF2VC7C1HJECS60XX"
	claims.Permissions = []string{perms.RemoveCollaborators}
	ctx = s.AuthContext(ctx, claims)

	// Should return an error if the user ID is invalid
	err = s.client.UserDelete(ctx, "invalid")
	s.CheckError(err, http.StatusBadRequest, "could not parse request")

	// Should return an error if the user does not exist
	err = s.client.UserDelete(ctx, "01234JSK7CZW0W282ZN3E9W86Z")
	s.CheckError(err, http.StatusNotFound, "user not found")

	// Should return an error if the user is not in the organization
	orgID := "01GKHJSK7CZW0W282ZN3E9W86Z"
	err = s.client.UserDelete(ctx, orgID)
	s.CheckError(err, http.StatusNotFound, "user not found")

	// Successfully deleting a user from an organization
	userID := "01GQYYKY0ECGWT5VJRVR32MFHM"
	err = s.client.UserDelete(ctx, userID)
	require.NoError(err, "could not delete user")

	// Ensure all organization API keys for the user were revoked
	keys, _, err := models.ListAPIKeys(context.Background(), ulids.MustParse(orgID), ulids.Null, ulids.MustParse(userID), nil)
	require.NoError(err, "could not list api keys")
	require.Empty(keys, "expected user keys to be revoked")

	// Ensure the organization mapping was removed
	_, err = models.GetOrgUser(context.Background(), userID, orgID)
	require.ErrorIs(err, models.ErrNotFound, "organization user mapping should not exist")

	// User should still exist
	_, err = models.GetUser(context.Background(), userID, ulids.Null)
	require.NoError(err, "user should still exist")

	// Remove the user from their last organization
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	err = s.client.UserDelete(ctx, userID)
	require.NoError(err, "could not delete user")

	// Ensure the organization API keys were revoked
	keys, _, err = models.ListAPIKeys(context.Background(), ulids.MustParse(claims.OrgID), ulids.Null, ulids.MustParse(userID), nil)
	require.NoError(err, "could not list api keys")
	require.Empty(keys, "expected user keys to be revoked")

	// Ensure the organization mapping was removed
	_, err = models.GetOrgUser(context.Background(), userID, claims.OrgID)
	require.ErrorIs(err, models.ErrNotFound, "organization user mapping should not exist")

	// Ensure the user was deleted since this was their last organization
	_, err = models.GetUser(context.Background(), userID, claims.OrgID)
	require.ErrorIs(err, models.ErrNotFound, "user should not exist")
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
