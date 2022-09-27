package interceptors

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/getsentry/sentry-go"
	config "github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Panic recovery logs the panic to Sentry if it is enabled and then converts the panic
// into a gRPC error to return to the client; this allows the server to stay online.
func UnaryRecovery(conf config.Config) grpc.UnaryServerInterceptor {
	useSentry := conf.UseSentry()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		panicked := true

		defer func() {
			// NOTE: recover only works for the current go routine so panics in any
			// go routine launched by the handler will not be recovered by this function
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

// Panic recovery logs the panic to Sentry if it is enabled and then converts the panic
// into a gRPC error to return to the client; this allows the server to stay online.
func StreamRecovery(conf config.Config) grpc.StreamServerInterceptor {
	useSentry := conf.UseSentry()
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			// NOTE: recover only works for the current go routine so panics in any
			// go routine launched by the handler will not be recovered by this function
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
