package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

func (s *Server) TenantList(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err error
		in  *api.Tenant
		out *api.Tenant
	)

	// TODO: Add authentication and authorization middleware

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
	}

	if in.ID == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrTenantIDRequired))
	}

	if in.TenantName == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant name is required"))
	}

	if in.EnvironmentType == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("environment type is required"))
	}

	// Add tenant to the database
	if err = db.CreateTenant(c.Request.Context(), &db.Tenant{}); err != nil {
		log.Error().Err(err).Msg("could not create tenant in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant"))
		return
	}

	// Build the response
	out = &api.Tenant{
		ID:              out.ID,
		TenantName:      out.TenantName,
		EnvironmentType: out.EnvironmentType,
	}

	c.JSON(http.StatusOK, out)
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
