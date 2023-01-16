package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/stretchr/testify/require"
)

func TestAuthorize(t *testing.T) {
	// TODO: create the test authentication server
	// srv, err := authtest.NewServer()
	// require.NoError(t, err, "could not start authtest server")
	// defer srv.Close()

	// TODO: create the authenticate middleware
	// authenticate, err := middleware.Authenticate(
	// 	middleware.WithJWKSEndpoint(srv.KeysURL()),
	// 	middleware.WithAudience(authtest.Audience),
	// 	middleware.WithIssuer(authtest.Issuer),
	// )
	// require.NoError(t, err, "could not create authenticate middleware")

	// Create Authorize middleware
	authorize := middleware.Authorize("foo:read", "foo:write")

	// Create a quick router for testing the middleware
	router := gin.Default()
	router.GET("/", authorize, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// If there is no authenticate middleware before authorize, expect a 401 response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// TODO: expect 401 with an authenticated request without the permissions

	// TODO: expect 200 with an authenticated request with the permissions
}
