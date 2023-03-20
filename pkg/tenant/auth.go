package tenant

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

// Lifetimes for CSRF token cookies that the server generates
const (
	protectLoginCSRFLifetime = time.Minute * 10
	authCSRFLifetime         = time.Hour * 12 // Should be longer than the access token lifetime
)

// Register is a publically accessible endpoint that allows new users to create an
// account via Quarterdeck by providing an email address and password.
//
// Route: POST /v1/register
func (s *Server) Register(c *gin.Context) {
	var err error
	ctx := c.Request.Context()

	// Parse the request body
	params := &api.RegisterRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse register request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Filter bad requests before they reach Quarterdeck
	// Note: This is a simple check to ensure that all required fields are present.
	if err = params.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Make the register request to Quarterdeck
	projectID := ulids.New()
	req := &qd.RegisterRequest{
		ProjectID:    projectID.String(),
		Name:         params.Name,
		Email:        params.Email,
		Password:     params.Password,
		PwCheck:      params.PwCheck,
		Organization: params.Organization,
		Domain:       params.Domain,
		AgreeToS:     params.AgreeToS,
		AgreePrivacy: params.AgreePrivacy,
	}

	var reply *qd.RegisterReply
	if reply, err = s.quarterdeck.Register(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Create member model for the new user
	member := &db.Member{
		ID:    reply.ID,
		OrgID: reply.OrgID,
		Name:  req.Name,
		Role:  reply.Role,
	}

	// Create a default tenant and project for the new user
	// Note: This task will error if the member model is invalid
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		return db.CreateUserResources(ctx, projectID, req.Organization, member)
	}), tasks.WithRetries(3),
		tasks.WithBackoff(backoff.NewExponentialBackOff()),
		tasks.WithError(fmt.Errorf("could not create default tenant and project for new user %s", reply.ID.String())),
	)

	// Add to SendGrid Ensign Marketing list in go routine
	// TODO: use worker queue to limit number of go routines for tasks like this
	// TODO: test in live integration tests to make sure this works
	hub := sentrygin.GetHubFromContext(c).Clone()
	go func() {
		contact := &sendgrid.Contact{
			Email: params.Email,
		}
		contact.ParseName(params.Name)

		if err := s.sendgrid.AddContact(contact); err != nil {
			log.Error().Err(err).Msg("could not add newly registered user to sendgrid ensign marketing list")
			if hub != nil {
				hub.CaptureException(err)
			}
		}
	}()

	// Return the response from Quarterdeck
	c.Status(http.StatusNoContent)
}

// Login is a publically accessible endpoint that allows users to login into their
// account via Quarterdeck and receive access and refresh tokens for future requests.
//
// Route: POST /v1/login
func (s *Server) Login(c *gin.Context) {
	var err error

	// Parse the request body
	params := &api.LoginRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Validate that required fields were provided
	if params.Email == "" || params.Password == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing email/password for login"))
		return
	}

	// Make the login request to Quarterdeck
	req := &qd.LoginRequest{
		Email:    params.Email,
		Password: params.Password,
	}

	var reply *qd.LoginReply
	if reply, err = s.quarterdeck.Login(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// TODO: Add user state checks and create required resources for first logins
	// (tenants, projects, etc.)

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(authCSRFLifetime)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set csrf protection cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set cookies"))
		return
	}

	// Return the access and refresh tokens from Quarterdeck
	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		LastLogin:    reply.LastLogin,
	}
	c.JSON(http.StatusOK, out)
}

// Refresh is a publicly accessible endpoint that allows users to refresh their
// access token using their refresh token. This enables frontend clients to provide a
// seamless login experience for the user.
//
// Route: POST /v1/refresh
func (s *Server) Refresh(c *gin.Context) {
	var err error

	// Parse the request body
	params := &api.RefreshRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse refresh request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Validate that required fields were provided
	if params.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing refresh token"))
		return
	}

	// Make the refresh request to Quarterdeck
	req := &qd.RefreshRequest{
		RefreshToken: params.RefreshToken,
	}
	var reply *qd.LoginReply
	if reply, err = s.quarterdeck.Refresh(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(authCSRFLifetime)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set csrf protection cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set cookies"))
		return
	}

	// Return the access and refresh tokens from Quarterdeck
	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		LastLogin:    reply.LastLogin,
	}
	c.JSON(http.StatusOK, out)
}

// ProtectLogin prepares the front-end for login by setting the double cookie
// tokens for CSRF protection.
func (s *Server) ProtectLogin(c *gin.Context) {
	expiresAt := time.Now().Add(protectLoginCSRFLifetime)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set csrf login protection cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set cookies"))
		return
	}
	c.JSON(http.StatusOK, &api.Reply{Success: true})
}

// VerifyEmail is a publicly accessible endpoint that allows users to verify their
// email address by supplying a token that was sent to their email address. If the
// token has already been verified, this endpoint returns a 202 Accepted response.
//
// Route: POST /v1/verify
func (s *Server) VerifyEmail(c *gin.Context) {
	var (
		err    error
		params *api.VerifyRequest
	)

	// Parse the request body
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse verify email request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Validate that required fields were provided
	if params.Token == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing token in request"))
		return
	}

	// Make the verify request to Quarterdeck
	req := &qd.VerifyRequest{
		Token: params.Token,
	}
	if err = s.quarterdeck.VerifyEmail(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Note: This obscures 202 Accepted responses as 200 OK responses which prevents
	// the user from being able to tell if they were already verified. To allow the
	// user to distinguish between the two we would have to return an error or modify
	// the response body to include that information.
	c.JSON(http.StatusOK, &api.Reply{Success: true})
}
