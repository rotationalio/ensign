package quarterdeck

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/gravatar"
	"github.com/rotationalio/ensign/pkg/utils/metrics"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

const (
	UserHuman   = "human"
	UserMachine = "machine"
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
// the user, although the field is optional. This endpoint does not handle adding users
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
		sentry.Warn(c).Err(err).Msg("could not parse register request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryRegisterAgain))
		return
	}

	if err = in.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// ProjectID is an optional field
	var projectID ulid.ULID
	if in.ProjectID != "" {
		if projectID, err = ulid.Parse(in.ProjectID); err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryRegisterAgain))
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
		sentry.Error(c).Err(err).Msg("could not create derived key for user password")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Create a new organization for the user
	org := &models.Organization{
		Name:   in.Organization,
		Domain: in.Domain,
	}

	// Create a verification token to send to the user
	if err = user.CreateVerificationToken(); err != nil {
		sentry.Error(c).Err(err).Msg("could not create verification token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Create the user, they are the owner in their own organization
	if err = user.Create(ctx, org, permissions.RoleOwner); err != nil {
		// Handle constraint errors
		var dberr *models.ConstraintError
		if errors.As(err, &dberr) {
			c.JSON(http.StatusConflict, api.ErrorResponse(api.ErrUserExists))
			return
		}

		sentry.Error(c).Err(err).Msg("could not create user in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Verification emails should happen asynchronously because sending emails can be
	// slow and waiting for SendGrid to send the email could cause the request to time
	// out even though the user was successfully created.
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		return s.SendVerificationEmail(user)
	}),
		tasks.WithRetries(3),
		tasks.WithBackoff(backoff.NewExponentialBackOff()),
		tasks.WithError(fmt.Errorf("could not send verification email to user %s", user.ID.String())),
	)

	// If a project ID is provided then link the user's organization to the project by
	// creating a database record. This allows a path for the client to create a
	// default project for new users without having to go through a separate,
	// authenticated request. See ProjectCreate for more details.
	if !ulids.IsZero(projectID) {
		// WARNING: Failure to save this record will create an inconsistent situation
		// where the user and organization were created but the project was not linked
		// to the organization.
		// TODO: ensure this is added to a transaction context somehow.
		op := &models.OrganizationProject{
			OrgID:     org.ID,
			ProjectID: projectID,
		}
		if err = op.Save(ctx); err != nil {
			// WARNING: Errors in saving the organization project are very bad!
			sentry.Fatal(c).Err(err).Msg("user and organization created but project not linked to the organization")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
			return
		}
	}

	// increment registered users (can happen in this function and via invite)
	metrics.Registered.WithLabelValues(ServiceName).Inc()

	// increment registered organizations (can only happen in this function)
	metrics.Organizations.WithLabelValues(ServiceName).Inc()

	// Prepare response to return to the registering user.
	out = &api.RegisterReply{
		ID:        user.ID,
		OrgID:     org.ID,
		Email:     user.Email,
		OrgName:   org.Name,
		OrgDomain: org.Domain,
		Message:   "Welcome to Ensign!",
		Role:      permissions.RoleOwner,
		Created:   user.Created,
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
		sentry.Warn(c).Err(err).Msg("could not parse login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "unparseable").Inc()
		return
	}

	if in.Email == "" || in.Password == "" {
		log.Debug().Msg("missing email or password from login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "missing credentials").Inc()
		return
	}

	// Only one of orgID or invite token is accepted
	if !ulids.IsZero(in.OrgID) && in.InviteToken != "" {
		log.Debug().Msg("both orgID and invite token provided")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "both orgID and invite token provided").Inc()
		return
	}

	// Retrieve the user by email (read-only transaction)
	// For the invite case, the "default" organization is used to retrieve the user
	if user, err = models.GetUserEmail(c.Request.Context(), in.Email, in.OrgID); err != nil {
		// handle user not found error with a 403.
		if errors.Is(err, models.ErrNotFound) {
			log.Debug().Msg("could not find user by email address")
			c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrTryLoginAgain))

			// increment failure count
			metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "user not found").Inc()
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve the user from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "could not get user").Inc()
		return
	}

	// User must be verified to log in.
	if !user.EmailVerified {
		log.Debug().Msg("user has not verified their email address")
		api.Unverified(c)

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "unverified email").Inc()
		return
	}

	// Check that the password supplied by the user is correct.
	if verified, err := passwd.VerifyDerivedKey(user.Password, in.Password); err != nil || !verified {
		log.Debug().Err(err).Msg("invalid login credentials")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrTryLoginAgain))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "invalid password").Inc()
		return
	}

	// If an invite token was provided, accept the invite. This method modifies the
	// data at the user pointer and handles the logging and error responses.
	if in.InviteToken != "" && s.acceptInvite(c, user, in.InviteToken) != nil {
		return
	}

	// Create the access and refresh tokens and return them to the user.
	var claims *tokens.Claims
	if claims, err = user.NewClaims(c.Request.Context()); err != nil {
		sentry.Error(c).Err(err).Msg("could not create claims for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "could not create claims for user").Inc()
		return
	}

	out = &api.LoginReply{
		LastLogin: user.LastLogin.String,
	}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		sentry.Error(c).Err(err).Msg("could not create access and refresh token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "jwt token error").Inc()
		return
	}

	// Update the users last login in a Go routine so it doesn't block
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		return user.UpdateLastLogin(ctx)
	}), tasks.WithError(fmt.Errorf("could not update last login timestamp for user %s", user.ID.String())))

	// increment active users (in grafana we will divide by 24 hrs to get daily active)
	metrics.Active.WithLabelValues(ServiceName, UserHuman).Inc()
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
		sentry.Warn(c).Err(err).Msg("could not parse authentication request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "unparseable").Inc()
		return
	}

	if in.ClientID == "" || in.ClientSecret == "" {
		log.Debug().Msg("missing client id or secret from authentication request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing credentials"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "missing client id or secret").Inc()
		return
	}

	// Retrieve the API key by the client ID (read-only transaction)
	if apikey, err = models.GetAPIKey(c.Request.Context(), in.ClientID); err != nil {
		// handle apikey not found with a 403.
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))

			// increment failure count
			metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "api key not found").Inc()
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve apikey from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete request"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "api key-client id mismatch").Inc()
		return
	}

	// Check that the client secret supplied by the user is correct.
	if verified, err := passwd.VerifyDerivedKey(apikey.Secret, in.ClientSecret); err != nil || !verified {
		log.Debug().Err(err).Msg("invalid api key credentials")
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "invalid client secret").Inc()
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
		sentry.Error(c).Err(err).Msg("could not get permissions from model")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "permissions not found").Inc()
		return
	}

	out = &api.LoginReply{
		LastLogin: apikey.LastUsed.String,
	}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		sentry.Error(c).Err(err).Msg("could not create access and refresh token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "jwt token error").Inc()
		return
	}

	// Update the api keys last authentication in a Go routine so it doesn't block.
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		return apikey.UpdateLastUsed(ctx)
	}), tasks.WithError(fmt.Errorf("could not update last seen timestamp for api key %s", apikey.ID.String())))

	// increment active users (in grafana we will divide by 24 hrs to get daily active)
	metrics.Active.WithLabelValues(ServiceName, UserMachine).Inc()
	c.JSON(http.StatusOK, out)
}

