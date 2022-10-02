package quarterdeck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (s *Server) ProjectList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) ProjectCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) ProjectDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) ProjectUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) ProjectDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}
