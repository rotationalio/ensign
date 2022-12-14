package quarterdeck

import (
	"database/sql"
	"net/http"

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

// TODO: add documentation
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
