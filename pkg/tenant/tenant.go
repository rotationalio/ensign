package tenant

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	mw "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/tenant/db"
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

	// Adds the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

const ServiceName = "tenant"

func New(conf config.Config) (s *Server, err error) {
	// Loads the default configuration from the environment if the config is empty.
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Creates the server and prepares to serve
	s = &Server{
		Server: *service.New(conf.BindAddr, service.WithMode(conf.Mode)),
		conf:   conf,
	}

	s.Server.Register(s)
	return s, nil
}

// Server implements the service.Service interface and provides handlers to respond to
// Tenant-specific API routes and requests.
type Server struct {
	service.Server
	conf        config.Config        // server configuration
	ensign      *EnsignClient        // client to issue requests to Ensign
	quarterdeck qd.QuarterdeckClient // client to issue requests to Quarterdeck
	sendgrid    *emails.EmailManager // send emails and manage contacts
	tasks       *tasks.TaskManager   // task manager for performing background tasks
	topics      *TopicSubscriber     // consume topic updates from Ensign
	wg          *sync.WaitGroup      // waitgroup for go routines
}

// Setup the server before the routes are configured.
func (s *Server) Setup() (err error) {
	// Sets up logging config first
	zerolog.SetGlobalLevel(s.conf.GetLogLevel())
	if s.conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Configures Sentry
	if s.conf.Sentry.UseSentry() {
		if err = sentry.Init(s.conf.Sentry); err != nil {
			return fmt.Errorf("could not init sentry: %w", err)
		}
	}

	// Connect to services when not in maintenance mode
	if !s.conf.Maintenance {
		s.tasks = tasks.New(4, 64, time.Second)
		log.Debug().Int("workers", 4).Int("queue_size", 64).Msg("task manager started")

		// Connect to the trtl database
		if err = db.Connect(s.conf.Database); err != nil {
			return fmt.Errorf("could not connect to db: %w", err)
		}

		// Connect to Ensign
		if s.ensign, err = NewEnsignClient(s.conf.Ensign); err != nil {
			return fmt.Errorf("could not create ensign client: %w", err)
		}

		// Initialize the email manager
		if s.sendgrid, err = emails.New(s.conf.SendGrid); err != nil {
			return fmt.Errorf("could no init sendgrid: %w", err)
		}

		// Initialize the quarterdeck client
		if s.quarterdeck, err = s.conf.Quarterdeck.Client(); err != nil {
			return fmt.Errorf("could not create quarterdeck client: %w", err)
		}

		// Wait for specified duration until Quarterdeck is online and ready.
		ctx, cancel := context.WithTimeout(context.Background(), s.conf.Quarterdeck.WaitForReady)
		defer cancel()

		if err = s.quarterdeck.WaitForReady(ctx); err != nil {
			// Could not connect to Quarterdeck with the specified timeout, cannot start Tenant.
			return err
		}

		// Wait until ensign is ready
		if attempts, err := s.ensign.WaitForReady(); err != nil {
			return fmt.Errorf("could not connect to ensign after %d attempts: %w", attempts, err)
		}

		// Start the metatopic subscriber as a go routine
		s.topics = NewTopicSubscriber(s.ensign)
		s.wg = &sync.WaitGroup{}
		if err = s.topics.Run(s.wg); err != nil {
			return fmt.Errorf("could not start metatopic subscriber: %w", err)
		}
	}

	return nil
}

// Called when the server has been started and is ready.
func (s *Server) Started() (err error) {
	if s.conf.Maintenance {
		log.Warn().Msg("starting tenant server in maintenance mode")
	}

	// Startup services that cannot be started in maintenance mode.
	if !s.conf.Maintenance {
		if !s.conf.SendGrid.Enabled() {
			log.Warn().Msg("sendgrid is not enabled")
		}
	}

	log.Info().Str("listen", s.URL()).Str("version", pkg.Version()).Msg("tenant server started")
	return nil
}

// Cleanup when the server is being shutdown. Note that in tests you should call
// Shutdown() to ensure the server stops and not this method.
func (s *Server) Stop(context.Context) (err error) {
	log.Info().Msg("gracefully shutting down the tenant server")

	// Shutdown the running services
	if !s.conf.Maintenance {
		s.tasks.Stop()
		s.topics.Stop()

		// Wait for all go routines to finish
		s.wg.Wait()

		if err = db.Close(); err != nil {
			return fmt.Errorf("could not gracefully shutdown connection to trtldb: %w", err)
		}
	}

	// Flush sentry errors
	if s.conf.Sentry.UseSentry() {
		sentry.Flush(2 * time.Second)
	}

	log.Debug().Msg("successfully shutdown the tenant server")
	return nil
}

