package quarterdeck

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/gravatar"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

// Register creates a new user in the database with the specified password, allowing the
// user to login to Quarterdeck. This endpoint requires a "strong" password and a valid
// register request, otherwise a 400 reply is returned. The password is stored in the
// database as an argon2 derived key so it is impossible for a hacker to get access to
// raw passwords.
//
// An organization is created for the user registering based on the organization data
// in the register request and the user is assigned the Owner role. A project ID can be
// provided in the request to allow the client to safely create a default project for
// the user, alhough the field is optional. This endpoint does not handle adding users
// to existing organizations through collaborator invites.
// TODO: add rate limiting to ensure that we don't get spammed with registrations
func (s *Server) Register(c *gin.Context) {
	var (
		err error
		in  *api.RegisterRequest
		out *api.RegisterReply
	)

	ctx := c.Request.Context()

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse register request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if err = in.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// ProjecID is an optional field
	var projectID ulid.ULID
	if in.ProjectID != "" {
		if projectID, err = ulid.Parse(in.ProjectID); err != nil {
			log.Error().Err(err).Str("project_id", in.ProjectID).Msg("could not parse project ID")
			c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project ID in request"))
			return
		}
	}

	// Create a user model to insert into the database.
	user := &models.User{
		Name:  in.Name,
		Email: in.Email,
	}

	// Set the user agreement fields
	user.SetAgreement(in.AgreeToS, in.AgreePrivacy)

	// Create password derived key so that we're not storing raw passwords
	if user.Password, err = passwd.CreateDerivedKey(in.Password); err != nil {
		log.Error().Err(err).Msg("could not create password derived key")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	// Create a new organization to associate with the user since this is not an invite.
	org := &models.Organization{
		Name:   in.Organization,
		Domain: in.Domain,
	}

	// Create a verification token to send to the user
	if err = user.CreateVerificationToken(); err != nil {
		log.Error().Err(err).Msg("could not create verification token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	if err = user.Create(ctx, org, permissions.RoleOwner); err != nil {
		// Handle constraint errors
		var dberr *models.ConstraintError
		if errors.As(err, &dberr) {
			c.JSON(http.StatusConflict, api.ErrorResponse(api.ErrUserExists))
			return
		}

		log.Error().Err(err).Msg("could not insert user into database during registration")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	// Verification emails should happen asynchronously because sending emails can be
	// slow and waiting for SendGrid to send the email could cause the request to time
	// out even though the user was successfully created.
	s.tasks.Queue(tasks.TaskFunc(func(ctx context.Context) {
		if err := s.SendVerificationEmail(user); err != nil {
			log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not send verification email to user")
		}
	}))

	// If a project ID is provided then link the user's organization to the project by
	// creating a database record. This allows a path for the client to create a
	// default project for new users without having to go through a separate,
	// authenticated request. See ProjectCreate for more details.
	if !ulids.IsZero(projectID) {
		// TODO: Failure to save this record will create an inconsistent situation
		// where the user and organization were created but the project was not linked
		// to the organization.
		op := &models.OrganizationProject{
			OrgID:     org.ID,
			ProjectID: projectID,
		}
		if err = op.Save(ctx); err != nil {
			log.Error().Err(err).Msg("could not link organization to project")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
			return
		}
	}

	// Prepare response to return to the registering user.
	out = &api.RegisterReply{
		ID:      user.ID,
		OrgID:   org.ID,
		Email:   user.Email,
		Message: "Welcome to Ensign!",
		Role:    permissions.RoleOwner,
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
	if user, err = models.GetUserEmail(c.Request.Context(), in.Email, in.OrgID); err != nil {
		// handle user not found error with a 403.
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusForbidden, api.ErrorResponse("invalid login credentials"))
			return
		}

		log.Error().Err(err).Msg("could not find user by email")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete request"))
		return
	}

	// Check that the password supplied by the user is correct.
	if verified, err := passwd.VerifyDerivedKey(user.Password, in.Password); err != nil || !verified {
		log.Debug().Err(err).Msg("invalid login credentials")
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid login credentials"))
		return
	}

	// User must be verified to log in.
	if !user.EmailVerified {
		log.Debug().Msg("user has not verified their email address")
		c.JSON(http.StatusForbidden, api.ErrorResponse("email address not verified"))
		return
	}

	// Create the access and refresh tokens and return them to the user.
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.String(),
		},
		Name:    user.Name,
		Email:   user.Email,
		Picture: gravatar.New(user.Email, nil),
	}

	// Add the orgID to the claims
	var orgID ulid.ULID
	if orgID, err = user.OrgID(); err != nil {
		log.Error().Err(err).Msg("could not get orgID from user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}
	claims.OrgID = orgID.String()

	// Add the user permissions to the claims.
	// NOTE: these should have been fetched on the first query.
	if claims.Permissions, err = user.Permissions(c.Request.Context(), false); err != nil {
		log.Error().Err(err).Msg("could not fetch user permissions")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	out = &api.LoginReply{
		LastLogin: user.LastLogin.String,
	}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		log.Error().Err(err).Msg("could not create jwt tokens on login")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Update the users last login in a Go routine so it doesn't block
	s.tasks.Queue(tasks.TaskFunc(func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := user.UpdateLastLogin(ctx); err != nil {
			log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not update last login timestamp")
		}
	}))
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
		// handle apikey not found with a 403.
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
			return
		}

		log.Error().Err(err).Msg("could not retrieve api key by client id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete request"))
		return
	}

	// Check that the client secret supplied by the user is correct.
	if verified, err := passwd.VerifyDerivedKey(apikey.Secret, in.ClientSecret); err != nil || !verified {
		log.Debug().Err(err).Msg("invalid api key credentials")
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
		return
	}

	// Create the access and refresh tokens and return them.
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: apikey.ID.String(),
		},
		OrgID:     apikey.OrgID.String(),
		ProjectID: apikey.ProjectID.String(),
	}

	// Add the key permissions to the claims.
	// NOTE: these should have been fetched on the first query and cached.
	if claims.Permissions, err = apikey.Permissions(c.Request.Context(), false); err != nil {
		log.Error().Err(err).Msg("could not fetch api key permissions")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	out = &api.LoginReply{
		LastLogin: apikey.LastUsed.String,
	}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		log.Error().Err(err).Msg("could not create jwt tokens on authenticate")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Update the api keys last authentication in a Go routine so it doesn't block.
	s.tasks.Queue(tasks.TaskFunc(func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := apikey.UpdateLastUsed(ctx); err != nil {
			log.Error().Err(err).Str("api_key_id", apikey.ID.String()).Msg("could not update last seen timestamp")
		}
	}))
	c.JSON(http.StatusOK, out)
}

