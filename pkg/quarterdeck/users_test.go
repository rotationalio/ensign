package quarterdeck_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

func (s *quarterdeckTestSuite) TestUserDetail() {
	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test passing empty ULID results in StatusNotFound error
	user, err := s.client.UserDetail(ctx, "")
	s.CheckError(err, http.StatusNotFound, "resource not found")
	require.Nil(user, "expected no data returned after an error")

	// Test passing invalid ULID results in StatusUnauthorized error
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
	fmt.Println(user)
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
	s.CheckError(err, http.StatusBadRequest, "invalid user claims")
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
