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

// TenantDetail retrieves a summary detail of a specified tenant.
func (s *Server) TenantDetail(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	// Get the tenant ID from the URL
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Get the specified tenant from the database
	if _, err = db.RetrieveTenant(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not retrieve tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// TenantDelete deletes a tenant from a user's request.
// This returns a 200 OK response, not an error response.
func (s *Server) TenantDelete(c *gin.Context) {
	var (
		err    error
		tenant *api.Reply
	)

	// Get the tenant ID from the URL
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Delete the tenant from the database
	if err = db.DeleteTenant(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Str("tenantID", tenantID.String()).Msg("could not delete tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete tenant"))
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (s *Server) TenantUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
