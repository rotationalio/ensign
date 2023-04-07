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
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

func (s *Server) UserDetail(c *gin.Context) {
	var (
		err    error
		userID ulid.ULID
		model  *models.User
		claims *tokens.Claims
	)

	// Retrieve ID component from the URL and parse it.
	if userID, err = ulid.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("id")).Msg("could not parse user id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	//retrieve the orgID from the claims and check if it is valid
	orgID := claims.ParseOrgID()
	if ulids.IsZero(orgID) {
		sentry.Warn(c).Msg("invalid user claims sent in request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	if model, err = models.GetUser(c.Request.Context(), userID, orgID); err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		case errors.Is(err, models.ErrUserOrganization):
			c.Error(err)
			log.Warn().Msg("attempt to fetch user from different organization")
			c.JSON(http.StatusForbidden, api.ErrorResponse("requester is not authorized to access this user"))
		default:
			sentry.Error(c).Err(err).Msg("could not get user from database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}
		return
	}

	// Populate the response from the model
	c.JSON(http.StatusOK, model.ToAPI())

}

func (s *Server) UserUpdate(c *gin.Context) {
	//TODO: add functionality to update email
	var (
		err    error
		userID ulid.ULID
		user   *api.User
		model  *models.User
		claims *tokens.Claims
	)

	// Retrieve ID component from the URL and parse it.
	if userID, err = ulid.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("id")).Msg("could not parse user id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		return
	}

	if err = c.BindJSON((&user)); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse update user request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Sanity check: the URL endpoint and the user ID on the model match.
	if !ulids.IsZero(user.UserID) && user.UserID.Compare(userID) != 0 {
		c.Error(api.ErrModelIDMismatch)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrModelIDMismatch))
		return
	}

	// Validate the request from the API side.
	if err = user.ValidateUpdate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	//retrieve the orgID and userID from the claims and check if they are valid
	orgID := claims.ParseOrgID()
	requesterID := claims.ParseUserID()
	if ulids.IsZero(orgID) || ulids.IsZero(requesterID) {
		sentry.Warn(c).Msg("invalid user claims sent in request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Create a thin model to update in the database
	model = &models.User{
		ID:   user.UserID,
		Name: user.Name,
	}

	// Attempt to update the name in the database
	if err = model.Update(c.Request.Context(), orgID); err != nil {
		// Check if the error is a not found error or a validation error.
		var verr *models.ValidationError

		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		case errors.As(err, &verr):
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			sentry.Error(c).Err(err).Msg("could not update user in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}
		return
	}

	// Populate the response from the model
	c.JSON(http.StatusOK, model.ToAPI())
}

func (s *Server) UserList(c *gin.Context) {
	var (
		err                error
		orgID              ulid.ULID
		users              []*models.User
		nextPage, prevPage *pagination.Cursor
		claims             *tokens.Claims
		out                *api.UserList
	)

	query := &api.UserPageQuery{}
	if err = c.BindQuery(query); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse user page query request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	if query.NextPageToken != "" {
		if prevPage, err = pagination.Parse(query.NextPageToken); err != nil {
			sentry.Warn(c).Err(err).Msg("could not parse next page token")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	} else {
		prevPage = pagination.New("", "", int32(query.PageSize))
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Warn(c).Msg("invalid user claims sent in request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	if users, nextPage, err = models.ListUsers(c.Request.Context(), orgID, prevPage); err != nil {
		// Check if the error is a not found error or a validation error.
		var verr *models.ValidationError

		switch {
		case errors.Is(err, models.ErrNotFound):
			// TODO: can this error happen or is an empty page returned?
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		case errors.As(err, &verr):
			c.Error(verr)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			sentry.Error(c).Err(err).Msg("could not list users in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}
		return
	}

	// Prepare response
	out = &api.UserList{
		Users: make([]*api.User, 0, len(users)),
	}

	for _, user := range users {
		out.Users = append(out.Users, user.ToAPI())
	}

	// If a next page token is available, add it to the response.
	if nextPage != nil {
		if out.NextPageToken, err = nextPage.NextPageToken(); err != nil {
			sentry.Error(c).Err(err).Msg("could not create next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

// Delete a user by their ID.  This endpoint allows admins to delete a user from the organization
// TODO: determine all the components of this process (billing, removal of organization, etc)
func (s *Server) UserDelete(c *gin.Context) {
	sentry.Warn(c).Msg("user delete is not implemented")
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

// Invite a user to the organization. This is an authenticated endpoint that sends an
// invitation email to an email address. The link in the email contains a token that
// can only be verified by Quarterdeck. New users are not created in the database until
// they accept the email invitation.
func (s *Server) UserInvite(c *gin.Context) {
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
