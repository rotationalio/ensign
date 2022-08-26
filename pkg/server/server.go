/*
Package server implements the Ensign single node server.
*/
package server

import (
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-multierror"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	// Initialize zerolog for server-side process logging
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// An Ensign server implements the Ensign service as defined by the wire protocol.
type Server struct {
	api.UnimplementedEnsignServer
	srv     *grpc.Server  // The gRPC server that handles incoming requests in individual go routines
	conf    config.Config // Primary source of truth for server configuration
	started time.Time     // The timestamp that the server was started (for uptime)
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

	s = &Server{
		conf:  conf,
		echan: make(chan error, 1),
	}

	// Prepare to receive gRPC requests and configure RPCs
	opts := make([]grpc.ServerOption, 0, 2)
	opts = append(opts, grpc.Creds(insecure.NewCredentials()))
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

	if len(errs) > 0 {
		log.Debug().Int("n_errs", len(errs)).Msg("could not successfully shutdown ensign server")
		return multierror.Append(err, errs...)
	}

	log.Debug().Msg("successfully shutdown ensign server")
	return nil
}
