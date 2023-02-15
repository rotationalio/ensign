/*
Uptime is a small, lightweight service that tracks the status of Rotational services
and watches for outages. This package is intended to be a stand-alone service and has
limited tests since it is not a mission critical application.
*/
package uptime

import (
	"context"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/uptime/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/service"
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
	// Load the default configuration from the environment if config is empty.
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Create the service and register it with the server.
	s = &Server{
		Server: *service.New(conf.BindAddr, service.WithMode(conf.Mode)),
		conf:   conf,
	}

	s.Server.Register(s)
	return s, nil
}

type Server struct {
	service.Server
	conf config.Config
}

// Setup the server before the routes are configured.
func (s *Server) Setup() (err error) {
	// Setup our logging config first thing
	zerolog.SetGlobalLevel(s.conf.GetLogLevel())
	if s.conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return nil
}

// Called when the server has been started and is ready.
func (s *Server) Started() (err error) {
	log.Info().Str("listen", s.URL()).Str("version", pkg.Version()).Msg("uptime server started")
	return nil
}

// Cleanup when the server is being shutdown. Note that in tests you should call
// Shutdown() to ensure the server is gracefully closed and not this method.
func (s *Server) Stop(ctx context.Context) (err error) {
	log.Info().Msg("gracefully shutting down the uptime server")
	return nil
}

// Setup the server's middleware and routes.
func (s *Server) Routes(router *gin.Engine) (err error) {
	// Setup CORS configuration
	corsConf := cors.Config{
		AllowMethods: []string{"GET", "HEAD"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
		AllowOrigins: s.conf.AllowOrigins,
		MaxAge:       12 * time.Hour,
	}

	// Application Middleware
	// NOTE: ordering is important to how middleware is handled
	middlewares := []gin.HandlerFunc{
		// Logging should be on the outside so we can record the correct latency of requests
		// NOTE: logging panics will not recover
		logger.GinLogger("uptime"),

		// Panic recovery middleware
		gin.Recovery(),

		// CORS configuration allows the front-end to make cross-origin requests
		cors.New(corsConf),
	}

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			router.Use(middleware)
		}
	}

	// NotFound and NotAllowed routes
	router.NoRoute(api.NotFound)
	router.NoMethod(api.NotAllowed)
	return nil
}
