package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

func (s *Server) TenantProjectList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TenantProjectCreate adds a new tenant project to the database
// and returns a 201 StatusCreated response.
//
// Route: /tenant/:tenantID/projects
func (s *Server) TenantProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
	)

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&project); err != nil {
		log.Warn().Err(err).Msg("could not bind tenant project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a project ID does not exist and return a 400 response
	// if the project ID exists.
	if project.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project id cannot be specified on create"))
		return
	}

	// Verify that a project names has been provided and return a 400 response
	// if the project name does not exist.
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project name is required"))
		return
	}

	// Add project to the database and return a 500 response if it cannot be added.
	if err = db.CreateProject(c.Request.Context(), &db.Project{Name: project.Name}); err != nil {
		log.Error().Err(err).Msg("could not create tenant project in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant project"))
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (s *Server) ProjectList(c *gin.Context) {
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

// ProjectCreate adds a new project to the database and returns
// a 201 StatusCreated response.
//
// Route: /project
func (s *Server) ProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
	)

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&project); err != nil {
		log.Warn().Err(err).Msg("could not bind project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a project ID does not exist and return a 400 response
	// if the project ID exists.
	if project.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project ID cannot be specified on create"))
		return
	}

	// Verify that a project name has been provided and return a 400 response
	// if the project name does not exist.
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project name is required"))
		return
	}

	// Add project to the database and return a 500 response if not successful.
	if err = db.CreateProject(c.Request.Context(), &db.Project{Name: project.Name}); err != nil {
		log.Error().Err(err).Msg("could not create project in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add project"))
		return
	}

	c.JSON(http.StatusCreated, project)
}

// ProjectDetail retrieves a summary detail of a project by its
// ID and returns a 200 OK response.
//
// Route: /project/:projectID
func (s *Server) ProjectDetail(c *gin.Context) {
	var (
		err   error
		reply *api.Project
	)

	// Get the project ID from the URL and return a 400 response
	// if the project does not exist.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Error().Err(err).Msg("could not parse project ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project ulid"))
		return
	}

	// Get the specified project from the database and return a 404 response
	// if it cannot be retrieved.
	var project *db.Project
	if project, err = db.RetrieveProject(c.Request.Context(), projectID); err != nil {
		log.Error().Err(err).Str("projectID", projectID.String()).Msg("could not retrieve project")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not retrieve project"))
		return
	}

	reply = &api.Project{
		ID:   project.ID.String(),
		Name: project.Name,
	}

	c.JSON(http.StatusOK, reply)
}

// ProjectUpdate updates the record of a project with a given ID
// and returns a 200 OK response.
//
// Route: /project/:projectID
func (s *Server) ProjectUpdate(c *gin.Context) {
	var (
		err     error
		project *api.Project
	)

	// Get the project ID from the URL and return a 400 response if
	// the project ID is not a ULID.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Error().Err(err).Msg("could not parse project ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project ulid"))
		return
	}

	// Bind the user request with JSON and return a 400 response
	// if binding is not successful.
	if err = c.BindJSON(&project); err != nil {
		log.Warn().Err(err).Msg("could not parse project update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind user request"))
		return
	}

	// Verify that the project name exists and return a 400 response if it doesn't.
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project name is required"))
		return
	}

	// Get the specified project from the database and return a 404 response if
	// it cannot be retrieved.
	var p *db.Project
	if p, err = db.RetrieveProject(c.Request.Context(), projectID); err != nil {
		log.Error().Err(err).Str("projectID", projectID.String()).Msg("could not retrieve project")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Update project in the database and return a 500 response if the
	// project record cannot be updated.
	if err = db.UpdateProject(c.Request.Context(), p); err != nil {
		log.Error().Err(err).Str("projectID", projectID.String()).Msg("could not save project")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update project"))
		return
	}
	c.JSON(http.StatusOK, project)
}

// ProjectDelete deletes a project from a user's request with a given ID
// and returns a 200 OK response instead of an error response.
//
// Route: /project/:projectID
func (s *Server) ProjectDelete(c *gin.Context) {
	var (
		err error
	)

	// Get the project ID from the URL and return a 400 response
	// if the project does not exist.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Error().Err(err).Msg("could not parse project ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project ulid"))
		return
	}

	// Delete the project and return a 404 response if it cannot be removed.
	if err = db.DeleteProject(c.Request.Context(), projectID); err != nil {
		log.Error().Err(err).Str("projectID", projectID.String()).Msg("could not delete project")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not delete project"))
		return
	}
	c.Status(http.StatusOK)
}
