package interceptors

import (
	"context"

	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const statusEndpoint = "/ensign.v1beta1.Ensign/Status"

// The maintenance interceptor only allows Status endpoint to be queried and returns a
// service unavailable error otherwise. If the server is not in maintenance mode when
// the interceptor is created this method returns nil.
func UnaryMaintenance(conf config.Config) grpc.UnaryServerInterceptor {
	// Do not return an interceptor if we're not in maintenance mode.
	if !conf.Maintenance {
		return nil
	}

	// This interceptor will supercede all following interceptors.
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Allow the Status endpoint through, otherwise return unavailable
		if info.FullMethod == statusEndpoint {
			return handler(ctx, req)
		}

		err = status.Error(codes.Unavailable, "the Ensign server is currently in maintenance mode")
		log.Debug().Err(err).Str("method", info.FullMethod).Msg("ensign service unavailable during maintenance")
		return nil, err
	}
}

// The stream maintenance interceptor simply returns an unavailable error. If the server
// is not in maintenance mode when the interceptor is created this method returns nil.
func StreamMaintenance(conf config.Config) grpc.StreamServerInterceptor {
	// Do not return an interceptor if we're not in maintenance mode.
	if !conf.Maintenance {
		return nil
	}

	// This interceptor will supercede all following interceptors
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = status.Error(codes.Unavailable, "the Ensign server is currently in maintenance mode")
		log.Debug().Err(err).Str("method", info.FullMethod).Msg("ensign service unavailable during maintenance")
		return err
	}
}
