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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Prometheus namespaces for the collectors defined in this package.
const (
	NamespaceEnsign = "ensign"
	NamespaceGRPC   = "grpc"
)

var (
	// All Ensign specific collectors for observability are defined here.
	Events            *prometheus.CounterVec
	OnlinePublishers  prometheus.Gauge
	OnlineSubscribers prometheus.Gauge

	// Generic gRPC collectors for observability defined here.
	RPCStarted    *prometheus.CounterVec
	RPCHandled    *prometheus.CounterVec
	RPCDuration   *prometheus.HistogramVec
	StreamMsgSent *prometheus.CounterVec
	StreamMsgRecv *prometheus.CounterVec
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

		// If not enabled, simply return here
		if !cfg.Enabled {
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

// Initializes all gRPC metrics with their appropriate null value for a gRPC service.
// This is useful to ensure that all metrics exist when collecting and querying without
// having to wait for an RPC to be called. This method should not be called before the
// Serve method is called otherwise it will panic.
func PreRegisterGRPCMetrics(srv *grpc.Server) {
	var allCodes = []codes.Code{
		codes.OK, codes.Canceled, codes.Unknown, codes.InvalidArgument, codes.DeadlineExceeded,
		codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
		codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange,
		codes.Unimplemented, codes.Internal, codes.Unavailable, codes.DataLoss,
	}

	for service, info := range srv.GetServiceInfo() {
		for _, mInfo := range info.Methods {
			method := mInfo.Name
			stype := "unary"
			if mInfo.IsClientStream || mInfo.IsServerStream {
				stype = "stream"
			}

			// These are references (not increments) to create labels but not set values
			RPCStarted.GetMetricWithLabelValues(stype, service, method)
			RPCDuration.GetMetricWithLabelValues(stype, service, method)

			if stype == "stream" {
				StreamMsgSent.GetMetricWithLabelValues(stype, service, method)
				StreamMsgRecv.GetMetricWithLabelValues(stype, service, method)
			}

			for _, code := range allCodes {
				RPCHandled.GetMetricWithLabelValues(stype, service, method, code.String())
			}
		}
	}
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
	collectors := make([]prometheus.Collector, 0, 8)

	// Ensign Collectors
	Events = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceEnsign,
		Name:      "events",
		Help:      "count the number of events handled by the ensign system",
	}, []string{"node", "region"})
	collectors = append(collectors, Events)

	OnlinePublishers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NamespaceEnsign,
		Name:      "online_publishers",
		Help:      "the number of publisher streams currently connected to the node",
	})
	collectors = append(collectors, OnlinePublishers)

	OnlineSubscribers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NamespaceEnsign,
		Name:      "online_subscribers",
		Help:      "the number of subscribe streams currently connected to the node",
	})
	collectors = append(collectors, OnlineSubscribers)

	// Generic GRPC collectors
	RPCStarted = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceGRPC,
		Name:      "server_started_total",
		Help:      "count the total number of RPCs started on the server",
	}, []string{"type", "service", "method"})
	collectors = append(collectors, RPCStarted)

	RPCHandled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceGRPC,
		Name:      "server_handled_total",
		Help:      "count the total number of RPCs completed on the server regardless of success or failure",
	}, []string{"type", "service", "method", "code"})
	collectors = append(collectors, RPCHandled)

	RPCDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: NamespaceGRPC,
		Name:      "server_handler_duration",
		Help:      "response latency (in seconds) of the application handler for the rpc method",
	}, []string{"type", "service", "method"})
	collectors = append(collectors, RPCDuration)

	StreamMsgSent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceGRPC,
		Name:      "server_stream_messages_sent",
		Help:      "total number of streaming messages sent by the server",
	}, []string{"type", "service", "method"})
	collectors = append(collectors, StreamMsgSent)

	StreamMsgRecv = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NamespaceGRPC,
		Name:      "server_stream_messages_recv",
		Help:      "total number of streaming messages received by the server",
	}, []string{"type", "service", "method"})
	collectors = append(collectors, StreamMsgRecv)

	// Register all the collectors
	for _, collector := range collectors {
		if err = prometheus.Register(collector); err != nil {
			log.Debug().Err(err).Msg("could not register collector")
			return err
		}
	}

	return nil
}
