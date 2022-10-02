package quarterdeck_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (suite *quarterdeckTestSuite) TestRegister() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.RegisterRequest{}
	_, err := suite.client.Register(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestLogin() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.LoginRequest{}
	_, err := suite.client.Login(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestAuthenticate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.APIAuthentication{}
	_, err := suite.client.Authenticate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestRefresh() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := suite.client.Refresh(ctx)
	require.Error(err, "expected unimplemented error")
}
