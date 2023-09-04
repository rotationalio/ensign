package tenant

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
)

// InvitePreview returns "preview" information about an invite given a token. This
// endpoint must not be authenticated because unauthorized users should be able to
// accept organization invitations. Frontends should use this endpoint to validate an
// invitation token after the user has clicked on an invitation link in their email.
// The preview must contain enough information so the user knows which organization
// they are joining and also whether or not the email address is already registered to
// an account. This allows frontends to know whether or not to prompt the user to
// login or to create a new account.
//
// Route: /invites/:token
func (s *Server) InvitePreview(c *gin.Context) {
	var err error

	token := c.Param("token")

	// Call Quarterdeck to retrieve the invite preview.
	var rep *qd.UserInvitePreview
	if rep, err = s.quarterdeck.InvitePreview(c.Request.Context(), token); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Create the preview response
	out := &api.MemberInvitePreview{
		Email:       rep.Email,
		OrgName:     rep.OrgName,
		InviterName: rep.InviterName,
		Role:        rep.Role,
		HasAccount:  rep.UserExists,
	}
	c.JSON(http.StatusOK, out)
}

// InviteAccept is an authenticated endpoint to accept an invitation to join an
// organization. The invitation token must be provided in the request body, and the
// email in the user claims must match the email address in the token. If the
// invitation is invalid this endpoint returns a 404. If successful, the user is logged
// into the organization and credentials are set as cookies. Frontends should use this
// endpoint when a user is already logged in and is accepting an invitation. If the
// user is not logged in, the Login endpoint should be used instead.
//
// Route: /invites/accept
func (s *Server) InviteAccept(c *gin.Context) {
	var (
		ctx context.Context
		req *api.MemberInviteToken
		err error
	)

	// User credentials are required for the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user credentials from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
	}

	// Parse the token in the request body
	if err = c.BindJSON(&req); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse accept invite token request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	if req.Token == "" {
		sentry.Warn(c).Msg("missing token in accept invite request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	// Create the Quarterdeck request with the token
	acceptRequest := &qd.UserInviteToken{
		Token: req.Token,
	}

	// Call Quarterdeck to accept the invite
	var rep *qd.LoginReply
	if rep, err = s.quarterdeck.InviteAccept(ctx, acceptRequest); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Set the access and refresh tokens as cookies for the front-end
	if err := middleware.SetAuthCookies(c, rep.AccessToken, rep.RefreshToken, s.conf.Auth.CookieDomain); err != nil {
		sentry.Error(c).Err(err).Msg("could not set access and refresh token cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Return 200 response
	c.Status(http.StatusOK)
}
