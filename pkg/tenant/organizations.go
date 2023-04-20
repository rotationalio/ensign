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
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

// OrganizationList fetches the list of organizations the authenticated user is a part of from Quarterdeck.
//
// Route: GET /v1/organization
func (s *Server) OrganizationList(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// User credentials are required to make the Quarterdeck request.
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Parse query parameters.
	query := &api.PageQuery{}
	if err = c.ShouldBindQuery(query); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse page query request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Build the Quarterdeck request.
	req := &qd.OrganizationPageQuery{
		PageSize:      int(query.PageSize),
		NextPageToken: query.NextPageToken,
	}

	// Request a page of organizations from Quarterdeck.
	var reply *qd.OrganizationList
	if reply, err = s.quarterdeck.OrganizationList(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Return page of organizations.
	out := &api.OrganizationPage{
		Organizations: make([]*api.Organization, 0),
		NextPageToken: reply.NextPageToken,
	}

	for _, org := range reply.Organizations {
		orgs := &api.Organization{
			ID:        org.ID.String(),
			Name:      org.Name,
			Domain:    org.Domain,
			Created:   db.TimeToString(org.Created),
			LastLogin: db.TimeToString(org.LastLogin),
		}
		out.Organizations = append(out.Organizations, orgs)
	}

	c.JSON(http.StatusOK, out)
}

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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Fetch the orgID from the claims
	var claimsID ulid.ULID
	if claimsID = orgIDFromContext(c); ulids.IsZero(claimsID) {
		return
	}

	// Parse the orgID passed in from the URL
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(c.Param("orgID")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("orgID")).Msg("could not parse org id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// User can only list their own organization
	if claimsID.Compare(orgID) != 0 {
		sentry.Warn(c).Str("user_org", claimsID.String()).Str("params_org", orgID.String()).Msg("user cannot access this organization")
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// Fetch the organization from Quarterdeck
	var org *qd.Organization
	if org, err = s.quarterdeck.OrganizationDetail(ctx, orgID.String()); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Build the response from the Quarterdeck response
	out := &api.Organization{
		ID:       org.ID.String(),
		Name:     org.Name,
		Domain:   org.Domain,
		Projects: org.Projects,
		Created:  org.Created.Format(time.RFC3339Nano),
		Modified: org.Modified.Format(time.RFC3339Nano),
	}

	// Get the organization owner
	if out.Owner, err = getOwner(ctx, org); err != nil {
		sentry.Error(c).Err(err).Str("org", org.ID.String()).Msg("could not retrieve organization owner")
	}

	c.JSON(http.StatusOK, out)
}

// Helper to fetch the owner of the organization. Since an organization can have
// multiple owners, this method returns the first owner found.
func getOwner(ctx context.Context, org *qd.Organization) (_ string, err error) {
	// List the members in the organization
	// TODO: Create list method that will not require pagination for this endpoint.
	// Set page size to return all projects and topics.
	getAll := &pg.Cursor{StartIndex: "", EndIndex: "", PageSize: 100}
	var members []*db.Member
	if members, _, err = db.ListMembers(ctx, org.ID, getAll); err != nil {
		return "", err
	}

	// Return the first owner found
	// TODO: Once user invites are implemented, this may need to be updated to list all
	// the owners or the original owner.
	for _, member := range members {
		if member.Role == permissions.RoleOwner {
			return member.Name, nil
		}
	}

	// Organizations should have at least one owner
	return "", errors.New("organization has no owners")
}

// Helper to fetch the orgID from the gin context. This method also logs and returns
// any errors to allow endpoints to have consistent error handling.
func orgIDFromContext(c *gin.Context) (orgID ulid.ULID) {
	var (
		claims *tokens.Claims
		err    error
	)
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return ulid.ULID{}
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Error(c).Err(err).Msg("could not parse orgID from claims")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("invalid user claims"))
		return ulid.ULID{}
	}

	return orgID
}
