package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) ProjectAPIKeyList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) ProjectAPIKeyCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyList(c *gin.Context) {
	// The following TODO task items will need to be
	// implemented for each endpoint.

	// TODO: Add authentication and authorization middleware
	// TODO: Identify top-level info
	// TODO: Parse and validate user input
	// TODO: Perform work on the request, e.g. database interactions,
	// sending notifications, accessing other services, etc.

	// Return response with the correct status code

	// TODO: Replace StatusNotImplemented with StatusOk and
	// replace "not yet implemented" message.
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
