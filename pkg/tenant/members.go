package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// TenantMemberList retrieves all members assigned to a tenant
// and returns a 200 OK response.
//
// Route: tenant/:tenantID/member
func (s *Server) TenantMemberList(c *gin.Context) {
	var (
		err error
	)

	// Get the member's tenant ID from the URL and return a 400 response
	// if the tenant ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant ulid"))
		return
	}

	// Get members from the database and return a 500 response
	// if not successful.
	var members []*db.Member
	if members, err = db.ListMembers(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not fetch members from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch members from the database"))
		return
	}

	// Build the response.
	out := &api.TenantMemberPage{
		TenantID:      tenantID.String(),
		TenantMembers: make([]*api.Member, 0),
	}

	// Loop over member. For each db.Member inside the array, create a tenantMember
	// which will be an api.Member{} and assign the ID and Name fetched from db.Member
	// to that struct and then append to the out.TenantMembers array.
	for _, dbMember := range members {
		tenantMember := &api.Member{
			ID:   dbMember.ID.String(),
			Name: dbMember.Name,
			Role: dbMember.Role,
		}
		out.TenantMembers = append(out.TenantMembers, tenantMember)
	}
	c.JSON(http.StatusOK, out)
}

// / TenantMemberCreate adds a new member to a tenant in the database
// and returns a 201 StatusCreated response.
//
// Route: /tenant/:tenantID/members
func (s *Server) TenantMemberCreate(c *gin.Context) {
	var (
		err    error
		member *api.Member
		out    *api.Member
	)

	// Get the tenant ID from the URL and return a 400 if the tenant does not exist.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Bind the user request with JSON and return a 400 response if
	// binding is not successful.
	if err = c.BindJSON(&member); err != nil {
		log.Warn().Err(err).Msg("could not bind tenant member create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a member ID does not exist and return a 400 response if
	// the member id exists.
	if member.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member id cannot be specified on create"))
		return
	}

	// Verify that a member name has been provided and return a 400 repsonse if
	// the member name does not exist.
	if member.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant member name is required"))
		return
	}

	// Verify that a member role has been provided and return a 400 response if
	// the member role does not exist.
	if member.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("tenant member role is required"))
		return
	}

	tmember := &db.Member{
		TenantID: tenantID,
		Name:     member.Name,
		Role:     member.Role,
	}

	if err = db.CreateTenantMember(c.Request.Context(), tmember); err != nil {
		log.Error().Err(err).Msg("could not create tenant member in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant member"))
		return
	}

	out = &api.Member{
		ID:   tmember.ID.String(),
		Name: member.Name,
		Role: member.Role,
	}

	c.JSON(http.StatusCreated, out)
}

// MemberList retrieves all members assigned to an organization
// and returns a 200 OK response.
//
// Route: /member
func (s *Server) MemberList(c *gin.Context) {
	// TODO: Fetch the member's tenant ID from key.
	var tenantID ulid.ULID

	// Get members from the database and return a 500 response
	// if not succesful.
	if _, err := db.ListMembers(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not fetch members from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch members from database"))
		return
	}

	// Build the response.
	out := &api.MemberPage{Members: make([]*api.Member, 0)}

	member := &api.Member{}

	out.Members = append(out.Members, member)

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
		out    *api.Member
	)

	// TODO: Add authentication middleware to fetch the organization ID.

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&member); err != nil {
		log.Warn().Err(err).Msg("could not bind member create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a member id does not exist and return a 400 response if
	// the member id exists.
	if member.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member id cannot be specified on create"))
		return
	}

	// Verify that a member name exists and return a 400 response if
	// the member name does not exist.
	if member.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member name is required"))
		return
	}

	if member.Role == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("member role is required"))
		return
	}

	m := &db.Member{
		Name: member.Name,
		Role: member.Role,
	}

	if err = db.CreateMember(c.Request.Context(), m); err != nil {
		log.Error().Err(err).Msg("could not create member in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add member"))
		return
	}

	out = &api.Member{
		ID:   m.ID.String(),
		Name: member.Name,
		Role: member.Role,
	}

	c.JSON(http.StatusCreated, out)
}

// MemberDetail retrieves a summary detail of a member by its ID
// and returns a 200 OK response.
//
// Route: /member/:memberID
func (s *Server) MemberDetail(c *gin.Context) {
	var (
		err   error
		reply *api.Member
	)

	// Get the member ID from the URL and return a 400 if the member does not exist.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		log.Error().Err(err).Msg("could not parse member ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse member id"))
		return
	}

	// Get the specified member from the database and return a 404 response
	// if it cannot be retrieved.
	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), memberID); err != nil {
		log.Error().Err(err).Str("memberID", memberID.String()).Msg("could not retrieve member")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not retrieve member"))
		return
	}

	reply = &api.Member{
		ID:   member.ID.String(),
		Name: member.Name,
		Role: member.Role,
	}

	c.JSON(http.StatusOK, reply)
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

	// Get the member ID from the URL and return a 400 if the
	// member ID is not a ULID.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		log.Error().Err(err).Msg("could not parse member id")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse member id"))
		return
	}

	// Bind the user request with JSON and return a 400 response
	// if binding is not successful.
	if err = c.BindJSON(&member); err != nil {
		log.Warn().Err(err).Msg("could not parse member update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind user request"))
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

	// Get the specified member from the database and return a 404 response
	// if it cannot be retrieved.
	var m *db.Member
	if m, err = db.RetrieveMember(c.Request.Context(), memberID); err != nil {
		log.Error().Err(err).Str("memberID", memberID.String()).Msg("could not retrieve member")
		c.JSON(http.StatusNotFound, api.ErrorResponse("member not found"))
		return
	}

	// Update member in the database and return a 500 response if the
	// member record cannot be updated.
	if err = db.UpdateMember(c.Request.Context(), m); err != nil {
		log.Error().Err(err).Str("memberID", memberID.String()).Msg("could not save member")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update member"))
		return
	}
	c.JSON(http.StatusOK, member)
}

// MemberDelete deletes a member from a user's request with a given
// ID and returns a 200 OK response instead of an an error response.
//
// Route: /member/:memberID
func (s *Server) MemberDelete(c *gin.Context) {
	var (
		err error
	)

	// Get the member ID from the URL and return a 400 response
	// if the member does not exist.
	var memberID ulid.ULID
	if memberID, err = ulid.Parse(c.Param("memberID")); err != nil {
		log.Error().Err(err).Msg("could not parse member ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse member id"))
		return
	}

	// Delete the member and return a 404 response if it cannot be removed.
	if err = db.DeleteMember(c.Request.Context(), memberID); err != nil {
		log.Error().Err(err).Str("memberID", memberID.String()).Msg("could not delete member")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not delete member"))
		return
	}
	c.Status(http.StatusOK)
}
