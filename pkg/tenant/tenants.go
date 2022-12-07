package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// CreateTenant creates a new tenant in the database
func (s *Server) CreateTenant(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	if err = c.BindJSON(&tenant); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if tenant.ID != "" {
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

	// Add Tenant to the database
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

func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	// TODO: Add authentication and authorization middleware
	if err = c.BindJSON(&tenant); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
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

	// TODO: Add tenant to the database
	if err = db.CreateTenant(c.Request.Context(), tenant); err != nil {
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
	var (
		err error
		out *api.Reply
	)

	tenantID := c.Param("tenantID")

	if tenantID == "" {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	// TODO: Replace uuid with ulid
	req := &db.Tenant{
		ID: ulid.Make(),
	}

	// Delete Tenant and
	if err = db.DeleteTenant(c.Request.Context(), req.ID); err != nil {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	c.JSON(http.StatusOK, out)
}
