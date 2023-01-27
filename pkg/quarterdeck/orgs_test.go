package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

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
