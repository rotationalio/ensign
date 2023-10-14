/*
Package ensign implements the Ensign single node server.
*/
package ensign

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-multierror"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/info"
	"github.com/rotationalio/ensign/pkg/ensign/interceptors"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	quarterdeck "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	health "github.com/rotationalio/ensign/pkg/utils/probez/grpc/v1"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
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
	health.ProbeServer
	api.UnimplementedEnsignServer

	srv     *grpc.Server                // The gRPC server that handles incoming requests in individual go routines
	conf    config.Config               // Primary source of truth for server configuration
	auth    *interceptors.Authenticator // Fetches public keys from Quarterdeck to authenticate token requests
	broker  *broker.Broker              // Brokers all incoming events from publishers and queues them to subscribers
	infog   *info.TopicInfoGatherer     // Gathers topic information in a background go routine
	data    store.EventStore            // Storage for event data - writing to this store must happen as fast as possible
	meta    store.MetaStore             // Storage for metadata such as topics and placement
	tasks   *tasks.TaskManager          // Manager for performing background tasks
	started time.Time                   // The timestamp that the server was started (for uptime)
	echan   chan error                  // Sending errors down this channel stops the server (is fatal)
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
		conf:  conf,
		echan: make(chan error, 1),
	}

	// Perform setup tasks if we're not in maintenance mode.
	if !conf.Maintenance {
		// Open local data storage for Ensign
		if s.data, s.meta, err = store.Open(conf.Storage); err != nil {
			return nil, err
		}

		// Create the authenticator
		if s.auth, err = interceptors.NewAuthenticator(conf.Auth.AuthOptions()...); err != nil {
			return nil, err
		}

		// Create the broker with access to the data stores
		s.broker = broker.New(s.data)

		// Create the topic info gatherer
		s.infog = info.New(s.data, s.meta)

		// Create the background task manager
		s.tasks = tasks.New(4, 64, time.Second)
		log.Debug().Int("workers", 4).Int("queue_size", 64).Msg("task manager started")
	}

	// Prepare to receive gRPC requests and configure RPCs
	opts := make([]grpc.ServerOption, 0, 4)
	opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	opts = append(opts, grpc.ChainUnaryInterceptor(s.UnaryInterceptors()...))
	opts = append(opts, grpc.ChainStreamInterceptor(s.StreamInterceptors()...))
	s.srv = grpc.NewServer(opts...)

	// Initialize the Ensign service
	api.RegisterEnsignServer(s.srv, s)
	health.RegisterHealthServer(s.srv, s)
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

	// Set the server to a not serving state
	s.NotHealthy()

	// Setup non-maintenance mode actions
	if !s.conf.Maintenance {
		// Wait for Quarterdeck before being ready to serve requests
		if err = s.WaitForQuarterdeck(); err != nil {
			sentry.Error(nil).Err(err).Msg("could not connect to quarterdeck")
			return err
		}

		// Start the broker to handle publish and subscribe
		s.broker.Run(s.echan)

		// Start the info gathering routine
		s.infog.Run()
	}

	// Run monitoring and metrics server
	if err = o11y.Serve(s.conf.Monitoring); err != nil {
		sentry.Error(nil).Err(err).Msg("could not start monitoring server")
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
		sentry.Error(nil).Err(err).Str("bindaddr", s.conf.BindAddr).Msg("could not listen on given bindaddr")
		return err
	}

	go s.Run(sock)
	log.Info().Str("listen", s.conf.BindAddr).Msg("ensign server started")

	// Now that the server is running set the start time to track uptime
	s.started = time.Now()

	// Set the server to ready and serving requests
	s.Healthy()

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
	// Set the server to a not serving state
	s.NotHealthy()

	errs := make([]error, 0)
	log.Info().Msg("gracefully shutting down ensign server")
	s.srv.GracefulStop()

	// Shutdown running services if not in maintenance mode
	if !s.conf.Maintenance {
		// Shutdown the running broker and finalize all events
		if err = s.broker.Shutdown(); err != nil {
			errs = append(errs, err)
		}

		// Shutdown the topic info gatherer
		if err = s.infog.Shutdown(); err != nil {
			errs = append(errs, err)
		}

		// Shutdown the task manger
		s.tasks.Stop()

		// Gracefully close the data stores.
		if err = s.meta.Close(); err != nil {
			errs = append(errs, err)
		}

		if err = s.data.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if err = o11y.Shutdown(context.Background()); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		log.Debug().Int("n_errs", len(errs)).Msg("could not successfully shutdown ensign server")
		return multierror.Append(err, errs...)
	}

	// Flush sentry errors
	if s.conf.Sentry.UseSentry() {
		sentry.Flush(2 * time.Second)
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
	opts := make([]grpc.UnaryServerInterceptor, 0, 4)

	// If we're in maintenance mode only return the maintenance mode interceptor and
	// the panic recovery interceptor (just in case). Otherwise continue to build chain.
	if maintenance := interceptors.UnaryMaintenance(s.conf); maintenance != nil {
		opts = append(opts, maintenance)
		opts = append(opts, interceptors.UnaryRecovery(s.conf.Sentry))
		return opts
	}

	opts = append(opts, interceptors.UnaryMonitoring(s.conf))
	opts = append(opts, interceptors.UnaryRecovery(s.conf.Sentry))
	opts = append(opts, sentry.UnaryInterceptor(s.conf.Sentry))
	opts = append(opts, s.auth.Unary())
	return opts
}

