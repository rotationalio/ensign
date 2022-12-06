package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (suite *tenantTestSuite) TestCreateTenant() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := suite.client.TenantCreate(ctx, &api.Tenant{})
	require.Error(err, http.StatusBadRequest, "tenant id is required")

	req := &api.Tenant{
		ID: "001",
	}

	tenant, err := suite.client.TenantCreate(ctx, req)
	require.NoError(err, "could not add tenant")
	require.Equal(req.ID, tenant.ID, "tenant id should match")
}
