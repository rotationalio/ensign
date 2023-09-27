package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

// MemberList retrieves members assigned to a specified organization
// and returns a 200 OK response.
//
// Route: /member
func (s *Server) MemberList(c *gin.Context) {
	var (
		err        error
		orgID      ulid.ULID
		next, prev *pg.Cursor
	)

	// Members exist in organizations
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	query := &api.PageQuery{}
	if err = c.BindQuery(query); err != nil {
		log.Error().Err(err).Msg("could not parse query")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query"))
		return
	}

	if query.NextPageToken != "" {
		if prev, err = pg.Parse(query.NextPageToken); err != nil {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse next page token"))
			return
		}
	} else {
		prev = pg.New("", "", int32(query.PageSize))
	}

	// Get members from the database and return a 500 response if not succesful.
	var members []*db.Member
	if members, next, err = db.ListMembers(c.Request.Context(), orgID, prev); err != nil {
		sentry.Error(c).Err(err).Msg("could not list members")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list members"))
		return
	}

	// Build the response.
	out := &api.MemberPage{Members: make([]*api.Member, 0)}

	// Loop over db.Member and retrieve each member.
	for _, dbMember := range members {
		out.Members = append(out.Members, dbMember.ToAPI())
	}

	if next != nil {
		if out.NextPageToken, err = next.NextPageToken(); err != nil {
			log.Error().Err(err).Msg("could not set next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list members"))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

// MemberCreate starts the team member invitation process by forwarding the request to
// Quarterdeck. If successful, an invitation email is sent to the email address in the
// request and a unverified member is created in Trtl, returning a 201 Created
// response.
//
// Route: /member
func (s *Server) MemberCreate(c *gin.Context) {
	var (
		err    error
		ctx    context.Context
		member *api.Member
		orgID  ulid.ULID
	)

	// User credentials are required to make the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Members exist in organizations
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&member); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse member create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify that a member id does not exist and return a 400 response if it does.
	if member.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member id cannot be specified on create"))
		return
	}

	// Verify that a member email exists.
	if member.Email == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member email is required"))
		return
	}

	// Verify that a member role exists and return a 400 response if it does not.
	if member.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member role is required"))
		return
	}

	// Validate user role
	if !perms.IsRole(member.Role) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown member role"))
		return
	}

	// Email address must be unique in the organization.
	if err = db.VerifyMemberEmail(c.Request.Context(), orgID, member.Email); err != nil {
		if errors.Is(err, db.ErrMemberExists) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("team member already exists with this email address"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not check team member existence")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add team member"))
		return
	}

	// Call Quarterdeck to create and send the invite email.
	req := &qd.UserInviteRequest{
		Email: member.Email,
		Role:  member.Role,
	}

	var reply *qd.UserInviteReply
	if reply, err = s.quarterdeck.InviteCreate(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Create the pending record in the database.
	dbMember := &db.Member{
		OrgID:        reply.OrgID,
		ID:           reply.UserID,
		Email:        reply.Email,
		Name:         reply.Name,
		Role:         reply.Role,
		Organization: reply.Organization,
		Workspace:    reply.Workspace,
		Invited:      true,
	}

	if err = db.CreateMember(c.Request.Context(), dbMember); err != nil {
		sentry.Error(c).Err(err).Msg("could not create member in database after invitation")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add team member"))
		return
	}

	c.JSON(http.StatusCreated, dbMember.ToAPI())
}

// MemberDetail retrieves a summary detail of a member by its ID
// and returns a 200 OK response.
//
// Route: /member/:memberID
func (s *Server) MemberDetail(c *gin.Context) {
	var err error

	// Members exist on organizations
	// This method handles the logging and error responses
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the member ID from the URL and return a 400 if the member does not exist.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("memberID")).Msg("could not parse member id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
		return
	}

	// Get the specified member from the database
	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve member from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member"))
		return
	}

	c.JSON(http.StatusOK, member.ToAPI())
}

// MemberUpdate updates the record of a member with a given ID. This endpoint is used
// to update metadata for team members but does not allow user profile information to
// be updated. Multiple errors may be returned if there are multiple errors in the
// profile.
//
// route: /member/:memberID
func (s *Server) MemberUpdate(c *gin.Context) {
	var (
		err error
		req *api.Member
	)

	// Members exist on organizations
	// This method handles the logging and error responses
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the member ID from the URL and return a 400 if the
	// member ID is not a ULID.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("memberID")).Msg("could not parse member id")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrMemberNotFound))
		return
	}

	// Bind the user request with JSON and return a 400 response
	// if binding is not successful.
	if err = c.BindJSON(&req); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse member update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrMemberNotFound))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve member from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Only the user name can be updated for now.
	req.Normalize()
	member.Name = req.Name

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

