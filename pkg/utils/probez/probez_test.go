package probez_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/utils/probez"
	"github.com/stretchr/testify/require"
)

func TestServerURL(t *testing.T) {
	srv := probez.New()
	err := srv.Serve(":0")
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	uri, err := url.Parse(srv.URL())
	require.NoError(t, err)
	require.Equal(t, "http", uri.Scheme)
	require.Equal(t, "127.0.0.1", uri.Hostname())
	require.NotEmpty(t, uri.Port())
}

func TestZeroServer(t *testing.T) {
	// A zero-valued server should still serve but return 503 by default.
	srv := &probez.Server{}
	err := srv.Serve(":0")
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	probe, err := probez.NewProbe(srv.URL())
	require.NoError(t, err)

	ok, status, err := probe.Healthy(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)

	ok, status, err = probe.Live(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)

	ok, status, err = probe.Ready(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)
}

func TestNewServer(t *testing.T) {
	// A new server should be healthy but not ready.
	srv := probez.New()
	err := srv.Serve(":0")
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	probe, err := probez.NewProbe(srv.URL())
	require.NoError(t, err)

	ok, status, err := probe.Healthy(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	ok, status, err = probe.Live(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	ok, status, err = probe.Ready(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)
}

func TestHTTPServer(t *testing.T) {
	// Should be able to change healthy and ready status and get different responses.
	// A new server should be healthy but not ready.
	srv := probez.New()
	err := srv.Serve(":0")
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	probe, err := probez.NewProbe(srv.URL())
	require.NoError(t, err)

	srv.Healthy()
	ok, status, err := probe.Healthy(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	ok, status, err = probe.Live(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	srv.NotHealthy()
	ok, status, err = probe.Healthy(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)

	ok, status, err = probe.Live(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)

	srv.Ready()
	ok, status, err = probe.Ready(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	srv.NotReady()
	ok, status, err = probe.Ready(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)
}

func TestGinServer(t *testing.T) {
	// Should be able to use a gin router to serve healthy and ready requests.
	router := gin.Default()
	srv := probez.New()
	srv.Handle(router)

	ts := httptest.NewServer(router)
	defer ts.Close()

	probe, err := probez.NewProbe(ts.URL)
	require.NoError(t, err)

	srv.Healthy()
	ok, status, err := probe.Healthy(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	ok, status, err = probe.Live(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	srv.NotHealthy()
	ok, status, err = probe.Healthy(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)

	ok, status, err = probe.Live(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)

	srv.Ready()
	ok, status, err = probe.Ready(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.True(t, ok)

	srv.NotReady()
	ok, status, err = probe.Ready(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, status)
	require.False(t, ok)
}

func TestConcurrency(t *testing.T) {
	// Multiple threads should not race to set values on the server.
	var wg sync.WaitGroup
	srv := probez.New()
	err := srv.Serve(":0")
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				srv.Healthy()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				srv.NotHealthy()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				srv.Ready()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				srv.NotReady()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			probe, err := probez.NewProbe(srv.URL())
			require.NoError(t, err)

			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				_, _, err = probe.Healthy(context.Background())
				require.NoError(t, err)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			probe, err := probez.NewProbe(srv.URL())
			require.NoError(t, err)

			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				_, _, err = probe.Ready(context.Background())
				require.NoError(t, err)
			}
		}()
	}

	wg.Wait()
}
