package tenant

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
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

// MemberCreate adds a new member to an organization in the database
// and returns a 201 StatusCreated response.
//
// Route: /member
func (s *Server) MemberCreate(c *gin.Context) {
	var (
		err    error
		member *api.Member
		orgID  ulid.ULID
	)

	const MemberConfirmed = "Confirmed"

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

	// Verify that a member name exists and return a 400 response if it does not.
	if member.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member name is required"))
		return
	}

	// Verify that a member role exists and return a 400 response if it does not.
	if member.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member role is required"))
		return
	}

	dbMember := &db.Member{
		OrgID:  orgID,
		Email:  member.Email,
		Name:   member.Name,
		Role:   member.Role,
		Status: MemberConfirmed,
	}

	if err = db.CreateMember(c.Request.Context(), dbMember); err != nil {
		sentry.Error(c).Err(err).Msg("could not create member in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add member"))
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

	// Verify the member name exists and return a 400 responsoe if it doesn't.
	if member.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member name is required"))
		return
	}

	// Verify the member role exists and return a 400 response if it doesn't.
	if member.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member role is required"))
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
	var err error

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

	// TODO: Add org verification

	// Bind the user request with JSON.
	params := &api.UpdateMemberParams{}
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse member update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify member role exists.
	if params.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member role is required"))
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
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member role"))
		return
	}

	// Check to ensure the memberID from the URL matches the member ID from the database.
	if memberID != member.ID {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member id does not match id in URL"))
	}

	// Update member role.
	member.Role = params.Role

	// Update member in the database.
	if err = db.UpdateMember(c.Request.Context(), member); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not update member in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member"))
		return
	}

	// Get members from the database and set page size to return all members.
	// TODO: Create list method that will not require pagination for this endpoint.
	getAll := &pg.Cursor{StartIndex: "", EndIndex: "", PageSize: 100}
	var members []*db.Member
	if members, _, err = db.ListMembers(c.Request.Context(), orgID, getAll); err != nil {
		sentry.Error(c).Err(err).Msg("could not list members")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member role"))
		return
	}

	// Loop over dbMember and count the number of members whose role is Owner to verify that at least one Owner remains in the organization.
	count := 0
	for _, dbMember := range members {
		if dbMember.Role == perms.RoleOwner {
			count++
		}
	}

	if count < 1 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("organization must have at least one owner"))
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
