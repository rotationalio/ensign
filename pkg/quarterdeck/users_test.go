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
