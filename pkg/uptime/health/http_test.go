package health_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/stretchr/testify/require"
)

var serviceStatus *health.ServiceStatus

func TestOnline(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	mon, err := health.NewHTTPMonitor(ts.URL)
	require.NoError(t, err)

	status, err := mon.Status(context.Background())
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Online, status.Status())
	require.False(t, status.CheckedAt().IsZero())
	require.NotEmpty(t, status.Hash())
}

func TestTooManyRedirects(t *testing.T) {
	// Create a constantly redirecting server
	redirects := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirects++
		loc := r.URL.ResolveReference(&url.URL{Path: fmt.Sprintf("/%d", redirects)})

		w.Header().Set("Location", loc.String())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}))
	defer ts.Close()

	mon, err := health.NewHTTPMonitor(ts.URL)
	require.NoError(t, err)

	status, err := mon.Status(context.Background())
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Unhealthy, status.Status())
}

func TestTimeout(t *testing.T) {
	// Create a server that does not return a response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	mon, err := health.NewHTTPMonitor(ts.URL)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	status, err := mon.Status(ctx)
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Degraded, status.Status())
}

func TestOffline(t *testing.T) {
	mon, err := health.NewHTTPMonitor("http://localhost:40784")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	status, err := mon.Status(ctx)
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Offline, status.Status())
}

func TestHTTPServiceStatus(t *testing.T) {
	status := &health.HTTPServiceStatus{}
	require.Implements(t, serviceStatus, status)

	status.StatusCode = http.StatusOK
	status.Timestamp = time.Now().Truncate(1 * time.Millisecond)

	require.Equal(t, health.Online, status.Status())

	// Should be able to marshal and unmarshal
	data, err := status.Marshal()
	require.NoError(t, err, "could not marshal status")

	other := &health.HTTPServiceStatus{}
	err = other.Unmarshal(data)
	require.NoError(t, err, "could not unmarshal status")
	require.Equal(t, status, other, "expected the unmarshaled type to be equal")

	// Hashing should be static
	require.NotEmpty(t, status.Hash())
	require.True(t, bytes.Equal(status.Hash(), other.Hash()))
}

func TestHTTPServiceStatusHash(t *testing.T) {
	testCases := []struct {
		status   *health.HTTPServiceStatus
		expected string
	}{
		{
			status:   &health.HTTPServiceStatus{},
			expected: "66ad387c16757277b806e89d2e80f03d",
		},
		{
			status: &health.HTTPServiceStatus{
				Error: "connection refused",
			},
			expected: "9bb093a320101d0bb8b62e26678f89b9",
		},
		{
			status: &health.HTTPServiceStatus{
				Error:     "context deadline exceeded",
				ErrorType: health.Degraded,
			},
			expected: "b7da6725854e78afdb9ebe93da4be389",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusOK,
			},
			expected: "66ae36d87e757277b806e89d96d4d6e5",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusNoContent,
			},
			expected: "66ae24ad52757277b806e89d8f6121d9",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusTooManyRequests,
			},
			expected: "66aeb17e5d757277b806e89dc923608f",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusRequestTimeout,
			},
			expected: "66af10dc18757277b806e89df03fcfdc",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusNotFound,
			},
			expected: "66af230744757277b806e89df7b384e8",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusServiceUnavailable,
			},
			expected: "66ad615faf757277b806e89d3f46cb31",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusBadGateway,
			},
			expected: "66ad65e58e757277b806e89d4120b142",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusGatewayTimeout,
			},
			expected: "66ad5ccff8757277b806e89d3d66d6bc",
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusInternalServerError,
			},
			expected: "66ad6efb24757277b806e89d44da8bc8",
		},
	}

	for i, tc := range testCases {
		tc.status.Timestamp = time.Now()
		sum := tc.status.Hash()
		expected, _ := hex.DecodeString(tc.expected)
		require.True(t, bytes.Equal(expected, sum), "test case %d failed, sum was %x", i, sum)

		// Check equality -- assumes that all the test cases have different hashes.
		for j, otc := range testCases {
			var assert require.BoolAssertionFunc
			assert = require.False
			if i == j {
				assert = require.True
			}
			assert(t, health.Equal(tc.status, otc.status))
		}
	}
}

func TestHTTPServiceStatusStatus(t *testing.T) {
	testCases := []struct {
		status   *health.HTTPServiceStatus
		expected health.Status
	}{
		{
			status:   &health.HTTPServiceStatus{},
			expected: health.Unknown,
		},
		{
			status: &health.HTTPServiceStatus{
				Error: "connection refused",
			},
			expected: health.Offline,
		},
		{
			status: &health.HTTPServiceStatus{
				Error:     "context deadline exceeded",
				ErrorType: health.Degraded,
			},
			expected: health.Degraded,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusOK,
			},
			expected: health.Online,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusNoContent,
			},
			expected: health.Online,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusTooManyRequests,
			},
			expected: health.Degraded,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusRequestTimeout,
			},
			expected: health.Degraded,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusNotFound,
			},
			expected: health.Unhealthy,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusServiceUnavailable,
			},
			expected: health.Maintenance,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusBadGateway,
			},
			expected: health.Offline,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusGatewayTimeout,
			},
			expected: health.Offline,
		},
		{
			status: &health.HTTPServiceStatus{
				StatusCode: http.StatusInternalServerError,
			},
			expected: health.Unhealthy,
		},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, tc.status.Status(), "test case %d failed", i)
	}
}
