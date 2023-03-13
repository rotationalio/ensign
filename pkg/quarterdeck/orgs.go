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
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
)

const OneTimeAccessDuration = 5 * time.Minute

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

	// Initialize prometheus collectors (this function has a sync.Once so it's safe to call more than once)
	metrics.Setup()

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

	// Increment total number of projects in prometheus
	metrics.Projects.WithLabelValues("quarterdeck").Inc()
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

	// Check that the orgID claim is allowed access to the project
	model := &models.OrganizationProject{OrgID: claims.ParseOrgID(), ProjectID: project.ProjectID}
	if exists, err = model.Exists(c.Request.Context()); !exists || err != nil {
		if !exists {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown project id"))
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete project access request"))
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

	// Add only the user permissions related to topics to these claims -- whatever
	// access to topics the user has, so to will the one time access claims.
	for _, permission := range claims.Permissions {
		if permissions.InGroup(permission, permissions.PrefixTopics) {
			ota.Permissions = append(ota.Permissions, permission)
		}
	}

	creds = &api.LoginReply{}
	token := s.tokens.CreateToken(ota)
	if creds.AccessToken, err = s.tokens.Sign(token); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "could not complete project access request")
		return
	}

	c.JSON(http.StatusOK, creds)
}
