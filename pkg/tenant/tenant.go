package tenant

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
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/tenant/db"
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

	// Adds the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

func New(conf config.Config) (s *Server, err error) {
	// Loads the default configuration from the environment if the config is empty.
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Sets up logging config first
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Configures Sentry
	if conf.Sentry.UseSentry() {
		if err = sentry.Init(conf.Sentry); err != nil {
			return nil, err
		}
	}

	// Creates the server and prepares to serve
	s = &Server{
		conf: conf,
		errc: make(chan error, 1),
	}

	// Connect to services when not in maintenance mode
	if !s.conf.Maintenance {
		// Connect to the trtl database
		if err = db.Connect(s.conf.Database); err != nil {
			return nil, err
		}
	}

	// Creates the router
	gin.SetMode(conf.Mode)
	s.router = gin.New()
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Creates the http server
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
	conf    config.Config // server configuration
	srv     *http.Server  // http server that handles requests
	router  *gin.Engine   // router that defines the http handler
	started time.Time     // time that the server was started
	healthy bool          // states if we're online or shutting down
	url     string        // external url of the server from the socket
	errc    chan error    // any errors sent to this channel are fatal
}

// Serves API requests while listening on the specified bind address.
func (s *Server) Serve() (err error) {
	// Catches OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.errc <- s.Shutdown()
	}()

	// Sets health of the service to true unless in maintenance mode
	s.SetHealth(!s.conf.Maintenance)
	if s.conf.Maintenance {
		log.Warn().Msg("starting tenant server in maintenance mode")
	}

	// Startup services that cannot be started in maintenance mode.
	if !s.conf.Maintenance {
		if !s.conf.SendGrid.Enabled() {
			log.Warn().Msg("sendgrid is not enabled")
		}
	}

	// Creates a socket to listen on and infer the final URL.
	// NOTE: if the bindaddr is 127.0.0.1:0 for testing, a random port will be assigned,
	// manually creating the listener will allow us to determine which port.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %s", s.conf.BindAddr, err)
	}

	// Sets URL from the listener
	s.SetURL("http://" + sock.Addr().String())
	s.started = time.Now()

	// Listens for HTTP requests and handles them.
	go func() {
		if err = s.srv.Serve(sock); err != nil && err != http.ErrServerClosed {
			s.errc <- err
		}
		// If there isn't an error, return nil so that this function exits if
		// Shutdown is called manually.
		s.errc <- nil
	}()

	log.Info().Str("listen", s.url).Str("version", pkg.Version()).Msg("tenant server started")

	//Listens for any errors that might have occurred and waits for all go routines to stop
	return <-s.errc
}

// Shuts down the server gracefully
func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down the tenant server")

	s.SetHealth(false)
	s.srv.SetKeepAlivesEnabled(false)

	// Close connection to the trtl database
	if err = db.Close(); err != nil {
		log.Warn().Err(err).Msg("could not gracefully shutdown connection to trtl database")
	}

	// Requires shutdown occurs in 30 seconds without blocking.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Debug().Msg("successfully shutdown the tenant server")
	return nil
}

// Sets up the server's middleware and routes
func (s *Server) setupRoutes() error {
	// Instantiates Sentry Handlers
	var tags gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": "tenant"}
		tags = sentry.TrackPerformance(tagmap)
	}

	var tracing gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": "tenant"}
		tracing = sentry.TrackPerformance(tagmap)
	}

	// Sets up CORS configuration
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
	middlewares := []gin.HandlerFunc{
		// Logging should be on the outside so that we can record the correct latency of requests
		logger.GinLogger("tenant"),

		// Panic recovery middleware
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),

		// Adds searchable tags to sentry context
		tags,

		// Tracing helps measure performance metrics with Sentry
		tracing,

		// CORS configuration allows the front-end to make cross origin requests
		cors.New(corsConf),

		// Maintenance mode handling - should not require authentication
		s.Available(),
	}
	// Adds middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Adds the v1 API routes
	v1 := s.router.Group("v1")
	{
		// Heartbeat route (authentication not required)
		v1.GET("/status", s.Status)

		// Notification signups (authentication not required)
		v1.POST("/notifications/signup", s.SignUp)

		// Adds tenant to the API routes
		// Routes to tenants
		v1.GET("/tenant", s.TenantList)
		v1.GET("/tenant/:tenantID", s.TenantDetail)
		v1.POST("/tenant", s.TenantCreate)
		v1.PUT("/tenant/:tenantID", s.TenantUpdate)
		v1.DELETE("/tenant/:tenantID", s.TenantDelete)

		// Routes to members
		v1.GET("/tenant/:tenantID/members", s.TenantMemberList)
		v1.POST("/tenant/:tenantID/members", s.TenantMemberCreate)

		v1.GET("/members", s.MemberList)
		v1.GET("/members/:memberID", s.MemberDetail)
		v1.POST("/members", s.MemberCreate)
		v1.PUT("/members/:memberID", s.MemberUpdate)
		v1.DELETE("/members/:memberID", s.MemberDelete)

		// Routes to projects
		v1.GET("/tenant/:tenantID/projects", s.TenantProjectList)
		v1.POST("/tenant/:tenantID/projects", s.TenantProjectCreate)

		v1.GET("/projects", s.ProjectList)
		v1.GET("/projects/:projectID", s.ProjectDetail)
		v1.POST("/projects", s.ProjectCreate)
		v1.PUT("/projects/:projectID", s.ProjectUpdate)
		v1.DELETE("/projects/:projectID", s.ProjectDelete)

		// Routes to topics
		v1.GET("/projects/:projectID/topics", s.ProjectTopicList)
		v1.POST("/projects/:projectID/topics", s.ProjectTopicCreate)

		v1.GET("/topics", s.TopicList)
		v1.POST("/topics", s.TopicCreate)
		v1.GET("/topics/:topicID", s.TopicDetail)
		v1.PUT("/topics/:topicID", s.TopicUpdate)
		v1.DELETE("/topics/:topicID", s.TopicDelete)

		// Routes to APIKeys
		v1.GET("/projects/:projectID/aoikeys", s.ProjectAPIKeyList)
		v1.POST("/projects/:projectID/apikeys", s.ProjectAPIKeyCreate)

		v1.GET("/apikeys", s.APIKeyList)
		v1.GET("/apikeys/:apiKeyID", s.APIKeyDetail)
		v1.POST("/apikeys", s.APIKeyCreate)
		v1.PUT("/apikeys/:apiKeyID", s.APIKeyUpdate)
		v1.DELETE("/apikeys/:apiKeyID", s.APIKeyDelete)
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
