package quarterdeck

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Initializes zerolog with our default logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment if the config is empty.
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Setup our logging config first thing
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Configure Sentry
	if conf.Sentry.UseSentry() {
		if err = sentry.Init(conf.Sentry); err != nil {
			return nil, err
		}
	}

	// Create the server and prepare to serve
	s = &Server{
		conf: conf,
		errc: make(chan error, 1),
	}

	// TODO: handle maintenance mode

	// Create the router
	gin.SetMode(conf.Mode)
	s.router = gin.New()
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	s.srv = &http.Server{
		Addr:         s.conf.BindAddr,
		Handler:      s.router,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s, nil
}

// Server implements the API router and handlers.
type Server struct {
	sync.RWMutex
	conf    config.Config // the server configuration
	srv     *http.Server  // the http server to handle requests on
	router  *gin.Engine   // the router that defines the http handler
	started time.Time     // the time that the server was started
	healthy bool          // if we're online or shutting down
	url     string        // the external url of the server from the socket
	errc    chan error    // any errors sent on this channel are fatal
}

// Serve API requests while listening on the specified bind address.
func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.errc <- s.Shutdown()
	}()

	// Set the health of the service to true unless we're in maintenance mode
	s.SetHealth(!s.conf.Maintenance)
	if s.conf.Maintenance {
		log.Warn().Msg("starting quarterdeck server in maintenance mode")
	}

	// Create a socket to listen on and infer the final URL.
	// NOTE: if the bindaddr is 127.0.0.1:0 for testing, a random port will be assigned,
	// manually creating the listener will allow us to determine which port.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %s", s.conf.BindAddr, err)
	}

	// Set the URL from the listener
	s.SetURL("http://" + sock.Addr().String())
	s.started = time.Now()

	// Listen for HTTP requests and handle them.
	go func() {
		if err = s.srv.Serve(sock); err != nil && err != http.ErrServerClosed {
			s.errc <- err
		}

		// If there is no error, return nil so this function exits if Shutdown is
		// called manually (e.g. not from an OS signal).
		s.errc <- nil
	}()

	log.Info().Str("listen", s.url).Str("version", pkg.Version()).Msg("quarterdeck server started")

	// Listen for any errors that might have occurred and wait for all go routines to stop
	return <-s.errc
}

// Shutdown the server gracefully (usually called by OS signal).
func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down the quarterdeck server")

	s.SetHealth(false)
	s.srv.SetKeepAlivesEnabled(false)

	// Require the shutdown occurs in 30 seconds without blocking
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Debug().Msg("successfully shutdown the quarterdeck server")
	return nil
}

// Setup the server's middleware and routes (done once in New).
func (s *Server) setupRoutes() error {
	// Instantiate Sentry Handlers
	var tags gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": "quarterdeck"}
		tags = sentry.UseTags(tagmap)
	}

	var tracing gin.HandlerFunc
	if s.conf.Sentry.UsePerformanceTracking() {
		tagmap := map[string]string{"service": "quarterdeck"}
		tracing = sentry.TrackPerformance(tagmap)
	}

	// Setup CORS configuration
	corsConf := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN", "sentry-trace", "baggage"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	if s.conf.AllowAllOrigins() {
		corsConf.AllowAllOrigins = true
	} else {
		corsConf.AllowOrigins = s.conf.AllowOrigins
	}

	// Application Middleware
	// NOTE: ordering is important to how middleware is handled
	middlewares := []gin.HandlerFunc{
		// Logging should be on the outside so we can record the correct latency of requests
		// NOTE: logging panics will not recover
		logger.GinLogger("quarterdeck"),

		// Panic recovery middleware
		// NOTE: gin middleware needs to be added before sentry
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),

		// Add searchable tags to sentry context
		tags,

		// Tracing helps us measure performance metrics with Sentry
		tracing,

		// CORS configuration allows the front-end to make cross-origin requests
		cors.New(corsConf),

		// Maintenance mode handling - should not require authentication
		s.Available(),
	}

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Add the v1 API routes
	v1 := s.router.Group("/v1")
	{
		// Heartbeat route (no authentication required)
		v1.GET("/status", s.Status)
	}

	// NotFound and NotAllowed routes
	s.router.NoRoute(api.NotFound)
	s.router.NoMethod(api.NotAllowed)
	return nil
}

func (s *Server) SetHealth(health bool) {
	s.Lock()
	s.healthy = health
	s.Unlock()
	log.Debug().Bool("healthy", health).Msg("server health set")
}

func (s *Server) Healthy() bool {
	s.RLock()
	defer s.RUnlock()
	return s.healthy
}

func (s *Server) SetURL(url string) {
	s.Lock()
	s.url = url
	s.Unlock()
	log.Debug().Str("url", url).Msg("server url set")
}

func (s *Server) URL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url
}
