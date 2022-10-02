package quarterdeck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (s *Server) APIKeyList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) APIKeyCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) APIKeyDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) APIKeyUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) APIKeyDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}
