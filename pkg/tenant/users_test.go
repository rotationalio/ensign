package tenant_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (suite *tenantTestSuite) TestTenantCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.Tenant{
		ID:              "001",
		TenantName:      "tenant01",
		EnvironmentType: "Prod",
	}

	_, err := suite.client.TenantCreate(ctx, req)
	require.NoError(err, "could not add tenant")
}
