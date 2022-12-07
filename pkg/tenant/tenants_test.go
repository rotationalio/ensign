package tenant_test

import (
	"context"
	"net/http"
	"time"

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
		return &pb.PutReply{
			Success: true,
		}, nil
	}
	_, err := suite.client.TenantCreate(ctx, &db.Tenant{})
	require.Error(err, http.StatusBadRequest, "expected unimplemented error")

	req := &db.Tenant{
		Name:            "tenant01",
		EnvironmentType: "prod",
	}

	tenant, err := suite.client.TenantCreate(ctx, req)
	require.NoError(err, "could not add tenant")
	require.NotEmpty(tenant.ID, "tenant id should not be zero")
	require.Equal(req.Name, tenant.Name, "tenant name should match")
	require.Equal(req.EnvironmentType, tenant.EnvironmentType, "tenant envrionment type should match")
}
