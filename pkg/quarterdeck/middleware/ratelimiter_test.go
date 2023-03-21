package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/stretchr/testify/require"
)

func TestRatelimiter(t *testing.T) {
	// ////////////////////// Test 1 /////////////////////////////
	// Test that setting the Burst to 0 results in a 429 error code
	router := gin.New()

	// setting the Burst to 0 means the token bucket will be empty
	// therefore, all requests will be rejected
	conf := config.RateLimitConfig{PerSecond: 1, Burst: 0, TTL: 1}
	router.GET("/", middleware.RateLimiter(conf), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv := httptest.NewServer(router)
	defer srv.Close()
	client := srv.Client()
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	require.NoError(t, err)
	rep, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, rep.StatusCode)
	require.Equal(t, "0.00", rep.Header.Get("X-RateLimit-Remaining"))

	// ////////////////////// Test 2 /////////////////////////////
	// Test that setting the Limit to 1 and Burst to 3 will result in a 200 code
	router = gin.New()

	// token bucket is full, so the first request will be allowed
	conf = config.RateLimitConfig{PerSecond: 1, Burst: 3, TTL: 1}
	router.GET("/", middleware.RateLimiter(conf), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv = httptest.NewServer(router)
	defer srv.Close()
	client = srv.Client()
	req, err = http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	require.NoError(t, err)
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)
	require.Equal(t, "2.00", rep.Header.Get("X-RateLimit-Remaining"))

	// ////////////////////// Test 3 /////////////////////////////
	// Test submission of multiple requests over the Burst amount results in a 429 error code
	router = gin.New()

	conf = config.RateLimitConfig{PerSecond: 1, Burst: 3, TTL: 1}
	router.GET("/", middleware.RateLimiter(conf), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv = httptest.NewServer(router)
	defer srv.Close()
	client = srv.Client()
	req, err = http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	require.NoError(t, err)
	req.RemoteAddr = "1.2.3.5"

	for i := 0; i < 5; i++ {
		rep, err := client.Do(req)
		require.NoError(t, err)
		// the first two requests will be allowed and rate limit remaining will be greater than 0
		if i < 2 {
			require.Equal(t, http.StatusOK, rep.StatusCode)
			require.NotEqual(t, "0.00", rep.Header.Get("X-RateLimit-Remaining"))
			// the third request will be allowed but rate limit remaining will be equal to zero
		} else if i == 2 {
			require.Equal(t, http.StatusOK, rep.StatusCode)
			require.Equal(t, "0.00", rep.Header.Get("X-RateLimit-Remaining"))
			// beyond the third request all requests will be rejected
		} else {
			require.Equal(t, http.StatusTooManyRequests, rep.StatusCode)
			require.Equal(t, "0.00", rep.Header.Get("X-RateLimit-Remaining"))
		}

	}
}
