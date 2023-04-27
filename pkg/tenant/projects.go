package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

// TenantProjectList retrieves projects assigned to a specified
// tenant and returns a 200 OK response.
//
// Route: /tenant/:tenantID/projects
func (s *Server) TenantProjectList(c *gin.Context) {
	var (
		err        error
		next, prev *pg.Cursor
	)

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

	// Get the project's tenant ID from the URL and return a 404 response
	// if the tenant ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", c.Param("tenantID")).Msg("could not parse tenant id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Get projects from the database
	var projects []*db.Project
	if projects, next, err = db.ListProjects(c.Request.Context(), tenantID, prev); err != nil {
		sentry.Error(c).Err(err).Msg("could not list projects from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list projects"))
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

	if next != nil {
		if out.NextPageToken, err = next.NextPageToken(); err != nil {
			log.Error().Err(err).Msg("could not set next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list projects"))
			return
		}
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
		project *api.Project
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to create the project
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the project's tenant ID from the URL and return a 400 response
	// if the tenant ID is not a ULID.
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(c.Param("tenantID")); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", c.Param("tenantID")).Msg("could not parse tenant id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
		return
	}

	// Verify tenant exists in the organization.
	if err = db.VerifyOrg(c, orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
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

	if project.Description != "" {
		tproject.Description = project.Description
	}

	// Create the project in the database and register it with Quarterdeck.
	// TODO: Distinguish between trtl errors and quarterdeck errors.
	// TODO: it is now even more important to distinguish between these errors!
	if err = s.createProject(ctx, tproject); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	c.JSON(http.StatusCreated, tproject.ToAPI())
}

// ProjectList retrieves projects assigned to a specified
// organization and returns a 200 OK response.
//
// Route: /projects
func (s *Server) ProjectList(c *gin.Context) {
	var (
		err        error
		orgID      ulid.ULID
		next, prev *pg.Cursor
	)

	// org ID is required to list the projects
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

	// Get projects from the database.
	var projects []*db.Project
	if projects, next, err = db.ListProjects(c.Request.Context(), orgID, prev); err != nil {
		sentry.Error(c).Err(err).Msg("could not list projects from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list projects"))
		return
	}

	// Build the response.
	out := &api.ProjectPage{Projects: make([]*api.Project, 0)}

	//Loop over db.Project and retrieve each project.
	for _, dbProject := range projects {
		out.Projects = append(out.Projects, dbProject.ToAPI())
	}

	if next != nil {
		if out.NextPageToken, err = next.NextPageToken(); err != nil {
			log.Error().Err(err).Msg("could not set next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list projects"))
			return
		}
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
		project *api.Project
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to create the project
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Bind the user request and return a 400 response if binding
	// is not successful.
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
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

	// Parse the tenant ID from the request
	var tenantID ulid.ULID
	if tenantID, err = ulid.Parse(project.TenantID); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse tenant ID in project")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse tenant id"))
		return
	}

	// Verify tenant exists in the organization.
	if err = db.VerifyOrg(c, orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("tenant not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	dbProject := &db.Project{
		OrgID:    orgID,
		TenantID: tenantID,
		Name:     project.Name,
	}

	// Create the project in the database and register it with Quarterdeck.
	// TODO: Distinguish between trtl errors and quarterdeck errors.
	if err = s.createProject(ctx, dbProject); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dbProject.ToAPI())
}

// ProjectDetail retrieves a summary detail of a project by its
// ID and returns a 200 OK response.
//
// Route: /project/:projectID
func (s *Server) ProjectDetail(c *gin.Context) {
	var (
		err   error
		orgID ulid.ULID
	)

	// orgID is required to check ownership of the project
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the project ID from the URL and return a 400 response
	// if the project does not exist.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Warn().Err(err).Str("projectID", c.Param("projectID")).Msg("could not parse project id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Get the specified project from the database
	var project *db.Project
	if project, err = db.RetrieveProject(c.Request.Context(), projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve project from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve project"))
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
		orgID   ulid.ULID
	)

	// orgID is required to check ownership of the project
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the project ID from the URL and return a 400 response if
	// the project ID is not a ULID.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Warn().Err(err).Str("projectID", c.Param("projectID")).Msg("could not parse project id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Bind the user request with JSON and return a 400 response
	// if binding is not successful.
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify that the project name exists and return a 400 response if it doesn't.
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("project name is required"))
		return
	}

	// Get the specified project from the database
	var p *db.Project
	if p, err = db.RetrieveProject(c.Request.Context(), projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve project from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update project"))
		return
	}

	// Update all user provided fields
	p.Name = project.Name

	// Update project in the database
	if err = db.UpdateProject(c.Request.Context(), p); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not update project in database")
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
		err   error
		orgID ulid.ULID
	)

	// orgID is required to check ownership of the project
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the project ID from the URL and return a 400 response
	// if the project does not exist.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Warn().Err(err).Str("projectID", c.Param("projectID")).Msg("could not parse project id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Delete the project from the database
	if err = db.DeleteProject(c.Request.Context(), projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not delete project from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete project"))
		return
	}
	c.Status(http.StatusNoContent)
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
