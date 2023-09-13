package quarterdeck

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

const OneTimeAccessDuration = 5 * time.Minute

func (s *Server) OrganizationList(c *gin.Context) {
	var (
		err                error
		userID             ulid.ULID
		orgs               []*models.Organization
		nextPage, prevPage *pagination.Cursor
		claims             *tokens.Claims
		out                *api.OrganizationList
	)

	query := &api.OrganizationPageQuery{}
	if err = c.BindQuery(query); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse api page query")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	if query.NextPageToken != "" {
		if prevPage, err = pagination.Parse(query.NextPageToken); err != nil {
			sentry.Warn(c).Err(err).Msg("could not parse next page token")
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid next page token"))
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

	if userID = claims.ParseUserID(); ulids.IsZero(userID) {
		sentry.Warn(c).Msg("invalid user claims sent in request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	if orgs, nextPage, err = models.ListOrgs(c.Request.Context(), userID, prevPage); err != nil {
		// Check if the error is a not found error or a validation error.
		var verr *models.ValidationError

		switch {
		case errors.Is(err, models.ErrNotFound):
			// TODO: can this error happen or is an empty page returned?
			c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		case errors.As(err, &verr):
			c.Error(verr)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			sentry.Error(c).Err(err).Msg("could not list organizations in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list organizations"))
		}
		return
	}

	// Prepare response
	out = &api.OrganizationList{
		Organizations: make([]*api.Organization, 0, len(orgs)),
	}

	for _, org := range orgs {
		out.Organizations = append(out.Organizations, org.ToAPI())
	}

	// If a next page token is available, add it to the response.
	if nextPage != nil {
		if out.NextPageToken, err = nextPage.NextPageToken(); err != nil {
			sentry.Error(c).Err(err).Msg("could not create next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list organizations"))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

// Retrieve an organization by ID. Users are only allowed to retrieve their own
// organization.
// TODO: Eventually allow users to retrieve other organizations they are a part of.
func (s *Server) OrganizationDetail(c *gin.Context) {
	var (
		err    error
		orgID  ulid.ULID
		org    *models.Organization
		claims *tokens.Claims
	)

	// Parse the orgID passed in from the URL
	if orgID, err = ulid.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("id")).Msg("could not parse org id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// User claims are required to verify the user's organization
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// User must be a member of the organization
	if claims.OrgID != orgID.String() {
		sentry.Warn(c).Msg("user attempted to access organization they don't belong to")
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// Fetch the organization from the database
	if org, err = models.GetOrg(c.Request.Context(), orgID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not get organization from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		return
	}

	// Build the response from the model
	c.JSON(http.StatusOK, org.ToAPI())
}

// Update an organization by ID. Users are only allowed to update their own
// organization.
func (s *Server) OrganizationUpdate(c *gin.Context) {
	var (
		err    error
		orgID  ulid.ULID
		req    *api.Organization
		claims *tokens.Claims
	)

	// Parse the orgID passed in from the URL
	if orgID, err = ulid.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("id")).Msg("could not parse org id")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrOrganizationNotFound))
		return
	}

	// User claims are required to verify the user's organization
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// User must be a member of the organization
	if claims.OrgID != orgID.String() {
		sentry.Warn(c).Msg("user attempted to access organization they don't belong to")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrOrganizationNotFound))
		return
	}

	// Parse the organization from the request body
	if err = c.BindJSON(&req); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse organization from request body")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryOrganizationAgain))
		return
	}

	// Validate the organization request
	req.Name = strings.TrimSpace(req.Name)
	req.Domain = strings.TrimSpace(req.Domain)
	if err = req.ValidateUpdate(); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Create a database organization model for the update
	org := &models.Organization{
		ID:     req.ID,
		Name:   req.Name,
		Domain: req.Domain,
	}

	// Save the organization to the database
	if err = org.Save(c.Request.Context()); err != nil {
		var verr *models.ValidationError

		switch {
		case errors.Is(err, models.ErrNotFound):
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrOrganizationNotFound))
		case errors.Is(err, models.ErrDuplicate):
			c.Error(err)
			c.JSON(http.StatusConflict, api.ErrorResponse(responses.ErrDomainAlreadyExists))
		case errors.As(err, &verr):
			c.Error(verr)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			sentry.Error(c).Err(err).Msg("could not update organization in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}
		return
	}

	// Build the response from the model
	c.JSON(http.StatusOK, org.ToAPI())
}

// Lookup an organization's workspace by domain slug.
func (s *Server) WorkspaceLookup(c *gin.Context) {
	var (
		err    error
		in     *api.WorkspaceQuery
		out    *api.Workspace
		claims *tokens.Claims
		org    *models.Organization
	)

	in = &api.WorkspaceQuery{}
	if err = c.BindQuery(in); err != nil {
		sentry.Error(c).Err(err).Msg("could not bind workspace query")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryOrganizationAgain))
		return
	}

	if in.Domain == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrBadWorkspaceLookup))
		return
	}

	out = &api.Workspace{Domain: in.Domain, IsAvailable: false}

	if org, err = models.LookupWorkspace(c.Request.Context(), in.Domain); err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			if in.CheckAvailable {
				out.IsAvailable = true
				c.JSON(http.StatusOK, out)
				return
			}

			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrWorkspaceNotFound))
			return
		default:
			sentry.Error(c).Err(err).Msg("could not lookup workspace domain")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
			return
		}
	}

	// User claims are required to verify the user's organization
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// If the user is a member of the looked up organization, then supply extra details.
	if claims.OrgID == org.ID.String() {
		out.OrgID = org.ID
		out.Name = org.Name
	}

	c.JSON(http.StatusOK, out)
}
