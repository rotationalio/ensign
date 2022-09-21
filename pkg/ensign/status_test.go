package ensign_test

import (
	"context"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
)

func (s *serverTestSuite) TestStatus() {
	var err error
	require := s.Require()

	var rep *api.ServiceState
	rep, err = s.client.Status(context.Background(), &api.HealthCheck{})
	require.NoError(err, "could not make status request")

	require.Equal(api.ServiceState_HEALTHY, rep.Status)
	require.NotEmpty(rep.Version)
	require.NotEmpty(rep.Uptime)
	require.NotEmpty(rep.NotBefore)
	require.NotEmpty(rep.NotAfter)
}
