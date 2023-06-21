package tenant_test

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *tenantTestSuite) TestTenantList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	defer cancel()

	tenants := []*db.Tenant{
		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GQ38QWNR7MYQXSQ682PJQM7T"),
			Name:            "tenant001",
			EnvironmentType: "prod",
			Created:         time.Unix(1668660681, 0),
			Modified:        time.Unix(1668661302, 0),
		},

		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GQ38QMW7FGKG7AN1TVJTGHJA"),
			Name:            "tenant002",
			EnvironmentType: "staging",
			Created:         time.Unix(1673659941, 0),
			Modified:        time.Unix(1673659941, 0),
		},

		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GQ38QBN8XYA2S0KTW8AHPXHR"),
			Name:            "tenant003",
			EnvironmentType: "dev",
			Created:         time.Unix(1674073941, 0),
			Modified:        time.Unix(1674073941, 0),
		},
		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GVZQ12C17KB13GJP0490T23H"),
			Name:            "tenant004",
			EnvironmentType: "prod",
			Created:         time.Unix(1674073941, 0),
			Modified:        time.Unix(1674073941, 0),
		},
		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GVZQ6C2DET9C5QJYWN85RE3Z"),
			Name:            "tenant005",
			EnvironmentType: "staging",
			Created:         time.Unix(1674073941, 0),
			Modified:        time.Unix(1674073941, 0),
		},
	}

	prefix := orgID[:]
	namespace := "tenants"

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, tenant := range tenants {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, tenant.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := tenant.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       tenant.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	req := &api.PageQuery{}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMBVR86186E0EKCHQK4ESJB1",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := suite.client.TenantList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have the correct permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadOrganizations}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Retrieve all tenants.
	rep, err := suite.client.TenantList(ctx, req)
	require.NoError(err, "could not list tenants")
	require.Len(rep.Tenants, 5, "expected 5 tenants")
	require.Empty(rep.NextPageToken, "next page token should not be set since there isn't a next page")

	// Verify tenant data has been populated.
	for i := range tenants {
		require.Equal(tenants[i].ID.String(), rep.Tenants[i].ID, "tenant id should match")
		require.Equal(tenants[i].Name, rep.Tenants[i].Name, "tenant name should match")
		require.Equal(tenants[i].EnvironmentType, rep.Tenants[i].EnvironmentType, "tenant environment type should match")
		require.Equal(tenants[i].Created.Format(time.RFC3339Nano), rep.Tenants[i].Created, "tenant created timestamp should match")
		require.Equal(tenants[i].Modified.Format(time.RFC3339Nano), rep.Tenants[i].Modified, "tenant modified timestamp should match")
	}

	// Set page size and test pagination.
	req.PageSize = 2
	rep, err = suite.client.TenantList(ctx, req)
	require.NoError(err, "could not list tenants")
	require.Len(rep.Tenants, 2, "expected 2 tenants")
	require.NotEmpty(rep.NextPageToken, "next page token expected")

	// Test next page token.
	req.NextPageToken = rep.NextPageToken
	rep2, err := suite.client.TenantList(ctx, req)
	require.NoError(err, "could not list tenants")
	require.Len(rep2.Tenants, 2, "expected 2 tenants")
	require.NotEqual(rep.Tenants[0].ID, rep2.Tenants[0].ID, "should not have same tenant ID")
	require.NotEqual(rep.Tenants[1].ID, rep2.Tenants[1].ID, "should not have same tenant ID")
	require.NotEmpty(rep2.NextPageToken, "next page token expected")
	require.NotEqual(rep.NextPageToken, rep2.NextPageToken, "should have new next page token")

	req.NextPageToken = rep2.NextPageToken
	rep3, err := suite.client.TenantList(ctx, req)
	require.NoError(err, "could not list tenants")
	require.Len(rep3.Tenants, 1, "expected 1 tenant")
	require.NotEqual(rep2.Tenants[0].ID, rep3.Tenants[0].ID, "should not have same tenant ID")
	require.Empty(rep3.NextPageToken, "should be empty when a next page does not exist")

	// Limit maximum number of requests to 5, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for i := 0; i < 5; i++ {
		page, err := suite.client.TenantList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Tenants)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 3, "expected 5 results in 3 pages")
	require.Equal(nResults, 5, "expected 5 results in 3 pages")

	// Set test fixture.
	test := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "",
		Permissions: []string{perms.ReadOrganizations},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.TenantList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")
}

