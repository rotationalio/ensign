/*
Package service unifies how our HTTP services run to deduplicate code between different
microservices. The common Service struct should be embedded into custom server
implementations and it should be registered with a Server that implements the routes
and handler functionality. The Service struct handles log initialization, the
construction and management of the http Server, the construction and management of a
gin Engine for multiplexing, and the liveness and readiness probes for Kubernetes.
*/
package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/rotationalio/ensign/pkg/utils/probez"
)

// Service specifies the methods that a service implementation for a specific
// microservice needs to implement in order to be registered with the http server. The
// simplest services only need to specify the Service interface in order to create a
// quickly working service!
type Service interface {
	// Routes is called once the server is up and running and is in a live state but
	// before it is in a ready state. The intended use case of this method is for the
	// Service to specify its own middleware, routes, and handlers on the router, which
	// is the majority of what differentiates a service from the main server.
	// If an error is returned, the error will shutdown the server and stop startup.
	Routes(router *gin.Engine) error

	// Stop is called after the http server has gracefully concluded serving and the
	// service is in a not ready, not healthy state. This method should be used to
	// clean up connections to databases, close resources properly, etc. Note that the
	// context may have a deadline on it before the server is terminated.
	// If an error is returned, the error will be reported to the caller of Serve().
	Stop(ctx context.Context) error
}

// Initializer services want to do work before the server is able to serve liveness
// probes (e.g. there is no http server running at all). It is strongly recommended that
// services specify Setup() instead of Initialize() unless work is needed to done in a
// non-live state for some reason.
type Initializer interface {
	// Initialize is called before the server is in a live state and the server is not
	// serving any requests (unhealthy, no liveness probes) and is the first possible
	// entry point to the service interface. If an error is returned from Initialize,
	// the server will not start at all.
	Initialize() error
}

// Preparer services want to do work before the server is started up such as connecting
// to databases or other external resources.
type Preparer interface {
	// Setup is called after initialize and after the server has started serving
	// liveness probes (e.g. it is healthy but not ready). It is called prior to routes
	// and is intended for more general setup and connection to databases and other
	// resources. If an error is returned, the server will shutdown and will not proceed
	// with its startup. Retry logic to delay readiness should be added to this method.
	Setup() error
}

// Starter services want to do work (usually logging or booting up background routines)
// when the server is started up and ready.
type Starter interface {
	// Started is called after the server is serving and ready, after the Routes method
	// has been called. Errors returned from this function will cause the server to
	// shutdown will will have the effect of calling the services Shutdown method as
	// well. Any errors returned from shutdown or this method are returned as a
	// multierror to report everything that happened in the error process. This method
	// is generally used for logging or booting up background routines and is a good
	// choice for work that needs to be done that is error prone but may not affect the
	// service or the ability of the service to handle requests.
	Started() error
}

// Server implements common functionality for all service implementations that use Gin.
// Users should create a new server then register their service with the server in order
// to serve the server to respond to requests. After the server is created when the
// user calls its serve method, the server uses the registered service as it transitions
// through startup and shutdown states. Briefly those states are:
//
// Initialize: called before the server is started (unhealthy, no liveness probes)
// Setup: called after the server is started but before routes are setup (healthy, not ready)
// Routes: called after setup and is intended to setup service routes (healthy, not ready)
// Started: called after the server is in a ready state (healthy, ready)
// Shutdown: called after the server is shutdown to cleanup the service (unhealthy, no liveness probes)
type Server struct {
	sync.RWMutex
	service Service
	srv     *http.Server
	router  *gin.Engine
	status  *probez.Handler
	started time.Time
	url     *url.URL
	errc    chan error
}

func New(addr string, opts ...Option) *Server {
	options := newOptions(opts...)
	srv := &Server{
		srv:    options.server,
		router: options.router,
		status: probez.New(),
		errc:   make(chan error, 1),
	}

	// Create the gin router if not specified by the user
	if options.mode != "" {
		gin.SetMode(options.mode)
	}

	// Create and configure the gin router
	if srv.router == nil {
		srv.router = gin.New()
		srv.router.RedirectTrailingSlash = true
		srv.router.RedirectFixedPath = false
		srv.router.HandleMethodNotAllowed = true
		srv.router.ForwardedByClientIP = true
		srv.router.UseRawPath = false
		srv.router.UnescapePathValues = true
	}

	// Create the http server if it was not specified by the user
	if srv.srv == nil {
		srv.srv = &http.Server{
			Addr:              addr,
			Handler:           srv.router,
			ErrorLog:          nil,
			ReadHeaderTimeout: 20 * time.Second,
			WriteTimeout:      20 * time.Second,
			IdleTimeout:       30 * time.Second,
		}
	}

	// Ensure the Addr is set on the server even if specified by the user.
	srv.srv.Addr = addr
	return srv
}

