package ensign

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: move this to its own package for better organization
// Prepares the interceptors (middleware) for the unary RPC endpoings of the server.
// The first interceptor will be the outer most, while the last interceptor will be the
// inner most wrapper around the real call. All unary interceptors returned by this
// method should be chained using grpc.ChainUnaryInterceptor().
func (s *Server) UnaryInterceptors() []grpc.UnaryServerInterceptor {
	opts := make([]grpc.UnaryServerInterceptor, 0, 1)

	opts = append(opts, s.UnaryRecovery())
	return opts
}

// Prepares the interceptors (middleware) for the unary RPC endpoings of the server.
// The first interceptor will be the outer most, while the last interceptor will be the
// inner most wrapper around the real call. All stream interceptors returned by this
// method should be chained using grpc.ChainStreamInterceptor().
func (s *Server) StreamInterceptors() []grpc.StreamServerInterceptor {
	opts := make([]grpc.StreamServerInterceptor, 0, 1)

	opts = append(opts, s.StreamRecovery())
	return opts
}

// Panic recovery logs the panic to Sentry if it is enabled and then converts the panic
// into a gRPC error to return to the client; this allows the server to stay online.
func (s *Server) UnaryRecovery() grpc.UnaryServerInterceptor {
	useSentry := s.conf.Sentry.UseSentry()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				if useSentry {
					sentry.CurrentHub().Recover(r)
				}

				log.WithLevel(zerolog.PanicLevel).
					Err(fmt.Errorf("%v", r)).
					Bool("panicked", panicked).
					Str("stack_trace", string(debug.Stack())).
					Msg("ensign server has recovered from a panic")
				err = status.Error(codes.Internal, "an unhandled exception occurred")
			}
		}()

		rep, err := handler(ctx, req)
		panicked = false
		return rep, err
	}
}

func (s *Server) StreamRecovery() grpc.StreamServerInterceptor {
	useSentry := s.conf.Sentry.UseSentry()
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				if useSentry {
					sentry.CurrentHub().Recover(r)
				}

				log.WithLevel(zerolog.PanicLevel).
					Err(fmt.Errorf("%v", r)).
					Bool("panicked", panicked).
					Str("stack_trace", string(debug.Stack())).
					Msg("ensign server has recovered from a panic")
				err = status.Error(codes.Internal, "an unhandled exception occurred")
			}
		}()

		err = handler(srv, stream)
		panicked = false
		return err
	}
}