func (suite *tenantTestSuite) TestTenantCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to a mock trtl database
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
		OrgID:       "01GMBVR86186E0EKCHQK4ESJB1",
		Permissions: []string{"create:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set client csrf protection")
	_, err := suite.client.TenantCreate(ctx, &api.Tenant{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantCreate(ctx, &api.Tenant{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditOrganizations}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if tenant id exists.
	_, err = suite.client.TenantCreate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "tenant01", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusBadRequest, "tenant id cannot be specified on create", "expected error when tenant id exists")

	// Should return an error if tenant name does not exist.
	_, err = suite.client.TenantCreate(ctx, &api.Tenant{ID: "", Name: "", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusBadRequest, "tenant name is required", "expected error when tenant name does not exist")

	// Should return an error if tenant environment type does not exist.
	_, err = suite.client.TenantCreate(ctx, &api.Tenant{ID: "", Name: "tenant01", EnvironmentType: ""})
	suite.requireError(err, http.StatusBadRequest, "tenant environment type is required", "expected error when tenant environment type does not exist")

	// Create a tenant test fixture
	req := &api.Tenant{
		Name:            "tenant01",
		EnvironmentType: "prod",
	}

	rep, err := suite.client.TenantCreate(ctx, req)
	require.NoError(err, "could not add tenant")
	require.NotEmpty(rep.ID, "expected non-zero ulid to be populated")
	require.Equal(req.Name, rep.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, rep.EnvironmentType, "tenant environment type should match")
	require.NotEmpty(rep.Created, "expected non-zero created timestamp to be populated")
	require.NotEmpty(rep.Modified, "expected non-zero modified timestamp to be populated")

	// Create a test fixture.
	test := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "",
		Permissions: []string{perms.EditOrganizations},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.TenantCreate(ctx, &api.Tenant{})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")
}

func (suite *tenantTestSuite) TestTenantDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	fixture := &db.Tenant{
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "example-staging",
		EnvironmentType: "prod",
		Created:         time.Now().Add(-time.Hour),
		Modified:        time.Now(),
	}

	// Marshal the data with msgpack
	data, err := fixture.MarshalValue()
	require.NoError(err, "could not marshal the tenant")

	// Call the OnGet method and return the tenant id if the namespace is "organizations".
	// If the namespace is not "organizations", return the JSON test data.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.TenantNamespace:
			return &pb.GetReply{Value: data}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: fixture.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.TenantDetail(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantDetail(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadOrganizations}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should error if the orgID is missing from the claims
	_, err = suite.client.TenantDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantDetail(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the tenant does not exist
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantDetail(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")

	// Create a tenant test fixture.
	req := &api.Tenant{
		ID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:            "example-staging",
		EnvironmentType: "prod",
	}

	reply, err := suite.client.TenantDetail(ctx, req.ID)
	require.NoError(err, "could not retrieve tenant")
	require.Equal(req.ID, reply.ID, "tenant id should match")
	require.Equal(req.Name, reply.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, reply.EnvironmentType, "tenant environment type should match")
	require.NotEmpty(reply.Created, "expected non-zero created timestamp to be populated")
	require.NotEmpty(reply.Modified, "expected non-zero modified timestamp to be populated")

	// Test not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: fixture.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	_, err = suite.client.TenantDetail(ctx, req.ID)
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")
}

func (suite *tenantTestSuite) TestTenantUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	fixture := &db.Tenant{
		OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "tenant001",
		EnvironmentType: "prod",
	}

	// Marshal the data with msgpack
	data, err := fixture.MarshalValue()
	require.NoError(err, "could not marshal the tenant")

	// Call the OnGet method and return the tenant id if the namespace is "organizations".
	// If the namespace is not "organizations", return the JSON test data.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.TenantNamespace:
			return &pb.GetReply{Value: data}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: fixture.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// OnPut should return a success reply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"create:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "example-staging", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "example-staging", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditOrganizations}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should error if the orgID is missing from the claims
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "invalid", Name: "example-staging", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "example-staging", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the tenant does not exist
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "invalid", Name: "tenant001", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")

	// Should return an error if the tenant name does not exist
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusBadRequest, "tenant name is required", "expected error when tenant name does not exist")

	// Should return an error if the tenant environment type does not exist
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "tenant001"})
	suite.requireError(err, http.StatusBadRequest, "tenant environment type is required", "expected error when tenant environent type does not exist")

	req := &api.Tenant{
		ID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:            "tenant001",
		EnvironmentType: "dev",
	}

	rep, err := suite.client.TenantUpdate(ctx, req)
	require.NoError(err, "could not update tenant")
	require.NotEqual(req.ID, "01GM8MEZ097ZC7RQRCWMPRPS0T", "tenant id should not match")
	require.Equal(req.Name, rep.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, rep.EnvironmentType, "tenant environment type should match")
	require.NotEmpty(rep.Created, "expected non-zero created timestamp to be populated")
	require.NotEmpty(rep.Modified, "expected non-zero modified timestamp to be populated")

	// Test not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: fixture.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	_, err = suite.client.TenantUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")
}

