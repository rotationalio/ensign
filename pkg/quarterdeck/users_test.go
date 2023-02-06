package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

func (s *quarterdeckTestSuite) TestUserUpdate() {
	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Need to be authorized to update a user
	in := &api.User{}
	out, err := s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	require.Nil(out, "expected no data returned after an error")

	// Updating a user requires the collaborators:edit permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
		OrgID: ulids.New().String(),
	}
	ctx = s.AuthContext(ctx, claims)
	out, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")
	require.Nil(out, "expected no data returned after an error")

	// validate incomplete user information returns error
	claims.Permissions = []string{perms.EditCollaborators}
	ctx = s.AuthContext(ctx, claims)
	out, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "missing required field: user_id")
	require.Nil(out, "expected no data returned after an error")

	// missing claims subject results in error
	in.UserID = ulids.New()
	in.Name = "Johnny Miller"
	out, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "invalid user claims")
	require.Nil(out, "expected no data returned after an error")

	// invalid user_id results in error
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	ctx = s.AuthContext(ctx, claims)
	out, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(out, "expected no data returned after an error")

	// passing in user from a different organization results in error
	in.UserID = ulid.MustParse("01GQFQ4475V3BZDMSXFV5DK6XX")
	out, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(out, "expected no data returned after an error")

	// invalid requester orgID results in error
	in.UserID = ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	out, err = s.client.UserUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(out, "expected no data returned after an error")

	// happy path test
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	out, err = s.client.UserUpdate(ctx, in)
	require.NoError(err, "should have been able to update the user")
	require.NotSame(in, out, "expected a different object to be returned")
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
	for i := 0; i < 4; i++ {
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
