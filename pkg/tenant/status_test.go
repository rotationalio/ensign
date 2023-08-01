package tenant_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/stretchr/testify/require"
)

func (suite *tenantTestSuite) TestStatus() {
	require := suite.Require()
	rep, err := suite.client.Status(context.Background())
	require.NoError(err, "could not execute client request")
	require.NotEmpty(rep, "expected a complete status reply")
	require.Equal("ok", rep.Status, "expected status to be ok")
	require.NotEmpty(rep.Uptime, "expected some value for uptime")
	require.NotEmpty(rep.Version, "expected some value for version")

	// Ensures that when the server is stopping we get back a stopping status
	suite.srv.SetHealth(false)
	defer suite.srv.SetHealth(true)

	rep, err = suite.client.Status(context.Background())
	require.NoError(err, "could not execute client request")
	require.Equal("stopping", rep.Status, "expected status to be ok")
	require.NotEmpty(rep.Uptime, "expected some value for uptime")
	require.NotEmpty(rep.Version, "expected some value for version")
}

func TestAvailableMaintenance(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	var err error
	logger.Discard()

	// Create a tenant server in maintenance mode and test the Available middleware
	// NOTE: this must be separate from the tenant test suite to run in maintenance mode
	// NOTE: specify a KeysURL that doesn't exist to ensure that maintenance mode does
	// not require a connection to a Quarterdeck server.
	conf, err := config.Config{
		Maintenance: true,
		BindAddr:    "127.0.0.1:0",
		Mode:        gin.TestMode,
		Auth: config.AuthConfig{
			KeysURL: "http://127.0.0.1:5000", // This server should not be running
		},
		Database: config.DatabaseConfig{
			Testing: true,
		},
	}.Mark()
	require.NoError(t, err, "could not create valid configuration for maintenance mode")

	srv, err := tenant.New(conf)
	require.NoError(t, err, "could not create tenant server in maintenance mode")

	stopped := make(chan bool)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
		<-stopped
		logger.ResetLogger()
	})
	go func() {
		if err := srv.Serve(); err != nil {
			t.Logf("could not serve tenant service: %s", err)
		}
		stopped <- true
	}()

	// Wait for 500ms to ensure the API server starts up
	time.Sleep(500 * time.Millisecond)

	// Expect that we get a 503 in maintenance mode for any RPC query
	paths := []string{
		"/", "/v1", "/v1/status", "/v1/tenants", "/v1/register", "/v1/login",
		"/v1/members", "/v1/members/foo", "/v1/projects", "/v1/projects/foo",
	}

	for _, path := range paths {
		rep, err := http.Get(srv.URL() + path)
		require.NoError(t, err, "could not execute http request with default client")
		require.Equal(t, http.StatusServiceUnavailable, rep.StatusCode, "expected status unavailable from maintenance mode server")

		// We expect a JSON response from the server
		status := &api.StatusReply{}
		err = json.NewDecoder(rep.Body).Decode(status)
		require.NoError(t, err, "could not decode JSON body from response")

		require.Equal(t, "maintenance", status.Status)
		require.NotEmpty(t, status.Uptime)
		require.NotEmpty(t, status.Version)
	}
}
