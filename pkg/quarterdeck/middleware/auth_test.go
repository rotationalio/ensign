package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	// Create the test authentication server
	srv, err := authtest.NewServer()
	require.NoError(t, err, "could not start authtest server")
	defer srv.Close()

	// Create the middleware
	authenticate, err := middleware.Authenticate(
		middleware.WithJWKSEndpoint(srv.KeysURL()),
		middleware.WithAudience(authtest.Audience),
		middleware.WithIssuer(authtest.Issuer),
	)
	require.NoError(t, err, "could not create authenticate middleware")

	// Create a quick router for testing with the handler and the middleware
	router := gin.Default()
	router.GET("/", authenticate, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Validate that an unauthenticated request is rejected
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// Validate that an authenticated request is accepted
	tks, err := srv.CreateAccessToken(&tokens.Claims{})
	require.NoError(t, err, "could not create access token")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tks)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAccessToken(t *testing.T) {
	// Create a test context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test when no authorization header is set
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := middleware.GetAccessToken(c)
	require.ErrorIs(t, err, middleware.ErrNoAuthorization)

	// Test when authorization header is set but empty
	c.Request.Header.Set("Authorization", "")
	_, err = middleware.GetAccessToken(c)
	require.ErrorIs(t, err, middleware.ErrNoAuthorization)

	// Test when authorization header is set but is not a bearer
	c.Request.Header.Set("Authorization", "Basic ZGVtbzpwQDU1dzByZA==")
	_, err = middleware.GetAccessToken(c)
	require.ErrorIs(t, err, middleware.ErrParseBearer)

	// Test when authorization header is set but is invalid
	c.Request.Header.Set("Authorization", "ZGVtbzpwQDU1dzByZA==")
	_, err = middleware.GetAccessToken(c)
	require.ErrorIs(t, err, middleware.ErrParseBearer)

	// Test when authorization header is set but Bearer is invalid
	c.Request.Header.Set("Authorization", "Bearer     ")
	_, err = middleware.GetAccessToken(c)
	require.ErrorIs(t, err, middleware.ErrParseBearer)

	// Test correct authorization header
	c.Request.Header.Set("Authorization", "Bearer JWT")
	tks, err := middleware.GetAccessToken(c)
	require.NoError(t, err)
	require.Equal(t, "JWT", tks)
}

func TestDefaultAuthOptions(t *testing.T) {
	// Should be able to create a default auth options with no extra input.
	conf := middleware.NewAuthOptions()
	require.NotZero(t, conf, "a zero valued configuration was returned")
	require.Equal(t, middleware.DefaultKeysURL, conf.KeysURL)
	require.Equal(t, middleware.DefaultAudience, conf.Audience)
	require.Equal(t, middleware.DefaultIssuer, conf.Issuer)
	require.Equal(t, middleware.DefaultMinRefreshInterval, conf.MinRefreshInterval)
	require.NotZero(t, conf.Context, "no context was created")
}

func TestAuthOptions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	conf := middleware.NewAuthOptions(
		middleware.WithJWKSEndpoint("http://localhost:8088/.well-known/jwks.json"),
		middleware.WithAudience("http://localhost:3000"),
		middleware.WithIssuer("http://localhost:8088"),
		middleware.WithMinRefreshInterval(67*time.Minute),
		middleware.WithContext(ctx),
	)

	cancel()
	require.NotZero(t, conf, "a zero valued configuration was returned")
	require.Equal(t, "http://localhost:8088/.well-known/jwks.json", conf.KeysURL)
	require.Equal(t, "http://localhost:3000", conf.Audience)
	require.Equal(t, "http://localhost:8088", conf.Issuer)
	require.Equal(t, 67*time.Minute, conf.MinRefreshInterval)
	require.ErrorIs(t, conf.Context.Err(), context.Canceled)
}

func TestAuthOptionsValidator(t *testing.T) {
	validator := &tokens.MockValidator{}
	conf := middleware.NewAuthOptions(middleware.WithValidator(validator))
	require.NotZero(t, conf, "a zero valued configuration was returned")

	actual, err := conf.Validator()
	require.NoError(t, err, "could not create default validator")
	require.Same(t, validator, actual, "conf did not reutn the same validator")
}
