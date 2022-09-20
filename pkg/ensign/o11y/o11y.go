/*
Package o11y (a numeronym for "observability") exports server-specific metric
collectors to Prometheus for monitoring of the service and Ensign nodes. The package
also manages the http server that presents the metric collectors to the Prometheus
scraper in an on demand fashion.

At least once before use, an external caller needs to call the package Serve() method
with the desired configuration -- this will ensure the collector server starts up and
that the metrics collected are available for use in external packages. To clean up the
server, external callers should call the Shutdown() method.
*/
package o11y

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rs/zerolog/log"
)

// Prometheus namespace for the collectors defined in this package.
const Namespace = "ensign"

// All Ensign specific collectors for observability are defined here.
var (
	Events            *prometheus.CounterVec
	OnlinePublishers  prometheus.Gauge
	OnlineSubscribers prometheus.Gauge
)

// Internal package variables for serving the collectors to the Prometheus scraper.
var (
	srv   *http.Server
	cfg   config.MonitoringConfig
	setup sync.Once
	mu    sync.Mutex
	err   error
)

// Serve the prometheus metric collectors server so that the Prometheus scraper can
// collect metrics from the node. This method must be called at least once before any of
// the metrics in this package can be used. Calling serve multiple times will not cause
// problems but only the first call to Serve (and therefore the first config) will be
// used. It is possible to call Serve again after calling Shutdown.
func Serve(conf config.MonitoringConfig) error {
	// Guard against concurrent Serve and Shutdown
	mu.Lock()
	defer mu.Unlock()

	// Ensure that the initialization of the metrics and the server occurs only once.
	setup.Do(func() {
		// Register the collectors
		cfg = conf
		if err = registerCollectors(); err != nil {
			return
		}

		// Setup the prometheus handler and collectors server.
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		srv = &http.Server{
			Addr:         cfg.BindAddr,
			Handler:      mux,
			ErrorLog:     nil,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		// Serve the metrics server in its own go routine
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error().Err(err).Msg("o11y server shutdown prematurely")
			}
		}()

		log.Info().Str("addr", fmt.Sprintf("http://%s/metrics", conf.BindAddr)).Msg("o11y server started and ready for prometheus collector")
	})

	return err
}

// Shutdown the prometheus metrics collectors server and reset the package. This method
// should be called at least once by outside callers before the process shuts down to
// ensure that system resources are cleaned up correctly.
func Shutdown(ctx context.Context) error {
	// Guard against concurrent Serve and Shutdown
	mu.Lock()
	defer mu.Unlock()

	// If we're already shutdown don't panic
	if srv == nil {
		return nil
	}

	// Ensure that no matter what happens we reset the package so it can be served again.
	defer func() {
		srv = nil
		cfg = config.MonitoringConfig{}
		err = nil
		setup = sync.Once{}
	}()

	// Ensure there is a shutdown deadline so we don't block forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Initializes and registers the metric collectors in Prometheus. This function should
// only be called once from the Serve function. All new metrics must be defined in this
// function so that they can be used.
func registerCollectors() (err error) {
	// Track all collectors to make it easier to register them at the end of this
	// function. When adding new collectors make sure to increase the capacity.
	collectors := make([]prometheus.Collector, 0, 3)

	Events = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "events",
		Help:      "count the number of events handled by the ensign system",
	}, []string{"node", "region"})
	collectors = append(collectors, Events)

	OnlinePublishers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "online_publishers",
		Help:      "the number of publisher streams currently connected to the node",
	})
	collectors = append(collectors, OnlinePublishers)

	OnlineSubscribers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "online_subscribers",
		Help:      "the number of subscribe streams currently connected to the node",
	})
	collectors = append(collectors, OnlineSubscribers)

	// Register all the collectors
	for _, collector := range collectors {
		if err = prometheus.Register(collector); err != nil {
			log.Debug().Err(err).Msg("could not register collector")
			return err
		}
	}

	return nil
}
