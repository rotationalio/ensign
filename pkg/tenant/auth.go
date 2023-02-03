package tenant

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// Set the maximum age of login protection cookies.
const doubleCookiesMaxAge = time.Minute * 10

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

	// Validate that required fields were provided
	if params.Name == "" || params.Email == "" || params.Password == "" || params.PwCheck == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing required fields for registration"))
		return
	}

	// Simple validation of the provided password
	// Note: Quarterdeck also checks this along with password strength, but this allows
	// us to filter some bad requests before they reach Quarterdeck.
	if params.Password != params.PwCheck {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("passwords do not match"))
		return
	}

	// Make the register request to Quarterdeck
	req := &qd.RegisterRequest{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
		PwCheck:  params.PwCheck,
	}

	var reply *qd.RegisterReply
	if reply, err = s.quarterdeck.Register(ctx, req); err != nil {
		log.Error().Err(err).Msg("could not register user")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not complete registration"))
		return
	}

	// TODO: Send verification email to the provided email address

	// Create a partial member record for the new user
	// Note: Tenant ID is not populated because it hasn't been created yet
	member := &db.Member{
		OrgID: reply.OrgID,
		Name:  req.Name,
		Role:  reply.Role,
	}

	// Create a default tenant and project for the new user
	// Note: This method returns an error if the member model is invalid
	if err = db.CreateUserResources(ctx, member); err != nil {
		log.Error().Str("user_id", reply.ID.String()).Err(err).Msg("could not create default tenant and project for new user")
		// TODO: Does this leave the user in a bad state? Can they still use the app?
	}

	// Add to SendGrid Ensign Marketing list in go routine
	// TODO: use worker queue to limit number of go routines for tasks like this
	// TODO: test in live integration tests to make sure this works
	if s.conf.SendGrid.Enabled() {
		go func() {
			name := strings.Split(params.Name, "")
			contact := &sgContact{
				Email:     params.Email,
				FirstName: name[0],
			}
			if len(name) > 1 {
				contact.LastName = strings.Join(name[1:], " ")
			}

			if err := s.AddContactToSendGrid(contact); err != nil {
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
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not complete login"))
		return
	}

	// TODO: Add user state checks and create required resources for first logins
	// (tenants, projects, etc.)

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(doubleCookiesMaxAge)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set cookies on login reply")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set cookies"))
		return
	}

	// Return the access and refresh tokens from Quarterdeck
	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}
	c.JSON(http.StatusOK, out)
}

// Refresh is a publically accessible endpoint that allows users to refresh their
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
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not complete refresh"))
		return
	}

	// Return the access and refresh tokens from Quarterdeck
	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}
	c.JSON(http.StatusOK, out)
}

// ProtectLogin prepares the front-end for login by setting the double cookie
// tokens for CSRF protection.
func (s *Server) ProtectLogin(c *gin.Context) {
	expiresAt := time.Now().Add(doubleCookiesMaxAge)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set cookies"))
		return
	}
	c.JSON(http.StatusOK, &api.Reply{Success: true})
}
