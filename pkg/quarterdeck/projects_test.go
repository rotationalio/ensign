package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/responses"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *quarterdeckTestSuite) TestProjectList() {
	require := s.Require()
	ctx := context.Background()

	// Listing Projects requires authentication
	req := &api.PageQuery{}
	_, err := s.client.ProjectList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Listing Projects requires the projects:read permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.ProjectList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	claims.OrgID = "01GQFQ14HXF2VC7C1HJECS60XX"
	claims.Permissions = []string{perms.ReadProjects}
	ctx = s.AuthContext(ctx, claims)

	// Should be able to list all projects for the specified organization
	page, err := s.client.ProjectList(ctx, req)
	require.NoError(err, "could not fetch projects list")
	require.Len(page.Projects, 10, "expected 10 projects back from the fixtures")
	require.Empty(page.NextPageToken, "expected no next page token in response")

	// Should be able to paginate the request
	req.PageSize = 8
	page, err = s.client.ProjectList(ctx, req)
	require.NoError(err, "could not fetch projects list")
	require.Len(page.Projects, 8, "expected 8 projects back from the fixtures")
	require.NotEmpty(page.NextPageToken, "expected next page token in response")

	// Test fetching the next page with the next page token
	req.NextPageToken = page.NextPageToken
	page2, err := s.client.ProjectList(ctx, req)
	require.NoError(err, "Could not fetch projects list second page")
	require.Empty(page2.NextPageToken, "expected no more pages")
	require.Len(page2.Projects, 2, "expected 1 project back from fixtures")
	require.NotEqual(page.Projects[0].ProjectID, page2.Projects[0].ProjectID)
}

