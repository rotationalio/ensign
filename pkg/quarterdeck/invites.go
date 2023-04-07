package quarterdeck

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

// InvitePreview returns details for a user invitation. This is a publicly accessible
// endpoint because it has to be accessed by unauthenticated users without any claims.
// Although the endpoint is read-only, it could return sensitive information about an
// organization so this endpoint relies on the token being difficult to guess.
func (s *Server) InvitePreview(c *gin.Context) {
	var err error

	// Get the token from the URL param
	token := c.Param("token")

	// Retrieve the invite from the database
	var invite *models.UserInvitation
	if invite, err = models.GetUserInvite(c.Request.Context(), token); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid invitation"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve user invite")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve invitation"))
		return
	}

	// The following checks are security checks to make sure that someone is not trying
	// to guess our invite token structure; even though they are submitting the invite to view
	// it, it must be a validly issued invite from Quarterdeck. 
	// Ensure the role is a recognized role
	if !perms.IsRole(invite.Role) {
		sentry.Warn(c).Str("role", invite.Role).Msg("invalid role for user invite")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid invitation"))
		return
	}

	// Ensure the invite is valid and not expired
	if err = invite.Validate(invite.Email); err != nil {
		if errors.Is(err, db.ErrTokenExpired) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invitation has expired"))
			return
		}

		sentry.Error(c).Err(err).Msg("bad user invite token")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid invitation"))
		return
	}

	// Attempt to lookup the user by their email address with their "default" organization
	var user *models.User
	if user, err = models.GetUserEmail(c.Request.Context(), invite.Email, ulid.ULID{}); err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			sentry.Error(c).Err(err).Msg("could not retrieve user")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve invitation"))
			return
		}

		// Don't return an error to allow new users to create an account
	}

	// Fetch the invite organization from the database
	var org *models.Organization
	if org, err = models.GetOrg(c.Request.Context(), invite.OrgID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid invitation"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve organization in user invite")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve invitation"))
		return
	}

	// Fetch the inviter user from the database
	var inviter *models.User
	if inviter, err = models.GetUser(c.Request.Context(), invite.CreatedBy, invite.OrgID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid invitation"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve inviter in user invite")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve invitation"))
		return
	}

	// Return the invite preview
	out := &api.UserInvitePreview{
		Email:       invite.Email,
		OrgName:     org.Name,
		InviterName: inviter.Name,
		Role:        invite.Role,
		UserExists:  user != nil,
	}

	c.JSON(http.StatusOK, out)
}

// Invite a user to the organization. This is an authenticated endpoint that sends an
// invitation email to an email address. The link in the email contains a token that
// can only be verified by Quarterdeck. New users are not created in the database until
// they accept the email invitation.
func (s *Server) InviteCreate(c *gin.Context) {
	var err error

	// Parse the invite request
	req := &api.UserInviteRequest{}
	if err := c.BindJSON(req); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse user invite request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.MissingField("email")))
		return
	}

	if req.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.MissingField("role")))
		return
	}

	if !perms.IsRole(req.Role) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnknownUserRole))
		return
	}

	// Fetch the user claims from the context
	var claims *tokens.Claims
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	var orgID ulid.ULID
	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Warn(c).Msg("invalid user claims sent in request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Fetch the organization from the database
	var org *models.Organization
	if org, err = models.GetOrg(c.Request.Context(), orgID); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch organization from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not invite user"))
		return
	}

	// Fetch the requesting user to create the invitation
	var user *models.User
	if user, err = models.GetUser(c.Request.Context(), claims.ParseUserID(), orgID); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch requesting user from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not invite user"))
		return
	}

	// Create the invite in the database
	var invite *models.UserInvitation
	if invite, err = user.CreateInvite(c.Request.Context(), req.Email, req.Role); err != nil {
		sentry.Error(c).Err(err).Msg("could not create invite in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not invite user"))
		return
	}

	// Send the user invite with the token
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		return s.SendInviteEmail(user, org, invite)
	}),
		tasks.WithRetries(3),
		tasks.WithBackoff(backoff.NewExponentialBackOff()),
		tasks.WithError(fmt.Errorf("could not send invite email to user %s", user.ID.String())),
	)

	out := &api.UserInviteReply{
		UserID:    invite.UserID,
		OrgID:     invite.OrgID,
		Email:     invite.Email,
		Role:      invite.Role,
		ExpiresAt: invite.Expires,
		CreatedBy: user.ID,
		Created:   invite.Created,
	}
	c.JSON(http.StatusOK, out)
}
