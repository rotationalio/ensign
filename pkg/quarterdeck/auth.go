package quarterdeck

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
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
// again. The refresh token not before overlaps with the access token to provide a
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
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := user.UpdateLastLogin(ctx); err != nil {
			log.Error().Err(err).Str("user_id", user.ID.String()).Msg("could not update last login timestamp")
		}
	}()
	c.JSON(http.StatusOK, out)
}

// TODO: add documentation
// TODO: simplify the code in this handler
func (s *Server) Authenticate(c *gin.Context) {
	var (
		err error
		in  *api.APIAuthentication
		out *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse authenticate request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if in.ClientID == "" || in.ClientSecret == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing credentials"))
		return
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: true}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process authentication"))
		return
	}
	defer tx.Rollback()

	// Fetch the derived key from the database to perform authentication
	var (
		keyID      int64
		derivedKey string
	)
	lookup := `SELECT id, secret FROM api_keys WHERE key_id=$1;`
	if err = tx.QueryRow(lookup, in.ClientID).Scan(&keyID, &derivedKey); err != nil {
		// TODO: more graceful handling of error codes and failures
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
		return
	}

	if verified, err := passwd.VerifyDerivedKey(derivedKey, in.ClientSecret); err != nil || !verified {
		// TODO: more graceful handling of error and failures
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid credentials"))
		return
	}

	// Create the access and refresh tokens and return them to the user.
	var rows *sql.Rows
	query := `SELECT p.name FROM api_key_permissions ak JOIN permissions p ON p.id=ak.permission_id WHERE ak.api_key_id=$1`
	if rows, err = tx.Query(query, keyID); err != nil {
		log.Error().Err(err).Msg("could not fetch permissions")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process authentication"))
		return
	}
	defer rows.Close()

	claims := &tokens.Claims{Permissions: make([]string, 0)}
	claims.RegisteredClaims.Subject = in.ClientID
	for rows.Next() {
		var permission string
		if err = rows.Scan(&permission); err != nil {
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
			return
		}
		claims.Permissions = append(claims.Permissions, permission)
	}

	var atk, rtk *jwt.Token
	if atk, err = s.tokens.CreateAccessToken(claims); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	if rtk, err = s.tokens.CreateRefreshToken(atk); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	out = &api.LoginReply{}
	if out.AccessToken, err = s.tokens.Sign(atk); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}
	if out.RefreshToken, err = s.tokens.Sign(rtk); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, out)
}

func (s *Server) Refresh(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}
