package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

// TenantList retrieves tenants assigned to a specified organization and
// returns a 200 OK response.
//
// Route: /tenant
func (s *Server) TenantList(c *gin.Context) {
	var (
		err        error
		orgID      ulid.ULID
		next, prev *pg.Cursor
	)

	// Tenants exist on organizations
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	query := &api.PageQuery{}
	if err = c.BindQuery(query); err != nil {
		log.Error().Err(err).Msg("could not parse query request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query request"))
		return
	}

	if query.NextPageToken != "" {
		if prev, err = pg.Parse(query.NextPageToken); err != nil {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse next page token"))
			return
		}
	} else {
		prev = pg.New("", "", int32(query.PageSize))
	}

	// Get tenants from the database and return a 500 response if not successful.
	var tenants []*db.Tenant
	if tenants, next, err = db.ListTenants(c.Request.Context(), orgID, prev); err != nil {
		sentry.Error(c).Err(err).Msg("could not list tenants in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list tenants"))
		return
	}

	// Build the response
	out := &api.TenantPage{Tenants: make([]*api.Tenant, 0)}

	// Loop over db.Tenant and retrieve each tenant.
	for _, dbTenant := range tenants {
		out.Tenants = append(out.Tenants, dbTenant.ToAPI())
	}

	if next != nil {
		if out.NextPageToken, err = next.NextPageToken(); err != nil {
			log.Error().Err(err).Msg("could not set next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list tenants"))
			return
		}
	}
	c.JSON(http.StatusOK, out)
}

// TenantCreate adds a new tenant to the database
// and returns a 201 StatusCreated response.
//
// Route: /tenant
func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err   error
		t     *api.Tenant
		orgID ulid.ULID
	)

	// Tenants exist on organizations
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Bind the user request with JSON and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&t); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify that a tenant ID does not exist and return a 400 response if the
	// tenant id exists.
	if t.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant id cannot be specified on create"))
		return
	}

	// Verify that a tenant name has been provided and return a 400 response
	// if the tenant name does not exist.
	if t.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant name is required"))
		return
	}

	// Verify that an environment type has been provided and return a 400 response
	// if the tenant environment type does not exist.
	if t.EnvironmentType == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant environment type is required"))
		return
	}

	tenant := &db.Tenant{
		OrgID:           orgID,
		Name:            t.Name,
		EnvironmentType: t.EnvironmentType,
	}

	if err = db.CreateTenant(c.Request.Context(), tenant); err != nil {
		sentry.Error(c).Err(err).Msg("could not create tenant in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant"))
		return
	}

	c.JSON(http.StatusCreated, tenant.ToAPI())
}

// TenantDetail retrieves a summary detail of a tenant by its ID and
// returns a 200 OK response.
//
// Route: /tenant/:tenantID
func (s *Server) TenantDetail(c *gin.Context) {
	var err error

	// Tenants exist in organizations
	// This method handles the logging and error response
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the tenant ID from the URL and return a 400 if the
	// tenant does not exist.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", c.Param("tenantID")).Msg("could not parse tenant id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Get the specified tenant from the database and return a 404 response
	// if it cannot be retrieved.
	var tenant *db.Tenant
	if tenant, err = db.RetrieveTenant(c.Request.Context(), orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve tenant from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	c.JSON(http.StatusOK, tenant.ToAPI())
}

// TenantUpdate will update a tenants record and
// returns a 200 OK response.
//
// Route: /tenant/:tenantID
func (s *Server) TenantUpdate(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	// Tenants exist in organizations
	// This method handles the logging and error response
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the tenant ID from the URL and return a 400 if the tenant
	// ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", c.Param("tenantID")).Msg("could not parse tenant id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Bind the user request with JSON and return a 400 response if
	// binding is not successful.
	if err = c.BindJSON(&tenant); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse update tenant request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify the tenant name exists and return a 400 response if it does not exist.
	if tenant.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant name is required"))
		return
	}

	// Verify the tenant environment type exists and return a 400 response if it does
	// not exist.
	if tenant.EnvironmentType == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant environment type is required"))
		return
	}

	// Get the specified tenant from the database.
	var t *db.Tenant
	if t, err = db.RetrieveTenant(c.Request.Context(), orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve tenant from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update tenant"))
		return
	}

	// Update all user provided fields.
	t.Name = tenant.Name
	t.EnvironmentType = tenant.EnvironmentType

	// Update tenant in the database.
	if err = db.UpdateTenant(c.Request.Context(), t); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("Could not update tenant in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update tenant"))
		return
	}

	c.JSON(http.StatusOK, t.ToAPI())
}

