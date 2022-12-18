package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

func (s *Server) TenantMemberList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TenantMemberCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) MemberList(c *gin.Context) {
	// The following TODO task items will need to be
	// implemented for each endpoint.

	// TODO: Add authentication and authorization middleware
	// TODO: Identify top-level info
	// TODO: Parse and validate user input
	// TODO: Perform work on the request, e.g. database interactions,
	// sending notifications, accessing other services, etc.

	// Return response with the correct status code

	// TODO: Replace StatusNotImplemented with StatusOk and
	// replace "not yet implemented" message.
	c.JSON(http.StatusNotImplemented, "not implemented yet")
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
		log.Debug().Err(err).Msg("could not parse member ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse member id"))
		return
	}

	// Get the specified member from the database and return a 404 response
	// if it cannot be retrieved.
	var member *db.Member
	if member, err = db.RetrieveMember(c.Request.Context(), memberID); err != nil {
		log.Error().Err(err).Msg("could not retrieve member")
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

func (s *Server) MemberCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) MemberUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
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
