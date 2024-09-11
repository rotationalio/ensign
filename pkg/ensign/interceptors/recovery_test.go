package interceptors_test

import (
	"context"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/interceptors"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestRecovery(t *testing.T) {
	conf := sentry.Config{
		DSN: "https://test.sentry.io/1234",
	}

	// Create the mock Ensign server to test the interceptors with
	opts := make([]grpc.ServerOption, 0, 2)
	opts = append(opts, grpc.UnaryInterceptor(interceptors.UnaryRecovery(conf)))
	opts = append(opts, grpc.StreamInterceptor(interceptors.StreamRecovery(conf)))
	srv := mock.New(nil, opts...)

	// Create client to trigger requests
	client, err := srv.Client()
	require.NoError(t, err, "could not connect client to mock")

	t.Run("UnaryPanic", func(t *testing.T) {
		t.Cleanup(srv.Reset)
		srv.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			panic("something very bad happened")
		}

		rep, err := client.Status(context.Background(), &api.HealthCheck{})
		require.EqualError(t, err, "rpc error: code = Internal desc = an unhandled exception occurred", "expected an unhandled exception error after the panic")
		require.Nil(t, rep, "expected a nil response after the panic")
	})

	t.Run("StreamPanic", func(t *testing.T) {
		t.Cleanup(srv.Reset)
		srv.OnPublish = func(api.Ensign_PublishServer) error {
			panic("something very bad happened")
		}

		stream, err := client.Publish(context.Background())
		require.NoError(t, err, "should not error on establishing the stream")

		msg, err := stream.Recv()
		require.EqualError(t, err, "rpc error: code = Internal desc = an unhandled exception occurred", "expected an unhandled exception error after the panic")
		require.Nil(t, msg, "expected a nil msg after the panic")
	})

	t.Run("UnaryNoPanic", func(t *testing.T) {
		t.Cleanup(srv.Reset)
		srv.UseError(mock.StatusRPC, codes.NotFound, "status not found")

		_, err := client.Status(context.Background(), &api.HealthCheck{})
		require.EqualError(t, err, "rpc error: code = NotFound desc = status not found", "expected a not found error instead of a panic")
	})

	t.Run("StreamNoPanic", func(t *testing.T) {
		t.Cleanup(srv.Reset)
		srv.UseError(mock.PublishRPC, codes.DataLoss, "missing data from stream")

		stream, err := client.Publish(context.Background())
		require.NoError(t, err, "should not error on establishing the stream")

		msg, err := stream.Recv()
		require.EqualError(t, err, "rpc error: code = DataLoss desc = missing data from stream", "expected a data loss error instead of a panic")
		require.Nil(t, msg, "expected a nil msg after the panic")
	})
}
