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

	tenantID := c.Param("tenantID")

	if tenantID == "" {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve tenant"))
		return
	}

	// TODO: Replace uuid with ulid
	req := &db.Tenant{
		ID: uuid.MustParse("6efd40b4-7035-47c1-afc5-d8142760e36c"),
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
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("temamt does not have an ID"))
		return
	}

	// Delete Tenant and
	if err = db.DeleteTenant(c.Request.Context(), uuid.MustParse("6efd40b4-7035-47c1-afc5-d8142760e36c")); err != nil {
		log.Error().Err(err).Msg(db.ErrMissingID.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete tenant"))
		return
	}
	c.JSON(http.StatusOK, out)
}
