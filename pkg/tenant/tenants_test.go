package tenant_test

import (
	"context"
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
