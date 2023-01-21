package quarterdeck

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rs/zerolog/log"
)

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
		err error
		key *api.APIKey
	)

	if err = c.BindJSON(&key); err != nil {
		log.Warn().Err(err).Msg("could not parse create api key request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	// TODO: get user claims

	// TODO: add key validation

	// Create the API Key database model and generate key material.
	model := &models.APIKey{
		Name:      key.Name,
		ProjectID: key.ProjectID,
	}

	// Add permissions to the database model
	if err = model.SetPermissions(key.Permissions...); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}

	if err = model.Create(c.Request.Context()); err != nil {
		// TODO: handle constraint violation errors with a 400
		c.Error(err)
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create API Key"))
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
