package tenant_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
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
	var err error
	logger.Discard()

	// Maintenance mode does not require authentication but we still create the
	// authentication middleware, so the authtest server needs to be runnning.
	auth, err := authtest.NewServer()
	require.NoError(t, err, "could not start the authtest server")
	t.Cleanup(func() {
		auth.Close()
	})

	// Create a tenant server in maintenance mode and test the Available middleware
	// NOTE: this must be separate from the tenant test suite to run in maintenance mode
	conf, err := config.Config{
		Maintenance:  true,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		AllowOrigins: []string{"http://localhost:3000"},
		Auth: config.AuthConfig{
			Audience:     authtest.Audience,
			Issuer:       authtest.Issuer,
			KeysURL:      auth.KeysURL(),
			CookieDomain: "localhost",
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
		srv.Shutdown()
		<-stopped
		logger.ResetLogger()
	})
	go func() {
		srv.Serve()
		stopped <- true
	}()

	// Wait for 500ms to ensure the API server starts up
	time.Sleep(500 * time.Millisecond)

	// Expect that we get a 503 in maintenance mode for any RPC query
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	require.NoError(t, err, "could not create http request")
	rep, err := http.DefaultClient.Do(req)
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
