package quarterdeck

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
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
		c.Error(err)
		c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	//retrieve the orgID from the claims and check if it is valid
	orgID := claims.ParseOrgID()
	if ulids.IsZero(orgID) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid user claims"))
		return
	}

	if model, err = models.GetUser(c.Request.Context(), userID, orgID); err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		case errors.Is(err, models.ErrUserOrganization):
			log.Warn().Msg("attempt to fetch user from different organization")
			c.JSON(http.StatusForbidden, api.ErrorResponse("requester is not authorized to access this user"))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}
		c.Error(err)
		return
	}

	// Populate the response from the model
	c.JSON(http.StatusOK, model.ToAPI(c.Request.Context()))

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
		c.Error(err)
		c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		return
	}

	if err = c.BindJSON((&user)); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
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
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	//retrieve the orgID and userID from the claims and check if they are valid
	orgID := claims.ParseOrgID()
	requesterID := claims.ParseUserID()
	if ulids.IsZero(orgID) || ulids.IsZero(requesterID) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid user claims"))
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
			c.JSON(http.StatusNotFound, api.ErrorResponse("user id not found"))
		case errors.As(err, &verr):
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}

		c.Error(err)
		return
	}

	// Populate the response from the model
	c.JSON(http.StatusOK, model.ToAPI(c.Request.Context()))
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
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query"))
		return
	}

	if query.NextPageToken != "" {
		if prevPage, err = pagination.Parse(query.NextPageToken); err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	} else {
		prevPage = pagination.New("", "", int32(query.PageSize))
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	if users, nextPage, err = models.ListUsers(c.Request.Context(), orgID, prevPage); err != nil {
		// Check if the error is a not found error or a validation error.
		var verr *models.ValidationError

		switch {
		case errors.Is(err, models.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse("user not found"))
		case errors.As(err, &verr):
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}

		c.Error(err)
		return
	}

	// Prepare response
	out = &api.UserList{
		Users: make([]*api.User, 0, len(users)),
	}

	for _, user := range users {
		out.Users = append(out.Users, user.ToAPI(c.Request.Context()))
	}

	// If a next page token is available, add it to the response.
	if nextPage != nil {
		if out.NextPageToken, err = nextPage.NextPageToken(); err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}
