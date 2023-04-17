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
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/emails/mock"
)

func (s *quarterdeckTestSuite) TestInvitePreview() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Should return an error if the invite token is not found
	_, err := s.client.InvitePreview(ctx, "invalid-token")
	s.CheckError(err, http.StatusBadRequest, "invalid invitation")

	// Should return an error if the invite token is expired
	_, err = s.client.InvitePreview(ctx, "s6jsNBizyGh_C_ZsUSuJsquONYa--gpcfzorN8DsdjIA")
	s.CheckError(err, http.StatusBadRequest, "invitation has expired")

	// Should return an error if the invite token is invalid
	_, err = s.client.InvitePreview(ctx, "pUqQaDxWrqSGZzkxFDYNfCMSMlB--gpcfzorN8DsdjIA")
	s.CheckError(err, http.StatusBadRequest, "invalid invitation")

	// Successfully retrieving an invite for a new user
	invite, err := s.client.InvitePreview(ctx, "s6jsNBizyGh_C_ZsUSuJsquONYa-KH_2cmoJZd-jnIk")
	require.NoError(err, "expected valid invite for new user")
	require.Equal("joe@checkers.io", invite.Email, "email did not match")
	require.Equal("Checkers", invite.OrgName, "organization did not match")
	require.Equal("Edison Edgar Franklin", invite.InviterName, "inviter name did not match")
	require.Equal("Member", invite.Role, "role did not match")
	require.False(invite.UserExists, "user should not exist")

	// Successfully retrieving an invite for an existing user
	invite, err = s.client.InvitePreview(ctx, "pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA")
	require.NoError(err, "expected valid invite for existing user")
	require.Equal("eefrank@checkers.io", invite.Email, "email did not match")
	require.Equal("Testing", invite.OrgName, "organization did not match")
	require.Equal("Jannel P. Hudson", invite.InviterName, "inviter name did not match")
	require.Equal("Admin", invite.Role, "role did not match")
	require.True(invite.UserExists, "user should exist")
}

func (s *quarterdeckTestSuite) TestInviteCreate() {
	require := s.Require()
	defer s.ResetDatabase()
	defer s.ResetTasks()
	defer mock.Reset()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Inviting users requires authentication
	req := &api.UserInviteRequest{}
	_, err := s.client.InviteCreate(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Inviting users requires the collaborators:add permission
	claims := &tokens.Claims{
		Name:  "Edison Edgar Franklin",
		Email: "eefrank@checkers.io",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.InviteCreate(ctx, req)
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
	_, err = s.client.InviteCreate(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: email")

	// Inviting a user requires a role
	userID := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	req.Email = "jannel@example.com"
	req.Role = ""
	_, err = s.client.InviteCreate(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: role")

	// Should return an error if the role is invalid
	req.Role = "not-a-valid-role"
	_, err = s.client.InviteCreate(ctx, req)
	s.CheckError(err, http.StatusBadRequest, api.ErrUnknownUserRole.Error())

	// Valid request - invited user already has an account
	req.Role = perms.RoleMember
	sent := time.Now()
	rep, err := s.client.InviteCreate(ctx, req)
	require.NoError(err, "could not invite user")
	require.Equal(userID, rep.UserID, "expected user ID to match")
	require.Equal(orgID, rep.OrgID, "expected org ID to match")
	require.Equal(req.Email, rep.Email, "expected email to match")
	require.Equal(req.Role, rep.Role, "expected role to match")
	require.Equal(subjectID, rep.CreatedBy, "expected created by to match")
	require.NotEmpty(rep.Created, "expected created at to be set")
	require.NotEmpty(rep.ExpiresAt, "expected expires at to be set")
	require.Equal("Jannel P. Hudson", rep.Name, "expected invited user name to match")

	// Valid request - invited user does not have an account
	req.Email = "gon@hunters.com"
	rep, err = s.client.InviteCreate(ctx, req)
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
			Subject:   fmt.Sprintf(emails.InviteRE, "Edison Edgar Franklin"),
			Timestamp: sent,
		},
		{
			To:        "gon@hunters.com",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   fmt.Sprintf(emails.InviteRE, "Edison Edgar Franklin"),
			Timestamp: sent,
		},
	}
	mock.CheckEmails(s.T(), messages)
}
