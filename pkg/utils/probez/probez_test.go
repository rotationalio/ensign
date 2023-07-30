package probez_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/probez"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logger.Discard()
	exitVal := m.Run()
	logger.ResetLogger()
	os.Exit(exitVal)
}

func TestHandlerState(t *testing.T) {
	srv := probez.New()
	require.True(t, srv.IsHealthy(), "handler should be healthy when instantiated")
	require.False(t, srv.IsReady(), "handler should not be ready when instantiated")

	srv.Healthy()
	srv.Ready()
	require.True(t, srv.IsHealthy())
	require.True(t, srv.IsReady())

	srv.NotHealthy()
	srv.NotReady()
	require.False(t, srv.IsHealthy())
	require.False(t, srv.IsReady())
}

func TestGinRouter(t *testing.T) {
	// Should be able to use a gin router to serve healthy and ready requests.
	router := gin.Default()
	srv := probez.New()
	srv.Use(router)

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

func TestServeMux(t *testing.T) {
	// Should be able to use a serve mux to serve healthy and ready requests.
	mux := http.NewServeMux()
	srv := probez.New()
	srv.Mux(mux)

	ts := httptest.NewServer(mux)
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

func TestHandler(t *testing.T) {
	// Should be able to use a Handler directly to serve healthy and ready requests.
	srv := probez.New()
	ts := httptest.NewServer(srv)
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

func TestHandlerConcurrency(t *testing.T) {
	// Multiple threads should not race to set values on the handler
	var wg sync.WaitGroup
	srv := probez.New()

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
			for i := 0; i < 10; i++ {
				time.Sleep(2 * time.Millisecond)
				srv.IsHealthy()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				srv.IsReady()
			}
		}()
	}

	wg.Wait()
}

func TestServeHTTPConcurrency(t *testing.T) {
	// Multiple threads should not race when using ServeHTTP
	// Multiple threads should not race to set values on the server.
	var wg sync.WaitGroup
	srv := probez.New()
	ts := httptest.NewServer(srv)
	defer ts.Close()

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
			probe, err := probez.NewProbe(ts.URL)
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
			probe, err := probez.NewProbe(ts.URL)
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
