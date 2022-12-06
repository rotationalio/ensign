package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (suite *tenantTestSuite) TestTenantDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trtl := db.GetMock()
	defer trtl.Reset()

	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{}, nil
	}

	_, err := suite.client.TenantDetail(ctx, "001")
	require.Error(err, http.StatusBadRequest, "tenant id is required")

	req := &api.Tenant{
		ID: "001",
	}

	tenant, err := suite.client.TenantDetail(ctx, "001")
	require.NoError(err, http.StatusBadRequest, "could not retrieve tenant")
	require.Equal(req.ID, tenant.ID, "tenant id should match")
}
func (suite *tenantTestSuite) TestTenantDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trtl := db.GetMock()
	defer trtl.Reset()

	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{}, nil
	}

	err := suite.client.TenantDelete(ctx, "001")
	require.Error(err, http.StatusBadRequest, "tenant id is required")

	err = suite.client.TenantDelete(ctx, "001")
	require.NoError(err, "could not delete tenant")

}
