package tenant

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/rs/zerolog/log"
)

// ProjectAPIKeyList lists API keys in the specified project by forwarding the request
// to Quarterdeck.
//
// Route: GET /v1/projects/:projectID/apikeys
func (s *Server) ProjectAPIKeyList(c *gin.Context) {
	var (
		claims *tokens.Claims
		ctx    context.Context
		err    error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// User claims are required to check ownership of the project
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not get user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch claims for authenticated user"))
		return
	}

	// Parse the params from the GET request
	params := &api.PageQuery{}
	if err = c.ShouldBindQuery(params); err != nil {
		log.Warn().Err(err).Msg("could not parse query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query params"))
		return
	}

	// Parse the project ID from the URL
	paramID := c.Param("projectID")
	var projectID ulid.ULID
	if projectID, err = ulids.Parse(paramID); err != nil {
		log.Warn().Str("id", paramID).Err(err).Msg("could not parse project id")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project id"))
		return
	}

	// Retrieve the project from the database
	var project *db.Project
	if project, err = db.RetrieveProject(ctx, projectID); err != nil {
		log.Error().Str("id", paramID).Err(err).Msg("could not retrieve project from database")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// User should not be able to list API keys in another organization
	if claims.OrgID != project.OrgID.String() {
		log.Warn().Str("user_org", claims.OrgID).Str("project_org", project.OrgID.String()).Msg("user cannot list API keys in this project")
		c.JSON(http.StatusForbidden, api.ErrorResponse("user is not authorized to access this project"))
		return
	}

	// Build the Quarterdeck request from the params
	req := &qd.APIPageQuery{
		ProjectID:     paramID,
		PageSize:      int(params.PageSize),
		NextPageToken: params.NextPageToken,
	}

	// Request a page of API keys from Quarterdeck
	var reply *qd.APIKeyList
	if reply, err = s.quarterdeck.APIKeyList(ctx, req); err != nil {
		log.Error().Err(err).Msg("could not list API keys")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not list API keys"))
		return
	}

	// Return the page of API keys
	out := &api.ProjectAPIKeyPage{
		ProjectID:     req.ProjectID,
		PrevPageToken: req.NextPageToken,
		NextPageToken: reply.NextPageToken,
		APIKeys:       make([]*api.APIKey, 0),
	}
	for _, key := range reply.APIKeys {
		out.APIKeys = append(out.APIKeys, &api.APIKey{
			ID:       key.ID.String(),
			ClientID: key.ClientID,
			Name:     key.Name,
		})
	}

	c.JSON(http.StatusOK, out)
}

// ProjectAPIKeyCreate creates a new API key in a project by forwarding the request to
// Quarterdeck.
//
// Route: POST /v1/projects/:projectID/apikeys
func (s *Server) ProjectAPIKeyCreate(c *gin.Context) {
	var (
		claims *tokens.Claims
		ctx    context.Context
		err    error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// User claims are required to check ownership of the project
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not get user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch claims for authenticated user"))
		return
	}

	// Parse the params from the POST request
	params := &api.APIKey{}
	if err = c.BindJSON(params); err != nil {
		log.Warn().Err(err).Msg("could not parse API key params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse API key params"))
		return
	}

	// Name is required
	if params.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key name is required"))
		return
	}

	// Permissions are required
	if len(params.Permissions) == 0 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key permissions are required"))
		return
	}

	// Build the Quarterdeck request
	// See ValidateCreate() for required fields
	req := &qd.APIKey{
		Name:        params.Name,
		Permissions: params.Permissions,
	}

	// ProjectID is required
	projectID := c.Param("projectID")
	if req.ProjectID, err = ulid.Parse(projectID); err != nil {
		log.Warn().Err(err).Str("projectID", projectID).Msg("could not parse project ID")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid project ID"))
		return
	}

	// Retrieve the Project from the database
	var project *db.Project
	if project, err = db.RetrieveProject(ctx, req.ProjectID); err != nil {
		log.Error().Err(err).Str("projectID", projectID).Msg("could not retrieve project from database")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// User should not be able to create API keys in another organization
	if claims.OrgID != project.OrgID.String() {
		log.Warn().Str("user_org", claims.OrgID).Str("project_org", project.OrgID.String()).Msg("user cannot create API keys in this project")
		c.JSON(http.StatusForbidden, api.ErrorResponse("user is not authorized to access this project"))
		return
	}

	// TODO: Add source to request

	// Create the API key with Quarterdeck
	var key *qd.APIKey
	if key, err = s.quarterdeck.APIKeyCreate(ctx, req); err != nil {
		log.Error().Err(err).Msg("could not create API key")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not create API key"))
		return
	}

	// Return the API key
	out := &api.APIKey{
		ID:           key.ID.String(),
		ClientID:     key.ClientID,
		ClientSecret: key.ClientSecret,
		Name:         key.Name,
		Owner:        key.CreatedBy.String(),
		Permissions:  key.Permissions,
		Created:      key.Created.Format(time.RFC3339Nano),
		Modified:     key.Modified.Format(time.RFC3339Nano),
	}

	c.JSON(http.StatusCreated, out)
}

