package quarterdeck

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rs/zerolog/log"
)

const (
	DefaultRole = "Member"
)

// Register creates a new user in the database with the specified password, allowing the
// user to login to Quarterdeck. This endpoint requires a "strong" password and a valid
// register request, otherwise a 400 reply is returned. The password is stored in the
// database as an argon2 derived key so it is impossible for a hacker to get access to
// raw passwords. By default the user is given the Member role, unless an organization
// is being created for the user, in which case the user is assigned the Owner role.
// TODO: add rate limiting to ensure that we don't get spammed with registrations
// TODO: review and ensure the register methodology is what we want
// TODO: handle organizations and invites (e.g. with role association).
func (s *Server) Register(c *gin.Context) {
	var (
		err error
		in  *api.RegisterRequest
		out *api.RegisterReply
	)

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse register request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if err = in.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Create a user model to insert into the database with the default role.
	// TODO: ensure role can be associated with the model directly.
	user := &models.User{
		Name:  in.Name,
		Email: in.Email,
	}

	// Create password derived key so that we're not storing raw passwords
	if user.Password, err = passwd.CreateDerivedKey(in.Password); err != nil {
		log.Error().Err(err).Msg("could not create password derived key")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	if err = user.Create(c.Request.Context(), DefaultRole); err != nil {
		// TODO: handle database constraint errors (e.g. unique email address)
		log.Error().Err(err).Msg("could not insert user into database during registration")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	// Prepare response to return to the registering user.
	out = &api.RegisterReply{
		ID:      user.ID.String(),
		Email:   user.Email,
		Message: "Welcome to Ensign!",
		Role:    "Member",
		Created: user.Created,
	}
	c.JSON(http.StatusCreated, out)
}

// Login is oriented towards human users who use their email and password for
// authentication (whereas authenticate is used for machine access using API keys).
// Login verifies the password submitted for the user is correct by looking up the user
// by email and using the argon2 derived key verification process to confirm the
// password matches. Upon authentication an access token and a refresh token with the
// authorized claims of the user (based on role) are returned. The user can use the
// access token to authenticate to Ensign systems and the claims within for
// authorization. The access token has an expiration and the refresh token can be used
// with the refresh endpoint to get a new access token without the user having to log in
// again. The refresh token overlaps with the access token to provide a
// seamless authentication experience and the user can refresh their access token so
// long as the refresh token is valid.
//
// This method primarily uses read queries (fetching the user from the database and
// fetching the user permissions from the database). It does update the user's last
// logged in timestamp in the database but should be highly available without
// Quarterdeck Raft replication in most cases.
// TODO: add rate limiting on a per-user basis to prevent Quarterdeck DOS.
func (s *Server) Login(c *gin.Context) {
	var (
		err  error
		user *models.User
		in   *api.LoginRequest
		out  *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if in.Email == "" || in.Password == "" {
		log.Debug().Msg("missing email or password from login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing credentials"))
		return
	}

	// Retrieve the user by email (read-only transaction)
	if user, err = models.GetUserEmail(c.Request.Context(), in.Email); err != nil {
		// TODO: handle user not found error with a 403.
		log.Error().Err(err).Msg("could not find user by email")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete request"))
		return
	}

	// Check that the password supplied by the user is correct.
	if verified, err := passwd.VerifyDerivedKey(user.Password, in.Password); err != nil || !verified {
		// TODO: more graceful handling of error and failures
		log.Debug().Err(err).Msg("invalid login credentials")
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid login credentials"))
		return
	}

	// Create the access and refresh tokens and return them to the user.
	// TODO: add organization ID and project ID to the claims
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.String(),
		},
		Name:  user.Name,
		Email: user.Email,
	}

	// Add the user permissions to the claims.
	// NOTE: these should have been fetched on the first query.
	if claims.Permissions, err = user.Permissions(c.Request.Context(), false); err != nil {
		log.Error().Err(err).Msg("could not fetch user permissions")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		log.Error().Err(err).Msg("could not create jwt tokens on login")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Update the users last login in a Go routine so it doesn't block
	// TODO: create a channel and workers to update last login to limit the num of go routines
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := user.UpdateLastLogin(ctx); err != nil {
			log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not update last login timestamp")
		}
	}()
	c.JSON(http.StatusOK, out)
}