// Refresh re-authenticates users and api keys using a refresh token rather than requiring a username
// and password or API key credentials a second time and returns a new access and refresh token pair
// with the current credentials of the user. This endpoint is intended to facilitate long-running
// connections to ensign systems that last longer than the duration of an access token; e.g. long
// sessions on the Beacon UI or (especially) long running publishers and subscribers (machine users)
// that need to stay authenticated semi-permanently.
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
		c.JSON(http.StatusForbidden, api.ErrorResponse("could not verify refresh token"))
		return
	}

	// verify that the token is indeed a refresh token
	if !claims.VerifyAudience(s.tokens.RefreshAudience(), true) {
		log.Debug().Msg("token does not contain refresh audience")
		c.JSON(http.StatusForbidden, api.ErrorResponse("could not verify refresh token"))
		return
	}

	// get the user from the database using the ID
	user, err := models.GetUser(c, claims.Subject, claims.OrgID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
			return
		}

		log.Error().Err(err).Msg("could not retrieve user from claims")
		c.JSON(http.StatusForbidden, api.ErrorResponse("could not retrieve user from claims"))
		return
	}

	// Create a new claims object using the user retrieved from the database
	// Create the access and refresh tokens and return them to the user.
	refreshClaims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.String(),
		},
		Name:    user.Name,
		Email:   user.Email,
		Picture: gravatar.New(user.Email, nil),
	}

	// Add the orgID to the claims
	var orgID ulid.ULID
	if orgID, err = user.OrgID(); err != nil {
		log.Error().Err(err).Msg("could not get orgID from user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}
	refreshClaims.OrgID = orgID.String()

	// Add the user permissions to the claims.
	// NOTE: these should have been fetched on the first query.
	if refreshClaims.Permissions, err = user.Permissions(c.Request.Context(), false); err != nil {
		log.Error().Err(err).Msg("could not fetch user permissions")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Create a new access token/refresh token pair
	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(refreshClaims); err != nil {
		log.Error().Err(err).Msg("could not create jwt tokens on refresh")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// Add the issued timestamp as the last login on the request to avoid DB lookups. It
	// is possible that the user/api key has logged in since the last time but that is
	// probably not valid information on a refresh.
	out.LastLogin = claims.IssuedAt.Format(time.RFC3339Nano)

	// Update the users last login in a Go routine so it doesn't block
	s.tasks.Queue(tasks.TaskFunc(func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := user.UpdateLastLogin(ctx); err != nil {
			log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not update last login timestamp")
		}
	}))
	c.JSON(http.StatusOK, out)
}

// VerifyEmail verifies a user's email address by validating the token in the request.
// This endpoint is intended to be called by frontend applications after the user has
// followed the link in the verification email. If the token is already verified this
// endpoint returns a 202 Accepted response.
func (s *Server) VerifyEmail(c *gin.Context) {
	var (
		req *api.VerifyRequest
		err error
	)

	// Get the token from the POST request
	if err = c.BindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("could not parse verify email request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing token in request"))
		return
	}

	// Look up the user by the token
	var user *models.User
	if user, err = models.GetUserByToken(c, req.Token); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid token"))
			return
		}

		log.Error().Err(err).Msg("could not lookup user by token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
		return
	}

	// Return 202 if user is already verified
	if user.EmailVerified {
		c.Status(http.StatusAccepted)
		return
	}

	// Construct the user token from the database fields
	token := &db.VerificationToken{
		Email: user.Email,
	}
	if token.ExpiresAt, err = user.GetVerificationExpires(); err != nil {
		log.Error().Err(err).Msg("could not get token expiration")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
		return
	}

	// Verify the token with the stored secret
	if err = token.Verify(user.GetVerificationToken(), user.EmailVerificationSecret); err != nil {
		if errors.Is(err, db.ErrTokenExpired) {
			// If expired, create a new token for the user
			if err = user.CreateVerificationToken(); err != nil {
				log.Error().Err(err).Msg("could not create new verification token to replace expired token")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
				return
			}

			if err = user.Save(c.Request.Context()); err != nil {
				log.Error().Err(err).Msg("could not save user with new verification token")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
				return
			}

			// Send the new token to the user
			s.tasks.Queue(tasks.TaskFunc(func(ctx context.Context) {
				if err := s.SendVerificationEmail(user); err != nil {
					log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not send verification email to user")
				}
			}))

			c.JSON(http.StatusGone, api.ErrorResponse("token expired, a new verification token has been sent to the email associated with the account"))
			return
		}

		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid token"))
		return
	}

	// Mark user as verified so they can login
	user.EmailVerified = true
	if err = user.Save(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("could not save user with verified email")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
		return
	}

	c.Status(http.StatusNoContent)
}
