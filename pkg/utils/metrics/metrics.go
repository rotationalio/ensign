/*
metrics tracks request behavior inside tenant (which manages Ensign user accounts,
projects, and topics) and quarterdeck (which manages authentication and authorization)
*/
package metrics

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// Prometheus namespaces for the collectors defined in this package.
const (
	NamespaceHTTPMetrics = "http_stats"
	NamespaceUserMetrics = "user_metrics"
)

var (
	// Total HTTP requests, disaggregated by service (e.g. "tenant" or "quarterdeck", status code, query path)
	RequestsHandled *prometheus.CounterVec

	// HTTP request duration, disaggregated by service (e.g. "tenant" or "quarterdeck", status code, query path)
	RequestDuration *prometheus.HistogramVec

	// Daily active users, collected via Quarterdeck usage, disaggregated by type (e.g. human or machine)
	Active *prometheus.CounterVec

	// Failed logins, collected via Quarterdeck usage, disaggregated by user type (e.g. human or machine) and cause of failure
	FailedLogins *prometheus.CounterVec

	// Verified users, collected via Quarterdeck usage
	Verified *prometheus.CounterVec

	// Registered users, collected via Quarterdeck usage
	Registered *prometheus.CounterVec

	// Registered organizations, collected via Quarterdeck usage
	Organizations *prometheus.CounterVec

	// Projects, collected via Quarterdeck usage
	Projects *prometheus.CounterVec

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

func Routes(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// Initializes and registers the metric collectors in Prometheus.
// This function should only be called once (e.g. from the GinLogger).
// All new metrics must be defined in this function to be used.
func initCollectors() (err error) {

	// Track all collectors to register at the end of the function.
	// When adding new collectors make sure to increase the capacity.
	collectors := make([]prometheus.Collector, 0, 8)

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

	Active = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceUserMetrics,
		Name:      "active_users",
		Help:      "daily active users, collected via quarterdeck usage, by type (human v machine)",
	}, []string{"service", "type"})
	collectors = append(collectors, Active)

	FailedLogins = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceUserMetrics,
		Name:      "failed_logins",
		Help:      "failed logins, collected via quarterdeck usage, by user type and cause",
	}, []string{"service", "type", "cause"})
	collectors = append(collectors, FailedLogins)

	Verified = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceUserMetrics,
		Name:      "verified_users",
		Help:      "verified users, collected via quarterdeck usage",
	}, []string{"service"})
	collectors = append(collectors, Verified)

	Registered = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceUserMetrics,
		Name:      "registered_users",
		Help:      "registered users, collected via quarterdeck usage",
	}, []string{"service"})
	collectors = append(collectors, Registered)

	Organizations = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceUserMetrics,
		Name:      "organizations",
		Help:      "registered organizations, collected via quarterdeck usage",
	}, []string{"service"})
	collectors = append(collectors, Organizations)

	Projects = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceUserMetrics,
		Name:      "projects",
		Help:      "projects, collected via quarterdeck usage",
	}, []string{"service"})
	collectors = append(collectors, Projects)

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
