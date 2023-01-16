package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/stretchr/testify/require"
)

func TestDoubleCookies(t *testing.T) {
	// Test both the DoubleCookie middleware and the SetDoubleCookieTokens handler
	router := gin.New()

	// Add a route that sets the cookies
	router.GET("/protect", func(c *gin.Context) {
		err := middleware.SetDoubleCookieToken(c, "", time.Now().Add(10*time.Minute))
		require.NoError(t, err)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Add a route that requires double cookie submit
	router.POST("/action", middleware.DoubleCookie(), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"success": true})
	})

	// Create a tls test server with the CSRF protected router
	srv := httptest.NewTLSServer(router)
	defer srv.Close()

	// Create an https client with a cookie jar
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := srv.Client()
	client.Jar = jar

	// Atttempt to make a request that is not CSRF protected
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	require.NoError(t, err)

	// Ensure the request is Forbidden
	rep, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, rep.StatusCode)

	// Check the data in the response
	data, err := readJSON(rep)
	require.NoError(t, err, "could not parse response")
	require.Contains(t, data, "error")
	require.Equal(t, middleware.ErrCSRFVerification.Error(), data["error"].(string))

	// Login and set the cookies
	req, err = http.NewRequest(http.MethodGet, srv.URL+"/protect", nil)
	require.NoError(t, err)

	// Execute the protect request
	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)

	// Check that we got back two cookies in the response
	cookies := rep.Cookies()
	require.Len(t, cookies, 2)

	testCases := make([]*http.Request, 0, 3)

	// Attempt to send a request with the cookies but no X-CSRF-TOKEN Header
	reqa, err := http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	require.NoError(t, err)
	testCases = append(testCases, reqa)

	// Send a request with the cookies but an empty X-CSRF-TOKEN Header
	reqb, err := http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	reqb.Header.Set(middleware.CSRFHeader, "")
	require.NoError(t, err)
	testCases = append(testCases, reqb)

	// Send a request with the cookies but an incorrect X-CSRF-TOKEN Header
	reqc, err := http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	reqc.Header.Set(middleware.CSRFHeader, "not a valid csrf token")
	require.NoError(t, err)
	testCases = append(testCases, reqc)

	for i, req := range testCases {
		// Ensure the request is Forbidden
		rep, err = client.Do(req)
		require.NoError(t, err, "bad request %d failed", i)
		require.Equal(t, http.StatusForbidden, rep.StatusCode, "bad request %d failed", i)

		// Check data in the response
		data, err = readJSON(rep)
		require.NoError(t, err, "bad request %d failed", i)
		require.Contains(t, data, "error", "bad request %d failed", i)
		require.Equal(t, middleware.ErrCSRFVerification.Error(), data["error"].(string), "bad request %d failed", i)

	}

	// Send a valid request with the double cookie protection intact
	var cookieToken string
	for _, cookie := range cookies {
		if cookie.Name == middleware.CSRFCookie {
			cookieToken = cookie.Value
		}
	}

	require.NotEmpty(t, cookieToken, "could not find cookie in response")
	req, err = http.NewRequest(http.MethodPost, srv.URL+"/action", nil)
	req.Header.Set(middleware.CSRFHeader, cookieToken)
	require.NoError(t, err)

	rep, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rep.StatusCode)

	data, err = readJSON(rep)
	require.NoError(t, err)
	require.NotContains(t, data, "error")
	require.Contains(t, data, "success")
	require.True(t, data["success"].(bool))

}

func readJSON(rep *http.Response) (gin.H, error) {
	defer rep.Body.Close()
	data := make(gin.H)
	if err := json.NewDecoder(rep.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func TestGenerateCSRFToken(t *testing.T) {
	// Generate 100 tokens and ensure that they do not equal the previous tokens
	tokens := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		token, err := middleware.GenerateCSRFToken()
		require.NoError(t, err, "could not generate CSRF token")
		tokens = append(tokens, token)
	}

	for i, token := range tokens {
		for j, other := range tokens {
			if i != j {
				require.NotEqual(t, token, other, "all tokens generated should be unique")
			}
		}
	}
}
