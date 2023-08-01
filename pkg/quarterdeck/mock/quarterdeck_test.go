package mock_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/stretchr/testify/require"
)

func TestMock(t *testing.T) {
	// Start the mock server
	quarterdeck, err := mock.NewServer()
	require.NoError(t, err, "could not start the mock quarterdeck server")
	t.Cleanup(func() {
		quarterdeck.Close()
	})
	client := quarterdeck.Client()

	// Unconfigured path should return 404
	req, err := http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/invalid", nil)
	require.NoError(t, err, "could not create request")
	rep, err := client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusNotFound, rep.StatusCode, "expected status code to be 404")

	// Try a request with a default handler
	quarterdeck.OnStatus()
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/status", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusOK, rep.StatusCode, "expected status code to be 200")

	// Try request with a different status code
	quarterdeck.OnStatus(mock.UseStatus(http.StatusServiceUnavailable))
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/status", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusServiceUnavailable, rep.StatusCode, "expected status code to be 503")

	// Try request with a fixture
	fixture := &api.StatusReply{
		Status: "ok",
	}
	quarterdeck.OnStatus(mock.UseJSONFixture(fixture))
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/status", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusOK, rep.StatusCode, "expected status code to be 200")
	actual := &api.StatusReply{}
	require.NoError(t, json.NewDecoder(rep.Body).Decode(actual), "could not decode response")
	require.Equal(t, fixture, actual, "expected response to match fixture")

	// Try request with a custom handler
	quarterdeck.OnStatus(mock.UseHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/status", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusServiceUnavailable, rep.StatusCode, "expected status code to be 503")

	// Endpoint with a path parameter
	quarterdeck.OnAPIKeysDetail("somekey")
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/apikeys/somekey", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusOK, rep.StatusCode, "expected status code to be 201")

	// Verify that the handlers were called the expected number of times
	require.Equal(t, 4, quarterdeck.StatusCount(), "expected status handler to be called 4 times")
	require.Equal(t, 1, quarterdeck.APIKeysDetailCount("somekey"), "expected apikeys handler to be called 1 time")
}
