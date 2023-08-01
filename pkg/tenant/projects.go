package tenant

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	pb "github.com/rotationalio/go-ensign/api/v1beta1"
	mt "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"github.com/rs/zerolog/log"
)

const maxQueryResults = 10

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
		// Ensure the project owner info is populated
		// TODO: Use a member cache to avoid multiple DB calls
		var owner *db.Member
		if owner, err = dbProject.Owner(c.Request.Context()); err != nil {
			sentry.Error(c).Err(err).Str("member_id", dbProject.OwnerID.String()).Msg("could not fetch project owner info")
			continue
		}

		// Return only the fields that are required for list
		// TODO: Return data storage, which should have units
		project := &api.Project{
			ID:          dbProject.ID.String(),
			Name:        dbProject.Name,
			Description: dbProject.Description,
			Owner: api.Member{
				ID:      owner.ID.String(),
				Name:    owner.Name,
				Picture: owner.Picture(),
			},
			Status:       dbProject.Status(),
			ActiveTopics: dbProject.Topics,
			DataStorage: api.StatValue{
				Value: 0,
				Units: "GB",
			},
			Created: db.TimeToString(dbProject.Created),
		}
		out.TenantProjects = append(out.TenantProjects, project)
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
		claims  *tokens.Claims
		project *api.Project
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Get the user claims to populate the owner info
	if claims, err = middleware.GetClaims(c); err != nil {
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
		if len(project.Description) > db.MaxDescriptionLength {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("project description is too long"))
			return
		}

		tproject.Description = project.Description
	}

	if err = tproject.SetOwnerFromClaims(claims); err != nil {
		sentry.Error(c).Err(err).Msg("could not set project owner from user claims")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
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

// TenantProjectPatch applies a partial update to a project identified by the tenantID
// and projectID in the URL.
//
// Route: /tenant/:tenantID/projects/:projectID
func (s *Server) TenantProjectPatch(c *gin.Context) {
	var (
		err                        error
		orgID, tenantID, projectID ulid.ULID
	)

	// Get the orgID from the claims and handle errors
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Parse the tenantID from the URL
	tenantParam := c.Param("tenantID")
	if tenantID, err = ulid.Parse(tenantParam); err != nil {
		sentry.Warn(c).Err(err).Str("tenantID", tenantParam).Msg("could not parse tenant id from URL")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrTenantNotFound))
		return
	}

	// Parse the projectID from the URL
	projectParam := c.Param("projectID")
	if projectID, err = ulid.Parse(projectParam); err != nil {
		sentry.Warn(c).Err(err).Str("projectID", projectParam).Msg("could not parse project id from URL")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
		return
	}

	// Verify that the tenant exists in the user's organization.
	if err = db.VerifyOrg(c, orgID, tenantID); err != nil {
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, db.ErrOrgNotVerified) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrTenantNotFound))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Patch the project, this method handles the rest of the errors and responses.
	s.projectPatch(c, orgID, projectID)
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
		// Return only the fields that are required for list
		// TODO: Return data storage, which should have units
		project := dbProject.ToAPI()

		// Ensure the project owner info is populated
		// TODO: Use a member cache to avoid multiple DB calls
		var owner *db.Member
		if owner, err = dbProject.Owner(c.Request.Context()); err != nil {
			sentry.Error(c).Err(err).Str("member_id", dbProject.OwnerID.String()).Msg("could not fetch project owner info")
		}

		if owner != nil {
			project.Owner = api.Member{
				ID:      owner.ID.String(),
				Name:    owner.Name,
				Picture: owner.Picture(),
			}
		}

		out.Projects = append(out.Projects, project)
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
		claims  *tokens.Claims
		project *api.Project
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// Get the user claims to populate the owner info
	if claims, err = middleware.GetClaims(c); err != nil {
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

	if err = dbProject.SetOwnerFromClaims(claims); err != nil {
		sentry.Error(c).Err(err).Msg("could not set owner info from user claims")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
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

	// Ensure the project owner info is populated
	if _, err = project.Owner(c.Request.Context()); err != nil {
		sentry.Error(c).Err(err).Str("member_id", project.OwnerID.String()).Msg("could not retrieve project owner from database")
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

	if project.Description != "" {
		if len(project.Description) > db.MaxDescriptionLength {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("project description is too long"))
			return
		}

		p.Description = project.Description
	}

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

// ProjectPatch applies a partial update to a project identified by the project ID in
// the URL and the tenant ID in the request body.
//
// Route: /project/:projectID
func (s *Server) ProjectPatch(c *gin.Context) {
	var (
		err              error
		orgID, projectID ulid.ULID
	)

	// Get the orgID from the claims and handle errors
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Parse the projectID from the URL
	projectParam := c.Param("projectID")
	if projectID, err = ulid.Parse(projectParam); err != nil {
		sentry.Warn(c).Err(err).Str("projectID", projectParam).Msg("could not parse project id from URL")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
		return
	}

	// Patch the project, this method handles the rest of the responses and errors.
	s.projectPatch(c, orgID, projectID)
}

// projectPatch is a handler for multiple patch project routes that applies a partial
// update to a project. Only the provided fields are updated, and an error is returned
// if any field does not exist or if no updates were made.
func (s *Server) projectPatch(c *gin.Context, orgID, projectID ulid.ULID) {
	var (
		err     error
		project *db.Project
		patch   map[string]interface{}
	)

	// Parse the patch fields from the request body
	if err = c.BindJSON(&patch); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project patch request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Verify that the project exists in the user's organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, db.ErrOrgNotVerified) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Retrieve the project from the database.
	if project, err = db.RetrieveProject(c, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
			return
		}
		sentry.Error(c).Err(err).Msg("could not retrieve project")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Iterate over the patch fields and perform the updates.
	for field, value := range patch {
		var ok bool
		switch field {
		case "name":
			var name string
			if name, ok = value.(string); !ok {
				c.JSON(http.StatusBadRequest, api.ErrorResponse(api.FieldTypeError(field, "string")))
				return
			}

			project.Name = name
		case "description":
			var description string
			if description, ok = value.(string); !ok {
				c.JSON(http.StatusBadRequest, api.ErrorResponse(api.FieldTypeError(field, "string")))
				return
			}

			project.Description = description
		case "owner":
			// Parse the owner field, which should be a member object.
			var owner map[string]interface{}
			if owner, ok = value.(map[string]interface{}); !ok {
				c.JSON(http.StatusBadRequest, api.ErrorResponse(api.FieldTypeError(field, "member object")))
				return
			}

			// Ignore owner ID if not provided.
			var id interface{}
			if id, ok = owner["id"]; !ok {
				continue
			}

			var ownerID ulid.ULID
			if ownerID, err = ulids.Parse(id); err != nil {
				log.Warn().Err(err).Msg("could not parse owner as ULID")
				c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrMemberNotFound))
				return
			}

			// Ensure that the new owner is a member of the organization.
			if _, err = db.RetrieveMember(c.Request.Context(), orgID, ownerID); err != nil {
				if errors.Is(err, db.ErrNotFound) {
					c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrMemberNotFound))
					return
				}
				log.Error().Err(err).Msg("could not retrieve organization member")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
				return
			}

			project.OwnerID = ownerID
		default:
			// Ignore other nested fields for now to appease the Go client
			if _, ok = value.(map[string]interface{}); ok {
				continue
			}
			c.JSON(http.StatusBadRequest, api.ErrorResponse(api.InvalidFieldError(field)))
			return
		}
	}

	// Update the project in the database
	// TODO: Return better errors to the user
	if err = db.UpdateProject(c.Request.Context(), project); err != nil {
		switch err.(type) {
		case *db.ValidationError:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err.Error()))
		default:
			sentry.Error(c).Err(err).Msg("could not update project")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}
		return
	}

	c.JSON(http.StatusOK, project.ToAPI())
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

