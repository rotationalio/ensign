package tenant

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/radish"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

// ProjectAPIKeyList lists API keys in the specified project by forwarding the request
// to Quarterdeck.
//
// Route: GET /v1/projects/:projectID/apikeys
func (s *Server) ProjectAPIKeyList(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to check ownership of the project
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Parse the query parameters
	query := &api.PageQuery{}
	if err = c.ShouldBindQuery(query); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse page query request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Parse the project ID from the URL
	paramID := c.Param("projectID")
	var projectID ulid.ULID
	if projectID, err = ulids.Parse(paramID); err != nil {
		sentry.Warn(c).Err(err).Str("project_id", paramID).Msg("could not parse project id from query string")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project id"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Build the Quarterdeck request from the params
	req := &qd.APIPageQuery{
		ProjectID:     projectID.String(),
		PageSize:      int(query.PageSize),
		NextPageToken: query.NextPageToken,
	}

	// Request a page of API keys from Quarterdeck
	var reply *qd.APIKeyList
	if reply, err = s.quarterdeck.APIKeyList(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Return the page of API keys
	out := &api.ProjectAPIKeyPage{
		ProjectID:     req.ProjectID,
		PrevPageToken: req.NextPageToken,
		NextPageToken: reply.NextPageToken,
		APIKeys:       make([]*api.APIKeyPreview, 0),
	}
	for _, key := range reply.APIKeys {
		preview := &api.APIKeyPreview{
			ID:       key.ID.String(),
			ClientID: key.ClientID,
			Name:     key.Name,
			Status:   key.Status,
			LastUsed: db.TimeToString(key.LastUsed),
			Created:  db.TimeToString(key.Created),
			Modified: db.TimeToString(key.Modified),
		}

		// Return partial if permissions are missing, otherwise return full
		if key.Partial {
			preview.Permissions = api.PartialPermissions
		} else {
			preview.Permissions = api.FullPermissions
		}

		out.APIKeys = append(out.APIKeys, preview)
	}

	c.JSON(http.StatusOK, out)
}

// ProjectAPIKeyCreate creates a new API key in a project by forwarding the request to
// Quarterdeck.
//
// Route: POST /v1/projects/:projectID/apikeys
func (s *Server) ProjectAPIKeyCreate(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user credentials from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// User claims are required to validate the API key permissions
	var claims *tokens.Claims
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to check ownership of the project
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the user ID from the context
	var userID ulid.ULID
	if userID = userIDFromContext(c); ulids.IsZero(userID) {
		return
	}

	// Parse the params from the POST request
	params := &api.APIKey{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse create apikey request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Name is required
	if params.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key name is required"))
		return
	}

	// Permissions are required
	if len(params.Permissions) == 0 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key permissions are required."))
		return
	}

	// Ensure that the user can't create an API key with permissions they don't have
	// This is a lightweight check to prevent bad requests from going to Quarterdeck
	for _, permission := range params.Permissions {
		if perms.UserKeyPermission(permission) && !claims.HasPermission(permission) {
			// Do not allow the user to create an apikey with this permission
			sentry.Warn(c).Msg("user tried to create an API key with permissions they don't have")
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid permissions requested for API key"))
			return
		}
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
		sentry.Warn(c).Err(err).Str("projectID", projectID).Msg("could not parse project ID")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid project ID"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, req.ProjectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// TODO: Add source to request

	// Create the API key with Quarterdeck
	var key *qd.APIKey
	if key, err = s.quarterdeck.APIKeyCreate(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
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
	}

	// Update project stats in the background
	s.tasks.QueueContext(middleware.TaskContext(c), radish.Func(func(ctx context.Context) error {
		return s.UpdateProjectStats(ctx, userID, key.ProjectID)
	}), radish.WithErrorf("could not update stats for project %s", key.ProjectID.String()))

	c.JSON(http.StatusCreated, out)
}

// TODO: Implement by factoring out common code from ProjectAPIKeyCreate
func (s *Server) APIKeyList(c *gin.Context) {
	sentry.Warn(c).Msg("apikey list not implemented yet")
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TODO: Implement by factoring out common code from ProjectAPIKeyCreate
func (s *Server) APIKeyCreate(c *gin.Context) {
	sentry.Warn(c).Msg("apikey create not implemented yet")
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Parse the API key ID from the URL
	apiKeyID := c.Param("apiKeyID")

	// Get the API key from Quarterdeck
	var key *qd.APIKey
	if key, err = s.quarterdeck.APIKeyDetail(ctx, apiKeyID); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Parse the API key ID from the URL
	var id ulid.ULID
	apiKeyID := c.Param("apiKeyID")
	if id, err = ulid.Parse(apiKeyID); err != nil {
		sentry.Warn(c).Err(err).Str("apikeyID", apiKeyID).Msg("could not parse apikey id from url")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse API key ID from URL"))
		return
	}

	// Parse the request body
	params := &api.APIKey{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse apikey update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// ID should also be in the request body
	var paramsID ulid.ULID
	if paramsID, err = ulid.Parse(params.ID); err != nil {
		sentry.Warn(c).Err(err).Str("apikeyID", params.ID).Msg("could not parse apikey id from params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse API key ID from request body"))
		return
	}

	// Sanity check that the ID in the URL matches the ID in the request
	if !ulids.IsZero(paramsID) && id.Compare(paramsID) != 0 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("API key ID does not match key ID in request"))
		return
	}

	// Name is required
	if params.Name == "" {
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
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Get the user ID from the request context
	var userID ulid.ULID
	if userID = userIDFromContext(c); ulids.IsZero(userID) {
		return
	}

	// Parse the API key ID from the URL
	apiKeyID := c.Param("apiKeyID")

	// Figure out which project this key belongs to
	key := &qd.APIKey{}
	if key, err = s.quarterdeck.APIKeyDetail(ctx, apiKeyID); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Delete the API key using Quarterdeck
	if err = s.quarterdeck.APIKeyDelete(ctx, apiKeyID); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Update project stats in the background
	s.tasks.QueueContext(middleware.TaskContext(c), radish.Func(func(ctx context.Context) error {
		return s.UpdateProjectStats(ctx, userID, key.ProjectID)
	}), radish.WithErrorf("could not update stats for project %s", key.ProjectID.String()))

	c.Status(http.StatusNoContent)
}

// APIKeyPermissions returns the API key permissions available to the user by
// forwarding the request to Quarterdeck.
//
// Route: GET /v1/apikeys/permissions
func (s *Server) APIKeyPermissions(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Get the API key permissions from Quarterdeck
	var perms []string
	if perms, err = s.quarterdeck.APIKeyPermissions(ctx); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	c.JSON(http.StatusOK, perms)
}
