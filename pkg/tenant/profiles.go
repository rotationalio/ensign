package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

// ProfileDetail retrieves profile information for the authenticated user based on
// their current claims.
//
// Route: GET /profile
func (s *Server) ProfileDetail(c *gin.Context) {
	var (
		err             error
		memberID, orgID ulid.ULID
		claims          *tokens.Claims
		member          *db.Member
	)

	// Fetch the claims for the authenticated user
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claism from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// The memberID is the subject in the claims: their user ID
	if memberID = claims.ParseUserID(); ulids.IsZero(memberID) {
		sentry.Error(c).Msg("could not parse user ID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Get the orgID for the user's logged in organization
	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Error(c).Msg("could not parse orgID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Fetch the user's profile from the database
	if member, err = db.RetrieveMember(c, orgID, memberID); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve user profile")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	c.JSON(http.StatusOK, member.ToAPI())
}

// ProfileUpdate allows a user to update their own profile information within the
// context of their current logged in organization. This endpoint is also used to
// update profile information during the onboarding process, so it may make a request
// to Quarterdeck to update organization info for new users. Multiple errors may be
// returned if there are multiple errors in the onboarding information.
//
// Route: PUT /profile
func (s *Server) ProfileUpdate(c *gin.Context) {
	var (
		err             error
		claims          *tokens.Claims
		orgID, memberID ulid.ULID
		req             *api.Member
	)

	// Fetch the claims for the authenticated user
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claism from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// The memberID is the subject in the claims: their user ID
	if memberID = claims.ParseUserID(); ulids.IsZero(memberID) {
		sentry.Error(c).Msg("could not parse user ID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Get the orgID for the user's logged in organization
	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Error(c).Msg("could not parse orgID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Parse the request body into the member update
	if err = c.BindJSON(&req); err != nil {
		sentry.Error(c).Err(err).Msg("could not parse profile update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryProfileAgain))
		return
	}

	// This endpoint only allows users to update their own profile
	if req.ID != memberID.String() {
		sentry.Error(c).Msg("user attempted to update someone else's profile")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryProfileAgain))
		return
	}

	// Fetch the member record to be updated
	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), orgID, memberID); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve member from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Update all fields that are allowed to be updated
	req.Normalize()
	member.Name = req.Name
	member.ProfessionSegment = req.ProfessionSegment
	member.DeveloperSegment = req.DeveloperSegment

	if !member.Invited {
		member.Organization = req.Organization
		member.Workspace = req.Workspace
	}

	// Validate the member update, this is also validated in UpdateMember() but this
	// ensures that an invalid organization or workspace is not sent to Quarterdeck.
	if err = member.Validate(); err != nil {
		var verrs db.ValidationErrors
		switch {
		case errors.As(err, &verrs):
			// Return validation errors to the frontend with field names.
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verrs.ToAPI()))
		default:
			sentry.Error(c).Err(err).Msg("could not validate member update")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}
		return
	}

	if !member.Invited && member.IsOnboarded() {
		// If user is done onboarding, update the organization details in Quarterdeck.
		var ctx context.Context
		if ctx, err = middleware.ContextFromRequest(c); err != nil {
			sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
			c.JSON(http.StatusUnauthorized, api.ErrorResponse(responses.ErrTryLoginAgain))
			return
		}

		// Update the organization in Quarterdeck.
		org := &qd.Organization{
			ID:     orgID,
			Name:   member.Organization,
			Domain: member.Workspace,
		}
		if _, err = s.quarterdeck.OrganizationUpdate(ctx, org); err != nil {
			sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
			api.ReplyQuarterdeckError(c, err)
			return
		}
	}

	// Update member in the database.
	if err = db.UpdateMember(c.Request.Context(), member); err != nil {
		var verrs db.ValidationErrors
		switch {
		case errors.Is(err, db.ErrNotFound):
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrMemberNotFound))
		case errors.As(err, &verrs):
			// Return validation errors to the frontend with field names.
			c.JSON(http.StatusBadRequest, api.ErrorResponse(verrs.ToAPI()))
		default:
			sentry.Error(c).Err(err).Msg("could not update member in the database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}
		return
	}

	c.JSON(http.StatusOK, member.ToAPI())
}
