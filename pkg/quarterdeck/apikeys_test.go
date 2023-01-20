package quarterdeck_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

func (suite *quarterdeckTestSuite) TestAPIKeyList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: implement actual tests
	req := &api.PageQuery{}
	_, err := suite.client.APIKeyList(ctx, req)
	require.Error(err, "unauthorized requests should not return a response")

	// require.NoError(err, "should return an empty list")
	// require.Empty(rep.APIKeys)
	// require.Empty(rep.NextPageToken)
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

	req := &api.APIKey{ID: ulids.New()}
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