// Prepares the interceptors (middleware) for the unary RPC endpoints of the server.
// The first interceptor will be the outer most, while the last interceptor will be the
// inner most wrapper around the real call. All stream interceptors returned by this
// method should be chained using grpc.ChainStreamInterceptor().
func (s *Server) StreamInterceptors() []grpc.StreamServerInterceptor {
	// NOTE: if more interceptors are added, make sure to increase the capacity!
	opts := make([]grpc.StreamServerInterceptor, 0, 4)

	// If we're in maintenance mode only return the maintenance mode interceptor.
	if mainenance := interceptors.StreamMaintenance(s.conf); mainenance != nil {
		opts = append(opts, mainenance)
		return opts
	}

	opts = append(opts, interceptors.StreamMonitoring(s.conf))
	opts = append(opts, interceptors.StreamRecovery(s.conf.Sentry))
	opts = append(opts, sentry.StreamInterceptor(s.conf.Sentry))
	opts = append(opts, s.auth.Stream())
	return opts
}

// Creates a quarterdeck client connected to the same host as the Auth KeysURL and waits
// until Quarterdeck returns a healthy response or the exponential timeout backoff limit
// is reached before returning.
func (s *Server) WaitForQuarterdeck() (err error) {
	// Parse the auth url
	var qdu *url.URL
	if qdu, err = url.Parse(s.conf.Auth.KeysURL); err != nil {
		return err
	}

	var client quarterdeck.QuarterdeckClient
	if client, err = quarterdeck.New(qdu.String()); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err = client.WaitForReady(ctx); err != nil {
		return err
	}
	return nil
}

// StoreMock returns the underlying store for testing purposes. Can only be accessed in
// storage testing mode otherwise the method panics.
func (s *Server) StoreMock() *mock.Store {
	store, ok := s.data.(*mock.Store)
	if !ok {
		log.Panic().
			Str("data_store_type", fmt.Sprintf("%T", s.data)).
			Str("meta_store_type", fmt.Sprintf("%T", s.meta)).
			Bool("storage_testing", s.conf.Storage.Testing).
			Msg("store mock can only be retrieved in testing mode")
	}
	return store
}

// RunBroker runs the internal broker for testing purposes.
func (s *Server) RunBroker() {
	s.broker.Run(s.echan)
}