// Sets up the server's middleware and routes
func (s *Server) Routes(router *gin.Engine) (err error) {
	// If in maintenance mode, setup the maintenance routes and middleware.
	if s.conf.Maintenance {
		return s.MaintenanceRoutes(router)
	}

	// Set the authentication overrides from the configuration
	opts := mw.AuthOptions{
		Audience: s.conf.Auth.Audience,
		Issuer:   s.conf.Auth.Issuer,
		KeysURL:  s.conf.Auth.KeysURL,
	}

	// Creating the authenticator middleware requires a valid connection to Quarterdeck
	var authenticator gin.HandlerFunc
	if authenticator, err = mw.Authenticate(mw.WithAuthOptions(opts)); err != nil {
		return err
	}

	// Instantiate Sentry Handlers
	var tags gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": ServiceName}
		tags = sentry.TrackPerformance(tagmap)
	}

	var tracing gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": ServiceName}
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
		logger.GinLogger(ServiceName),

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
			router.Use(middleware)
		}
	}

	// CSRF protection is individually configured for POST, PUT, PATCH, and DELETE routes
	csrf := mw.DoubleCookie()

	// Initialize prometheus collectors (this function has a sync.Once so it's safe to call more than once)
	metrics.Setup()

	// Setup prometheus metrics (reserves the "/metrics" route)
	metrics.Routes(router)

	// Adds the v1 API routes
	v1 := router.Group("v1")
	{
		// Heartbeat route (authentication not required)
		v1.GET("/status", s.Status)

		// Set cookies for CSRF protection (authentication not required)
		v1.GET("/login", s.ProtectLogin)

		// User auth routes (authentication not required)
		v1.POST("/register", s.Register)
		v1.POST("/login", s.Login)
		v1.POST("/refresh", s.Refresh)
		v1.POST("/verify", s.VerifyEmail)
		v1.GET("/invites/:token", s.InvitePreview)

		// Authenticated routes
		v1.POST("/switch", authenticator, s.Switch)

		// Organization API routes must be authenticated
		organizations := v1.Group("/organization", authenticator)
		{
			organizations.GET("", mw.Authorize(perms.ReadOrganizations), s.OrganizationList)
			organizations.GET("/:orgID", mw.Authorize(perms.ReadOrganizations), s.OrganizationDetail)
		}

		// Tenant API routes must be authenticated
		tenant := v1.Group("/tenant", authenticator)
		{
			tenant.GET("", mw.Authorize(perms.ReadOrganizations), s.TenantList)
			tenant.POST("", csrf, mw.Authorize(perms.EditOrganizations), s.TenantCreate)
			tenant.GET("/:tenantID", mw.Authorize(perms.ReadOrganizations), s.TenantDetail)
			tenant.PUT("/:tenantID", csrf, mw.Authorize(perms.EditOrganizations), s.TenantUpdate)
			tenant.DELETE("/:tenantID", csrf, mw.Authorize(perms.DeleteOrganizations), s.TenantDelete)

			tenant.GET("/:tenantID/projects", mw.Authorize(perms.ReadProjects), s.TenantProjectList)
			tenant.POST("/:tenantID/projects", csrf, mw.Authorize(perms.EditProjects), s.TenantProjectCreate)
			tenant.PATCH("/:tenantID/projects/:projectID", csrf, mw.Authorize(perms.EditProjects), s.TenantProjectPatch)

			tenant.GET("/:tenantID/stats", mw.Authorize(perms.ReadOrganizations, perms.ReadProjects, perms.ReadTopics, perms.ReadAPIKeys), s.TenantStats)

			tenant.GET("/:tenantID/projects/stats", mw.Authorize(perms.ReadOrganizations, perms.ReadProjects, perms.ReadTopics, perms.ReadAPIKeys), s.TenantStats)

		}

		// Members API routes must be authenticated
		members := v1.Group("/members", authenticator)
		{
			members.GET("", mw.Authorize(perms.ReadCollaborators), s.MemberList)
			members.POST("", csrf, mw.Authorize(perms.AddCollaborators), s.MemberCreate)
			members.GET("/:memberID", mw.Authorize(perms.ReadCollaborators), s.MemberDetail)
			members.PUT("/:memberID", csrf, mw.Authorize(perms.EditCollaborators), s.MemberUpdate)
			members.POST("/:memberID", mw.Authorize(perms.EditCollaborators, perms.ReadCollaborators), s.MemberRoleUpdate)
			members.DELETE("/:memberID", csrf, mw.Authorize(perms.RemoveCollaborators), s.MemberDelete)
		}

		// Projects API routes must be authenticated
		projects := v1.Group("/projects", authenticator)
		{
			projects.GET("", mw.Authorize(perms.ReadProjects), s.ProjectList)
			projects.POST("", csrf, mw.Authorize(perms.EditProjects), s.ProjectCreate)
			projects.GET("/:projectID", mw.Authorize(perms.ReadProjects), s.ProjectDetail)
			projects.PUT("/:projectID", csrf, mw.Authorize(perms.EditProjects), s.ProjectUpdate)
			projects.PATCH("/:projectID", csrf, mw.Authorize(perms.EditProjects), s.ProjectPatch)
			projects.DELETE("/:projectID", csrf, mw.Authorize(perms.DeleteProjects), s.ProjectDelete)

			projects.GET("/:projectID/topics", mw.Authorize(perms.ReadTopics), s.ProjectTopicList)
			projects.POST("/:projectID/topics", csrf, mw.Authorize(perms.CreateTopics), s.ProjectTopicCreate)

			projects.GET("/:projectID/apikeys", mw.Authorize(perms.ReadAPIKeys), s.ProjectAPIKeyList)
			projects.POST("/:projectID/apikeys", csrf, mw.Authorize(perms.EditAPIKeys), s.ProjectAPIKeyCreate)
		}

		// Topics API routes must be authenticated
		topics := v1.Group("/topics", authenticator)
		{
			topics.GET("", mw.Authorize(perms.ReadTopics), s.TopicList)
			topics.POST("", csrf, mw.Authorize(perms.EditTopics), s.TopicCreate)
			topics.GET("/:topicID", mw.Authorize(perms.ReadTopics), s.TopicDetail)
			topics.PUT("/:topicID", csrf, mw.Authorize(perms.EditTopics), s.TopicUpdate)
			topics.DELETE("/:topicID", csrf, mw.Authorize(perms.DestroyTopics), s.TopicDelete)
		}

		// API key routes must be authenticated
		apikeys := v1.Group("/apikeys", authenticator)
		{
			apikeys.GET("", mw.Authorize(perms.ReadAPIKeys), s.APIKeyList)
			apikeys.POST("", csrf, mw.Authorize(perms.EditAPIKeys), s.APIKeyCreate)
			apikeys.GET("/:apiKeyID", mw.Authorize(perms.ReadAPIKeys), s.APIKeyDetail)
			apikeys.PUT("/:apiKeyID", csrf, mw.Authorize(perms.EditAPIKeys), s.APIKeyUpdate)
			apikeys.DELETE("/:apiKeyID", csrf, mw.Authorize(perms.DeleteAPIKeys), s.APIKeyDelete)
			apikeys.GET("/permissions", s.APIKeyPermissions)
		}
	}

	// NotFound and NotAllowed routes
	router.NoRoute(api.NotFound)
	router.NoMethod(api.NotAllowed)
	return nil
}