func (s *Server) Register(service Service) {
	s.service = service
}

// Serve API requests while listening on the specified bind address.
func (s *Server) Serve() (err error) {
	if s.service == nil {
		return ErrNoServiceRegistered
	}

	// Initialize the service if the service is an initializer
	if err = s.initialize(); err != nil {
		return err
	}

	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit

		// Require the shutdown occurs in 10 seconds without blocking
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.errc <- s.Shutdown(ctx)
	}()

	// Handle liveness and readiness probe requests
	s.status.Use(s.router)

	// Create a socket to listen on and infer the final URL.
	// NOTE: if the bindaddr is 127.0.0.1:0 for testing, a random port will be assigned,
	// manually creating the listener will allow us to determine which port.
	// When we start listening all incoming requests will be buffered until the server
	// actually starts up in its own go routine below.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.srv.Addr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %s", s.srv.Addr, err)
	}

	// Set the URL from the listener
	s.setURL(sock.Addr())
	s.status.Healthy()

	// Listen for HTTP requests and handle them.
	go func() {
		// Make sure we don't use the external err to avoid data races.
		if serr := s.srv.Serve(sock); !errors.Is(serr, http.ErrServerClosed) {
			s.errc <- serr
		}

		// If there is no error, return nil so this function exits if Shutdown is
		// called manually (e.g. not from an OS signal).
		s.errc <- nil
	}()

	// Call the setup method while we're healthy but not ready
	if err = s.setup(); err != nil {
		s.srv.Shutdown(context.Background())
		return err
	}

	// Call the service routes method to enable the service interface.
	// NOTE: this relies on the fact that a gin router is able to handle dynamic routes;
	// that is even though the server is running we can add any arbitrary handlers.
	if err = s.service.Routes(s.router); err != nil {
		s.srv.Shutdown(context.Background())
		return err
	}

	// Set the started timestamp and the service as ready
	s.setStartTime()
	s.status.Ready()

	// Call the startup method now that we're done setting up and initializing the server
	if err = s.startup(); err != nil {
		if serr := s.Shutdown(context.Background()); serr != nil {
			err = multierror.Append(err, serr)
		}
		return err
	}

	// Listen for any errors that might have occurred and wait for all go routines to stop
	return <-s.errc
}

// Shutdown the server gracefully (usually called by OS signal) but can be
// called by other triggers or manually during the test.
func (s *Server) Shutdown(ctx context.Context) (err error) {
	s.status.NotHealthy()
	s.status.NotReady()
	s.srv.SetKeepAlivesEnabled(false)

	if serr := s.service.Stop(ctx); serr != nil {
		err = multierror.Append(err, serr)
	}

	// NOTE: this must come last otherwise the ErrServerClosed will terminate the Serve thread.
	if serr := s.srv.Shutdown(ctx); serr != nil {
		err = multierror.Append(err, serr)
	}
	return err
}

// Set the URL from the TCPAddr when the server is started. Should be set by Serve().
func (s *Server) setURL(addr net.Addr) {
	s.Lock()
	defer s.Unlock()
	s.url = &url.URL{
		Scheme: "http",
		Host:   addr.String(),
	}

	if tcp, ok := addr.(*net.TCPAddr); ok && tcp.IP.IsUnspecified() {
		s.url.Host = fmt.Sprintf("127.0.0.1:%d", tcp.Port)
	}
}

// URL returns the URL of the server determined by the socket addr.
func (s *Server) URL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url.String()
}

// Set the StartTime in a thread safe manner.
func (s *Server) setStartTime() {
	s.Lock()
	defer s.Unlock()
	s.started = time.Now()
}

// StartTime returns the time that the server started.
func (s *Server) StartTime() time.Time {
	s.RLock()
	defer s.RUnlock()
	return s.started
}

// SetHealth is used by tests to set the health of the server
func (s *Server) SetHealth(health bool) {
	if health {
		s.status.Healthy()
	} else {
		s.status.NotHealthy()
	}
}

// IsHealthy returns whether the status probe is healthy or not.
func (s *Server) IsHealthy() bool {
	return s.status.IsHealthy()
}

// IsReady returns whether the status probe is ready or not.
func (s *Server) IsReady() bool {
	return s.status.IsReady()
}

func (s *Server) initialize() error {
	if initializer, ok := s.service.(Initializer); ok {
		return initializer.Initialize()
	}
	return nil
}

func (s *Server) setup() error {
	if preparer, ok := s.service.(Preparer); ok {
		return preparer.Setup()
	}
	return nil
}

func (s *Server) startup() error {
	if starter, ok := s.service.(Starter); ok {
		return starter.Started()
	}
	return nil
}
