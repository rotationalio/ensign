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

// Organization Detail fetches the details for an organization from Quarterdeck.
//
// Route: GET /v1/organizations/:orgID
func (s *Server) OrganizationDetail(c *gin.Context) {
	var (
		ctx    context.Context
		claims *tokens.Claims
		err    error
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// User claims are required to verify that the user is in the organization
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not get user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch claims for authenticated user"))
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
	if claims.OrgID != orgID.String() {
		log.Warn().Str("orgid", orgID.String()).Msg("user cannot access this organization")
		c.JSON(http.StatusForbidden, api.ErrorResponse("user is not authorized to access this organization"))
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
