package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

func TenantList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func TenantCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantDetail(c *gin.Context) {
	var (
		err    error
		tenant *api.Tenant
	)

	if tenant.ID == "" {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	// TODO: Replace uuid with ulid
	req := &db.Tenant{
		ID: uuid.UUID{},
	}

	if _, err = db.RetrieveTenant(c.Request.Context(), req.ID); err != nil {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func TenantUpdate(c *gin.Context) {
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
		ID: uuid.UUID{},
	}

	// Delete Tenant and
	if err = db.DeleteTenant(c.Request.Context(), req.ID); err != nil {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	c.JSON(http.StatusOK, out)
}
