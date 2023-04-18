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
		OrgID:  reply.OrgID,
		ID:     reply.UserID,
		Email:  reply.Email,
		Name:   reply.Name,
		Role:   reply.Role,
		Status: db.MemberStatusPending,
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

// MemberUpdate updates the record of a member with a given ID and
// returns a 200 OK response.
//
// route: /member/:memberID
func (s *Server) MemberUpdate(c *gin.Context) {
	var (
		err    error
		member *api.Member
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
		c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
		return
	}

	// Bind the user request with JSON and return a 400 response
	// if binding is not successful.
	if err = c.BindJSON(&member); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse member update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify the member email exists.
	if member.Email == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member email is required"))
		return
	}

	// Verify the member role exists and return a 400 response if it doesn't.
	if member.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member role is required"))
		return
	}

	// Verify the role provided is valid.
	if !perms.IsRole(member.Role) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown member role"))
		return
	}

	var m *db.Member
	if m, err = db.RetrieveMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve member from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member"))
		return
	}

	// Update all fields provided by the user
	m.Email = member.Email
	m.Name = member.Name
	m.Role = member.Role

	// Update member in the database.
	if err = db.UpdateMember(c.Request.Context(), m); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not update member in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve member"))
		return
	}

	c.JSON(http.StatusOK, m.ToAPI())
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

	// Check to ensure the memberID from the URL matches the member ID from the database.
	if memberID.Compare(member.ID) != 0 {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("member id does not match id in URL"))
		return
	}

	// If the member to be updated is an owner, loop over dbMember and break out of the loop if there are at least two owners.
	// If member is the only owner, their role cannot be changed.
	if member.Role == perms.RoleOwner {
		// Get members from the database and set page size to return all members.
		// TODO: Create list method that will not require pagination for this endpoint.
		getAll := &pg.Cursor{StartIndex: "", EndIndex: "", PageSize: 100}
		var members []*db.Member
		if members, _, err = db.ListMembers(c, orgID, getAll); err != nil {
			sentry.Error(c).Err(err).Msg("could not list members")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not get member from the database"))
			return
		}

		// Verify if org has more than one owner.
		count := orgOwnerCount(members)

		switch count {
		case 0:
			sentry.Warn(c).Err(err).Msg("could not find any owners")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not get member from the database"))
			return
		case 1:
			c.JSON(http.StatusBadRequest, api.ErrorResponse("organization must have at least one owner"))
			return
		}
	}

	// TOOD: Should we allow invitations to be updated?
	if member.Status == db.MemberStatusPending {
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

// MemberDelete deletes a member from a user's request with a given
// ID and returns a 200 OK response instead of an an error response.
//
// Route: /member/:memberID
func (s *Server) MemberDelete(c *gin.Context) {
	var err error

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

	// Check to ensure the memberID from the URL matches the member ID from the database.
	if memberID.Compare(member.ID) != 0 {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("member id does not match id in URL"))
		return
	}

	// Check to ensure member is not the only owner of the organization.
	if member.Role == perms.RoleOwner {
		// Get members from the database and set page size to return all members.
		// TODO: Create list method that will not require pagination for this endpoint.
		getAll := &pg.Cursor{StartIndex: "", EndIndex: "", PageSize: 100}
		var members []*db.Member
		if members, _, err = db.ListMembers(c, orgID, getAll); err != nil {
			sentry.Error(c).Err(err).Msg("could not list members")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not get member from the database"))
			return
		}

		// Verify if org has more than one owner.
		count := orgOwnerCount(members)

		switch count {
		case 0:
			sentry.Warn(c).Err(err).Msg("could not find any owners")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not get member from the database"))
			return
		case 1:
			c.JSON(http.StatusBadRequest, api.ErrorResponse("organization must have at least one owner"))
			return
		}
	}

	// Delete the member from the database
	if err = db.DeleteMember(c.Request.Context(), orgID, memberID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not delete member from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete member"))
		return
	}
	c.Status(http.StatusOK)
}

// InvitePreview returns "preview" information about an invite given a token. This
// endpoint must not be authenticated because unauthorized users should be able to
// accept organization invitations. Frontends should use this endpoint to validate an
// invitation token after the user has clicked on an invitation link in their email.
// The preview must contain enough information so the user knows which organization
// they are joining and also whether or not the email address is already registered to
// an account. This allows frontends to know whether or not to prompt the user to
// login or to create a new account.
//
// Route: /invites/:token
func (s *Server) InvitePreview(c *gin.Context) {
	var err error

	token := c.Param("token")

	// Call Quarterdeck to retrieve the invite preview.
	var rep *qd.UserInvitePreview
	if rep, err = s.quarterdeck.InvitePreview(c.Request.Context(), token); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Create the preview response
	out := &api.MemberInvitePreview{
		Email:       rep.Email,
		OrgName:     rep.OrgName,
		InviterName: rep.InviterName,
		Role:        rep.Role,
		HasAccount:  rep.UserExists,
	}
	c.JSON(http.StatusOK, out)
}

// Helper method to check if an organization has more than one owner.
func orgOwnerCount(members []*db.Member) (count uint8) {
	count = 0
	for _, dbMember := range members {
		if dbMember.Role == perms.RoleOwner {
			count++
			if count >= 2 {
				break
			}
		}
	}
	return count
}
