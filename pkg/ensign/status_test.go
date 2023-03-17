package ensign_test

import (
	"context"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/ensign"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/utils/bufconn"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *serverTestSuite) TestStatus() {
	var err error
	require := s.Require()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rep *api.ServiceState
	rep, err = s.client.Status(ctx, &api.HealthCheck{})
	require.NoError(err, "could not make status request")

	require.Equal(api.ServiceState_HEALTHY, rep.Status)
	require.NotEmpty(rep.Version)
	require.NotEmpty(rep.Uptime)
	require.NotEmpty(rep.NotBefore)
	require.NotEmpty(rep.NotAfter)
}

func TestMaintenanceMode(t *testing.T) {
	// Create an ensign server in maintenance mode
	// This configuration will run the ensign server as a fully functional gRPC service
	// on an in-memory socket allowing the testing of RPCs from the client perspective.
	conf, err := config.Config{
		Maintenance: true,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
		BindAddr:    "127.0.0.1:0",
		Storage: config.StorageConfig{
			ReadOnly: true,
			DataPath: t.TempDir(),
		},
	}.Mark()
	require.NoError(t, err, "could not mark test configuration in maintenance mode as valid")

	// Create the server and run it on a bufconn.
	srv, err := ensign.New(conf)
	require.NoError(t, err, "could not create server in maintenance mode with a test configuration")

	conn := bufconn.New()
	go srv.Run(conn.Sock())
	t.Cleanup(func() {
		srv.Shutdown()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	// Create a client for testing purposes
	cc, err := conn.Connect(ctx, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "could not connect to bufconn")
	client := api.NewEnsignClient(cc)

	// Should only be able to connect to the status endpoint
	rep, err := client.Status(ctx, &api.HealthCheck{})
	require.NoError(t, err, "could not make status request in maintenance mode")
	require.Equal(t, api.ServiceState_MAINTENANCE, rep.Status, "expected maintenance status not %q", rep.Status.String())
	require.NotEmpty(t, rep.Version)
	require.NotEmpty(t, rep.Uptime)
	require.NotEmpty(t, rep.NotBefore)
	require.NotEmpty(t, rep.NotAfter)

	// Should get unavailable errors from all other unary endpoints
	_, err = client.ListTopics(ctx, &api.PageInfo{})
	require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode")

	_, err = client.CreateTopic(ctx, &api.Topic{})
	require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode")

	_, err = client.DeleteTopic(ctx, &api.TopicMod{})
	require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode")

	// Should get unavailable errors from streaming endpoints
	// TODO: what happens on stream.Send()?
	pub, err := client.Publish(ctx)
	require.NoError(t, err, "should not get error from stream initialization")
	_, err = pub.Recv()
	require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode")

	sub, err := client.Subscribe(ctx)
	require.NoError(t, err, "should not get error from stream initialization")
	_, err = sub.Recv()
	require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode")
}
