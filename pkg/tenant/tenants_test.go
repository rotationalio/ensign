package tenant_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (suite *tenantTestSuite) TestCreateTenant() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.Tenant{}

	_, err := suite.client.TenantCreate(ctx, req)
	require.Error(err, "could not add tenant")
}
