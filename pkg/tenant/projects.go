package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// TenantProjectList retrieves all projects assigned to a tenant
// and returns a 200 OK response.
//
// Route: //tenant/:tenantID/projects
func (s *Server) TenantProjectList(c *gin.Context) {
	var (
		err error
	)
	// Get the project's tenant ID from the URL and return a 400 response
	// if the tenant ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant ulid"))
		return
	}

	// Get projects from the database and return a 500 response
	// if not successful.
	if _, err := db.ListProjects(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not fetch projects from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects from the database"))
		return
	}

	// Build the response.
	out := &api.TenantProjectPage{TenantProjects: make([]*api.Project, 0)}

	tenantProject := &api.Project{}

	out.TenantProjects = append(out.TenantProjects, tenantProject)

	c.JSON(http.StatusOK, out)
}

func (s *Server) TenantProjectCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// ProjectList retrieves all projects assigned to a tenant
// and returns a 200 OK response.
//
// Route: /projects
func (s *Server) ProjectList(c *gin.Context) {
	// TODO: Fetch the project's tenant ID from key.

	// Get projects from the database and return a 500 response
	// if not successful.
	if _, err := db.ListProjects(c.Request.Context(), ulid.ULID{}); err != nil {
		log.Error().Err(err).Msg("could not fetch projects from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects from the database"))
		return
	}

	// Build the response.
	out := &api.ProjectPage{Projects: make([]*api.Project, 0)}

	project := &api.Project{}

	out.Projects = append(out.Projects, project)

	c.JSON(http.StatusOK, out)
}

func (s *Server) ProjectCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
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
