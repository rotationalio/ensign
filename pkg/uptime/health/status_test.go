package health_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	statuses := []health.Status{
		health.Unknown, health.Maintenance, health.Stopping, health.Online,
		health.Degraded, health.Unhealthy, health.Offline, health.Outage,
	}

	for i, status := range statuses {
		s := status.String()
		require.NotEmpty(t, s, "could not stringify status %d", i)

		parsed, err := health.ParseStatus(s)
		require.NoError(t, err, "could not parse status %d", i)
		require.Equal(t, status, parsed, "unexpected parse for status %d", i)
	}

	// Should be able to parse ok
	parsed, err := health.ParseStatus("ok")
	require.NoError(t, err, "could not parse ok")
	require.Equal(t, health.Online, parsed, "ok did not parse as online")

	// Attempt to parse an unknown status
	parsed, err = health.ParseStatus("foo")
	require.Error(t, err, "should not have been able to parse foo")
	require.Equal(t, health.Unknown, parsed, "should have returned unknown status on error")
}
