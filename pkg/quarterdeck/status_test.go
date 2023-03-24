package quarterdeck_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/stretchr/testify/require"
)

func (s *quarterdeckTestSuite) TestStatus() {
	require := s.Require()
	rep, err := s.client.Status(context.Background())
	require.NoError(err, "could not execute client request")
	require.NotEmpty(rep, "expected a complete status reply")
	require.Equal("ok", rep.Status, "expected status to be ok")
	require.NotEmpty(rep.Uptime, "expected some value for uptime")
	require.NotEmpty(rep.Version, "expected some value for version")

	// Ensure when server is stopping we get back a stopping status
	s.srv.SetHealth(false)
	defer s.srv.SetHealth(true)

	rep, err = s.client.Status(context.Background())
	require.NoError(err, "could not execute client request")
	require.Equal("stopping", rep.Status, "expected status to be ok")
	require.NotEmpty(rep.Uptime, "expected some value for uptime")
	require.NotEmpty(rep.Version, "expected some value for version")
}

func TestAvailableMaintenance(t *testing.T) {
	// Create a quarterdeck server in maintenance mode and test the Available middleware
	// NOTE: this must be separate from the quarterdeck test suite to run in maintenance mode
	stopped := make(chan bool)
	logger.Discard()
	t.Cleanup(func() {
		<-stopped
		logger.ResetLogger()
	})

	conf, err := config.Config{
		Maintenance:  true,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		AllowOrigins: []string{"http://localhost:3000"},
		RateLimit: config.RateLimitConfig{
			PerSecond: 20.00,
			Burst:     120,
			TTL:       5 * time.Minute,
		},
	}.Mark()
	require.NoError(t, err, "could not create valid configuration for maintenance mode")

	srv, err := quarterdeck.New(conf)
	require.NoError(t, err, "could not create quarterdeck server in maintenance mode")

	go func() {
		srv.Serve()
		stopped <- true
	}()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
	})

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