// ProjectQuery executes simple queries on topics in a project using enSQL. This
// endpoint forwards the query to Ensign and a limited number of results to the client.
// Clients that require more results or complex queries should use the SDKs instead.
//
// Route: /projects/:projectID/query
func (s *Server) ProjectQuery(c *gin.Context) {
	var (
		err       error
		orgID     ulid.ULID
		projectID ulid.ULID
		in        *api.ProjectQueryRequest
	)

	// orgID is required to check ownership of the project
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	if err = c.BindJSON(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Get the project ID from the URL
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Warn().Err(err).Str("projectID", c.Param("projectID")).Msg("could not parse project id")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
		return
	}

	// Verify that the project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, db.ErrOrgNotVerified) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	in.Query = strings.TrimSpace(in.Query)
	if in.Query == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrMissingQueryField))
		return
	}

	if len(in.Query) > 2000 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrQueryTooLong))
		return
	}

	// Build the response with the query results
	out := &api.ProjectQueryResponse{
		Results: make([]*api.QueryResult, 0),
	}

	// TODO: Query events from Ensign rather than using the fixtures
	events := fixtureEvents()
	for _, event := range events {
		result := &api.QueryResult{
			Metadata: event.Metadata,
			Mimetype: event.Mimetype.MimeType(),
			Version:  fmt.Sprintf("%s v%d.%d.%d", event.Type.Name, event.Type.MajorVersion, event.Type.MinorVersion, event.Type.PatchVersion),
			Created:  event.Created.String(),
		}

		// Attempt to encode the event data.
		if result.Data, result.IsBase64Encoded, err = encodeToString(event.Data, event.Mimetype); err != nil {
			result.Data = "could not encode event data to string"
		}

		out.Results = append(out.Results, result)
		if len(out.Results) >= maxQueryResults {
			break
		}
	}

	c.JSON(http.StatusOK, out)
}

