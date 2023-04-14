package quarterdeck

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
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
	var user *api.User
	if user, err = model.ToAPI(); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize user model to api")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve user"))
		return
	}

	c.JSON(http.StatusOK, user)

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

	// Validate the request from the API side.
	if err = user.ValidateUpdate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Sanity check: the URL endpoint and the user ID on the model match.
	if !ulids.IsZero(user.UserID) && user.UserID.Compare(userID) != 0 {
		c.Error(api.ErrModelIDMismatch)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrModelIDMismatch))
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
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user"))
		}
		return
	}

	// Populate the response from the model
	if user, err = model.ToAPI(); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize user model to API")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user"))
		return
	}
	c.JSON(http.StatusOK, user)
}

// The UserRoleUpdate endpoint updates the role of a user in the organization. If the
// role is not valid or user already has the role, a 400 error is returned.
func (s *Server) UserRoleUpdate(c *gin.Context) {
	var (
		err    error
		userID ulid.ULID
		req    *api.UpdateRoleRequest
		user   *api.User
		model  *models.User
		claims *tokens.Claims
	)

	// Retrieve ID component from the URL and parse it.
	if userID, err = ulid.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("id")).Msg("could not parse user id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		return
	}

	if err = c.BindJSON((&req)); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse update user request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	if ulids.IsZero(userID) {
		sentry.Warn(c).Err(err).Msg("id in URL path is zero")
		c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
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

	// Verify that a valid role was provided
	if req.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.MissingField("role")))
		return
	}

	if !perms.IsRole(req.Role) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown user role"))
		return
	}

	// Fetch the user from the database
	if model, err = models.GetUser(c.Request.Context(), userID, orgID); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrUserOrganization) {
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not fetch user from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user role"))
		return
	}

	var role string
	if role, err = model.Role(); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch user role from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user role"))
		return
	}

	if role == req.Role {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user already has the specified role"))
		return
	}

	// Attempt to update the role in the database
	if err = model.ChangeRole(c.Request.Context(), orgID, req.Role); err != nil {
		// Check the error returned
		var verr *models.ValidationError
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		case errors.Is(err, models.ErrNoOwnerRole):
			c.Error(err)
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user role"))
		case errors.Is(err, models.ErrOwnerRoleConstraint):
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse("organization must have at least one owner"))
		case errors.As(err, &verr):
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			sentry.Error(c).Err(err).Msg("could not update user role in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user role"))
		}
		return
	}

	// Populate the response from the model
	if user, err = model.ToAPI(); err != nil {
		sentry.Error(c).Err(err).Msg("could not serialize user model to api")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user role"))
		return
	}

	c.JSON(http.StatusOK, user)
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

	for _, model := range users {
		var user *api.User
		if user, err = model.ToAPI(); err != nil {
			sentry.Error(c).Err(err).Msg("could not serialize user model to api")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list users"))
			return
		}

		out.Users = append(out.Users, user)
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

// Delete a user by their ID.  This endpoint allows admins to delete a user from the
// organization in the requesting user's claims. If the user does not exist in any
// other organization, their account will also be deleted.
// TODO: determine all the components of this process (billing, removal of organization, etc)
func (s *Server) UserDelete(c *gin.Context) {
	var (
		err    error
		userID ulid.ULID
		orgID  ulid.ULID
		claims *tokens.Claims
		user   *models.User
	)

	// Parse the user ID from the URL
	if userID, err = ulids.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse user id from request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Retrieve the orgID from the claims
	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Warn(c).Msg("invalid user claims sent in request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Retrieve the user to be deleted
	if user, err = models.GetUser(c.Request.Context(), userID, orgID); err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		case errors.Is(err, models.ErrUserOrganization):
			c.Error(err)
			log.Warn().Msg("attempt to fetch user from a different organization")
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		default:
			sentry.Error(c).Err(err).Msg("could not get user from database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete user"))
		}
		return
	}

	// Completely remove the user from the organization
	if err = user.RemoveOrganization(c.Request.Context(), orgID); err != nil {
		sentry.Error(c).Err(err).Msg("could not remove user from organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete user"))
		return
	}

	c.Status(http.StatusNoContent)
}
