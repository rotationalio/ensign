package quarterdeck

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/metrics"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
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
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
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
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
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

func (s *Server) ProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
		claims  *tokens.Claims
	)

	// Bind the Project request to the project data structure
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Validate the request from the API side.
	if err = project.Validate(); err != nil {
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

	// Create the OrganizationProjects mapping database model
	model := &models.OrganizationProject{
		OrgID:     claims.ParseOrgID(),
		ProjectID: project.ProjectID,
	}

	// Save the model to the database
	if err = model.Save(c.Request.Context()); err != nil {
		switch err.(type) {
		case *models.ValidationError:
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		default:
			sentry.Error(c).Err(err).Msg("could not create project")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
		}
		return
	}

	// Update the response to send to the user
	project.OrgID = model.OrgID
	project.ProjectID = model.ProjectID
	project.Created, _ = model.GetCreated()
	project.Modified, _ = model.GetModified()

	// Increment total number of projects in prometheus
	metrics.Projects.WithLabelValues(ServiceName).Inc()
	c.JSON(http.StatusOK, project)
}

func (s *Server) ProjectAccess(c *gin.Context) {
	var (
		err     error
		exists  bool
		project *api.Project
		claims  *tokens.Claims
		creds   *api.LoginReply
	)

	// Bind the Project request to the project data structure
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project access request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Validate the request from the API side.
	if err = project.Validate(); err != nil {
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

	// Check that the orgID claim is allowed access to the project
	model := &models.OrganizationProject{OrgID: claims.ParseOrgID(), ProjectID: project.ProjectID}
	if exists, err = model.Exists(c.Request.Context()); !exists || err != nil {
		if !exists {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown project id"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve organization project mapping from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete project access request"))
		return
	}

	// Prepare the credentials for the access token to make the one time access.
	now := time.Now()
	ota := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ulids.New().String(),
			Subject:   claims.Subject,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(OneTimeAccessDuration)),
		},
		OrgID:       claims.OrgID,
		ProjectID:   project.ProjectID.String(),
		Permissions: make([]string, 0, 4),
	}

	// Add only the user permissions related to topics to these claims -- whatever
	// access to topics the user has, so to will the one time access claims.
	for _, permission := range claims.Permissions {
		if permissions.InGroup(permission, permissions.PrefixTopics) {
			ota.Permissions = append(ota.Permissions, permission)
		}
	}

	creds = &api.LoginReply{}
	token := s.tokens.CreateToken(ota)
	if creds.AccessToken, err = s.tokens.Sign(token); err != nil {
		sentry.Error(c).Err(err).Msg("could not sign token")
		c.JSON(http.StatusInternalServerError, "could not complete project access request")
		return
	}

	c.JSON(http.StatusOK, creds)
}
