package tenant_test

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
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

	req := &db.Tenant{
		ID:   uuid.MustParse("1d4db493-a16f-4766-b328-62da380f28ec"),
		Name: "tenant-name",
	}

	// Replace string with ulid?
	tenant, err := suite.client.TenantDetail(ctx, "001")
	require.Error(err, http.StatusBadRequest, "tenant id is required")
	require.Equal(req.ID, tenant.ID, "tenant id should match")
	require.Equal(req.Name, tenant.TenantName, "tenant name should match")
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

	err := suite.client.TenantDelete(ctx, "4c3f75b2-b49a-4a4a-a207-ddd8d075c775")
	require.NoError(err, "could not delete tenant")
}
