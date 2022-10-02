package quarterdeck_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (suite *quarterdeckTestSuite) TestAPIKeyList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: implement actual tests
	req := &api.PageQuery{}
	rep, err := suite.client.APIKeyList(ctx, req)
	require.NoError(err, "should return an empty list")
	require.Empty(rep.APIKeys)
	require.Empty(rep.NextPageToken)
}

func (suite *quarterdeckTestSuite) TestAPIKeyCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.APIKey{}
	_, err := suite.client.APIKeyCreate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestAPIKeyDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := suite.client.APIKeyDetail(ctx, "42")
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestAPIKeyUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.APIKey{ID: 42}
	_, err := suite.client.APIKeyUpdate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestAPIKeyDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := suite.client.APIKeyDelete(ctx, "42")
	require.Error(err, "expected unimplemented error")
}
