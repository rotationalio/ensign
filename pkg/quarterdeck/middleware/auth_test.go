package middleware_test

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
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

func TestAuthorize(t *testing.T) {
	// create the test authentication server
	srv, err := authtest.NewServer()
	require.NoError(t, err, "could not start authtest server")
	defer srv.Close()

	// create the authenticate middleware
	authenticate, err := middleware.Authenticate(
		middleware.WithJWKSEndpoint(srv.KeysURL()),
		middleware.WithAudience(authtest.Audience),
		middleware.WithIssuer(authtest.Issuer),
	)
	require.NoError(t, err, "could not create authenticate middleware")

	// Create Authorize middleware
	authorize := middleware.Authorize("foo:read", "foo:write")

	// Create a quick router for testing the middleware
	router := gin.Default()
	router.GET("/noauth", authorize, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	router.GET("/auth", authenticate, authorize, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// If there is no authenticate middleware before authorize, expect a 401 response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/noauth", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// expect 401 with an authenticated request without the permissions
	tks, err := srv.CreateAccessToken(&tokens.Claims{Permissions: []string{"bar:read", "bar:write"}})
	require.NoError(t, err, "could not create access token")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/auth", nil)
	req.Header.Set("Authorization", "Bearer "+tks)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	// expect 200 with an authenticated request with the permissions
	tks, err = srv.CreateAccessToken(&tokens.Claims{Permissions: []string{"foo:read", "foo:write"}})
	require.NoError(t, err, "could not create access token")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/auth", nil)
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
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiXSwiZXhwIjoxNjgwNjE1MzMwLCJuYmYiOjE2ODA2MTE3MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AiLCJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6Impkb2VAZXhhbXBsZS5jb20iLCJvcmciOiIxMjMiLCJwcm9qZWN0IjoiYWJjIiwicGVybWlzc2lvbnMiOlsicmVhZDpkYXRhIiwid3JpdGU6ZGF0YSJdfQ.LLb6c2RdACJmoT3IFgJEwfu2_YJMcKgM2bF3ISF41A37gKTOkBaOe-UuTmjgZ7WEcuQ-cVkht0KI_4zqYYctB_WB9481XoNwff5VgFf3xrPdOYxS00YXQnl09RRqt6Fmca8nvd4mXfdO7uvpyNVuCIqNxBPXdSnRhreSoFB1GtFm42sBPAD7vF-MQUmU0c4PTsbiCfhR1_buH0NYEE1QFp3vYcgoiXOJHh9VStmRscqvLB12AQrcs26G9opdTCCORmvR2W3JLJ_hliHyp-d9lhXmCDFyiGkDEhTAUglqwBjqz5SO1UfAThWJO18PvZl4QPhb724oNT82VPh0DMDwfw"
	c.Request.Header.Set("Authorization", "Bearer "+token)
	tks, err := middleware.GetAccessToken(c)
	require.NoError(t, err)
	require.Equal(t, token, tks)

	// Test correct authorization header lowercase
	c.Request.Header.Set("authorization", "bearer "+token)
	tks, err = middleware.GetAccessToken(c)
	require.NoError(t, err)
	require.Equal(t, token, tks)
}

func TestGetAccessTokenCookie(t *testing.T) {
	// Test both the GetAccessToken Cookie and the SetAuthTokens handler
	router := gin.New()

	auth, err := authtest.NewServer()
	require.NoError(t, err, "could not start authentication server")
	defer auth.Close()

	// Add a route that sets the cookies
	router.GET("/login", func(c *gin.Context) {
		claims := &tokens.Claims{
			Name:  "Sylvia Brochen",
			Email: "sylvia@example.com",
		}

		access, refresh, err := auth.CreateTokenPair(claims)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		middleware.SetAuthCookies(c, access, refresh, "127.0.0.1")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Add a route that requires cookie-based authentication
	// create the authenticate middleware
	authenticate, err := middleware.Authenticate(
		middleware.WithJWKSEndpoint(auth.KeysURL()),
		middleware.WithAudience(authtest.Audience),
		middleware.WithIssuer(authtest.Issuer),
	)
	require.NoError(t, err, "could not create authenticate middleware")

	router.POST("/action", authenticate, func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"success": true})
	})

	// Create a tls test server with the cookie authenticated router
	srv := httptest.NewTLSServer(router)
	defer srv.Close()

	// Create an https client with a cookie jar
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := srv.Client()
	client.Jar = jar

	// Attempt to make a request that is not cookie authenticated
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	require.NoError(t, err)

	rep, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rep.StatusCode)

	// Login and set the cookies
	req, err = http.NewRequest(http.MethodGet, srv.URL+"/login", nil)
	require.NoError(t, err)
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)

	// Check that we got back two cookies in the response (access and refresh)
	cookies := rep.Cookies()
	require.Len(t, cookies, 2)

	// Should be able to do an authenticated request with the cookies
	req, err = http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	require.NoError(t, err)

	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rep.StatusCode)
}

func TestContextFromRequest(t *testing.T) {
	// Create a test context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Should error if there is no request on the context
	_, err := middleware.ContextFromRequest(c)
	require.ErrorIs(t, err, middleware.ErrNoRequest)

	// Test when no access token is set
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, err := middleware.ContextFromRequest(c)
	require.NoError(t, err, "could not get context from request")
	require.NotNil(t, ctx, "no context was returned")
	creds, _ := api.CredsFromContext(ctx)
	require.Nil(t, creds, "credentials should not have been set")

	// Test when an access token is set
	c.Set(middleware.ContextAccessToken, "JWT")
	ctx, err = middleware.ContextFromRequest(c)
	require.NoError(t, err, "could not get context from request")
	require.NotNil(t, ctx, "no context was returned")
	creds, _ = api.CredsFromContext(ctx)
	token, err := creds.AccessToken()
	require.NoError(t, err, "could not get access token from context")
	require.Equal(t, "JWT", token, "access token was not set correctly")
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

func TestAuthOptionsOverride(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	opts := middleware.AuthOptions{
		KeysURL:            "http://localhost:8088/.well-known/jwks.json",
		Audience:           "http://localhost:3000",
		Issuer:             "http://localhost:8088",
		MinRefreshInterval: 42 * time.Minute,
		Context:            ctx,
	}

	conf := middleware.NewAuthOptions(middleware.WithAuthOptions(opts))
	require.NotSame(t, opts, conf, "expected a new configuration object to be created")
	require.Equal(t, opts, conf, "expected the opts to override the configuration defaults")

	// Ensure the context is the same on the configuration
	cancel()
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