func (suite *tenantTestSuite) TestTenantDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	defer cancel()

	// Connect to a mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// OnGet passes the tenantID as a value.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: tenantID[:],
		}, nil
	}

	// Call the OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
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
	err := suite.client.TenantDelete(ctx, tenantID.String())
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.TenantDelete(ctx, tenantID.String())
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.DeleteOrganizations}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should error if the orgID is missing from the claims
	err = suite.client.TenantDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when user does not have permission")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.TenantDelete(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the tenant does not exist
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.TenantDelete(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")

	err = suite.client.TenantDelete(ctx, tenantID.String())
	require.NoError(err, "could not delete tenant")

	// Test not found path
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	err = suite.client.TenantDelete(ctx, tenantID.String())
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")
}

func (suite *tenantTestSuite) TestTenantStats() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to a mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	tenantID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	tenant := &db.Tenant{
		OrgID: ulid.MustParse(orgID),
		ID:    ulid.MustParse(tenantID),
	}

	var tenantData []byte
	tenantData, err := tenant.MarshalValue()
	require.NoError(err, "could not marshal tenant")

	// Trtl mock should return tenant id OnGet if the namespace is "organizations" and the tenant fixture if it is not.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.TenantNamespace:
			return &pb.GetReply{Value: tenantData}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: tenant.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	projects := []*db.Project{
		{
			OrgID:    tenant.OrgID,
			TenantID: tenant.ID,
			ID:       ulids.New(),
		},
		{
			OrgID:    tenant.OrgID,
			TenantID: tenant.ID,
			ID:       ulids.New(),
		},
	}
	projectPrefix := tenant.ID[:]

	topics := map[string][]*db.Topic{
		string(projects[0].ID[:]): {
			{
				OrgID:     projects[0].OrgID,
				ProjectID: projects[0].ID,
				ID:        ulids.New(),
			},
			{
				OrgID:     projects[0].OrgID,
				ProjectID: projects[0].ID,
				ID:        ulids.New(),
			},
		},
		string(projects[1].ID[:]): {
			{
				OrgID:     projects[1].OrgID,
				ProjectID: projects[1].ID,
				ID:        ulids.New(),
			},
		},
	}

	// Trtl mock should return projects and topics on Cursor
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		switch in.Namespace {
		case db.ProjectNamespace:
			if !bytes.Equal(in.Prefix, projectPrefix) {
				return status.Error(codes.FailedPrecondition, "unexpected prefix for cursor request")
			}
			for _, project := range projects {
				data, err := project.MarshalValue()
				require.NoError(err, "could not marshal project fixture data")
				stream.Send(&pb.KVPair{
					Key:       project.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		case db.TopicNamespace:
			require.Contains(topics, string(in.Prefix), "unexpected prefix for cursor request")
			for _, topic := range topics[string(in.Prefix)] {
				data, err := topic.MarshalValue()
				require.NoError(err, "could not marshal topic fixture data")
				stream.Send(&pb.KVPair{
					Key:       topic.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		default:
			return status.Error(codes.FailedPrecondition, "unexpected namespace for cursor request")
		}
		return nil
	}

	keys := &qd.APIKeyList{}

	// Initial quarterdeck mock expects authentication and returns 200 with no keys
	suite.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(keys), mock.RequireAuth())

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.TenantStats(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantStats(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadOrganizations, perms.ReadProjects, perms.ReadTopics, perms.ReadAPIKeys}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the orgID is missing from the claims
	_, err = suite.client.TenantStats(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when orgID is missing from claims")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantStats(ctx, tenantID)
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the tenant ID is not parseable
	claims.OrgID = orgID
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantStats(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant ID is not parseable")

	// Retrieving tenant stats without any keys
	claims.OrgID = orgID
	expected := []*api.StatValue{
		{
			Name:  "projects",
			Value: 2,
		},
		{
			Name:  "topics",
			Value: 3,
		},
		{
			Name:  "keys",
			Value: 0,
		},
		{
			Name:    "storage",
			Value:   0,
			Units:   "GB",
			Percent: 0,
		},
	}

	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	stats, err := suite.client.TenantStats(ctx, tenantID)
	require.NoError(err, "could not get tenant stats")
	require.Equal(expected, stats, "expected tenant stats to match")

	// Return two keys for each project
	// TODO: Testing multiple pages requires a more dynamic mock
	keys = &qd.APIKeyList{
		APIKeys: []*qd.APIKeyPreview{
			{
				ID: ulids.New(),
			},
			{
				ID: ulids.New(),
			},
		},
	}
	expected[2].Value = 4
	suite.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(keys), mock.RequireAuth())
	stats, err = suite.client.TenantStats(ctx, tenantID)
	require.NoError(err, "could not get tenant stats")
	require.Equal(expected, stats, "expected tenant stats to match")

	// Test that an error is returned if quarterdeck returns an error
	suite.quarterdeck.OnAPIKeys("", mock.UseError(http.StatusInternalServerError, "could not list API keys"), mock.RequireAuth())
	_, err = suite.client.TenantStats(ctx, tenantID)
	suite.requireError(err, http.StatusInternalServerError, "could not list API keys", "expected error when quarterdeck returns an error")

	// Test not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: tenant.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	_, err = suite.client.TenantStats(ctx, tenantID)
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")
}
