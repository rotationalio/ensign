package quarterdeck

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
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
	if rows, err = tx.Query(`SELECT k.id, k.key_id, k.name, p.slug, k.created, k.modified FROM api_keys k JOIN projects p ON p.id=k.project_id;`); err != nil {
		log.Error().Err(err).Msg("could not list api keys")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch api keys"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		k := &api.APIKey{}
		if err = rows.Scan(&k.ID, &k.ClientID, &k.Name, &k.Project, &k.Created, &k.Modified); err != nil {
			log.Error().Err(err).Msg("could not scan api key")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch api keys"))
			return
		}
		out.APIKeys = append(out.APIKeys, k)
	}

	tx.Commit()
	c.JSON(http.StatusOK, out)
}

// TODO: document
// TODO: actually implement this resource endpoint
// HACK: this is just a quick hack to get us going: it creates an api key quickly
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

	// TODO: add key validation

	// Create client ID and secret
	// HACK: generate better API keys
	key.ClientID = genkey(12)
	key.ClientSecret = genkey(32)

	// Insert derived key into the database rather than storing the secret directly
	var derivedKey string
	if derivedKey, err = passwd.CreateDerivedKey(key.ClientSecret); err != nil {
		log.Error().Err(err).Msg("could not create derived key from client secret")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create api key"))
		return
	}

	// Insert key into database
	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: false}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create api key"))
		return
	}
	defer tx.Rollback()

	// TODO: better error handling
	var result sql.Result
	if result, err = tx.Exec(`INSERT INTO api_keys (key_id, secret, name, project_id, created_by, created, modified) VALUES ($1, $2, $3, (SELECT id FROM projects WHERE slug=$4), (SELECT id FROM users WHERE email=$5), datetime('now'), datetime('now'));`, key.ClientID, derivedKey, key.Name, key.Project, key.Owner); err != nil {
		log.Error().Err(err).Msg("could not insert secret into the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create api key"))
		return
	}

	// HACK: should not ignore error here
	keyID, _ := result.LastInsertId()
	if err = tx.QueryRow(`SELECT id, created, modified FROM api_keys WHERE id=$1;`, keyID).Scan(&key.ID, &key.Created, &key.Modified); err != nil {
		log.Error().Err(err).Msg("could not insert secret into the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create api key"))
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, key)
}

func genkey(l int) string {
	data := make([]byte, l)
	rand.Read(data)
	return strings.ToLower(strings.Trim(base32.StdEncoding.EncodeToString(data), "="))
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
