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
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

const HeaderUserAgent = "User-Agent"

// List the API Keys for organization of the authenticated user, optionally filtered by
// project ID. The list response returns a subset of the fields in the APIKey object,
// to get more information about the API Key use the Detail endpoint. This endpoint
// returns a paginated response, limited by a default page size of 100 if one is not
// specified by the user (and a maximum page size of 5000). If there is another page of
// APIKeys the NextPageToken field will be populated, which can be used to make a
// subsequent request for the next page. Note that the page size or the projectID filter
// should not be changed between requests and that the NextPageToken will expire after
// 24 hours and can no longer be used.
//
// NOTE: the APIKey Secret should never be returned from this endpoint!
func (s *Server) APIKeyList(c *gin.Context) {
	var (
		err                error
		orgID, projectID   ulid.ULID
		keys               []*models.APIKey
		nextPage, prevPage *pagination.Cursor
		claims             *tokens.Claims
		out                *api.APIKeyList
	)

	query := &api.APIPageQuery{}
	if err = c.BindQuery(query); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query"))
		return
	}

	if query.ProjectID != "" {
		if projectID, err = ulid.Parse(query.ProjectID); err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(api.InvalidField("project_id")))
			return
		}
	}

	if query.NextPageToken != "" {
		if prevPage, err = pagination.Parse(query.NextPageToken); err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	} else {
		prevPage = pagination.New("", "", int32(query.PageSize))
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	if keys, nextPage, err = models.ListAPIKeys(c.Request.Context(), orgID, projectID, prevPage); err != nil {
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

	// Prepare response
	out = &api.APIKeyList{
		APIKeys: make([]*api.APIKey, 0, len(keys)),
	}

	for _, key := range keys {
		apikey := &api.APIKey{
			ID:        key.ID,
			ClientID:  key.KeyID,
			Name:      key.Name,
			OrgID:     key.OrgID,
			ProjectID: key.ProjectID,
		}
		apikey.LastUsed, _ = key.GetLastUsed()
		out.APIKeys = append(out.APIKeys, apikey)
	}

	// If a next page token is available, add it to the response.
	if nextPage != nil {
		if out.NextPageToken, err = nextPage.NextPageToken(); err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}
	}

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
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
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

	// NOTE: the OrgID and UserID MUST come from the user claims not from user input.
	// NOTE: we expect that the subject of the claims is the userID.
	model.OrgID = claims.ParseOrgID()
	model.CreatedBy = claims.ParseUserID()

	if ulids.IsZero(model.OrgID) || ulids.IsZero(model.CreatedBy) {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("invalid user claims"))
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
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
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
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Create a thin model to update in the database
	model = &models.APIKey{
		ID:   key.ID,
		Name: key.Name,
	}

	if model.OrgID = claims.ParseOrgID(); ulids.IsZero(model.OrgID) {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
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
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
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

// APIKeyPermissions returns the API key permissions available to the user.
func (s *Server) APIKeyPermissions(c *gin.Context) {
	var (
		tx     *sql.Tx
		err    error
		rows   *sql.Rows
		claims *tokens.Claims
	)

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Fetch all eligible API key claims
	out := make([]string, 0, 7)
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: true}); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch apikey permissions"))
		return
	}
	defer tx.Rollback()

	if rows, err = tx.Query("SELECT name FROM permissions WHERE allow_api_keys=true"); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch apikey permissions"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var permission string
		if err = rows.Scan(&permission); err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch apikey permissions"))
			return
		}
		out = append(out, permission)
	}

	// Filter other permissions based on the user's claims.
	outf := make([]string, 0, len(out))
	for _, permission := range out {
		// TODO: we'll need a better way to identify permissions that both the user and the API key can have.
		if perms.InGroup(permission, perms.PrefixTopics) || perms.InGroup(permission, perms.PrefixMetrics) {
			if !claims.HasPermission(permission) {
				// Do not return this permission
				continue
			}
		}
		outf = append(outf, permission)
	}

	c.JSON(http.StatusOK, outf)
}
