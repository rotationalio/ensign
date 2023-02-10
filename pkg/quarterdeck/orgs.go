package quarterdeck

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rs/zerolog/log"
)

// Retrieve an organization by ID. Users are only allowed to retrieve their own
// organization.
// TODO: Eventually allow users to retrieve other organizations they are a part of.
func (s *Server) OrganizationDetail(c *gin.Context) {
	var (
		err    error
		orgID  ulid.ULID
		org    *models.Organization
		claims *tokens.Claims
	)

	// Parse the orgID passed in from the URL
	if orgID, err = ulid.Parse(c.Param("id")); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// User claims are required to verify the user's organization
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	// User must be a member of the organization
	if claims.OrgID != orgID.String() {
		log.Warn().Str("orgid", orgID.String()).Msg("user cannot access this organization")
		c.JSON(http.StatusForbidden, api.ErrorResponse("user is not authorized to access this organization"))
		return
	}

	// Fetch the organization from the database
	if org, err = models.GetOrg(c.Request.Context(), orgID); err != nil {
		c.Error(err)
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		return
	}

	// Build the response from the model
	c.JSON(http.StatusOK, org.ToAPI())
}

func (s *Server) ProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
		claims  *tokens.Claims
	)

	// Bind the Project request to the project data structure
	if err = c.BindJSON(&project); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	// Validate the request from the API side.
	if err = project.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims unavailable"))
		return
	}

	// Create the OrganizationProjects mapping database model
	model := &models.OrganizationProject{
		OrgID:     claims.ParseOrgID(),
		ProjectID: project.ProjectID,
	}

	// Save the model to the database
	if err = model.Save(c.Request.Context()); err != nil {
		c.Error(err)
		switch err.(type) {
		case *models.ValidationError:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
		}
		return
	}

	// Update the response to send to the user
	project.OrgID = model.OrgID
	project.ProjectID = model.ProjectID
	project.Created, _ = model.GetCreated()
	project.Modified, _ = model.GetModified()
	c.JSON(http.StatusOK, project)
}
