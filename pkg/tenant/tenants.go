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

// TenantCreates adds a new tenant to the database
func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err error
		t   *api.Tenant
	)

	// TODO: Add authentication and authorization middleware

	if err = c.BindJSON(&t); err != nil {
		log.Warn().Err(err).Msg("could not bind tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	if t.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant id cannot be speicified on create"))
		return
	}

	if t.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant name is required"))
		return
	}

	if t.EnvironmentType == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("environment type is required"))
		return
	}

	tenant := &db.Tenant{
		Name:            t.Name,
		EnvironmentType: t.EnvironmentType,
	}

	// Add tenant to the database
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

func (s *Server) TenantDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
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

	// Bind the user request with JSON and return a 400 response if
	// binding is not successful.
	if err = c.BindJSON(&tenant); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Get the tenant ID from the URL and return a 400 if the tenant
	// does not exist.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
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

	// Prepare update request for insertion into the database.
	req := &db.Tenant{
		ID:              tenantID,
		Name:            tenant.Name,
		EnvironmentType: tenant.EnvironmentType,
	}

	// Update tenant in the database and return a 404 response if the
	// tenant record cannot be updated.
	if err := db.UpdateTenant(c.Request.Context(), req); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not update tenant"))
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (s *Server) TenantDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
