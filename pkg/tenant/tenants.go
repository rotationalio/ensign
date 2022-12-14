package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

func (s *Server) TenantList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TenantCreate adds a new tenant to the database
// and returns a 201 StatusCreated response.
//
// Route: /tenant
func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err error
		t   *api.Tenant
	)

	// TODO: Add authentication and authorization middleware

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
		c.JSON(http.StatusBadRequest, api.ErrorResponse("environment type is required"))
		return
	}

	tenant := &db.Tenant{
		Name:            t.Name,
		EnvironmentType: t.EnvironmentType,
	}

	if err = db.CreateTenant(c.Request.Context(), tenant); err != nil {
		log.Error().Err(err).Msg("could not create tenant in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant"))
		return
	}

	out := &api.Tenant{
		ID:              tenant.ID.String(),
		Name:            tenant.Name,
		EnvironmentType: tenant.EnvironmentType,
	}

	c.JSON(http.StatusCreated, out)
}

// TenantDetail retrieves a summary detail of a tenant by its id and
// returns a 200 OK response.
//
// Route: /tenant/:tenantID
func (s *Server) TenantDetail(c *gin.Context) {
	var (
		err   error
		reply *api.Tenant
	)

	// Get the tenant ID from the URL and return a 400 if the
	// tenant does not exist.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Get the specified tenant from the database and return a 500 response
	// if it cannot be retrieved.
	var tenant *db.Tenant
	if tenant, err = db.RetrieveTenant(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not retrieve tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	reply = &api.Tenant{
		ID:              tenant.ID.String(),
		Name:            tenant.Name,
		EnvironmentType: tenant.EnvironmentType,
	}
	c.JSON(http.StatusOK, reply)
}

// TenantDelete deletes a tenant from a user's request with a given
// id and returns a 200 OK response instead of an an error response.
//
// Route: /tenant/:tenantID
func (s *Server) TenantDelete(c *gin.Context) {
	var (
		err error
	)

	// Get the tenant ID from the URL and return a 400 if the
	// tenant does not exist.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Delete the tenant and return a 404 response if it cannot be removed.
	if err = db.DeleteTenant(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Str("tenantID", tenantID.String()).Msg("could not delete tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not delete tenant"))
		return
	}
	c.Status(http.StatusOK)
}

// TenantUpdate will update a tenants records and
// returns a 200 OK response.
//
// Route: /tenant/:tenantID
func (s *Server) TenantUpdate(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	// TODO: authentication and authorization middleware

	// Get the tenant ID from the URL and return a 400 if the tenant
	// ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
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

	// Get the specified tenant from the database and return a 500 response
	// if it cannot be retrieved.
	var t *db.Tenant
	if t, err = db.RetrieveTenant(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not retrieve tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	// Update tenant in the database and return a 404 response if the
	// tenant record cannot be updated.
	if err := db.UpdateTenant(c.Request.Context(), t); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update tenant"))
		return
	}

	c.JSON(http.StatusOK, tenant)
}
