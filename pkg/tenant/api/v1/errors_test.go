package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/stretchr/testify/require"
)

type IntString int

func (t IntString) String() string {
	return fmt.Sprintf("%04x", int(t))
}

type JMap map[string]string

func (j JMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string(j))
}

func TestErrorResponse(t *testing.T) {
	testCases := []struct {
		err      interface{}
		expected string
	}{
		{fmt.Errorf("something %s happened", "bad"), "something bad happened"},
		{errors.New("godzilla is here!"), "godzilla is here!"},
		{"this is a simple string", "this is a simple string"},
		{IntString(42), "002a"},
		{JMap{"color": "red"}, "{\"color\":\"red\"}"},
		{42, "unhandled error response"},
	}

	for _, tc := range testCases {
		rep := api.ErrorResponse(tc.err)
		require.False(t, rep.Success, "expected error reply to be success false")
		require.Equal(t, tc.expected, rep.Error, "unexpected result")
	}
}

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	api.NotFound(ctx)

	result := r.Result()
	defer result.Body.Close()
	require.Equal(t, result.StatusCode, http.StatusNotFound)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)

}

func TestNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	api.NotAllowed(ctx)

	result := r.Result()
	defer result.Body.Close()
	require.Equal(t, result.StatusCode, http.StatusMethodNotAllowed)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)
}

func TestReplyQuarterdeckError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// If err is nil then 200 OK is returned
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	api.ReplyQuarterdeckError(ctx, nil)
	responseEquals(t, r.Result(), http.StatusOK, api.Reply{Success: true})

	// If err is not a StatusError then 500 is returned
	r = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(r)
	api.ReplyQuarterdeckError(ctx, errors.New("something bad happened"))
	responseEquals(t, r.Result(), http.StatusInternalServerError, api.ErrorResponse("something bad happened"))

	// Empty StatusError should return a 500
	r = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(r)
	serr := &qd.StatusError{}
	api.ReplyQuarterdeckError(ctx, serr)
	responseEquals(t, r.Result(), http.StatusInternalServerError, api.Reply{Success: false})

	// Test StatusError is parsed correctly
	serr = &qd.StatusError{
		StatusCode: http.StatusNotFound,
		Reply:      qd.ErrorResponse("resource not found"),
	}
	r = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(r)
	api.ReplyQuarterdeckError(ctx, serr)
	responseEquals(t, r.Result(), http.StatusNotFound, api.ErrorResponse("resource not found"))
}

func responseEquals(t *testing.T, result *http.Response, code int, reply api.Reply) {
	defer result.Body.Close()
	require.Equal(t, code, result.StatusCode)
	require.Equal(t, result.Header.Get("Content-Type"), "application/json; charset=utf-8")

	var data api.Reply
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)
	require.Equal(t, reply, data)
}
