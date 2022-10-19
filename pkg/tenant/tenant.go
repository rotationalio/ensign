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

	// TODO: handle maintenance mode

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

		// Routes to members
		members := v1.Group("organization/members")
		{
			members.GET("", MemberList)
			members.GET("/:id", MemberDetail)
			members.POST("", MemberCreate)
			members.PUT("/:id", MemberUpdate)
			members.DELETE("/:id", MemberDelete)
		}

		// Routes to projects
		projects := v1.Group("organization/projects")
		{
			projects.GET("", ProjectList)
			projects.GET("/:id", ProjectDetail)
			projects.POST("", ProjectCreate)
			projects.PUT("/:id", ProjectUpdate)
			projects.DELETE("/:id", ProjectDelete)
		}

		// Routes to topics
		topics := v1.Group("organization/projects/topics")
		{
			topics.GET("", TopicList)
			topics.GET("/:id", TopicDetail)
			topics.POST("", TopicCreate)
			topics.PUT("/:id", TopicUpdate)
			topics.DELETE("/:id", TopicDelete)
		}

		// Routes to APIKeys
		apikeys := v1.Group("organization/projects/apikeys")
		{
			apikeys.GET("", APIKeyList)
			apikeys.GET("/:id", APIKeyDetail)
			apikeys.POST("", APIKeyCreate)
			apikeys.PUT("/:id", APIKeyUpdate)
			apikeys.DELETE("/:id", APIKeyDelete)
		}
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

// User handlers for Tenant stub. The below functions are
// currently listed below for testing purposes and will be moved
// to the users file once Tenant Client interfaces are updated
// in the api file.

func MemberList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func MemberDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func MemberCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func MemberUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func MemberDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func ProjectList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func ProjectDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func ProjectCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func ProjectUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func ProjectDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func TopicList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func TopicDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func TopicCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func TopicUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func TopicDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func APIKeyList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func APIKeyDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func APIKeyCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func APIKeyUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func APIKeyDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
