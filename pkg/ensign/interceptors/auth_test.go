package interceptors_test

import (
	"context"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/interceptors"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	health "github.com/rotationalio/ensign/pkg/utils/probez/grpc/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAuthenticator(t *testing.T) {
	// Create the test authentication server
	auth, err := authtest.NewServer()
	require.NoError(t, err, "could not start authtest server")
	defer auth.Close()

	// Create the interceptors and the mock gRPC server to test with
	authenticator, err := interceptors.NewAuthenticator(
		middleware.WithJWKSEndpoint(auth.KeysURL()),
		middleware.WithAudience(authtest.Audience),
		middleware.WithIssuer(authtest.Issuer),
	)
	require.NoError(t, err, "could not create authenticator interceptors")

	opts := make([]grpc.ServerOption, 0, 2)
	opts = append(opts, grpc.UnaryInterceptor(authenticator.Unary()))
	opts = append(opts, grpc.StreamInterceptor(authenticator.Stream()))
	srv := mock.New(nil, opts...)

	t.Run("Unary", func(t *testing.T) {
		t.Cleanup(srv.Reset)

		// Ensure that the unary endpoint returns a decent response
		srv.OnListTopics = func(ctx context.Context, _ *api.PageInfo) (*api.TopicsPage, error) {
			// Make sure that the claims are in the context, otherwise return an error.
			if _, ok := contexts.ClaimsFrom(ctx); !ok {
				return nil, status.Error(codes.PermissionDenied, "no claims in context")
			}
			return &api.TopicsPage{}, nil
		}

		// Create a client to trigger requests
		ctx := context.Background()
		client, err := srv.ResetClient(ctx)
		require.NoError(t, err, "could not connect client to mock")

		// Should not be able to connect to RPC without authentication
		_, err = client.ListTopics(ctx, &api.PageInfo{})
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = missing credentials")

		// Should not be able to connect with an invalid JWT token
		client, err = srv.ResetClient(ctx, mock.WithPerRPCToken("notarealjwtoken"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "could not connect client to mock")
		_, err = client.ListTopics(ctx, &api.PageInfo{})
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid credentials")

		// Should be able to connect with a valid auth token and claims should be in context
		token, err := auth.CreateAccessToken(&tokens.Claims{Email: "test@example.com"})
		require.NoError(t, err, "could not create access token")
		client, err = srv.ResetClient(ctx, mock.WithPerRPCToken(token), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "could not connect client to mock")

		_, err = client.ListTopics(ctx, &api.PageInfo{})
		require.NoError(t, err, "could not access endpoint with valid token")
	})

	t.Run("Stream", func(t *testing.T) {
		t.Cleanup(srv.Reset)

		// Create a client to trigger requests
		ctx := context.Background()
		client, err := srv.ResetClient(ctx)
		require.NoError(t, err, "could not connect client to mock")

		// Handle stream RPC
		srv.OnPublish = func(stream api.Ensign_PublishServer) error {
			// Make sure that the claims are in the context, otherwise return an error.
			if _, ok := contexts.ClaimsFrom(stream.Context()); !ok {
				return status.Error(codes.PermissionDenied, "no claims in context")
			}

			stream.Send(&api.PublisherReply{})
			return nil
		}

		// Should be able to connect to RPC without authentication
		stream, err := client.Publish(ctx)
		require.NoError(t, err, "expected to connect to stream without error")

		// Should not be able to send a message without authentication
		_, err = stream.Recv()
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = missing credentials")

		// Should not be able to connect with an invalid JWT token
		client, err = srv.ResetClient(ctx, mock.WithPerRPCToken("notarealjwtoken"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "could not connect client to mock")

		// Should be able to connect to RPC without authentication
		stream, err = client.Publish(ctx)
		require.NoError(t, err, "expected to connect to stream without error")

		// Should not be able to send a message without authentication
		_, err = stream.Recv()
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid credentials")

		// Should be able to connect with a valid auth token and claims should be in context
		token, err := auth.CreateAccessToken(&tokens.Claims{Email: "test@example.com"})
		require.NoError(t, err, "could not create access token")
		client, err = srv.ResetClient(ctx, mock.WithPerRPCToken(token), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "could not connect client to mock")

		// Should be able to connect to RPC without authentication
		stream, err = client.Publish(ctx)
		require.NoError(t, err, "expected to connect to stream without error")

		// Should not be able to send a message without authentication
		_, err = stream.Recv()
		require.NoError(t, err, "could not authenticate stream")
	})

	t.Run("Public", func(t *testing.T) {
		// Create a client to trigger requests
		ctx := context.Background()
		client, err := srv.ResetClient(ctx)
		require.NoError(t, err, "could not connect client to mock")

		probe, err := srv.HealthClient(ctx)
		require.NoError(t, err, "could not connect client to mock")

		srv.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			return &api.ServiceState{Status: api.ServiceState_HEALTHY}, nil
		}

		// The status endpoint should be unauthenticated
		rep, err := client.Status(ctx, &api.HealthCheck{})
		require.NoError(t, err, "could not make unauthenticated status request")
		require.Equal(t, api.ServiceState_HEALTHY, rep.Status)

		// The health check endpoint should be unauthenticated
		hb, err := probe.Check(ctx, &health.HealthCheckRequest{})
		require.NoError(t, err, "could not make unauthenticated health probe check")
		require.Equal(t, health.HealthCheckResponse_SERVING, hb.Status)

		// The watch health endpoint should be unauthenticated
		stream, err := probe.Watch(ctx, &health.HealthCheckRequest{})
		require.NoError(t, err, "could not initialize health probe watch stream")
		hb, err = stream.Recv()
		require.NoError(t, err, "could not fetch heartbeat from health probe watch")
		require.Equal(t, health.HealthCheckResponse_SERVING, hb.Status)
	})

}
