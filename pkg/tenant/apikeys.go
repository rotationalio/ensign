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
	"github.com/rs/zerolog/log"
)

func (s *Server) ProjectAPIKeyList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) ProjectAPIKeyCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyList(c *gin.Context) {
}

// APIKeyCreate creates a new API key by forwarding the request to Quarterdeck.
//
// Route: POST /v1/apikeys
func (s *Server) APIKeyCreate(c *gin.Context) {
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

	// The user's name is on the token claims
	var claims *tokens.Claims
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch user claims from context"))
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
	if req.ProjectID, err = ulid.Parse(params.ProjectID); err != nil {
		log.Warn().Err(err).Str("project_id", params.ProjectID).Msg("could not parse project ID")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid project ID"))
		return
	}

	// TODO: Add source to request

	// Create the API key with Quarterdeck
	var key *qd.APIKey
	if key, err = s.quarterdeck.APIKeyCreate(ctx, req); err != nil {
		log.Error().Err(err).Msg("could not create API key")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create API key"))
		return
	}

	// Return the API key
	out := &api.APIKey{
		ID:           key.ID.String(),
		ClientID:     key.ClientID,
		ClientSecret: key.ClientSecret,
		Name:         key.Name,
		ProjectID:    key.ProjectID.String(),
		Owner:        claims.Name,
		Permissions:  key.Permissions,
		Created:      key.Created.Format(time.RFC3339Nano),
		Modified:     key.Modified.Format(time.RFC3339Nano),
	}

	c.JSON(http.StatusCreated, out)
}

func (s *Server) APIKeyDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) APIKeyDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}
