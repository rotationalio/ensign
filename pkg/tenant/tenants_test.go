package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (suite *tenantTestSuite) TestTenantCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trtl := db.GetMock()
	defer trtl.Reset()

	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}
	_, err := suite.client.TenantCreate(ctx, &api.Tenant{})
	require.Error(err, http.StatusBadRequest, "expected unimplemented error")

	req := &api.Tenant{
		ID:              "001",
		Name:            "tenant01",
		EnvironmentType: "prod",
	}

	tenant, err := suite.client.TenantCreate(ctx, req)
	require.NoError(err, "could not add tenant")
	require.Equal(req.ID, tenant.ID, "tenant id should match")
	require.Equal(req.Name, tenant.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, tenant.EnvironmentType, "tenant id should match")
}

func (suite *tenantTestSuite) TestTenantUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if the tenant name does not exist
	_, err := suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", EnvironmentType: "prod"})
	suite.requireError(err, http.StatusBadRequest, "tenant name is required", "expected error when tenant name does not exist")

	// Should return an error if the tenant environment type does not exist
	_, err = suite.client.TenantUpdate(ctx, &api.Tenant{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "example-dev"})
	suite.requireError(err, http.StatusBadRequest, "tenant environment type is required", "expected error when tenant environent type does not exist")

	req := &api.Tenant{
		ID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:            "example-dev",
		EnvironmentType: "dev",
	}

	tenant, err := suite.client.TenantUpdate(ctx, req)
	require.NoError(err, "could not update tenant")
	require.Equal(req.ID, tenant.ID, "tenant id should match")
	require.Equal(req.Name, tenant.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, tenant.EnvironmentType, "tenant environment type should match")
}