// Refresh re-authenticates users and api keys using a refresh token rather than
// requiring a username and password or API key credentials a second time and returns a
// new access and refresh token pair with the current credentials of the user. This
// endpoint is intended to facilitate long-running connections to ensign systems that
// last longer than the duration of an access token; e.g. long sessions on the Beacon UI
// or (especially) long running publishers and subscribers (machine users) that need to
// stay authenticated semi-permanently.
func (s *Server) Refresh(c *gin.Context) {
	var (
		err error
		in  *api.RefreshRequest
		out *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse refresh request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrLogBackIn))
		return
	}

	// Check to see if the refresh token is included in the request
	if in.RefreshToken == "" {
		log.Debug().Msg("missing refresh token from request request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrLogBackIn))
		return
	}

	// verify the refresh token
	claims, err := s.tokens.Verify(in.RefreshToken)
	if err != nil {
		sentry.Warn(c).Err(err).Msg("could not verify refresh token")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))
		return
	}

	// verify that the token is indeed a refresh token
	if !claims.VerifyAudience(s.tokens.RefreshAudience(), true) {
		sentry.Warn(c).Msg("token does not contain refresh audience")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))
		return
	}

	// Refresh using the organization in the request, otherwise use the user's
	// currently selected organization.
	var orgID any
	if !ulids.IsZero(in.OrgID) {
		orgID = in.OrgID
	} else {
		orgID = claims.OrgID
	}

	// Identify if the subject is a user or an apikey
	var subject models.SubjectType
	if subject, err = models.IdentifySubject(c, claims.Subject, orgID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			// It is possible that either the user or API key was deleted.
			sentry.Warn(c).Err(err).Msg("could not identify subject in organization")
			c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))
			return
		}

		sentry.Warn(c).Err(err).Msg("could not identify subject")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	var (
		userType      string
		refreshClaims *tokens.Claims
	)

	switch subject {
	case models.UserSubject:
		userType = UserHuman
		if refreshClaims, err = s.refreshUser(c, claims.Subject, orgID); err != nil {
			// Error logging is handled in the refreshUser method
			return
		}
	case models.APIKeySubject:
		userType = UserMachine
		if refreshClaims, err = s.refreshAPIKey(c, claims.Subject, orgID); err != nil {
			// Error loggin is handled in the refreshAPIKey method
			return
		}
	default:
		sentry.Warn(c).Uint8("subject", uint8(subject)).Msg("unhandled subject type")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Create a new access token/refresh token pair
	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(refreshClaims); err != nil {
		sentry.Error(c).Err(err).Msg("could not create access and refresh tokens")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Add the issued timestamp as the last login on the request to avoid DB lookups. It
	// is possible that the user/api key has logged in since the last time but that is
	// probably not valid information on a refresh.
	out.LastLogin = claims.IssuedAt.Format(time.RFC3339Nano)

	// increment active users (in grafana we will divide by 24 hrs to get daily active)
	metrics.Active.WithLabelValues(ServiceName, userType).Inc()
	c.JSON(http.StatusOK, out)
}

func (s *Server) refreshUser(c *gin.Context, userID, orgID any) (_ *tokens.Claims, err error) {
	// Get the user from the database using the ID
	var user *models.User
	if user, err = models.GetUser(c, userID, orgID); err != nil {
		switch {
		case errors.Is(err, models.ErrUserOrganization):
			// The user is trying to log into an organization they don't belong to.
			sentry.Warn(c).Err(err).Msg("user is trying to log into an organization they don't belong to")
			c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))
		case errors.Is(err, models.ErrNotFound):
			// Not Found should only occur if the user was deleted after the refresh
			// tokens were issued, causing the token to be invalid.
			sentry.Warn(c).Err(err).Msg("user/organization in refresh token not found")
			c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))
		default:
			sentry.Warn(c).Err(err).Msg("could not retrieve user from database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}

		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "could not get user").Inc()
		return nil, err
	}

	// Create a new claims object using the user retrieved from the database
	refreshClaims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.String(),
		},
		Name:    user.Name,
		Email:   user.Email,
		Picture: gravatar.New(user.Email, nil),
	}

	// Add the orgID to the claims
	var refreshOrg ulid.ULID
	if refreshOrg, err = user.OrgID(); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch orgID from user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))

		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "organization not found").Inc()
		return nil, err
	}
	refreshClaims.OrgID = refreshOrg.String()

	// Add the user permissions to the claims.
	// NOTE: these should have been fetched on the first query.
	if refreshClaims.Permissions, err = user.Permissions(c.Request.Context(), false); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch permissions from user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))

		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "permissions not found").Inc()
		return nil, err
	}

	// Update the users last login in a Go routine so it doesn't block
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		return user.UpdateLastLogin(ctx)
	}), tasks.WithError(fmt.Errorf("could not update last login timestamp for user %s", user.ID.String())))
	return refreshClaims, nil
}

