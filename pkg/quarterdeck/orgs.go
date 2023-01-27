package quarterdeck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

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
