package tenant_test

import (
	"context"
	"net/http"
	"os"
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
	tenantID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	defer cancel()

	trtl := db.GetMock()
	defer trtl.Reset()

	data, err := os.ReadFile("testdata/tenant.json")
	if err != nil {
		return
	}

	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	req := &api.Tenant{
		ID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name:            "tenant-name",
		EnvironmentType: "prod",
	}

	tenant, err := suite.client.TenantDetail(ctx, tenantID)
	require.Error(err, http.StatusBadRequest, "could not get tenant")
	require.Equal(req, tenant, "tenant should match")
}

func (suite *tenantTestSuite) TestTenantDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"

	defer cancel()

	trtl := db.GetMock()
	defer trtl.Reset()

	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := suite.client.TenantDelete(ctx, tenantID)
	require.NoError(err, "could not delete tenant")
}
