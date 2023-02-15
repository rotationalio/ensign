package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

func (s *quarterdeckTestSuite) TestOrganizationDetail() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Retrieving an organization requires authentication
	_, err := s.client.OrganizationDetail(ctx, "invalid")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Retrieving an organization requires the read:organizations permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)
	orgID := ulids.New()
	_, err = s.client.OrganizationDetail(ctx, orgID.String())
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Valid ID is required in the URL
	claims.Permissions = []string{perms.ReadOrganizations}
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, "invalid")
	s.CheckError(err, http.StatusNotFound, "organization not found")

	// Specified organization must match the user's organization
	claims.OrgID = ulids.New().String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, orgID.String())
	s.CheckError(err, http.StatusForbidden, "user is not authorized to access this organization")

	// Organization must exist
	claims.OrgID = orgID.String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, orgID.String())
	s.CheckError(err, http.StatusNotFound, "organization not found")

	// Successfully retrieving the organization
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, claims.OrgID)
	require.NoError(err, "could not retrieve organization details")
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
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
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
	require.False(rep.Created.IsZero(), "no created returned in response")
	require.False(rep.Modified.IsZero(), "no modified returned in response")

	// Must specify a projectID
	_, err = s.client.ProjectCreate(ctx, &api.Project{})
	s.CheckError(err, http.StatusBadRequest, "missing required field: project_id")

	// Cannot specify an orgID
	_, err = s.client.ProjectCreate(ctx, &api.Project{OrgID: ulids.New(), ProjectID: ulids.New()})
	s.CheckError(err, http.StatusBadRequest, "field restricted for request: org_id")
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

	// Requesting one-time access to a project an Project requires the topics:read permission
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
	s.CheckError(err, http.StatusBadRequest, "unknown project id")
}
