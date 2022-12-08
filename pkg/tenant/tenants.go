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
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TenantUpdate will amend the name or environment type of a tenant
func (s *Server) TenantUpdate(c *gin.Context) {
	var (
		err error
		in  *db.Tenant
		out *api.Tenant
	)

	// TODO: authentication and authorization middleware

	// Bind JSON insertion
	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse tenant update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
	}

	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Debug().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Update tenant with a new ID and add to the database
	if err = db.UpdateTenant(c.Request.Context(), &db.Tenant{ID: tenantID}); err != nil {
		log.Error().Err(err).Msg("could not save tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update tenant"))
	}

	// Prepare what will go out
	out = &api.Tenant{
		ID:              out.ID,
		Name:            out.Name,
		EnvironmentType: out.EnvironmentType,
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) TenantDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