func (s *Server) refreshAPIKey(c *gin.Context, keyIDs, orgIDs any) (_ *tokens.Claims, err error) {
	// Parse the keyID and orgID into ULIDs
	var keyID, orgID ulid.ULID
	if keyID, err = ulids.Parse(keyIDs); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse apikey subject claims")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))

		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "invalid apikey id").Inc()
		return nil, err
	}

	if orgID, err = ulids.Parse(orgIDs); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse apikey orgID into ulid")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))

		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "invalid organization").Inc()
		return nil, err
	}

	// Get the APIKey from the database using the ID
	var apikey *models.APIKey
	if apikey, err = models.RetrieveAPIKey(c, keyID); err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			// Not Found should only occur if the apikey was deleted after the refresh
			// tokens were issued, causing the token to be invalid.
			sentry.Warn(c).Err(err).Msg("apikey/organization in refresh token not found")
			c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))
		default:
			sentry.Warn(c).Err(err).Msg("could not retrieve apikey from database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}

		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "could not get apikey").Inc()
		return nil, err
	}

	// Ensure that the orgID specified matches the orgID on the APIKey
	// TODO: this should be in the RetrieveAPIKey method
	if apikey.OrgID != orgID {
		sentry.Warn(c).Msg("requested organization does not match apikey organization")
		c.JSON(http.StatusForbidden, api.ErrorResponse(responses.ErrLogBackIn))

		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "apikey not in organization").Inc()
		return nil, models.ErrInvalidOrganization
	}

	// Create a new refreshClaims object using the apikey retrieved from the database
	refreshClaims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: apikey.ID.String(),
		},
		OrgID:     apikey.OrgID.String(),
		ProjectID: apikey.ProjectID.String(),
	}

	// Add the key permissions to the claims.
	// NOTE: these should have been fetched on the first query and cached.
	if refreshClaims.Permissions, err = apikey.Permissions(c.Request.Context(), false); err != nil {
		sentry.Error(c).Err(err).Msg("could not get permissions from model")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))

		metrics.FailedLogins.WithLabelValues(ServiceName, UserMachine, "permissions not found").Inc()
		return nil, err
	}

	// Update the api keys last authentication in a Go routine so it doesn't block.
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		return apikey.UpdateLastUsed(ctx)
	}), tasks.WithError(fmt.Errorf("could not update last seen timestamp for api key %s", apikey.ID.String())))
	return refreshClaims, nil
}

