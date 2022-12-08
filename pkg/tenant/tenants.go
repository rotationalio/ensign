package tenant

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// TenantCreates creates adds a new tenant to the database
func (s *Server) TenantCreate(c *gin.Context) {
	var (
		err error
		t   *api.Tenant
	)

	// TODO: Add authentication and authorization middleware

	if err = c.BindJSON(&t); err != nil {
		log.Warn().Err(err).Msg("could not bind tenant create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
	}

	if t.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrTenantIDRequired))
	}

	if t.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant name is required"))
	}

	if t.EnvironmentType == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("environment type is required"))
	}

	tenantID := ulid.Make()
	created := time.Now()
	modified := created

	tenant := &db.Tenant{
		ID:              tenantID,
		Name:            t.Name,
		EnvironmentType: t.EnvironmentType,
		Created:         created,
		Modified:        modified,
	}

	// Add tenant to the database
	if err = db.CreateTenant(c.Request.Context(), tenant); err != nil {
		log.Error().Err(err).Msg("could not create tenant in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant"))
	}

	c.JSON(http.StatusCreated, t)
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
