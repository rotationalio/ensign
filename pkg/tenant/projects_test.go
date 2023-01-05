package tenant_test

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (suite *tenantTestSuite) TestTenantProjectList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(cr *pb.CursorRequest, t pb.Trtl_CursorServer) error {
		return nil
	}

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	// Should return an error if the topic does not exist.
	_, err := suite.client.TenantProjectList(ctx, "invalid", req)
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant ulid", "expected error when tenant does not exist")

	projects, err := suite.client.TenantProjectList(ctx, "01GKKYAWC4PA72YC53RVXAEC67", req)
	require.NoError(err, "could not list tenant projects")
	require.Len(projects.TenantProjects, 1, "expected one project in the database")
}

func (suite *tenantTestSuite) TestProjectList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(cr *pb.CursorRequest, t pb.Trtl_CursorServer) error {
		return nil
	}

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	// TODO: Test length of values assigned to *api.ProjectPage
	_, err := suite.client.ProjectList(ctx, req)
	require.NoError(err, "could not list projects")
}

func (suite *tenantTestSuite) TestProjectDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
	}

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the project data with msgpack.
	data, err := project.MarshalValue()
	require.NoError(err, "could not marshal the project")

	// Unmarshal the project data with msgpack.
	other := &db.Project{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the project")

	// Call the OnGet method and return test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Should return an error if the project does not exist.
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse project ulid", "expected error when project does not exist")

	// Create a project test fixture.
	req := &api.Project{
		ID:   "01GKKYAWC4PA72YC53RVXAEC67",
		Name: "project-example",
	}

	rep, err := suite.client.ProjectDetail(ctx, req.ID)
	require.NoError(err, "could not retrieve project")
	require.Equal(req.ID, rep.ID, "expected project id to match")
	require.Equal(req.Name, rep.Name, "expected project name to match")

	// Should return an error if the project ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, errors.New("key not found")
	}

	_, err = suite.client.ProjectDetail(ctx, "01GKKYAWC4PA72YC53RVXAEC67")
	suite.requireError(err, http.StatusNotFound, "could not retrieve project", "expected error when project ID is not found")
}

func (suite *tenantTestSuite) TestProjectUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
	}

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the project data with msgpack.
	data, err := project.MarshalValue()
	require.NoError(err, "could not marshal the project")

	// Unmarshal the project data with msgpack.
	other := &db.Project{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the project")

	// Call the OnGet method and return the test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Should return an error if the project does not exist.
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse project ulid", "expected error when project does not exist")

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if the project name does not exist.
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67"})
	suite.requireError(err, http.StatusBadRequest, "project name is required", "expected error when project name does not exist")

	req := &api.Project{
		ID:   "01GKKYAWC4PA72YC53RVXAEC67",
		Name: "project-example",
	}

	rep, err := suite.client.ProjectUpdate(ctx, req)
	require.NoError(err, "could not update project")
	require.NotEqual(req.ID, "01GMTWFK4XZY597Y128KXQ4WHP", "project id should not match")
	require.Equal(rep.Name, req.Name, "expected project name to match")

	// Should return an error if the project ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, errors.New("key not found")
	}

	_, err = suite.client.ProjectDetail(ctx, "01GKKYAWC4PA72YC53RVXAEC67")
	suite.requireError(err, http.StatusNotFound, "could not retrieve project", "expected error when project ID is not found")
}

func (suite *tenantTestSuite) TestProjectDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projectID := "01GKKYAWC4PA72YC53RVXAEC67"

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{}, nil
	}

	// Should return an error if the project does not exist.
	err := suite.client.ProjectDelete(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse project ulid", "expected error when project does not exist")

	err = suite.client.ProjectDelete(ctx, projectID)
	require.NoError(err, "could not delete project")

	// Should return an error if the project ID is parsed but not found.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return nil, errors.New("key not found")
	}

	err = suite.client.ProjectDelete(ctx, "01GKKYAWC4PA72YC53RVXAEC67")
	suite.requireError(err, http.StatusNotFound, "could not delete project", "expected error when project ID is not found")
}
