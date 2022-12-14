package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

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

	// Create a tenant test fixture
	req := &api.Tenant{
		Name:            "tenant01",
		EnvironmentType: "prod",
	}

	tenant, err := suite.client.TenantCreate(ctx, req)
	require.NoError(err, "could not add tenant")
	require.Equal(req.Name, tenant.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, tenant.EnvironmentType, "tenant id should match")
}

func (suite *tenantTestSuite) TestTenantDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenant := &db.Tenant{
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "example-staging",
		EnvironmentType: "prod",
	}
	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the data with msgpack
	data, err := tenant.MarshalValue()
	require.NoError(err, "could not marshal the tenant")

	// Unmarshal the data with msgpack
	other := &db.Tenant{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the tenant")

	// Call the OnGet method and return the JSON test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Should return an error if the tenant does not exist
	_, err = suite.client.TenantDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant id", "expected error when tenant does not exist")

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
}

func (suite *tenantTestSuite) TestTenantUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenant := &db.Tenant{
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "example-staging",
		EnvironmentType: "prod",
	}

	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the data with msgpack
	data, err := tenant.MarshalValue()
	require.NoError(err, "could not marshal the tenant")

	// Unmarshal the data with msgpack
	other := &db.Tenant{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the tenant")

	// Call the OnGet method and return the JSON test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if the tenant does not exist
	_, err = suite.client.TenantDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant id", "expected error when tenant does not exist")

	// Should return an error if the tenant name does not exist
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusBadRequest, "tenant name is required", "expected error when tenant name does not exist")

	// Should return an error if the tenant environment type does not exist
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "example-dev"})
	suite.requireError(err, http.StatusBadRequest, "tenant environment type is required", "expected error when tenant environent type does not exist")

	req := &api.Tenant{
		ID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:            "example-dev",
		EnvironmentType: "dev",
	}

	rep, err := suite.client.TenantUpdate(ctx, req)
	require.NoError(err, "could not update tenant")
	require.NotEqual(req.ID, "01GM8MEZ097ZC7RQRCWMPRPS0T", "tenant id should not match")
	require.Equal(req.Name, rep.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, rep.EnvironmentType, "tenant environment type should match")
}

func (suite *tenantTestSuite) TestTenantDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

	defer cancel()

	// Connect to a mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
		return &pb.DeleteReply{}, nil
	}

	// Should return an error if the tenant does not exist
	err := suite.client.TenantDelete(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant id", "expected error when tenant does not exist")

	err = suite.client.TenantDelete(ctx, tenantID)
	require.NoError(err, "could not delete tenant")
}
