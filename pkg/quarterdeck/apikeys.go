package quarterdeck

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rs/zerolog/log"
)

const HeaderUserAgent = "User-Agent"

// TODO: document
// TODO: actually implement this resource endpoint
// TODO: implement pagination
// HACK: this is just a quick hack to get us going, it should filter api keys based on
// the authenticated user and organization instead of just returning everyting.
func (s *Server) APIKeyList(c *gin.Context) {
	var (
		err  error
		rows *sql.Rows
		out  *api.APIKeyList
	)

	// Fetch the api keys from the database
	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: true}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch api keys"))
		return
	}
	defer tx.Rollback()

	out = &api.APIKeyList{APIKeys: make([]*api.APIKey, 0)}
	if rows, err = tx.Query(`SELECT k.id, k.key_id, k.name, k.project_id, k.created, k.modified FROM api_keys k`); err != nil {
		log.Error().Err(err).Msg("could not list api keys")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch api keys"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		k := &api.APIKey{}
		if err = rows.Scan(&k.ID, &k.ClientID, &k.Name, &k.ProjectID, &k.Created, &k.Modified); err != nil {
			log.Error().Err(err).Msg("could not scan api key")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch api keys"))
			return
		}
		out.APIKeys = append(out.APIKeys, k)
	}

	tx.Commit()
	c.JSON(http.StatusOK, out)
}

// Create an API Key for the specified project with the specified permissions. Most of
// the fields on an APIKey cannot be updated (with the exception of the API Key name).
// This method is the only way a user can set a keys projectID, createdBy, source, and
// permissions fields. All other fields are managed by Quarterdeck.
//
// NOTE: a response to this request is the only time the key secret is exposed publicly.
// The secret is stored as an argon2 derived key so it is impossible for Quarterdeck to
// return the key to the user at any point after this method is called. The client must
// be responsible for recording the key and warning the user that this is the one time
// that it will be displayed. If the user loses the key, they will have to revoke
// (delete) the key and generate a new one.
func (s *Server) APIKeyCreate(c *gin.Context) {
	var (
		err    error
		key    *api.APIKey
		claims *tokens.Claims
	)

	// Bind the API request to the API Key
	if err = c.BindJSON(&key); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	// Validate the request from the API side. The Database Model also has a validation,
	// but the API validation should ensure users are sending (or not sending) the
	// correct input, where database validation ensures the data is correctly being put
	// into the database and that programatic constraints are observed.
	if err = key.ValidateCreate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Fetch the user-agent header from the request
	userAgent := c.GetHeader(HeaderUserAgent)

	// Create the API Key database model and generate key material.
	model := &models.APIKey{
		Name:      key.Name,
		ProjectID: key.ProjectID,
		Source: sql.NullString{
			Valid:  key.Source != "",
			String: key.Source,
		},
		UserAgent: sql.NullString{
			Valid:  userAgent != "",
			String: userAgent,
		},
	}

	if model.OrgID, err = ulid.Parse(claims.OrgID); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid user claims"))
		return
	}

	// NOTE: we expect that the subject of the claims is the userID.
	if model.CreatedBy, err = ulid.Parse(claims.Subject); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid user claims"))
		return
	}

	// Add permissions to the database model
	if err = model.SetPermissions(key.Permissions...); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	if err = model.Create(c.Request.Context()); err != nil {
		c.Error(err)
		// TODO: handle constraint violation errors with a 400
		switch err.(type) {
		case *models.ValidationError:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create API Key"))
		}
		return
	}

	// Update the response to send to the user
	key.ID = model.ID
	key.ClientID = model.KeyID
	key.ClientSecret = model.Secret
	key.Name = model.Name
	key.OrgID = model.OrgID
	key.ProjectID = model.ProjectID
	key.CreatedBy = model.CreatedBy
	key.Source = model.Source.String
	key.UserAgent = model.UserAgent.String
	key.LastUsed, _ = model.GetLastUsed()
	key.Permissions, _ = model.Permissions(c.Request.Context(), false)
	key.Created, _ = model.GetCreated()
	key.Modified, _ = model.GetModified()
	c.JSON(http.StatusCreated, key)
}

func (s *Server) APIKeyDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) APIKeyUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) APIKeyDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}
