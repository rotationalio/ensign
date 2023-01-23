package quarterdeck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
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
		KeyID:     keygen.KeyID(),
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

	// Create an APIKey but store it as a derived key in the database
	secret := keygen.Secret()
	if model.Secret, err = passwd.CreateDerivedKey(secret); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create API Key"))
		return
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
		switch err.(type) {
		case *models.ValidationError:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create API Key"))
		}
		return
	}

	// Update the response to send to the user
	key = model.ToAPI(c.Request.Context())
	key.ClientSecret = secret
	c.JSON(http.StatusCreated, key)
}

// Retrieve an APIKey by its ID. Most fields of the APIKey object are read-only, though
// some components, such as the APIKey secret, are not returned at all even on detail.
// An APIKey is returned if the ID can be parsed, it is found in the database, and the
// user OrgID claims match the organization the APIKey is assigned to. Otherwise this
// endpoint will return a 404 Not Found error if it cannot correctly retrieve the key.
//
// NOTE: the APIKey Secret should never be returned from this endpoint!
func (s *Server) APIKeyDetail(c *gin.Context) {
	var (
		err    error
		kid    ulid.ULID
		model  *models.APIKey
		claims *tokens.Claims
	)

	// Retrieve ID component from the URL and parse it.
	if kid, err = ulid.Parse(c.Param("id")); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Attempt to retrieve thekey from the database
	if model, err = models.RetrieveAPIKey(c.Request.Context(), kid); err != nil {
		// Check if the error is a not found error.
		c.Error(err)
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		return
	}

	// Ensure that the orgID on the claims matches the orgID on the APIKey
	if claims.OrgID != model.OrgID.String() {
		log.Warn().Msg("attempt to fetch key from different organization")
		c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
		return
	}

	// Populate the response from the model
	// NOTE: the secret should not be populated in the response!
	c.JSON(http.StatusOK, model.ToAPI(c.Request.Context()))
}

// Update an APIKey to change it's description. Most fields on the APIKey object are
// read-only; in order to "change" fields such as permissions it is necessary to delete
// the key and create a new one. The APIKey is updated if the ID can be parsed, it is
// found in the database, and the user OrgID claims match the organization the APIKey
// is assigned to. Otherwise this endpoint will return a 404 Not Found error.
//
// NOTE: the APIKey Secret should never be returned from this endpoint!
func (s *Server) APIKeyUpdate(c *gin.Context) {
	var (
		err    error
		kid    ulid.ULID
		key    *api.APIKey
		model  *models.APIKey
		claims *tokens.Claims
	)

	// Retrieve ID component from the URL and parse it.
	if kid, err = ulid.Parse(c.Param("id")); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
		return
	}

	// Bind the API request to the API Key
	if err = c.BindJSON(&key); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	// Sanity check: the URL endpoint and the key ID on the model match.
	if !ulids.IsZero(key.ID) && key.ID.Compare(kid) != 0 {
		c.Error(api.ErrModelIDMismatch)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrModelIDMismatch))
		return
	}

	// Validate the request from the API side. The Database Model also has a validation,
	// but the API validation should ensure users are sending (or not sending) the
	// correct input, where database validation ensures the data is correctly being put
	// into the database and that programatic constraints are observed.
	key.ID = kid
	if err = key.ValidateUpdate(); err != nil {
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

	// Create a thin model to update in the database
	model = &models.APIKey{
		ID:   key.ID,
		Name: key.Name,
	}

	if model.OrgID, err = ulid.Parse(claims.OrgID); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Attempt to retrieve thekey from the database
	if err = model.Update(c.Request.Context()); err != nil {
		// Check if the error is a not found error or a validation error.
		var verr *models.ValidationError

		switch {
		case errors.Is(err, models.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
		case errors.As(err, &verr):
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verr))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		}

		c.Error(err)
		return
	}

	// Populate the response from the model
	// NOTE: the secret should not be populated in the response!
	c.JSON(http.StatusOK, model.ToAPI(c.Request.Context()))
}

// Delete an APIKey by its ID. This endpoint allows user to revoke APIKeys so that they
// can no longer be used for authentication with Quarterdeck. The APIKey is deleted if
// its ID can be parsed, it is found in the database, and the user OrgID claims match
// the organization the APIKey is assigned to. Otherwise this endpoint will return a
// 404 Not Found error if it cannot correctly retrieve the key. If the API Key is
// successfully deleted, this endpoint returns a 204 No Content response.
func (s *Server) APIKeyDelete(c *gin.Context) {
	var (
		err        error
		kid, orgID ulid.ULID
		claims     *tokens.Claims
	)

	// Retrieve ID component from the URL and parse it.
	if kid, err = ulid.Parse(c.Param("id")); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	if orgID, err = ulid.Parse(claims.OrgID); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Delete the APIKey in the specified organization
	if err = models.DeleteAPIKey(c.Request.Context(), kid, orgID); err != nil {
		// Check if the error is a not found error.
		c.Error(err)
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("api key not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		return
	}

	c.Status(http.StatusNoContent)
}
