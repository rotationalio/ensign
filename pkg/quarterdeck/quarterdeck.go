package quarterdeck

import (
	"context"
	"errors"
	"os"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/metrics"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/service"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
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

const ServiceName = "quarterdeck"

func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment if the config is empty.
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Maintenance mode configuration checks
	if !conf.Maintenance {
		if len(conf.Token.Keys) == 0 {
			return nil, errors.New("invalid configuration: no token keys specified when not in maintenance mode")
		}
	}

	// Create the server and register it with the default service.
	s = &Server{
		Server: *service.New(conf.BindAddr, service.WithMode(conf.Mode)),
		conf:   conf,
	}

	s.Server.Register(s)
	return s, nil
}

// Server implements the service.Service interface and provides handlers to respond to
// Quarterdeck-specific API routes and requests.
type Server struct {
	service.Server
	conf     config.Config        // the server configuration
	tokens   *tokens.TokenManager // token manager for issuing JWT tokens for authentication
	tasks    *tasks.TaskManager   // task manager for performing background tasks
	sendgrid *emails.EmailManager // send emails and manage contacts
}

// Setup the server before the routes are configured.
func (s *Server) Setup() (err error) {
	// Setup our logging config first thing
	zerolog.SetGlobalLevel(s.conf.GetLogLevel())
	if s.conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Configure Sentry
	if s.conf.Sentry.UseSentry() {
		if err = sentry.Init(s.conf.Sentry); err != nil {
			return err
		}
	}

	// If the server is not in maintenance mode setup and configure required services.
	if !s.conf.Maintenance {
		if s.tokens, err = tokens.New(s.conf.Token); err != nil {
			return err
		}

		if s.sendgrid, err = emails.New(s.conf.SendGrid); err != nil {
			return err
		}

		s.tasks = tasks.New(4, 64)
		log.Debug().Int("workers", 4).Int("queue_size", 64).Msg("task manager started")

		if err = db.Connect(s.conf.Database.URL, s.conf.Database.ReadOnly); err != nil {
			return err
		}
		log.Debug().Bool("read-only", s.conf.Database.ReadOnly).Str("dsn", s.conf.Database.URL).Msg("connected to database")
	}

	return nil
}

// Called when the server has been started and is ready.
func (s *Server) Started() (err error) {
	if s.conf.Maintenance {
		log.Warn().Msg("starting quarterdeck server in maintenance mode")
	}

	log.Info().Str("listen", s.URL()).Str("version", pkg.Version()).Msg("quarterdeck server started")
	return nil
}

// Cleanup when the server is being shutdown. Note that in tests you should call
// Shutdown() to ensure the server is gracefully closed and not this method.
func (s *Server) Stop(ctx context.Context) (err error) {
	log.Info().Msg("gracefully shutting down the quarterdeck server")

	// Close the database connection
	if !s.conf.Maintenance {
		s.tasks.Stop()

		if err = db.Close(); err != nil {
			return err
		}
	}

	// Flush sentry errors
	if s.conf.Sentry.UseSentry() {
		sentry.Flush(2 * time.Second)
	}
	return nil
}

