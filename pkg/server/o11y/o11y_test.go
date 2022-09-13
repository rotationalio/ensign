package o11y_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rotationalio/ensign/pkg/config"
	"github.com/rotationalio/ensign/pkg/server/o11y"
	"github.com/stretchr/testify/require"
)

func TestO11yServer(t *testing.T) {
	// Ensure that we can setup the monitoring server and execute metrics.
	conf := config.MonitoringConfig{
		Enabled:  true,
		BindAddr: "127.0.0.1:48489",
		NodeID:   "testing-42",
	}

	err := o11y.Serve(conf)
	require.NoError(t, err, "could not serve the o11y metrics server")

	// Collect some metrics
	o11y.Events.WithLabelValues(conf.NodeID, "test").Inc()
	o11y.OnlinePublishers.Add(1)
	o11y.OnlineSubscribers.Add(3)

	// Attempt to collect the metrics
	rep, err := http.Get("http://127.0.0.1:48489/metrics")
	require.NoError(t, err, "could not make http request to o11y server")
	require.Equal(t, http.StatusOK, rep.StatusCode)

	err = o11y.Shutdown(context.Background())
	require.NoError(t, err, "could not shutdown the o11y metrics server")
}