// Switch re-authenticates users (human users only) using the access token the user
// posts in the headers in order to give the user new claims with a new organization ID.
// E.g. the user switches from being logged into one organization to being logged into
// another organization. The user must submit the orgID of the organization they wish
// to switch to and the user must belong to that organization otherwise an error is
// returned.
//
// NOTE: this endpoint cannot be used with api keys because api keys are only ever
// issued to one organization (and in fact, one project inside of one organization).
// Only human users can belong to multiple organizations.
func (s *Server) Switch(c *gin.Context) {
	var (
		err error
		in  *api.SwitchRequest
		out *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse refresh request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "unparseable").Inc()
		return
	}

	// Ensure the orgID is included in the request
	if ulids.IsZero(in.OrgID) {
		log.Debug().Msg("missing orgID in switch request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing organization id"))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "missing org_id").Inc()
		return
	}

	// Parse the claims from the access token in the request
	var claims *tokens.Claims
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "no access token").Inc()
		return
	}

	// If the orgID in the claims is the same as the requested orgID return an error
	// (the user must use the refresh endpoint to get claims for the same orgID)
	if orgID := claims.ParseOrgID(); orgID.Compare(in.OrgID) == 0 {
		sentry.Warn(c).Msg("user attempting to switch to the same organization")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot switch into the organization you are currently logged into"))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "switch to same org").Inc()
		return
	}

	// Fetch the user and the user's permissions from the database.
	// Ensure that the user is loaded in the supplied organization.
	var user *models.User
	if user, err = models.GetUser(c.Request.Context(), claims.Subject, in.OrgID); err != nil {
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, models.ErrUserOrganization) {
			sentry.Warn(c).Str("userID", claims.Subject).Str("orgID", in.OrgID.String()).Msg("user attempt to switch into organization they do not belong to")
			c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
			metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "invalid credentials").Inc()
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve organization user from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process switch request"))
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "unhandled error").Inc()
		return

	}

	// Create access and refresh tokens for new organization
	// NOTE: ensure that new claims are created and returned, not the old claims;
	// otherwise the user may receive incorrect permissions.
	newClaims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.String(),
		},
		Name:    user.Name,
		Email:   user.Email,
		Picture: gravatar.New(user.Email, nil),
		OrgID:   in.OrgID.String(),
	}

	// Add the user permissions to the claims
	if newClaims.Permissions, err = user.Permissions(c.Request.Context(), false); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user permissions from model")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "permissions not found").Inc()
		return
	}

	out = &api.LoginReply{
		LastLogin: user.LastLogin.String,
	}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(newClaims); err != nil {
		sentry.Error(c).Err(err).Msg("could not create access and refresh token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))

		// increment failure count
		metrics.FailedLogins.WithLabelValues(ServiceName, UserHuman, "jwt token error").Inc()
		return
	}

	// Update the user's last login in a Go routine so it doesn't block
	s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		return user.UpdateLastLogin(ctx)
	}), tasks.WithError(fmt.Errorf("could not update last login timestamp for user %s", user.ID.String())))

	// increment active users (in grafana we will divide by 24 hrs to get daily active)
	metrics.Active.WithLabelValues(ServiceName, UserHuman).Inc()
	c.JSON(http.StatusOK, out)
}

