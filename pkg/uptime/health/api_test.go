package health_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/stretchr/testify/require"
)

func TestAPIOnline(t *testing.T) {
	started := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rep := make(map[string]interface{})
		rep["status"] = "ok"
		rep["version"] = "1.0.1 (7b2ec52)"
		rep["uptime"] = time.Since(started)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rep)
	}))
	defer ts.Close()

	mon, err := health.NewAPIMonitor(ts.URL)
	require.NoError(t, err)

	status, err := mon.Status(context.Background())
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Online, status.Status())
	require.False(t, status.CheckedAt().IsZero())
	require.NotEmpty(t, status.Hash())
}

func TestAPIStopping(t *testing.T) {
	started := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rep := make(map[string]interface{})
		rep["status"] = "stopping"
		rep["version"] = "1.0.1 (7b2ec52)"
		rep["uptime"] = time.Since(started)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(rep)
	}))
	defer ts.Close()

	mon, err := health.NewAPIMonitor(ts.URL)
	require.NoError(t, err)

	status, err := mon.Status(context.Background())
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Stopping, status.Status())
	require.False(t, status.CheckedAt().IsZero())
	require.NotEmpty(t, status.Hash())
}

func TestAPIBodyParseError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`foo`))
	}))
	defer ts.Close()

	mon, err := health.NewAPIMonitor(ts.URL)
	require.NoError(t, err)

	status, err := mon.Status(context.Background())
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Online, status.Status())
	require.False(t, status.CheckedAt().IsZero())
	require.NotEmpty(t, status.Hash())
}

func TestAPINoResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("stop test server with EOF")
	}))
	defer ts.Close()

	mon, err := health.NewAPIMonitor(ts.URL)
	require.NoError(t, err)

	status, err := mon.Status(context.Background())
	require.NoError(t, err, "should have been able to execute the status check")
	require.Equal(t, health.Offline, status.Status())
	require.False(t, status.CheckedAt().IsZero())
	require.NotEmpty(t, status.Hash())
}

func TestAPIServiceStatus(t *testing.T) {
	status := &health.APIServiceStatus{}
	require.Implements(t, serviceStatus, status)

	status.StatusCode = http.StatusOK
	status.Content = map[string]interface{}{"status": "ok", "version": "1.0.0 (abcd123)"}
	status.Timestamp = time.Now().Truncate(1 * time.Millisecond)

	require.Equal(t, health.Online, status.Status())

	// Should be able to marshal and unmarshal
	data, err := status.Marshal()
	require.NoError(t, err, "could not marshal status")

	other := &health.APIServiceStatus{}
	err = other.Unmarshal(data)
	require.NoError(t, err, "could not unmarshal status")
	require.Equal(t, status, other, "expected the unmarshaled type to be equal")

	// Hashing should be static
	require.NotEmpty(t, status.Hash())
	require.True(t, bytes.Equal(status.Hash(), other.Hash()))

	// Should not be equal to http status
	httpStatus := &health.HTTPServiceStatus{StatusCode: http.StatusOK}
	require.False(t, bytes.Equal(status.Hash(), httpStatus.Hash()))

	require.Equal(t, "1.0.0 (abcd123)", status.Version())

	// Test unparsable status
	status.Content["status"] = "foo"
	_, err = status.ParseStatus()
	require.Error(t, err)

	status.Content["status"] = 1
	_, err = status.ParseStatus()
	require.ErrorIs(t, err, health.ErrUnparsableStatus)

	delete(status.Content, "status")
	_, err = status.ParseStatus()
	require.ErrorIs(t, err, health.ErrNoStatusResponse)

	// Test unparsable version
	status.Content["version"] = 1
	require.Empty(t, status.Version())

	delete(status.Content, "version")
	require.Empty(t, status.Version())
}

func TestAPIServiceStatusHash(t *testing.T) {
	testCases := []struct {
		status   *health.APIServiceStatus
		expected string
	}{
		{
			status:   &health.APIServiceStatus{},
			expected: "0a5ba753163c64bf6dc6a33db63bcb75",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
					Endpoint: "https://example.com/v1/status",
					Error:    "connection refused",
				},
			},
			expected: "e7544c04eb142f0b9337309b0581222e",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:  "https://example.com/v1/status",
					Error:     "context deadline exceeded",
					ErrorType: health.Degraded,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
			},
			expected: "28e77dec10090dd9d57c4a9054f6ba5e",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:   "https://example.com/v1/status",
					StatusCode: http.StatusOK,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
			},
			expected: "c2f93ca1d61e0f73ad00e2b246bd632a",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:   "https://example.com/v1/status",
					StatusCode: http.StatusNoContent,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
			},
			expected: "96ac3543ae3e2ebe0aa1cd83cb59476e",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:   "https://example.com/v1/status",
					StatusCode: http.StatusServiceUnavailable,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
				Content: map[string]interface{}{
					"status":  "stopping",
					"version": "123abcd",
					"uptime":  "2h42m14s",
				},
			},
			expected: "732af3bd5c65679db711973de5e874c1",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:   "https://example.com/v1/status",
					StatusCode: http.StatusOK,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
				Content: map[string]interface{}{
					"status":  "ok",
					"version": "1.0.0 (abcd1)",
					"uptime":  "12m",
				},
			},
			expected: "930691c810191c1d6cf8df5bd9d4f7c2",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:   "https://example.com/v1/status",
					StatusCode: http.StatusOK,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
				Content: map[string]interface{}{
					"status":  "ok",
					"version": "1.1.0 (21defa)",
					"uptime":  "18m",
				},
			},
			expected: "1e82f681fd2f2b0eb3359f47d8629c29",
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Endpoint:   "https://example.com/status",
					StatusCode: http.StatusOK,
					BaseStatus: health.BaseStatus{
						Timestamp: time.Now(),
					},
				},
				Content: map[string]interface{}{
					"status":  "ok",
					"version": "1.1.0 (21defa)",
					"uptime":  "18m",
				},
			},
			expected: "c0c5c6d016e9021857b5735351de2a41",
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

func TestAPIServiceStatusStatus(t *testing.T) {
	testCases := []struct {
		status   *health.APIServiceStatus
		expected health.Status
	}{
		{
			status:   &health.APIServiceStatus{},
			expected: health.Unknown,
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Error: "connection refused",
				},
			},
			expected: health.Offline,
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					Error:     "context deadline exceeded",
					ErrorType: health.Degraded,
				},
			},
			expected: health.Degraded,
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					StatusCode: http.StatusOK,
				},
			},
			expected: health.Online,
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					StatusCode: http.StatusServiceUnavailable,
				},
				Content: map[string]interface{}{
					"status":  "stopping",
					"version": "123abcd",
					"uptime":  "2h42m14s",
				},
			},
			expected: health.Stopping,
		},
		{
			status: &health.APIServiceStatus{
				HTTPServiceStatus: health.HTTPServiceStatus{
					StatusCode: http.StatusOK,
				},
				Content: map[string]interface{}{
					"status":  "maintenance",
					"version": "123abcd",
					"uptime":  "2h42m14s",
				},
			},
			expected: health.Maintenance,
		},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, tc.status.Status(), "test case %d failed", i)
	}
}
