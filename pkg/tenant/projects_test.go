package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *tenantTestSuite) TestTenantProjectList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
		Created:  time.Unix(1670424445, 0).In(time.UTC),
		Modified: time.Unix(1670424445, 0).In(time.UTC),
	}

	prefix := project.TenantID[:]
	namespace := "projects"

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i := 0; i < 7; i++ {
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     []byte(fmt.Sprintf("value %d", i)),
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	// Should return an error if the tenant does not exist.
	_, err := suite.client.TenantProjectList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant ulid", "expected error when tenant does not exist")

	rep, err := suite.client.TenantProjectList(ctx, project.TenantID.String(), &api.PageQuery{})
	require.NoError(err, "could not list tenant projects")
	require.Len(rep.TenantProjects, 7, "expected 7 projects")
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

	// TODO: Test length of values assigned to *api.ProjectPage
	_, err := suite.client.ProjectList(ctx, &api.PageQuery{})
	require.NoError(err, "could not list projects")
}

func (suite *tenantTestSuite) TestProjectDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
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
		Name: "project001",
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
		Name:     "project001",
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
		Name: "project001",
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

func (suite *tenantTestSuite) TestTenantProjectCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.Make().String()

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if tenant id is not a valid ULID.
	_, err := suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant id", "expected error when tenant id does not exist")

	// Should return an error if the project ID exists.
	_, err = suite.client.TenantProjectCreate(ctx, tenantID, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001"})
	suite.requireError(err, http.StatusBadRequest, "project id cannot be specified on create", "expected error when project id exists")

	// Should return an error if the project name does not exist.
	_, err = suite.client.TenantProjectCreate(ctx, tenantID, &api.Project{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "project name is required", "expected error when project name does not exist")

	// Create a project test fixture.
	req := &api.Project{
		Name: "project001",
	}

	project, err := suite.client.TenantProjectCreate(ctx, tenantID, req)
	require.NoError(err, "could not add project")
	require.Equal(req.Name, project.Name, "project name should match")
}

func (suite *tenantTestSuite) TestProjectCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if a project ID exists.
	_, err := suite.client.ProjectCreate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001"})
	suite.requireError(err, http.StatusBadRequest, "project id cannot be specified on create", "expected error when project id exists")

	// Should return an error if a project name does not exist.
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "project name is required", "expected error when project name does not exist")

	// Create a project test fixture.
	req := &api.Project{
		Name: "project001",
	}

	project, err := suite.client.ProjectCreate(ctx, req)
	require.NoError(err, "could not add project")
	require.Equal(req.Name, project.Name)
}