// TenantDelete deletes a tenant from a user's request with a given
// ID and returns a 200 OK response instead of an an error response.
//
// Route: /tenant/:tenantID
func (s *Server) TenantDelete(c *gin.Context) {
	var (
		err error
	)

	// Tenants exist in organizations
	// This method handles the logging and error response
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the tenant ID from the URL and return a 400 if the
	// tenant does not exist.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", c.Param("tenantID")).Msg("could not parse tenant id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Delete the tenant from the database.
	if err = db.DeleteTenant(c.Request.Context(), orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not delete tenant from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete tenant"))
		return
	}
	c.Status(http.StatusNoContent)
}

// TenantStats is a statistical view endpoint which returns high level counts of
// resources associated with a single Tenant.
//
// Route: /tenant/:tenantID/stats
func (s *Server) TenantStats(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to retrieve api keys from Quarterdeck
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Tenants exist in organizations
	// This method handles the logging and error response
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the tenantID from the URL
	id := c.Param("tenantID")
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(id); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", c.Param("tenantID")).Msg("could not parse tenant id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Retrieve the tenant from the database
	var tenant *db.Tenant
	if tenant, err = db.RetrieveTenant(ctx, orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve tenant from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	// Verify orgID from context matches the tenant orgID.
	db.VerifyOrg(ctx, orgID, tenant.OrgID)

	// TODO: Create list method that will not require pagination for this endpoint.
	// Set page size to return all projects and topics.
	getAll := &pg.Cursor{StartIndex: "", EndIndex: "", PageSize: 100}

	// Number of projects in the tenant
	var projects []*db.Project
	if projects, _, err = db.ListProjects(ctx, tenant.ID, getAll); err != nil {
		sentry.Error(c).Err(err).Msg("could not list projects in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant stats"))
		return
	}
	totalProjects := len(projects)

	// Count topics and api keys in each project
	var totalTopics, totalKeys int
	for _, project := range projects {
		var topics []*db.Topic
		if topics, _, err = db.ListTopics(ctx, project.ID, getAll); err != nil {
			sentry.Error(c).Err(err).Msg("could not list topics in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant stats"))
			return
		}
		totalTopics += len(topics)

		// API keys are stored in Quarterdeck
		req := &qd.APIPageQuery{
			ProjectID: project.ID.String(),
			PageSize:  100,
		}

		// We will always retrieve at least one page; it's possible but unlikely for a
		// project to have more than 100 API keys.
		totalKeys = 0
	keysLoop:
		for {
			var page *qd.APIKeyList
			if page, err = s.quarterdeck.APIKeyList(ctx, req); err != nil {
				sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
				api.ReplyQuarterdeckError(c, err)
				return
			}
			totalKeys += len(page.APIKeys)

			if page.NextPageToken == "" {
				break keysLoop
			}
			req.NextPageToken = page.NextPageToken
		}
	}

	// Build the standardized stats response for the frontend
	// TODO: Add data usage stats
	out := []*api.StatValue{
		{
			Name:  "projects",
			Value: float64(totalProjects),
		},
		{
			Name:  "topics",
			Value: float64(totalTopics),
		},
		{
			Name:  "keys",
			Value: float64(totalKeys),
		},
		{
			Name:    "storage",
			Value:   0,
			Units:   "GB",
			Percent: 0,
		},
	}

	c.JSON(http.StatusOK, out)
}