// TODO: Implement by factoring out common code from ProjectAPIKeyCreate
func (s *Server) APIKeyList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TODO: Implement by factoring out common code from ProjectAPIKeyCreate
func (s *Server) APIKeyCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// APIKeyDetail returns details about a specific API key.
//
// Route: GET /v1/apikeys/:apiKeyID
func (s *Server) APIKeyDetail(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// Parse the API key ID from the URL
	apiKeyID := c.Param("apiKeyID")

	// Get the API key from Quarterdeck
	var key *qd.APIKey
	if key, err = s.quarterdeck.APIKeyDetail(ctx, apiKeyID); err != nil {
		log.Error().Err(err).Str("apiKeyID", apiKeyID).Msg("could not get API key")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not retrieve API key"))
		return
	}

	// Return everything but the client secret
	out := &api.APIKey{
		ID:          key.ID.String(),
		ClientID:    key.ClientID,
		Name:        key.Name,
		Owner:       key.CreatedBy.String(),
		Permissions: key.Permissions,
		Created:     key.Created.Format(time.RFC3339Nano),
		Modified:    key.Modified.Format(time.RFC3339Nano),
	}

	c.JSON(http.StatusOK, out)
}

// APIKeyUpdate updates an API key by forwarding the request to Quarterdeck.
//
// Route: PUT /v1/apikeys/:apiKeyID
func (s *Server) APIKeyUpdate(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// Parse the API key ID from the URL
	var id ulid.ULID
	apiKeyID := c.Param("apiKeyID")
	if id, err = ulid.Parse(apiKeyID); err != nil {
		log.Warn().Err(err).Str("apiKeyID", apiKeyID).Msg("could not parse API key ID")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse API key ID from URL"))
		return
	}

	// Parse the request body
	params := &api.APIKey{}
	if err = c.BindJSON(params); err != nil {
		log.Warn().Err(err).Msg("could not parse API key update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse API key update request"))
		return
	}

	// ID should also be in the request body
	var paramsID ulid.ULID
	if paramsID, err = ulid.Parse(params.ID); err != nil {
		log.Warn().Err(err).Str("paramsID", params.ID).Msg("could not parse API key ID from request body")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse API key ID from request body"))
		return
	}

	// Sanity check that the ID in the URL matches the ID in the request
	if !ulids.IsZero(paramsID) && id.Compare(paramsID) != 0 {
		log.Warn().Err(err).Str("apiKeyID", apiKeyID).Str("paramsID", params.ID).Msg("API key ID in URL does not match ID in request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key ID does not match key ID in request"))
		return
	}

	// Name is required
	if params.Name == "" {
		log.Warn().Str("apiKeyID", apiKeyID).Msg("API key name is required")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key name is required for update"))
		return
	}

	// Build the Quarterdeck request
	// See ValidateUpdate() for required and restricted fields
	req := &qd.APIKey{
		ID:   id,
		Name: params.Name,
	}

	// Update the API key with Quarterdeck
	var key *qd.APIKey
	if key, err = s.quarterdeck.APIKeyUpdate(ctx, req); err != nil {
		log.Error().Err(err).Str("apiKeyID", apiKeyID).Msg("could not update API key")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not update API key"))
		return
	}

	// Return the updated key
	out := &api.APIKey{
		ID:          key.ID.String(),
		ClientID:    key.ClientID,
		Name:        key.Name,
		Owner:       key.CreatedBy.String(),
		Permissions: key.Permissions,
		Created:     key.Created.Format(time.RFC3339Nano),
		Modified:    key.Modified.Format(time.RFC3339Nano),
	}
	c.JSON(http.StatusOK, out)
}

// APIKeyDelete deletes an API key by forwarding the request to Quarterdeck.
//
// Route: DELETE /v1/apikeys/:apiKeyID
func (s *Server) APIKeyDelete(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// Parse the API key ID from the URL
	apiKeyID := c.Param("apiKeyID")

	// Delete the API key using Quarterdeck
	if err = s.quarterdeck.APIKeyDelete(ctx, apiKeyID); err != nil {
		log.Error().Err(err).Str("apiKeyID", apiKeyID).Msg("could not delete API key")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not delete API key"))
		return
	}

	c.Status(http.StatusNoContent)
}