// VerifyEmail verifies a user's email address by validating the token in the request.
// This endpoint is intended to be called by frontend applications after the user has
// followed the link in the verification email. If the user is not verified and the
// token is valid then the user is logged in. If the user is already verified then a
// 204 response is returned.
func (s *Server) VerifyEmail(c *gin.Context) {
	var (
		req *api.VerifyRequest
		err error
	)

	// Get the token from the POST request
	if err = c.BindJSON(&req); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse verify email request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing token in request"))
		return
	}

	// Look up the user by the token
	var user *models.User
	if user, err = models.GetUserByToken(c, req.Token, req.OrgID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid token"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve user by email verification token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
		return
	}

	// Return 202 if user is already verified
	if user.EmailVerified {
		c.Status(http.StatusNoContent)
		return
	}

	// Construct the user token from the database fields
	token := &db.VerificationToken{
		Email: user.Email,
	}
	if token.ExpiresAt, err = user.GetVerificationExpires(); err != nil {
		sentry.Error(c).Err(err).Msg("could not get verification expiration")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
		return
	}

	// Verify the token with the stored secret
	if err = token.Verify(user.GetVerificationToken(), user.EmailVerificationSecret); err != nil {
		if errors.Is(err, db.ErrTokenExpired) {
			// If expired, create a new token for the user
			if err = user.CreateVerificationToken(); err != nil {
				sentry.Error(c).Err(err).Msg("could not create new email verification token")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
				return
			}

			if err = user.Save(c.Request.Context()); err != nil {
				sentry.Error(c).Err(err).Msg("could not save user")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
				return
			}

			// Send the new token to the user
			s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
				return s.SendVerificationEmail(user)
			}),
				tasks.WithRetries(3),
				tasks.WithBackoff(backoff.NewExponentialBackOff()),
				tasks.WithError(fmt.Errorf("could not send verification email to user %s", user.ID.String())),
			)

			c.JSON(http.StatusGone, api.ErrorResponse("token expired, a new verification token has been sent to the email associated with the account"))
			return
		}

		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid token"))
		return
	}

	// Mark user as verified so they can login
	user.EmailVerified = true
	if err = user.Save(c.Request.Context()); err != nil {
		sentry.Error(c).Err(err).Msg("could not save user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify email"))
		return
	}

	// increment verified users in prometheus
	metrics.Verified.WithLabelValues(ServiceName).Inc()

	// Issue claims to the user to log them in, this skips the password check so it
	// only happens the first time a user is verified.
	var claims *tokens.Claims
	if claims, err = user.NewClaims(c.Request.Context()); err != nil {
		sentry.Error(c).Err(err).Msg("could not create claims for user")
		c.Status(http.StatusNoContent)
		return
	}

	// Create a new access token/refresh token pair
	out := &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.CreateTokenPair(claims); err != nil {
		sentry.Error(c).Err(err).Msg("could not create access and refresh token")
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, out)
}

// ResendEmail accepts an email address via a POST request and always returns a 204
// response, no matter the input or result of the processing. This is to ensure that
// no secure information is leaked from this unauthenticated endpoint. If the email
// address belongs to a user who has not been verified, another verification email is
// sent. If the post request contains an orgID and the user is invited to that
// organization but hasn't accepted the invite, then the invite is resent.
func (s *Server) ResendEmail(c *gin.Context) {
	var (
		err error
		in  *api.ResendRequest
	)

	// If we cannot parse the request return a 400 error
	if err = c.BindJSON(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse resend email request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Email is required for this endpoint and must not be longer than 254 characters
	// NOTE: email length limit is 254 characters based on RFC 2821.
	in.Email = strings.TrimSpace(in.Email)
	if in.Email == "" || len(in.Email) > 254 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrInvalidField))
		return
	}

	// Any response after parsing the request should return a 204 No Content
	defer c.Status(http.StatusNoContent)

	var user *models.User
	if user, err = models.GetUserEmail(c, in.Email, in.OrgID); err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			sentry.Error(c).Err(err).Msg("could not retrieve user by email address")
		}
		return
	}

	// Resend the verification email if the user is not verified
	// NOTE: this will have the effect of invalidating any previous verification tokens.
	if !user.EmailVerified {
		// Create a new token for the user
		if err = user.CreateVerificationToken(); err != nil {
			sentry.Error(c).Err(err).Msg("could not create new email verification token")
			return
		}

		if err = user.Save(c); err != nil {
			sentry.Error(c).Err(err).Msg("could not save user")
			return
		}

		// Send the new token to the user
		s.tasks.QueueContext(sentry.CloneContext(c), tasks.TaskFunc(func(ctx context.Context) error {
			return s.SendVerificationEmail(user)
		}),
			tasks.WithRetries(3),
			tasks.WithBackoff(backoff.NewExponentialBackOff()),
			tasks.WithError(fmt.Errorf("could not send verification email to user %s", user.ID.String())),
		)
	}

	// TODO: implement resending invitations to specific organizations
	// If an organization is specified in the request, check if the user has been
	// invited to the organization and hasn't accepted the invitation yet.
	if !ulids.IsZero(in.OrgID) {
		sentry.Warn(c).Msg("invitation resend is not implemented yet")
	}
}
