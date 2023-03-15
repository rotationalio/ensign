package tenant

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
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
		log.Warn().Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse register request"))
		return
	}

	// Filter bad requests before they reach Quarterdeck
	// Note: This is a simple check to ensure that all required fields are present.
	if err = params.Validate(); err != nil {
		log.Warn().Err(err).Msg("missing required fields")
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
		log.Error().Err(err).Msg("could not register user")
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
	// Note: This method returns an error if the member model is invalid
	if err = db.CreateUserResources(ctx, projectID, req.Organization, member); err != nil {
		log.Error().Str("user_id", reply.ID.String()).Err(err).Msg("could not create default tenant and project for new user")
		// TODO: Does this leave the user in a bad state? Can they still use the app?
	}

	// Add to SendGrid Ensign Marketing list in go routine
	// TODO: use worker queue to limit number of go routines for tasks like this
	// TODO: test in live integration tests to make sure this works
	if s.conf.SendGrid.Enabled() {
		go func() {
			contact := &sendgrid.Contact{
				Email: params.Email,
			}
			contact.ParseName(params.Name)

			if err := s.sendgrid.AddContact(contact); err != nil {
				log.Warn().Err(err).Msg("could not add newly registered user to sendgrid ensign marketing list")
			}
		}()
	}

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
		log.Warn().Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse login request"))
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
		log.Error().Err(err).Msg("could not login user")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// TODO: Add user state checks and create required resources for first logins
	// (tenants, projects, etc.)

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(authCSRFLifetime)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set cookies on login reply")
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
		log.Warn().Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse refresh request"))
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
		log.Error().Err(err).Msg("could not refresh user access token")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(authCSRFLifetime)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set cookies on refresh reply")
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
		log.Error().Err(err).Msg("could not set cookies")
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
		log.Warn().Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse verify request"))
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
		log.Error().Err(err).Msg("could not verify email address")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Note: This obscures 202 Accepted responses as 200 OK responses which prevents
	// the user from being able to tell if they were already verified. To allow the
	// user to distinguish between the two we would have to return an error or modify
	// the response body to include that information.
	c.JSON(http.StatusOK, &api.Reply{Success: true})
}
