package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

const (
	ContextUserClaims = "quarterdeck_user_claims"
)

// Authorize is a middleware that requires specific permissions in an authenticated
// user's claims. If those permissions do not match or the request is unauthenticated
// the middleware returns a 401 response. The Authorize middleware must be chained
// following the Authenticate middleware.
// TODO: move this to the auth.go file
func Authorize(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetClaims(c)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoAuthorization))
			return
		}

		if !claims.HasAllPermissions(permissions...) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoPermission))
			return
		}

		c.Next()
	}
}

// GetClaims fetches and parses Quarterdeck claims from the gin context. Returns an
// error if no claims exist on the context; panics if the claims are not the correct
// type -- however the panic should be recovered by middleware.
func GetClaims(c *gin.Context) (*tokens.Claims, error) {
	// TODO: ensure these claims use the same context key as the Authenicate middleware
	claims, exists := c.Get(ContextUserClaims)
	if !exists {
		return nil, ErrNoClaims
	}
	return claims.(*tokens.Claims), nil
}