func (s *Server) MemberRoleUpdate(c *gin.Context) {
	var (
		err error
		ctx context.Context
	)

	// Quarterdeck request requires credentials in the context
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Warn(c).Err(err).Msg("could not retrieve credentials from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not retrieve credentials from request"))
		return
	}

	// Members exist on organizations
	// This method handles the logging and error responses
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the member ID from the URL and return a 400 if the member ID is not a ULID.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("memberID")).Msg("could not parse member id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
		return
	}

	// Bind the user request with JSON.
	params := &api.UpdateRoleParams{}
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse member update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify member role exists.
	if params.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("team member role is required"))
		return
	}

	// Verify the role provided is valid.
	if !perms.IsRole(params.Role) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown team member role"))
		return
	}

	// Check that the member can be updated.
	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("team member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve member from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update team member role"))
		return
	}

	// If the member to be updated is an owner, loop over dbMember and break out of the loop if there are at least two owners.
	// If member is the only owner, their role cannot be changed.
	if member.Role == perms.RoleOwner {

		// Verify if org has more than one owner.
		var count int
		if count, err = orgOwnerCount(c.Request.Context(), orgID); err != nil {
			sentry.Error(c).Err(err).Msg("could not list members")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member role"))
			return
		}

		switch count {
		case 0:
			sentry.Warn(c).Err(err).Msg("could not find any owners")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member role"))
			return
		case 1:
			c.JSON(http.StatusBadRequest, api.ErrorResponse("organization must have at least one owner"))
			return
		}
	}

	// TOOD: Should we allow pending members to be role updated?
	if member.OnboardingStatus() != db.MemberStatusActive {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot update role for pending team member"))
		return
	}

	// Ensure that the role can be updated.
	if member.Role == params.Role {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("team member already has the requested role"))
		return
	}

	// Update role in quarterdeck
	req := &qd.UpdateRoleRequest{
		ID:   memberID,
		Role: params.Role,
	}

	var reply *qd.User
	if reply, err = s.quarterdeck.UserRoleUpdate(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Update member role with the role returned from quarterdeck.
	member.Role = reply.Role
	if err = db.UpdateMember(c.Request.Context(), member); err != nil {
		sentry.Error(c).Err(err).Msg("could not update member in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update team member"))
		return
	}

	c.JSON(http.StatusOK, member.ToAPI())
}

// MemberDelete attempts to delete a team member from an organization by forwarding the
// request to Quarterdeck. If the deleted field is set to true in the Quarterdeck
// response, the team member is deleted from the Tenant database. If the deleted field
// is not set in the response, then additional confirmation is required from the user
// so this endpoint returns the confirmation details which includes a token. The token
// must be provided to the MemberDeleteConfirm endpoint to complete the delete.
// Otherwise, the team member is not deleted.
//
// Route: /member/:memberID
func (s *Server) MemberDelete(c *gin.Context) {
	var (
		err error
		ctx context.Context
	)

	// Quarterdeck request requires an authenticated context
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve credentials from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not retrieve credentials from request"))
		return
	}

	// Members exist on organizations
	// This method handles the logging and error responses
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the member ID from the URL and return a 404 response
	// if the member does not exist.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("memberID")).Msg("could not parse member id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
		return
	}

	// Retrieve member from the database.
	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve member from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete member"))
		return
	}

	// Check to ensure member is not the only owner of the organization.
	if member.Role == perms.RoleOwner {
		// Verify if org has more than one owner.
		var count int
		if count, err = orgOwnerCount(c.Request.Context(), member.OrgID); err != nil {
			sentry.Error(c).Err(err).Msg("could not list members")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete member"))
			return
		}

		switch count {
		case 0:
			sentry.Warn(c).Err(err).Msg("could not find any owners")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete member"))
			return
		case 1:
			c.JSON(http.StatusBadRequest, api.ErrorResponse("organization must have at least one owner"))
			return
		}
	}

	// Attempt to remove the user from the Quarterdeck organization.
	var reply *qd.UserRemoveReply
	if reply, err = s.quarterdeck.UserRemove(ctx, member.ID.String()); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	out := &api.MemberDeleteReply{
		APIKeys: reply.APIKeys,
		Token:   reply.Token,
		Deleted: reply.Deleted,
	}

	// If delete requires confirmation then just return the confirmation details.
	if !reply.Deleted {
		c.JSON(http.StatusOK, out)
		return
	}

	// If delete was successful then delete the member from the database
	if err = db.DeleteMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not delete member from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete member"))
		return
	}

	c.JSON(http.StatusOK, out)
}

// Helper method to check if an organization has more than one owner.
func orgOwnerCount(ctx context.Context, orgID ulid.ULID) (count int, err error) {
	// Get members from the database and set page size to return all members.
	// TODO: Create list method that will not require pagination for this endpoint.
	getAll := &pg.Cursor{StartIndex: "", EndIndex: "", PageSize: 100}
	var members []*db.Member
	if members, _, err = db.ListMembers(ctx, orgID, getAll); err != nil {
		return count, err
	}

	count = 0
	for _, dbMember := range members {
		if dbMember.Role == perms.RoleOwner {
			count++
			if count >= 2 {
				break
			}
		}
	}
	return count, nil
}
