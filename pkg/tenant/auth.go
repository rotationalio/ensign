package tenant

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
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

	// TODO: Handle error status codes returned by Quarterdeck
	var reply *qd.RegisterReply
	if reply, err = s.quarterdeck.Register(c.Request.Context(), req); err != nil {
		log.Error().Err(err).Msg("could not register user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete registration"))
		return
	}

	// TODO: Send verification email to the provided email address

	// Return the response from Quarterdeck
	c.JSON(http.StatusOK, &api.RegisterReply{
		Email:   reply.Email,
		Message: reply.Message,
		Role:    reply.Role,
	})
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

	// TODO: Handle error status codes returned by Quarterdeck
	var reply *qd.LoginReply
	if reply, err = s.quarterdeck.Login(c.Request.Context(), req); err != nil {
		log.Error().Err(err).Msg("could not login user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete login"))
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
