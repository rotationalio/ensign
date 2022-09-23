package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {

		fixture := &api.StatusReply{
			Status:  "fine",
			Uptime:  (2 * time.Second).String(),
			Version: "1.0.test",
		}

		// Create a Test Server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/status", r.URL.Path)

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fixture)
		}))
		defer ts.Close()

		// Create a client to execute tests against the test server
		client, err := api.New(ts.URL)
		require.NoError(t, err)

		out, err := client.Status(context.Background())
		require.NoError(t, err, "could not execute status request")
		require.Equal(t, fixture, out, "expected the fixture to be returned")
	})

	t.Run("Unavailable", func(t *testing.T) {
		fixture := &api.StatusReply{
			Status:  "ack!",
			Uptime:  (9 * time.Second).String(),
			Version: "1.0.panic",
		}

		// Create a Test Server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/status", r.URL.Path)

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(fixture)
		}))
		defer ts.Close()

		// Create a client to execute tests against the test server
		client, err := api.New(ts.URL)
		require.NoError(t, err)

		out, err := client.Status(context.Background())
		require.NoError(t, err, "could not execute status request")
		require.Equal(t, fixture, out, "expected the fixture to be returned")
	})
}
