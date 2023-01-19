package mock_test

import (
	"net/http"
	"testing"

	"github.com/rotationalio/ensign/pkg/tenant/mock"
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

	// Try a valid request
	quarterdeck.OnStatus(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/status", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusOK, rep.StatusCode, "expected status code to be 200")

	// Try request with a different handler
	quarterdeck.OnStatus(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	req, err = http.NewRequest(http.MethodGet, quarterdeck.URL()+"/v1/status", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusServiceUnavailable, rep.StatusCode, "expected status code to be 503")

	// Endpoint with a path parameter
	quarterdeck.OnAPIKeys("somekey", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	// POST request should return 201
	req, err = http.NewRequest(http.MethodPost, quarterdeck.URL()+"/v1/apikeys/somekey", nil)
	require.NoError(t, err, "could not create request")
	rep, err = client.Do(req)
	require.NoError(t, err, "could not execute request")
	require.Equal(t, http.StatusCreated, rep.StatusCode, "expected status code to be 201")
}
