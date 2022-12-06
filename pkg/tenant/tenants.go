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
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
