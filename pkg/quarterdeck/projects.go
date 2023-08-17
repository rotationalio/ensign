package quarterdeck

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/metrics"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *Server) ProjectList(c *gin.Context) {
	var (
		err      error
		orgID    ulid.ULID
		out      *api.ProjectList
		claims   *tokens.Claims
		projects []*models.Project
		prevPage *pagination.Cursor
		nextPage *pagination.Cursor
	)

	// Fetch the user claims from the reuqest
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Error(c).Msg("could not parse orgID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Process the next page token or page query options
	query := &api.PageQuery{}
	if err = c.BindQuery(query); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse api page query")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	if query.NextPageToken != "" {
		if prevPage, err = pagination.Parse(query.NextPageToken); err != nil {
			sentry.Warn(c).Err(err).Str("next_page_token", query.NextPageToken).Msg("could not parse next page token")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	} else {
		prevPage = pagination.New("", "", int32(query.PageSize))
	}

	if projects, nextPage, err = models.ListProjects(c.Request.Context(), orgID, prevPage); err != nil {
		sentry.Error(c).Err(err).Msg("could not list projects in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Prepare and return the response to the user
	out = &api.ProjectList{
		Projects: make([]*api.Project, 0, len(projects)),
	}

	for _, project := range projects {
		out.Projects = append(out.Projects, project.ToAPI())
	}

	if nextPage != nil {
		if out.NextPageToken, err = nextPage.NextPageToken(); err != nil {
			sentry.Error(c).Err(err).Msg("could not create next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) ProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
		claims  *tokens.Claims
	)

	// Bind the Project request to the project data structure
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryProjectAgain))
		return
	}

	// Validate the request from the API side.
	if err = project.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrFixProjectDetails))
		return
	}

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Create the OrganizationProjects mapping database model
	model := &models.OrganizationProject{
		OrgID:     claims.ParseOrgID(),
		ProjectID: project.ProjectID,
	}

	// Save the model to the database
	if err = model.Save(c.Request.Context()); err != nil {
		switch err.(type) {
		case *models.ValidationError:
			c.Error(err)
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryProjectAgain))
		default:
			sentry.Error(c).Err(err).Msg("could not create project")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		}
		return
	}

	// Update the response to send to the user
	project.OrgID = model.OrgID
	project.ProjectID = model.ProjectID
	project.Created, _ = model.GetCreated()
	project.Modified, _ = model.GetModified()

	// Increment total number of projects in prometheus
	metrics.Projects.WithLabelValues(ServiceName).Inc()
	c.JSON(http.StatusOK, project)
}

func (s *Server) ProjectAccess(c *gin.Context) {
	var (
		err     error
		exists  bool
		project *api.Project
		claims  *tokens.Claims
		creds   *api.LoginReply
	)

	// Bind the Project request to the project data structure
	if err = c.BindJSON(&project); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse project access request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryProjectAgain))
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Check that the orgID claim is allowed access to the project
	model := &models.OrganizationProject{OrgID: claims.ParseOrgID(), ProjectID: project.ProjectID}
	if exists, err = model.Exists(c.Request.Context()); !exists || err != nil {
		if !exists {
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryProjectAgain))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve organization project mapping from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Prepare the credentials for the access token to make the one time access.
	now := time.Now()
	ota := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ulids.New().String(),
			Subject:   claims.Subject,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(OneTimeAccessDuration)),
		},
		OrgID:       claims.OrgID,
		ProjectID:   project.ProjectID.String(),
		Permissions: make([]string, 0, 4),
	}

	// Add only the user permissions related to topics and metrics to these claims -- whatever
	// access to topics and metrics the user has, so to will the one time access claims.
	for _, permission := range claims.Permissions {
		if permissions.InGroup(permission, permissions.PrefixTopics) || permissions.InGroup(permission, permissions.PrefixMetrics) {
			ota.Permissions = append(ota.Permissions, permission)
		}
	}

	creds = &api.LoginReply{}
	token := s.tokens.CreateToken(ota)
	if creds.AccessToken, err = s.tokens.Sign(token); err != nil {
		sentry.Error(c).Err(err).Msg("could not sign token")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	c.JSON(http.StatusOK, creds)
}

func (s *Server) ProjectDetail(c *gin.Context) {
	var (
		err       error
		orgID     ulid.ULID
		projectID ulid.ULID
		project   *models.Project
		claims    *tokens.Claims
	)

	// Fetch the user claims from the request
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	if orgID = claims.ParseOrgID(); ulids.IsZero(orgID) {
		sentry.Error(c).Msg("could not parse orgID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Parse the projectID from the URL
	if projectID, err = ulid.Parse(c.Param("id")); err != nil {
		sentry.Warn(c).Err(err).Str("id", c.Param("id")).Msg("could not parse project id from url")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
		return
	}

	// Fetch the project from the database
	if project, err = models.FetchProject(c.Request.Context(), projectID, orgID); err != nil {
		// Check if the error is a not found error.
		if errors.Is(err, models.ErrNotFound) {
			c.Error(err)
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrProjectNotFound))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve project from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Return the API response to the user
	c.JSON(http.StatusOK, project.ToAPI())
}
