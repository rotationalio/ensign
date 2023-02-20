package tenant

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
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
		out.TenantProjects = append(out.TenantProjects, dbProject.ToAPI())
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
		ctx     context.Context
		claims  *tokens.Claims
		project *api.Project
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// Fetch member from the context.
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch member from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch member from context"))
		return
	}

	// Get the member's orgnaization ID and return a 500 response if it is not a ULID.
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(claims.OrgID); err != nil {
		log.Error().Err(err).Msg("could not parse org id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse org id"))
		return
	}

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
		OrgID:    orgID,
		TenantID: tenantID,
		Name:     project.Name,
	}

	// Create the project in the database and register it with Quarterdeck.
	// TODO: Distinguish between trtl errors and quarterdeck errors.
	if err = s.createProject(ctx, tproject); err != nil {
		log.Error().Err(err).Msg("could not create project")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not create project"))
		return
	}

	c.JSON(http.StatusCreated, tproject.ToAPI())
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
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch project from context"))
		return
	}

	// Get project's organization ID and return a 500 response if it is not a ULID.
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(project.OrgID); err != nil {
		log.Error().Err(err).Msg("could not parse org id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse org id"))
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
		out.Projects = append(out.Projects, dbProject.ToAPI())
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
		ctx     context.Context
		claims  *tokens.Claims
		project *api.Project
	)

	// Fetch project from the context.
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch project from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(err))
		return
	}

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// Get the project's organization ID and return a 500 response if it is not a ULID.
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(claims.OrgID); err != nil {
		log.Error().Err(err).Msg("could not parse org id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse org id"))
		return
	}

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&project); err != nil {
		log.Warn().Err(err).Msg("could not bind project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a project ID does not exist and return a 400 response if it does.
	if project.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project id cannot be specified on create"))
		return
	}

	// Verify that a project name has been provided and return a 400 response if it has not.
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project name is required"))
		return
	}

	dbProject := &db.Project{
		OrgID: orgID,
		Name:  project.Name,
	}

	// Create the project in the database and register it with Quarterdeck.
	// TODO: Distinguish between trtl errors and quarterdeck errors.
	if err = s.createProject(ctx, dbProject); err != nil {
		log.Error().Err(err).Msg("could not create project")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not create project"))
		return
	}

	c.JSON(http.StatusCreated, dbProject.ToAPI())
}

// ProjectDetail retrieves a summary detail of a project by its
// ID and returns a 200 OK response.
//
// Route: /project/:projectID
func (s *Server) ProjectDetail(c *gin.Context) {
	var err error

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

	c.JSON(http.StatusOK, project.ToAPI())
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

	c.JSON(http.StatusOK, p.ToAPI())
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

// createProject is a helper to create a project in the tenant database as well as
// register the orgid - projectid mapping in Quarterdeck in a single step. Any endpoint
// which allows a user to create a project should use this method to ensure that
// Quarterdeck is aware of the project. If this method returns an error, then the
// caller should return an error response to the client to indicate that the project
// creation has failed.
func (s *Server) createProject(ctx context.Context, project *db.Project) (err error) {
	// Create the project in the tenant database
	if err = db.CreateProject(ctx, project); err != nil {
		return err
	}

	// Only the project ID is required - Quarterdeck will extract the org ID from the
	// user claims.
	req := &qd.Project{
		ProjectID: project.ID,
	}

	// See Quarterdeck's ProjectCreate server method for more details
	if _, err = s.quarterdeck.ProjectCreate(ctx, req); err != nil {
		// TODO: Cleanup unused projects or delete them here
		return err
	}

	return nil
}
