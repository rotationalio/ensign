package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	conf := config.RatelimitConfig{PerSecond: 1, Burst: 0, TTL: 1}
	router.GET("/", middleware.RateLimiter(conf), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv := httptest.NewServer(router)
	defer srv.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
	require.Equal(t, "0.00", w.Header().Get("X-RateLimit-Remaining"))

	// ////////////////////// Test 2 /////////////////////////////
	// Test that setting the Limit to 1 and Burst to 3 will result in a 200 code
	router = gin.New()

	// token bucket is full, so the first request will be allowed
	conf = config.RatelimitConfig{PerSecond: 1, Burst: 3, TTL: 1}
	router.GET("/", middleware.RateLimiter(conf), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv = httptest.NewServer(router)
	defer srv.Close()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// ////////////////////// Test 3 /////////////////////////////
	// Test submission of multiple requests over the Burst amount results in a 429 error code
	router = gin.New()

	conf = config.RatelimitConfig{PerSecond: 1, Burst: 3, TTL: 1}
	router.GET("/", middleware.RateLimiter(conf), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv = httptest.NewServer(router)
	defer srv.Close()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/", nil)

	ticker := time.NewTicker(1 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				router.ServeHTTP(w, req)
				//fmt.Println(w.Code)
				//fmt.Println(w.Header().Get("X-RateLimit-Remaining"))
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	ticker.Stop()
	done <- true

	require.Equal(t, http.StatusTooManyRequests, w.Code)
	require.Equal(t, "0.00", w.Header().Get("X-RateLimit-Remaining"))
}
