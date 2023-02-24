/*
Uptime is a small, lightweight service that tracks the status of Rotational services
and watches for outages. This package is intended to be a stand-alone service and has
limited tests since it is not a mission critical application.
*/
package uptime

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/uptime/config"
	"github.com/rotationalio/ensign/pkg/uptime/db"
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
	conf    config.Config
	monitor *Monitor
}

// Setup the server before the routes are configured.
func (s *Server) Setup() (err error) {
	// Setup our logging config first thing
	zerolog.SetGlobalLevel(s.conf.GetLogLevel())
	if s.conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Open the levelDB database
	if err = db.Connect(s.conf.DataPath, false); err != nil {
		log.Error().Err(err).Str("path", s.conf.DataPath).Msg("could not open leveldb")
		return err
	} else {
		log.Debug().Str("path", s.conf.DataPath).Msg("leveldb opened")
	}

	// Create the uptime monitor
	if s.monitor, err = NewMonitor(s.conf.StatusInterval, s.conf.ServiceInfo); err != nil {
		log.Error().Err(err).Msg("could not create monitor")
		return err
	} else {
		log.Debug().Dur("interval", s.conf.StatusInterval).Str("infoPath", s.conf.ServiceInfo).Msg("uptime monitor created")
	}

	return nil
}

// Called when the server has been started and is ready.
func (s *Server) Started() (err error) {
	s.monitor.Start()
	log.Info().Str("listen", s.URL()).Str("version", pkg.Version()).Msg("uptime server started")
	return nil
}

// Cleanup when the server is being shutdown. Note that in tests you should call
// Shutdown() to ensure the server is gracefully closed and not this method.
func (s *Server) Stop(ctx context.Context) (err error) {
	log.Info().Msg("gracefully shutting down the uptime server")

	if serr := s.monitor.Stop(ctx); serr != nil {
		log.Error().Err(err).Msg("could not stop uptime monitor")
		err = multierror.Append(err, serr)
	}

	if serr := db.Close(); serr != nil {
		log.Error().Err(err).Str("path", s.conf.DataPath).Msg("could not close leveldb")
		err = multierror.Append(err, serr)
	}

	return err
}

// Setup the server's middleware and routes.
func (s *Server) Routes(router *gin.Engine) (err error) {
	// Setup HTML template renderer
	var html *template.Template
	if html, err = template.ParseFS(content, "templates/*.html"); err != nil {
		return err
	}
	router.SetHTMLTemplate(html)

	// Setup static content server
	var static fs.FS
	if static, err = fs.Sub(content, "static"); err != nil {
		return err
	}
	router.StaticFS("/static", http.FS(static))

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

	// Add index route
	router.GET("/", s.Index)

	// NotFound and NotAllowed routes
	router.NoRoute(s.NotFound)
	router.NoMethod(s.NotAllowed)
	return nil
}

func (s *Server) NotFound(c *gin.Context) {
	c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (s *Server) NotAllowed(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}
