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
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/rs/zerolog/log"
)

// Organization Detail fetches the details for an organization from Quarterdeck.
//
// Route: GET /v1/organizations/:orgID
func (s *Server) OrganizationDetail(c *gin.Context) {
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

	// Fetch the orgID from the claims
	var claimsID ulid.ULID
	if claimsID = orgIDFromContext(c); ulids.IsZero(claimsID) {
		return
	}

	// Parse the orgID passed in from the URL
	paramID := c.Param("orgID")
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(paramID); err != nil {
		log.Warn().Str("id", paramID).Err(err).Msg("could not parse orgID from URL")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse organization ID"))
		return
	}

	// User can only list their own organization
	if claimsID.Compare(orgID) != 0 {
		log.Warn().Str("user_org", claimsID.String()).Str("params_org", orgID.String()).Msg("user cannot access this organization")
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// Fetch the organization from Quarterdeck
	var org *qd.Organization
	if org, err = s.quarterdeck.OrganizationDetail(ctx, paramID); err != nil {
		log.Error().Err(err).Msg("could not fetch organization from Quarterdeck")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not detail organization"))
		return
	}

	// Build the response from the Quarter
	out := &api.Organization{
		ID:       org.ID.String(),
		Name:     org.Name,
		Domain:   org.Domain,
		Created:  org.Created.Format(time.RFC3339Nano),
		Modified: org.Modified.Format(time.RFC3339Nano),
	}
	c.JSON(http.StatusOK, out)
}

// Helper to fetch the orgID from the gin context. This method also logs and returns
// any errors to allow endpoints to have consistent error handling.
func orgIDFromContext(c *gin.Context) (orgID ulid.ULID) {
	var (
		claims *tokens.Claims
		err    error
	)
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not get user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return ulid.ULID{}
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		log.Error().Err(err).Msg("could not parse orgID from claims")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("invalid user claims"))
		return ulid.ULID{}
	}

	return orgID
}
