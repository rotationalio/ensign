package tenant_test

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (suite *tenantTestSuite) TestProjectDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
	}

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

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{tenant.ReadProjectPermission}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

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
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
	}

	// Marshal the project data with msgpack.
	data, err := project.MarshalValue()
	require.NoError(err, "could not marshal the project")

	// Unmarshal the project data with msgpack.
	other := &db.Project{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the project")

	// OnGet method should return the test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// OnPut method should return a success response.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "invalid"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "invalid"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{tenant.WriteProjectPermission}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the project ID is not parseable.
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "invalid"})
	suite.requireError(err, http.StatusBadRequest, "could not parse project ulid", "expected error when project does not exist")

	// Should return an error if the project name is missing.
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

	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001"})
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project ID is not found")
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

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	err := suite.client.ProjectDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	err = suite.client.ProjectDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{tenant.DeleteProjectPermission}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the project does not exist.
	err = suite.client.ProjectDelete(ctx, "invalid")
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

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err := suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{tenant.WriteTenantPermission}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if tenant id is not a valid ULID.
	_, err = suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})
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

	// Set the initial claims fixture.
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err := suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{tenant.WriteProjectPermission}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if a project ID exists.
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001"})
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
