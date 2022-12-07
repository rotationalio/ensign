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

func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	// TODO: Add authentication and authorization middleware

	if err = c.BindJSON(&tenant); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
	}

	if tenant.ID == "" {
		log.Error().Err(err).Msg("tenant does not have an ID")
	}

	if tenant.Name == "" {
		log.Warn().Err(err).Msg("tenant does not have an name")
	}

	if tenant.EnvironmentType == "" {
		log.Error().Err(err).Msg("tenant does not have an environment type")
	}

	if err = db.CreateTenant(c.Request.Context(), &db.Tenant{ID: ulid.Make()}); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant"))
	}

	c.JSON(http.StatusCreated, tenant)
}

func (s *Server) TenantDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantUpdate(c *gin.Context) {
	var (
		err error
		in  *api.Tenant
		out *api.Tenant
	)

	// TODO: authentication and authorization middleware

	// Get ID from param
	tenantID := c.Param("tenantID")

	// Bind JSON insertion
	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
	}

	// Confirm ID exists
	if tenantID == "" {
		log.Error().Err(err).Msg("tenant does not have an ID")
	}

	// Update tenant with a new ID and add to the database
	if err = db.UpdateTenant(c.Request.Context(), &db.Tenant{}); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update tenant"))
	}

	// Prepare what will go out
	out = &api.Tenant{
		ID:              in.ID,
		Name:            in.Name,
		EnvironmentType: in.EnvironmentType,
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) TenantDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
