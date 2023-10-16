package quarterdeck

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/metrics"
	"github.com/rotationalio/ensign/pkg/utils/radish"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
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
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return s.SendInviteEmail(user, org, invite)
	}),
		radish.WithRetries(3),
		radish.WithBackoff(backoff.NewExponentialBackOff()),
		radish.WithErrorf("could not send invite email to user %s", user.ID.String()),
		radish.WithContext(sentry.CloneContext(c)),
	)

	out := &api.UserInviteReply{
		UserID:       invite.UserID,
		OrgID:        invite.OrgID,
		Email:        invite.Email,
		Role:         invite.Role,
		Name:         invite.Name(),
		Organization: org.Name,
		Workspace:    org.Domain,
		ExpiresAt:    invite.Expires,
		CreatedBy:    user.ID,
		Created:      invite.Created,
	}
	c.JSON(http.StatusOK, out)
}

// InviteAccept accepts a user invitation. This is an authenticated endpoint, so in
// order to accept the invite the email address of the requesting user must match the
// email in the invitation. The user must also already exist in the database.
// Unauthenticated users can call the Login endpoint to accept an invitation if they do
// not have credentials.
func (s *Server) InviteAccept(c *gin.Context) {
	var (
		err error
		ctx context.Context
		req *api.UserInviteToken
	)

	// Fetch the user's claims from the context
	var claims *tokens.Claims
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Parse the invite token from the request
	if err = c.BindJSON(&req); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse accept invite request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	if req.Token == "" {
		sentry.Warn(c).Msg("missing invite token in request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	// Retrieve the user from the database
	var user *models.User
	ctx = c.Request.Context()
	if user, err = models.GetUser(ctx, claims.ParseUserID(), ulids.Null); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			log.Debug().Msg("could not find user by ID")
			c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrTryLoginAgain))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve user from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// User must be verified to accept the invite
	if !user.EmailVerified {
		log.Debug().Msg("user has not verified their email address")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrVerifyEmail))
		return
	}

	// Accept the invite and handle errors
	if s.acceptInvite(c, user, req.Token) != nil {
		return
	}

	// Construct the claims for the user
	var newClaims *tokens.Claims
	if newClaims, err = user.NewClaims(ctx); err != nil {
		sentry.Error(c).Err(err).Msg("could not create new claims for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	out := &api.LoginReply{
		LastLogin: user.LastLogin.String,
	}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(newClaims); err != nil {
		sentry.Error(c).Err(err).Msg("could not create access and refresh token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Update the user's last login in a Go routine
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return user.UpdateLastLogin(ctx)
	}), radish.WithErrorf("could not update last login timestamp for user %s", user.ID.String()),
		radish.WithContext(sentry.CloneContext(c)),
		radish.WithTimeout(1*time.Minute),
	)

	c.JSON(http.StatusOK, out)
}

// acceptInvite is a helper method that accepts an invitation and updates the user's
// organization membership in the database. This method handles the logging and error
// responses if the invitation is invalid or there is a database error.
func (s *Server) acceptInvite(c *gin.Context, user *models.User, token string) (err error) {
	// Retrieve the invitation from the database
	var invite *models.UserInvitation
	if invite, err = models.GetUserInvite(c.Request.Context(), token); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			log.Debug().Msg("could not find invite token")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrRequestNewInvite))
			metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "invite token not found").Inc()
			return err
		}

		sentry.Error(c).Err(err).Msg("could not retrieve the user invite from the database")
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "could not get invite token").Inc()
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return err
	}

	// Verify that the invite is for this email address
	if err = invite.Validate(user.Email); err != nil {
		switch {
		case errors.Is(err, db.ErrTokenExpired):
			log.Debug().Msg("invite token expired")
			metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "invite token expired").Inc()
		case errors.Is(err, models.ErrInvalidEmail):
			log.Debug().Msg("user email does not match invite")
			metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "invite token wrong email").Inc()
		default:
			sentry.Error(c).Err(err).Msg("invalid invite token")
			metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "invalid invite token").Inc()
		}

		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrRequestNewInvite))
		return err
	}

	// Add the user to the organization
	org := &models.Organization{
		ID: invite.OrgID,
	}
	if err = user.AddOrganization(c.Request.Context(), org, invite.Role); err != nil {
		sentry.Error(c).Err(err).Msg("could not add user to organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "could not add user to organization").Inc()
		return err
	}

	// Set the user's organization to the new one
	if err = user.SwitchOrganization(c.Request.Context(), org.ID); err != nil {
		sentry.Error(c).Err(err).Msg("could not switch user to new organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "could not switch user to new organization").Inc()
		return err
	}

	// At this point the user should be able to log into the org, so we can delete the invite
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return models.DeleteInvite(ctx, invite.Token)
	}), radish.WithErrorf("could not delete user invite with token %s", invite.Token),
		radish.WithContext(sentry.CloneContext(c)),
		radish.WithTimeout(1*time.Minute),
	)
	return nil
}
