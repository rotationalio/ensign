package quarterdeck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (s *Server) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) Authenticate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) Refresh(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}
