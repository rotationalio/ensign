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
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

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
				Endpoint: "https://example.com/v1/status",
				Error:    "connection refused",
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "f5a4a97a318f0068d37a32994b62d3c6",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:  "https://example.com/v1/status",
				Error:     "context deadline exceeded",
				ErrorType: health.Degraded,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "4566383c32f0116d094e02b9546ef736",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusOK,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "79506e16dae5a27089c37691448b0362",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/status",
				StatusCode: http.StatusOK,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "474f1b67f9d33c5f0b35eea57f16e8e8",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusNoContent,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "3f746a40d1c8730d80bb669254fb90e6",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusTooManyRequests,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "9f31276d7d0fd92208d1367b207eccc0",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusRequestTimeout,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "89aedc2035dac96705702834369a03b9",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusNotFound,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "0c023ccdeef77fce8845dc70b5cc8615",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusServiceUnavailable,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "6dcc66e7c2ee6b1e1ceaa3cd89bf53ee",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusBadGateway,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "3e5324ad3b2386660f05b8a2473b349b",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusGatewayTimeout,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "220adea78148fafd72c15f131a0fbc19",
		},
		{
			status: &health.HTTPServiceStatus{
				Endpoint:   "https://example.com/v1/status",
				StatusCode: http.StatusInternalServerError,
				BaseStatus: health.BaseStatus{
					Timestamp: time.Now(),
				},
			},
			expected: "e057a5a1f0e6a03208aabbdab6fb88f5",
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
