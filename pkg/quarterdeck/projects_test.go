package quarterdeck_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (suite *quarterdeckTestSuite) TestProjectList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: implement tests
	req := &api.PageQuery{}
	rep, err := suite.client.ProjectList(ctx, req)
	require.NoError(err, "should return an empty list")
	require.Empty(rep.Projects)
	require.Empty(rep.NextPageToken)
}

func (suite *quarterdeckTestSuite) TestProjectCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.Project{}
	_, err := suite.client.ProjectCreate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestProjectDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := suite.client.ProjectDetail(ctx, "42")
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestProjectUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.Project{ID: 42}
	_, err := suite.client.ProjectUpdate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestProjectDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := suite.client.ProjectDelete(ctx, "42")
	require.Error(err, "expected unimplemented error")
}
