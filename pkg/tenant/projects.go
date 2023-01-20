package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
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
	var projects []*db.Project
	if projects, err = db.ListProjects(c.Request.Context(), tenantID); err != nil {
		log.Error().Err(err).Msg("could not fetch projects from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects from the database"))
		return
	}

	// Build the response.
	out := &api.TenantProjectPage{
		TenantID:       tenantID.String(),
		TenantProjects: make([]*api.Project, 0),
	}

	// Loop over projects. For each db.Project inside the array, create a tenantProject
	// which will be an api.Project{} and assign the ID and Name fetched from db.Project
	// to that struct and then append to the out.TenantProjects array.
	for _, dbProject := range projects {
		tenantProject := &api.Project{
			ID:   dbProject.ID.String(),
			Name: dbProject.Name,
		}
		out.TenantProjects = append(out.TenantProjects, tenantProject)
	}

	c.JSON(http.StatusOK, out)
}

// TenantProjectCreate adds a new tenant project to the database
// and returns a 201 StatusCreated response.
//
// Route: /tenant/:tenantID/projects
func (s *Server) TenantProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
		out     *api.Project
	)

	// Get the project's tenant ID from the URL and return a 400 response
	// if the tenant ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		log.Error().Err(err).Msg("could not parse tenant ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

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

	tproject := &db.Project{
		TenantID: tenantID,
		Name:     project.Name,
	}

	// Add project to the database and return a 500 response if it cannot be added.
	if err = db.CreateTenantProject(c.Request.Context(), tproject); err != nil {
		log.Error().Err(err).Msg("could not create tenant project in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add tenant project"))
		return
	}

	out = &api.Project{
		ID:   tproject.ID.String(),
		Name: project.Name,
	}

	c.JSON(http.StatusCreated, out)
}

// ProjectList retrieves all projects assigned to an organization
// and returns a 200 OK response.
//
// Route: /projects
func (s *Server) ProjectList(c *gin.Context) {
	var (
		err     error
		project *tokens.Claims
	)

	// Fetch project from the context.
	if project, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch project from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch project from context"))
		return
	}

	// Get project's organization ID and return a 400 response if it is not a ULID.
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(project.OrgID); err != nil {
		log.Error().Err(err).Msg("could not parse org id")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse org id"))
		return
	}

	// Get projects from the database and return a 500 response if not successful.
	var projects []*db.Project
	if projects, err = db.ListProjects(c.Request.Context(), orgID); err != nil {
		log.Error().Err(err).Msg("could not fetch projects from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects from database"))
		return
	}

	// Build the response.
	out := &api.ProjectPage{Projects: make([]*api.Project, 0)}

	//Loop over db.Project and retrieve each project.
	for _, dbProject := range projects {
		project := &api.Project{
			ID:   dbProject.ID.String(),
			Name: dbProject.Name,
		}
		out.Projects = append(out.Projects, project)
	}

	c.JSON(http.StatusOK, out)
}

// ProjectCreate adds a new project to an organization in the database
// and returns a 201 StatusCreated response.
//
// Route: /project
func (s *Server) ProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
		out     *api.Project
	)

	// TODO: Add authentication middleware to fetch the organization ID.

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
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project id cannot be specified on create"))
		return
	}

	// Verify that a project name has been provided and return a 400 response
	// if the project name does not exist.
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project name is required"))
		return
	}

	p := &db.Project{
		Name: project.Name,
	}

	// Add project to the database and return a 500 response if not successful.
	if err = db.CreateProject(c.Request.Context(), p); err != nil {
		log.Error().Err(err).Msg("could not create project in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add project"))
		return
	}

	out = &api.Project{
		ID:   p.ID.String(),
		Name: project.Name,
	}

	c.JSON(http.StatusCreated, out)
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
