package quarterdeck

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rs/zerolog/log"
)

// TODO: add documentation
// TODO: review and ensure the register methodology is what we want
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

	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: false}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}
	defer tx.Rollback()

	// TODO: handle organizations and invites (e.g. with role associate).

	// Create password derived key so that we're not storing raw passwords
	var password string
	if password, err = passwd.CreateDerivedKey(in.Password); err != nil {
		log.Error().Err(err).Msg("could not create password derived key")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	insert := `INSERT INTO users (name, email, password, created, modified) VALUES ($1, $2, $3, datetime('now'), datetime('now'));`
	if _, err = tx.Exec(insert, in.Name, in.Email, password); err != nil {
		// TODO: how to check for constraint violations (e.g. unique or not null)
		log.Warn().Err(err).Msg("could not create new user")
		c.JSON(http.StatusConflict, api.ErrorResponse(err))
		return
	}

	// Prepare response
	out = &api.RegisterReply{
		Message: "Welcome to Ensign!",
		Role:    "Member",
	}

	// TODO: can we collect a time.Time from the database instead of a string?
	fetch := `SELECT id, email, created FROM users WHERE email=$1;`
	if err = tx.QueryRow(fetch, in.Email).Scan(&out.ID, &out.Email, &out.Created); err != nil {
		log.Error().Err(err).Msg("could not fetch newly created user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process registration"))
		return
	}

	// Assign the user the member role by default
	// TODO: assign the role from the invite if one is available
	role := `INSERT INTO user_roles (user_id, role_id, created, modified) VALUES ($1, (SELECT id FROM roles WHERE name='Member'), datetime('now'), datetime('now'));`
	if _, err = tx.Exec(role, out.ID); err != nil {
		log.Warn().Err(err).Msg("Could not assign user default role")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not finalize registration"))
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, out)
}

// TODO: add documentation
// TODO: simplify the code in this handler
func (s *Server) Login(c *gin.Context) {
	var (
		err error
		in  *api.LoginRequest
		out *api.LoginReply
	)

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not parse login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	if in.Email == "" || in.Password == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing credentials"))
		return
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: false}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process login"))
		return
	}
	defer tx.Rollback()

	// Fetch the derived key from the database to perform authentication
	var derivedKey string
	lookup := `SELECT password FROM users WHERE email=$1;`
	if err = tx.QueryRow(lookup, in.Email).Scan(&derivedKey); err != nil {
		// TODO: more graceful handling of error codes and failures
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid login credentials"))
		return
	}

	if verified, err := passwd.VerifyDerivedKey(derivedKey, in.Password); err != nil || !verified {
		// TODO: more graceful handling of error and failures
		c.JSON(http.StatusForbidden, api.ErrorResponse("invalid login credentials"))
		return
	}

	// Create the access and refresh tokens and return them to the user.
	claims := &tokens.Claims{Permissions: make([]string, 0)}
	querya := `SELECT name, email FROM users WHERE email=$1;`
	if err = tx.QueryRow(querya, in.Email).Scan(&claims.Name, &claims.Email); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}

	// TODO: this is not correct and needs an organization mapping to work
	// HACK: this quick hack expects that there is only one user role
	var rows *sql.Rows
	queryb := `SELECT p.name FROM user_roles ur JOIN users u ON u.id=ur.user_id JOIN roles r ON r.id=ur.role_id JOIN role_permissions rp ON r.id=rp.role_id JOIN permissions p ON p.id=rp.permission_id WHERE u.email=$1;`
	if rows, err = tx.Query(queryb, in.Email); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create credentials"))
		return
	}
	defer rows.Close()

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

	// Update the users last login
	update := `UPDATE users SET last_login=datetime('now') WHERE email=$1;`
	if _, err = tx.Exec(update, in.Email); err != nil {
		log.Warn().Err(err).Msg("could not update user last_login")
	}

	tx.Commit()
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