func (s *Server) MaintenanceRoutes(router *gin.Engine) (err error) {
	// Instantiate Sentry Handlers
	var tags gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": ServiceName}
		tags = sentry.TrackPerformance(tagmap)
	}

	var tracing gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": ServiceName}
		tracing = sentry.TrackPerformance(tagmap)
	}

	// Application Middleware
	middlewares := []gin.HandlerFunc{
		logger.GinLogger(ServiceName),
		tags,
		tracing,
		s.Available(),
	}

	// Adds middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			router.Use(middleware)
		}
	}

	// Add the status route
	router.GET("/v1/status", s.Status)

	// NotFound and NotAllowed routes
	router.NoRoute(api.NotFound)
	router.NoMethod(api.NotAllowed)
	return nil
}

//===========================================================================
// Accessor Methods
//===========================================================================

// Expose the Ensign client to the tests (only allowed in testing mode).
func (s *Server) GetEnsignClient() *EnsignClient {
	if s.conf.Mode == gin.TestMode {
		return s.ensign
	}
	log.Fatal().Msg("can only get ensign client in test mode")
	return nil
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
		s.tasks.Stop()
		if s.tasks.IsStopped() {
			s.tasks = tasks.New(4, 64, time.Second)
		}
		return
	}
	log.Fatal().Msg("can only reset task manager in test mode")
}
