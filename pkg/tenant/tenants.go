package tenant

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

// TenantList retrieves all tenants assigned to an organization and
// returns a 200 OK response.
//
// Route: /tenant
func (s *Server) TenantList(c *gin.Context) {
	var (
		err   error
		orgID ulid.ULID
	)

	// Tenants exist on organizations
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get tenants from the database and return a 500 response if not successful.
	var tenants []*db.Tenant
	if tenants, err = db.ListTenants(c.Request.Context(), orgID); err != nil {
		log.Error().Err(err).Msg("could not fetch tenants from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list tenants"))
		return
	}

	// Build the response
	out := &api.TenantPage{Tenants: make([]*api.Tenant, 0)}

	// Loop over db.Tenant and retrieve each tenant.
	for _, dbTenant := range tenants {
		out.Tenants = append(out.Tenants, dbTenant.ToAPI())
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
		log.Warn().Err(err).Msg("could not bind tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
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
		log.Error().Err(err).Msg("could not create tenant in database")
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
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Get the specified tenant from the database and return a 404 response
	// if it cannot be retrieved.
	var tenant *db.Tenant
	if tenant, err = db.RetrieveTenant(c.Request.Context(), orgID, tenantID); err != nil {
		log.Error().Err(err).Str("tenantID", tenantID.String()).Msg("could not retrieve tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
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
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Bind the user request with JSON and return a 400 response if
	// binding is not successful.
	if err = c.BindJSON(&tenant); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
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

	// Get the specified tenant from the database and return a 404 response
	// if it cannot be retrieved.
	var t *db.Tenant
	if t, err = db.RetrieveTenant(c.Request.Context(), orgID, tenantID); err != nil {
		log.Error().Err(err).Msg("could not retrieve tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Update tenant in the database and return a 500 response if the
	// tenant record cannot be updated.
	if err := db.UpdateTenant(c.Request.Context(), t); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
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
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Delete the tenant and return a 404 response if it cannot be removed.
	if err = db.DeleteTenant(c.Request.Context(), orgID, tenantID); err != nil {
		log.Error().Err(err).Str("tenantID", tenantID.String()).Msg("could not delete tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}
	c.Status(http.StatusOK)
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
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
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
		log.Error().Str("tenant_id", id).Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Retrieve the tenant from the database
	var tenant *db.Tenant
	if tenant, err = db.RetrieveTenant(ctx, orgID, tenantID); err != nil {
		log.Error().Err(err).Str("tenant_id", id).Msg("could not retrieve tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Number of projects in the tenant
	var projects []*db.Project
	if projects, err = db.ListProjects(ctx, tenant.ID); err != nil {
		log.Error().Err(err).Str("tenant_id", id).Msg("could not retrieve projects in tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant stats"))
		return
	}
	totalProjects := len(projects)

	// Count topics and api keys in each project
	var totalTopics, totalKeys int
	for _, project := range projects {
		var topics []*db.Topic
		if topics, err = db.ListTopics(ctx, project.ID); err != nil {
			log.Error().Err(err).Str("project_id", project.ID.String()).Msg("could not retrieve topics in project")
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
				log.Error().Err(err).Str("project_id", project.ID.String()).Msg("could not retrieve api keys in project")
				c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not retrieve tenant stats"))
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
