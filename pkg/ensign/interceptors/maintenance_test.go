package interceptors_test

import (
	"context"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/interceptors"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	health "github.com/rotationalio/ensign/pkg/utils/probez/grpc/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestMaintenance(t *testing.T) {
	conf := config.Config{
		Maintenance: false,
	}

	// Interceptors should return nil when maintenance is false
	require.Nil(t, interceptors.UnaryMaintenance(conf), "expected nil unary interceptor when not in maintenance mode")
	require.Nil(t, interceptors.StreamMaintenance(conf), "expected nil stream interceptor when not in maintenance mode")

	// Create maintenance mock Ensign server to test interceptors with
	conf.Maintenance = true
	opts := make([]grpc.ServerOption, 0, 2)
	opts = append(opts, grpc.UnaryInterceptor(interceptors.UnaryMaintenance(conf)))
	opts = append(opts, grpc.StreamInterceptor(interceptors.StreamMaintenance(conf)))
	srv := mock.New(nil, opts...)

	// Create client to trigger requests
	ctx := context.Background()
	client, err := srv.Client()
	require.NoError(t, err, "could not connect client to mock")

	// Create health probe
	probe, err := srv.HealthClient()
	require.NoError(t, err)

	t.Run("UnaryMaintenance", func(t *testing.T) {
		t.Cleanup(srv.Reset)

		// Should allow Status call through
		srv.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			return &api.ServiceState{
				Status: api.ServiceState_DANGER,
			}, nil
		}

		rep, err := client.Status(ctx, &api.HealthCheck{})
		require.NoError(t, err, "expected no unavailable error for the status endpoint")
		require.NotNil(t, rep, "expected status reply even in maintenance mode")
		require.Equal(t, api.ServiceState_DANGER, rep.Status)
		require.Equal(t, 1, srv.Calls[mock.StatusRPC])

		// Should allow health probe through
		hb, err := probe.Check(ctx, &health.HealthCheckRequest{})
		require.NoError(t, err, "could not probe health check endpoint")
		require.Equal(t, health.HealthCheckResponse_SERVING, hb.Status)

		// Should not allow any calls to other Unary RPCs
		_, err = client.ListTopics(ctx, &api.PageInfo{})
		require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode", "expected error from unary rpc")
		require.Zero(t, srv.Calls[mock.ListTopicsRPC])

		_, err = client.CreateTopic(ctx, &api.Topic{})
		require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode", "expected error from unary rpc")
		require.Zero(t, srv.Calls[mock.CreateTopicRPC])

		_, err = client.DeleteTopic(ctx, &api.TopicMod{})
		require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode", "expected error from unary rpc")
		require.Zero(t, srv.Calls[mock.DeleteTopicRPC])
	})

	t.Run("StreamMaintenance", func(t *testing.T) {
		t.Cleanup(srv.Reset)

		// Should allow health probe through
		hb, err := probe.Watch(ctx, &health.HealthCheckRequest{})
		require.NoError(t, err, "expected no error on stream initialization")
		update, err := hb.Recv()
		require.NoError(t, err, "could not recv heartbeat update")
		require.Equal(t, health.HealthCheckResponse_SERVING, update.Status)

		// Should not call publish stream handler
		pub, err := client.Publish(ctx)
		require.NoError(t, err, "expected no error on stream initialization")
		_, err = pub.Recv()
		require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode", "expected error when recv from publish stream")
		require.Zero(t, srv.Calls[mock.PublishRPC], "expected no calls to the mock server")

		// Should not call subscribe stream handler
		sub, err := client.Subscribe(ctx)
		require.NoError(t, err, "expected no error on stream initialization")
		_, err = sub.Recv()
		require.EqualError(t, err, "rpc error: code = Unavailable desc = the Ensign server is currently in maintenance mode", "expected error when recv from subscribe stream")
		require.Zero(t, srv.Calls[mock.SubscribeRPC], "expected no calls to the mock server")
	})

}
