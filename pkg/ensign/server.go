/*
Package ensign implements the Ensign single node server.
*/
package ensign

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-multierror"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/interceptors"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	// Initializes zerolog with our default logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg
	zerolog.DurationFieldInteger = false
	zerolog.DurationFieldUnit = time.Millisecond

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()

	// Disable gRPC logging to reduce logging verbosity
	logger.DisableGRPCLog()
}

// An Ensign server implements the Ensign service as defined by the wire protocol.
type Server struct {
	api.UnimplementedEnsignServer
	srv     *grpc.Server  // The gRPC server that handles incoming requests in individual go routines
	conf    config.Config // Primary source of truth for server configuration
	started time.Time     // The timestamp that the server was started (for uptime)
	pubsub  *PubSub       // An in-memory channel based buffer between publishers and subscribers
	echan   chan error    // Sending errors down this channel stops the server (is fatal)
}

// New creates a new ensign server with the given configuration. Most server setup is
// conducted in this method including setting up logging, connecting to databases, etc.
// If this method succeeds without an error, the server is ready to be served, but it
// will not listen to or handle requests until the Serve() method is called.
func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment if an empty config is supplied
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Configure logging (will modify logging globally for all packages!)
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Configure sentry for error and performance monitoring
	if conf.Sentry.UseSentry() {
		if err = sentry.Init(conf.Sentry); err != nil {
			return nil, err
		}
	}

	s = &Server{
		conf:   conf,
		echan:  make(chan error, 1),
		pubsub: NewPubSub(),
	}

	// Prepare to receive gRPC requests and configure RPCs
	opts := make([]grpc.ServerOption, 0, 4)
	opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	opts = append(opts, grpc.ChainUnaryInterceptor(s.UnaryInterceptors()...))
	opts = append(opts, grpc.ChainStreamInterceptor(s.StreamInterceptors()...))
	s.srv = grpc.NewServer(opts...)

	// TODO: perform setup tasks if we're not in maintenance mode.

	// Initialize the Ensign service
	api.RegisterEnsignServer(s.srv, s)
	return s, nil
}

// Serve RPC requests on the bind address specified in the configuration.
func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Run monitoring and metrics server
	if err = o11y.Serve(s.conf.Monitoring); err != nil {
		log.Error().Err(err).Msg("could not start monitoring server")
		return err
	}

	// Preregister gRPC metrics for prometheus if metrics are enabled to ensure that
	// Grafana dashboards are fully populated without waiting for requests.
	if s.conf.Monitoring.Enabled {
		o11y.PreRegisterGRPCMetrics(s.srv)
	}

	// Listen for TCP requests (other sockets such as bufconn for tests should use Run()).
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		log.Error().Err(err).Str("bindaddr", s.conf.BindAddr).Msg("could not listen on given bindaddr")
		return err
	}

	go s.Run(sock)
	log.Info().Str("listen", s.conf.BindAddr).Msg("ensign server started")

	// Now that the server is running set the start time to track uptime
	s.started = time.Now()

	// Listen for any fatal errors on the error channel, blocking while the server go
	// routine does its work. If the error is nil we expect a graceful shutdown.
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

// Run the gRPC server on the specified socket. This method can be used to serve TCP
// requests or to connect to a bufconn for testing purposes. This method blocks while
// the server is running so it should be run in a go routine.
func (s *Server) Run(sock net.Listener) {
	defer sock.Close()
	if err := s.srv.Serve(sock); err != nil {
		s.echan <- err
	}
}

// Shutdown stops the ensign server and all long running processes gracefully. May
// return a multierror if there were multiple problems during shutdown but it will
// attempt to close all open services and processes.
func (s *Server) Shutdown() (err error) {
	errs := make([]error, 0)
	log.Info().Msg("gracefully shutting down ensign server")
	s.srv.GracefulStop()

	if err = o11y.Shutdown(context.Background()); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		log.Debug().Int("n_errs", len(errs)).Msg("could not successfully shutdown ensign server")
		return multierror.Append(err, errs...)
	}

	log.Debug().Msg("successfully shutdown ensign server")
	return nil
}

// Prepares the interceptors (middleware) for the unary RPC endpoings of the server.
// The first interceptor will be the outer most, while the last interceptor will be the
// inner most wrapper around the real call. All unary interceptors returned by this
// method should be chained using grpc.ChainUnaryInterceptor().
func (s *Server) UnaryInterceptors() []grpc.UnaryServerInterceptor {
	// NOTE: if more interceptors are added, make sure to increase the capacity!
	opts := make([]grpc.UnaryServerInterceptor, 0, 2)

	// If we're in maintenance mode only return the maintenance mode interceptor and
	// the panic recovery interceptor (just in case). Otherwise continue to build chain.
	if maintenace := interceptors.UnaryMaintenance(s.conf); maintenace != nil {
		opts = append(opts, maintenace)
		opts = append(opts, interceptors.UnaryRecovery(s.conf.Sentry))
		return opts
	}

	opts = append(opts, interceptors.UnaryMonitoring(s.conf))
	opts = append(opts, interceptors.UnaryRecovery(s.conf.Sentry))
	return opts
}

// Prepares the interceptors (middleware) for the unary RPC endpoints of the server.
// The first interceptor will be the outer most, while the last interceptor will be the
// inner most wrapper around the real call. All stream interceptors returned by this
// method should be chained using grpc.ChainStreamInterceptor().
func (s *Server) StreamInterceptors() []grpc.StreamServerInterceptor {
	// NOTE: if more interceptors are added, make sure to increase the capacity!
	opts := make([]grpc.StreamServerInterceptor, 0, 2)

	// If we're in maintenance mode only return the maintenance mode interceptor.
	if mainenance := interceptors.StreamMaintenance(s.conf); mainenance != nil {
		opts = append(opts, mainenance)
		return opts
	}

	opts = append(opts, interceptors.StreamMonitoring(s.conf))
	opts = append(opts, interceptors.StreamRecovery(s.conf.Sentry))
	return opts
}