// Encode event data into a string for the response. Returns true if the data was
// base64 encoded.
func encodeToString(data []byte, mime mt.MIME) (encoded string, isBase64Encoded bool, err error) {
	switch mime {
	case mt.TextPlain, mt.TextCSV, mt.TextHTML:
		return string(data), false, nil
	case mt.ApplicationJSON:
		if encoded, err = responses.RemarshalJSON(data); err != nil {
			return "", false, err
		}
		return encoded, false, nil
	default:
		return base64.StdEncoding.EncodeToString(data), true, nil
	}
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

// UpdateProjectStats updates the stat fields on a project by performing readonly
// queries to Quarterdeck and Ensign. Because this requires a few RPCs, it should be
// called in a background task where possible to avoid blocking user requests. The
// context passed to this method must contain authentication credentials in order to
// query Quarterdeck which must include the topics:read and projects:read permissions.
// TODO: This data can be updated asynchronously once the Ensign "meta" topics are up
// and running.
func (s *Server) UpdateProjectStats(ctx context.Context, userID, projectID ulid.ULID) (err error) {
	// Go routine to fetch the number of project API keys from Quarterdeck.
	countAPIKeys := func(ctx context.Context) (_ uint64, err error) {
		var reply *qd.Project
		if reply, err = s.quarterdeck.ProjectDetail(ctx, projectID.String()); err != nil {
			return 0, err
		}

		return uint64(reply.APIKeysCount), nil
	}

	// Go routine to fetch the number of topics from Ensign.
	countTopics := func(ctx context.Context) (_ uint64, err error) {
		// Get the one-time access token to make the Ensign request.
		var token string
		if token, err = s.EnsignProjectToken(ctx, userID, projectID); err != nil {
			return 0, err
		}

		// Make the Ensign request with the one-time access token.
		var project *pb.ProjectInfo
		if project, err = s.ensign.InvokeOnce(token).Info(ctx); err != nil {
			return 0, err
		}

		// Count the number of non-readonly topics.
		if project.ReadonlyTopics > project.Topics {
			return 0, nil
		}

		return project.Topics - project.ReadonlyTopics, nil
	}

	// Run both RPCs in parallel and capture errors.
	var wg sync.WaitGroup
	errs := make([]error, 2)
	res := make([]uint64, 2)
	rpc := func(fn func(context.Context) (uint64, error), i int) {
		defer wg.Done()
		res[i], errs[i] = fn(ctx)
	}
	wg.Add(2)
	go rpc(countAPIKeys, 0)
	go rpc(countTopics, 1)
	wg.Wait()

	// Return if both RPCs errored.
	merr := multierror.Append(err, errs...)
	if len(merr.Errors) == len(errs) {
		return merr
	}

	// Retrieve the project from trtl.
	var project *db.Project
	if project, err = db.RetrieveProject(ctx, projectID); err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}

	// TODO: A project write in between here could cause updates to be stomped (e.g. if
	// the project name is updated it will be overwritten here).

	// Update the project stats.
	if errs[0] == nil {
		project.APIKeys = res[0]
	}
	if errs[1] == nil {
		project.Topics = res[1]
	}
	if err = db.UpdateProject(ctx, project); err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}

	return merr.ErrorOrNil()
}
