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

func (s *quarterdeckTestSuite) TestAccountUpdate() {
	require := s.Require()
	defer s.ResetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Need to be authorized to update a user
	in := &api.User{}
	user, err := s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	require.Nil(user, "expected no data returned after an error")

	// Updating a user requires the collaborators:edit permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
		OrgID: ulids.New().String(),
	}
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")
	require.Nil(user, "expected no data returned after an error")

	// validate incomplete user information returns error
	claims.Permissions = []string{perms.ReadCollaborators}
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "missing required field: user_id")
	require.Nil(user, "expected no data returned after an error")

	// missing claims subject results in error
	in.UserID = ulids.New()
	in.Name = "Johnny Miller"
	user, err = s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "invalid user claims")
	require.Nil(user, "expected no data returned after an error")

	// mismatch between requester_id and user_id results in error
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "resource id does not match id of endpoint")
	require.Nil(user, "expected no data returned after an error")

	// passing in a different user id results in error
	in.UserID = ulid.MustParse("01GQFQ4475V3BZDMSXFV5DK6XX")
	user, err = s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusBadRequest, "resource id does not match id of endpoint")
	require.Nil(user, "expected no data returned after an error")

	// invalid requester orgID results in error
	in.UserID = ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	user, err = s.client.AccountUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "user id not found")
	require.Nil(user, "expected no data returned after an error")

	// happy path test: valid orgID, requesterID and user_id match
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	user, err = s.client.AccountUpdate(ctx, in)
	require.NoError(err, "should have been able to update the user")
	require.NotSame(in, user, "expected a different object to be returned")
}
