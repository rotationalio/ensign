package quarterdeck_test

import "context"

func (suite *quarterdeckTestSuite) TestStatus() {
	require := suite.Require()
	rep, err := suite.client.Status(context.Background())
	require.NoError(err, "could not execute client request")
	require.NotEmpty(rep, "expected a complete status reply")
	require.Equal("ok", rep.Status, "expected status to be ok")
	require.NotEmpty(rep.Uptime, "expected some value for uptime")
	require.NotEmpty(rep.Version, "expected some value for version")

	// Ensure when server is stopping we get back a stopping status
	suite.srv.SetHealth(false)
	defer suite.srv.SetHealth(true)

	rep, err = suite.client.Status(context.Background())
	require.NoError(err, "could not execute client request")
	require.Equal("stopping", rep.Status, "expected status to be ok")
	require.NotEmpty(rep.Uptime, "expected some value for uptime")
	require.NotEmpty(rep.Version, "expected some value for version")

	// Ensure that when the server is in maintenance mode we get back maintenance
}
