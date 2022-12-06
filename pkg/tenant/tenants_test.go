package tenant_test

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (suite *tenantTestSuite) TestTenantDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	_, err = suite.client.TenantDetail(ctx, "001")
	require.Error(err, http.StatusBadRequest, "tenant id is required")
}

func (suite *tenantTestSuite) TestTenantDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trtl := db.GetMock()
	defer trtl.Reset()

	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := suite.client.TenantDelete(ctx, "001")
	require.NoError(err, "could not delete tenant")

}