func (s *quarterdeckTestSuite) TestProjectCreate() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Creating a project requires authentication
	req := &api.Project{ProjectID: ulids.New()}
	_, err := s.client.ProjectCreate(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Creating an Project requires the projects:edit permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.ProjectCreate(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	id := ulid.MustParse("01GKHJSK7CZW0W282ZN3E9W86Z")
	claims.Subject = id.String()
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	claims.Permissions = []string{perms.EditProjects}
	ctx = s.AuthContext(ctx, claims)

	// Test Happy Path
	rep, err := s.client.ProjectCreate(ctx, req)
	require.NoError(err, "could not create project after authentication")
	require.NotEmpty(rep, "expected a project response from the server")

	// Validate the response returned by the server
	require.False(ulids.IsZero(rep.OrgID), "no orgID returned in response")
	require.Equal(req.ProjectID, rep.ProjectID, "expected project id to be identical in response")
	require.Equal(id, rep.Owner.ID, "expected owner id to be identical in response")
	require.Equal(claims.Name, rep.Owner.Name, "expected owner name to be identical in response")
	require.Equal(claims.Email, rep.Owner.Email, "expected owner email to be identical in response")
	require.False(rep.Created.IsZero(), "no created returned in response")
	require.False(rep.Modified.IsZero(), "no modified returned in response")

	// Must specify a projectID
	_, err = s.client.ProjectCreate(ctx, &api.Project{})
	s.CheckError(err, http.StatusBadRequest, responses.ErrFixProjectDetails)

	// Cannot specify an orgID
	_, err = s.client.ProjectCreate(ctx, &api.Project{OrgID: ulids.New(), ProjectID: ulids.New()})
	s.CheckError(err, http.StatusBadRequest, responses.ErrFixProjectDetails)
}

func (s *quarterdeckTestSuite) TestProjectAccess() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Requesting one-time access to a project requires authentication
	req := &api.Project{ProjectID: ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT")}
	_, err := s.client.ProjectAccess(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Requesting one-time access to a project requires the topics:read permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.ProjectAccess(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	claims.Permissions = []string{perms.ReadAPIKeys, perms.ReadTopics, perms.DeleteAPIKeys, perms.CreateTopics, perms.EditAPIKeys}
	ctx = s.AuthContext(ctx, claims)

	// Test Happy Path
	rep, err := s.client.ProjectAccess(ctx, req)
	require.NoError(err, "could not create project after authentication")
	require.NotEmpty(rep, "expected a project response from the server")

	// Validate the response returned by the server
	require.NotEmpty(rep.AccessToken, "no access token returned in response")
	require.Empty(rep.RefreshToken, "no refresh token should have been returned in the response")
	require.Empty(rep.LastLogin, "no last login timestamp should have been returned in the response")

	// Validate the claims returned by the server
	ota := &tokens.Claims{}
	parser := &jwt.Parser{SkipClaimsValidation: true}

	_, _, err = parser.ParseUnverified(rep.AccessToken, ota)
	require.NoError(err, "could not parse access token")

	require.NotEqual(claims.ID, ota.ID)
	require.Equal(claims.Subject, ota.Subject)
	require.Equal(claims.OrgID, ota.OrgID)
	require.Equal(req.ProjectID.String(), ota.ProjectID)
	require.Equal([]string{perms.ReadTopics, perms.CreateTopics}, ota.Permissions)
	require.Greater(time.Until(ota.ExpiresAt.Time), 1*time.Minute)
	require.Less(time.Until(ota.ExpiresAt.Time), 10*time.Minute)

	// Must specify a projectID
	_, err = s.client.ProjectAccess(ctx, &api.Project{})
	s.CheckError(err, http.StatusBadRequest, "missing required field: project_id")

	// Cannot specify an orgID
	_, err = s.client.ProjectAccess(ctx, &api.Project{OrgID: ulids.New(), ProjectID: ulids.New()})
	s.CheckError(err, http.StatusBadRequest, "field restricted for request: org_id")

	// Must specify a projectID that belongs to the organization
	_, err = s.client.ProjectAccess(ctx, &api.Project{ProjectID: ulid.MustParse("01GQFQCFC9P3S7QZTPYFVBJD7F")})
	s.CheckError(err, http.StatusBadRequest, responses.ErrTryProjectAgain)
}

func (s *quarterdeckTestSuite) TestProjectDetail() {
	require := s.Require()
	ctx := context.Background()

	// Retrieving a Project requires authentication
	_, err := s.client.ProjectDetail(ctx, "01GQFQCFC9P3S7QZTPYFVBJD7F")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Retrieving an API Key requires the apikeys:read permission
	claims := &tokens.Claims{
		Name:  "Tom Riddle",
		Email: "voldy@example.com",
		OrgID: ulids.New().String(),
	}

	ctx = s.AuthContext(ctx, claims)
	project, err := s.client.ProjectDetail(ctx, "01GQFQCFC9P3S7QZTPYFVBJD7F")
	require.Nil(project, "no reply should be returned")
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Cannot retrieve a key that is not in the same organization
	claims.Permissions = []string{perms.ReadProjects}
	ctx = s.AuthContext(ctx, claims)

	project, err = s.client.ProjectDetail(ctx, "01GQFQCFC9P3S7QZTPYFVBJD7F")
	require.Nil(project, "no project should be returned")
	s.CheckError(err, http.StatusNotFound, responses.ErrProjectNotFound)

	// Test happy path and fetch project
	claims = &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01GKHJSK7CZW0W282ZN3E9W86Z",
		},
		Name:        "Edison Edgar Franlin",
		Email:       "eefrank@checkers.io",
		OrgID:       "01GQFQ14HXF2VC7C1HJECS60XX",
		Permissions: []string{perms.ReadProjects},
	}
	ctx = s.AuthContext(ctx, claims)

	project, err = s.client.ProjectDetail(ctx, "01GQFQCFC9P3S7QZTPYFVBJD7F")
	require.NoError(err, "could not fetch project")

	require.Equal("01GQFQ14HXF2VC7C1HJECS60XX", project.OrgID.String())
	require.Equal("01GQFQCFC9P3S7QZTPYFVBJD7F", project.ProjectID.String())
	require.Equal(3, project.APIKeysCount)
	require.Equal(0, project.RevokedCount)
	require.Equal("2023-01-23T16:22:54Z", project.Created.Format(time.RFC3339))
	require.Equal("2023-01-23T16:22:54Z", project.Modified.Format(time.RFC3339))

	// Test cannot parse ULID returns not found
	_, err = s.client.ProjectDetail(ctx, "notaulid")
	s.CheckError(err, http.StatusNotFound, responses.ErrProjectNotFound)

	// Test database not found
	_, err = s.client.ProjectDetail(ctx, ulids.New().String())
	s.CheckError(err, http.StatusNotFound, responses.ErrProjectNotFound)
}