// Authenticate is oriented to machine users that have an API key with a client ID and
// secret for authentication (whereas login is used for human access using an email and
// password). Authenticate verifies the client secret submitted is correct by looking
// up the api key by the key ID and using the argon2 derived key verification process
// to confirm the secret matches. Upon authentication, an access and refresh token with
// the authorized claims of the keys are returned. These tokens can be used to
// authenticate with ensign systems and the claims used for authorization. The access
// and refresh tokens work the same way the user tokens work and the refresh token can
// be used to fetch a new key pair without having to transmit a secret again.
//
// This method primarily uses read queries so should be highly available. The only write
// is the update of the last time the key was used, but it does this in a go routine to
// ensure that this endpoint is not blocked by Quarterdeck Raft replication.
// TODO: add rate limiting on a per-ip basis to prevent Quarterdeck DOS.
func (s *Server) Authenticate(c *gin.Context) {
	var (
		err    error
		apikey *models.APIKey
		in     *api.APIAuthentication
		out    *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse authenticate request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if in.ClientID == "" || in.ClientSecret == "" {
		log.Debug().Msg("missing client id or secret from authentication request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing credentials"))
		return
	}

	// Retrieve the API key by the client ID (read-only transaction)
	if apikey, err = models.GetAPIKey(c.Request.Context(), in.ClientID); err != nil {
		// TODO: handle apikey not found with a 404.
		log.Error().Err(err).Msg("could not find api key by client id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete request"))
		return
	}

	// Check that the client secret supplied by the user is correct.
	if verified, err := passwd.VerifyDerivedKey(apikey.Secret, in.ClientSecret); err != nil || !verified {
		// TODO: more graceful handling of error and failures
		log.Debug().Err(err).Msg("invalid api key credentials")
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
		return
	}

	// Create the access and refresh tokens and return them.
	// TODO: add the organization ID to the claims
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: apikey.ID.String(),
		},
		ProjectID: apikey.ProjectID.String(),
	}

	// Add the key permissions to the claims.
	// NOTE: these should have been fetched on the first query.
	if claims.Permissions, err = apikey.Permissions(c.Request.Context(), false); err != nil {
		log.Error().Err(err).Msg("could not fetch api key permissions")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		log.Error().Err(err).Msg("could not create jwt tokens on authenticate")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Update the api keys last authentication in a Go routine so it doesn't block.
	// TODO: create a channel and workers to update last seen to limit the num of go routines
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := apikey.UpdateLastUsed(ctx); err != nil {
			log.Error().Err(err).Str("api_key_id", apikey.ID.String()).Msg("could not update last seen timestamp")
		}
	}()
	c.JSON(http.StatusOK, out)
}

// Refresh is oriented to machine users that have a refresh token and need to obtain
// a new access token and refresh token to continue using ensign systems
func (s *Server) Refresh(c *gin.Context) {
	var (
		err error
		in  *api.RefreshRequest
		out *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse refresh request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	// Check to see if the refresh token is included in the request
	if in.RefreshToken == "" {
		log.Debug().Msg("missing refresh token from request request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing credentials"))
		return
	}

	// verify the refresh token
	claims, err := s.tokens.Verify(in.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("could not verify refresh token")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify refresh token"))
		return
	}

	// get the user from the database using the ID
	user, err := models.GetUser(c, claims.Subject)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve user from claims")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not retrieve user from claims"))
		return
	}

	// create a new access token/refresh token pair
	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		log.Error().Err(err).Msg("could not create jwt tokens on refresh")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Update the users last login in a Go routine so it doesn't block
	// TODO: create a channel and workers to update last login to limit the num of go routines
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := user.UpdateLastLogin(ctx); err != nil {
			log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not update last login timestamp")
		}
	}()
	c.JSON(http.StatusOK, out)
}