// Setup the server's middleware and routes.
func (s *Server) Routes(router *gin.Engine) (err error) {
	// Instantiate Sentry Handlers
	var tags gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": ServiceName}
		tags = sentry.UseTags(tagmap)
	}

	var tracing gin.HandlerFunc
	if s.conf.Sentry.UsePerformanceTracking() {
		tagmap := map[string]string{"service": ServiceName}
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
		logger.GinLogger(ServiceName),

		// Panic recovery middleware
		// NOTE: gin middleware needs to be added before sentry
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),
		//TODO: add ratelimiter
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
			router.Use(middleware)
		}
	}

	// Instantiate per-route middleware.
	var authenticate gin.HandlerFunc
	if authenticate, err = middleware.Authenticate(middleware.WithValidator(s.tokens)); err != nil {
		return err
	}

	//router.Use(middleware.RateLimiter())

	// Initialize prometheus collectors (this function has a sync.Once so it's safe to call more than once)
	metrics.Setup()

	// Setup prometheus metrics (reserves the "/metrics" route)
	metrics.Routes(router)

	// Add the v1 API routes
	v1 := router.Group("/v1")
	{
		// Heartbeat route (no authentication required)
		v1.GET("/status", s.Status)

		// Unauthenticated access routes
		v1.POST("/register", s.Register)
		v1.POST("/login", s.Login)
		v1.POST("/authenticate", s.Authenticate)
		v1.POST("/refresh", s.Refresh)
		v1.POST("/verify", s.VerifyEmail)

		// Organizations Resource
		orgs := v1.Group("/organizations", authenticate)
		{
			orgs.GET("/:id", middleware.Authorize(perms.ReadOrganizations), s.OrganizationDetail)
		}

		// API Keys Resource
		apikeys := v1.Group("/apikeys", authenticate)
		{
			apikeys.GET("", middleware.Authorize(perms.ReadAPIKeys), s.APIKeyList)
			apikeys.POST("", middleware.Authorize(perms.EditAPIKeys), s.APIKeyCreate)
			apikeys.GET("/:id", middleware.Authorize(perms.ReadAPIKeys), s.APIKeyDetail)
			apikeys.PUT("/:id", middleware.Authorize(perms.EditAPIKeys), s.APIKeyUpdate)
			apikeys.DELETE("/:id", middleware.Authorize(perms.DeleteAPIKeys), s.APIKeyDelete)
			apikeys.GET("/permissions", s.APIKeyPermissions)
		}

		// Projects Resource
		projects := v1.Group("/projects", authenticate)
		{
			projects.POST("", middleware.Authorize(perms.EditProjects), s.ProjectCreate)
			projects.POST("/access", middleware.Authorize(perms.ReadTopics), s.ProjectAccess)
		}

		// Users Resource
		users := v1.Group("/users", authenticate)
		{
			users.GET("/:id", middleware.Authorize(perms.ReadCollaborators), s.UserDetail)
			users.PUT("/:id", middleware.Authorize(perms.EditCollaborators), s.UserUpdate)
			users.GET("", middleware.Authorize(perms.ReadCollaborators), s.UserList)
			users.DELETE("/:id", middleware.Authorize(perms.RemoveCollaborators), s.UserDelete)
		}

		// Accounts Resource - endpoint for users to manage their own account
		accounts := v1.Group("/accounts", authenticate)
		{
			accounts.PUT("/:id", s.AccountUpdate)
		}
	}

	// The "well known" routes expose client security information and credentials.
	wk := router.Group("/.well-known")
	{
		wk.GET("/jwks.json", s.JWKS)
		wk.GET("/security.txt", s.SecurityTxt)
		wk.GET("/openid-configuration", s.OpenIDConfiguration)
	}

	// NotFound and NotAllowed routes
	router.NoRoute(api.NotFound)
	router.NoMethod(api.NotAllowed)
	return nil
}

// AccessToken returns a token that can be used in tests and is only available if the
// server is in testing mode, otherwise an empty string is returned.
func (s *Server) AccessToken(claims *tokens.Claims) string {
	if s.conf.Mode == gin.TestMode {
		token, err := s.tokens.CreateAccessToken(claims)
		if err != nil {
			panic(err)
		}

		atks, err := s.tokens.Sign(token)
		if err != nil {
			panic(err)
		}
		return atks
	}
	return ""
}

// Return an access and refresh token that can be used in tests and is only available
// if the server is in testing mode, otherwise empty strings are returned. It is preferred
// to use the AccessToken() function for most tests, use this function if a refresh
// token is required for testing.
func (s *Server) CreateTokenPair(claims *tokens.Claims) (string, string) {
	if s.conf.Mode == gin.TestMode {
		accessToken, refreshToken, err := s.tokens.CreateTokenPair(claims)
		if err != nil {
			panic(err)
		}
		return accessToken, refreshToken
	}
	return "", ""
}

// VerifyToken extracts the claims from an access or refresh token returned by the
// server. This is only available if the server is in testing mode.
func (s *Server) VerifyToken(tks string) (*tokens.Claims, error) {
	if s.conf.Mode == gin.TestMode {
		return s.tokens.Verify(tks)
	}
	log.Fatal().Msg("can only verify tokens in test mode")
	return nil, errors.New("can only use this method in test mode")
}

// Expose the task manager to the tests (only allowed in testing mode).
func (s *Server) GetTaskManager() *tasks.TaskManager {
	if s.conf.Mode == gin.TestMode {
		return s.tasks
	}
	log.Fatal().Msg("can only get task manager in test mode")
	return nil
}

// Reset the task manager from the tests (only allowed in testing mode)
func (s *Server) ResetTaskManager() {
	if s.conf.Mode == gin.TestMode {
		if s.tasks.IsStopped() {
			s.tasks = tasks.New(4, 64)
		}
		return
	}
	log.Fatal().Msg("can only reset task manager in test mode")
}
