package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// TenantCreates creates adds a new tenant to the database
func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err    error
		tenant *db.Tenant
	)

	if err = c.BindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Create an error for invalid tenant field
	if tenant.ID.Compare(ulid.ULID{}) != 0 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrInvalidTenantField))
		return
	}

	if tenant.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant name is required"))
		return
	}

	if tenant.EnvironmentType == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("environment type is required"))
		return
	}

	tenant = &db.Tenant{
		ID:              tenant.ID,
		Name:            tenant.Name,
		EnvironmentType: tenant.EnvironmentType,
		Created:         tenant.Created,
		Modified:        tenant.Modified,
	}

	// Create a tenant to be passed into the database
	if err = db.CreateTenant(c.Request.Context(), tenant); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant"))
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (s *Server) TenantList(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
