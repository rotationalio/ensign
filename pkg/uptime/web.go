package uptime

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

// content holds our static web server content.
//
//go:embed all:templates
//go:embed all:static
var content embed.FS

func (s *Server) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
