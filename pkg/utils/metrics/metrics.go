/*
metrics tracks request behavior inside tenant (which manages Ensign user accounts,
projects, and topics) and quarterdeck (which manages authentication and authorization)
*/
package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// Prometheus namespaces for the collectors defined in this package.
const (
	NamespaceHTTPMetrics = "http_stats"
)

var (
	// Total HTTP requests, disaggregated by service (e.g. "tenant" or "quarterdeck", status code, query path)
	RequestsHandled *prometheus.CounterVec

	// HTTP request duration, disaggregated by service (e.g. "tenant" or "quarterdeck", status code, query path)
	RequestDuration *prometheus.HistogramVec

	// Protection from unintentionally re-registering collectors
	setup sync.Once
	err   error
)

func Setup() {
	// Ensure that the initialization of the metrics and the server occurs only once.
	setup.Do(func() {
		// Register the collectors
		if err = initCollectors(); err != nil {
			return
		}
	})
}

// Initializes and registers the metric collectors in Prometheus.
// This function should only be called once (e.g. from the GinLogger).
// All new metrics must be defined in this function to be used.
func initCollectors() (err error) {

	// Track all collectors to register at the end of the function.
	// When adding new collectors make sure to increase the capacity.
	collectors := make([]prometheus.Collector, 0, 2)

	RequestsHandled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceHTTPMetrics,
		Name:      "requests_handled",
		Help:      "total requests handled, disaggregated by service, http status code, and path",
	}, []string{"service", "code", "path"})
	collectors = append(collectors, RequestsHandled)

	RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: NamespaceHTTPMetrics,
		Name:      "request_duration",
		Help:      "duration of requests, disaggregated by service, http status code, and path",
	}, []string{"service", "code", "path"})
	collectors = append(collectors, RequestDuration)

	// Register the collectors
	registerCollectors(collectors)
	return nil
}

// Helper function to ensure collectors are registered. Will emit a log error if the
// caller attempts to re-register an already-registered collector, but will not fail
func registerCollectors(collectors []prometheus.Collector) {
	var err error
	// Register the collectors
	for _, collector := range collectors {
		if err = prometheus.Register(collector); err != nil {
			err = fmt.Errorf("cannot register %s", collector)
			log.Warn().Err(err).Msg("collector already registered")
		}
	}
}

// TODO do we want to pre-register any HTTP status code labels?
// if so, do we need to know all the possible query paths to initialize those as well?
