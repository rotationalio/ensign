package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
)

// UnaryInterceptor for gRPC services that are using Sentry. Ensures that Sentry is used
// in a thread-safe manner and that performance, panics, and errors are correctly
// tracked with gRPC method names.
func UnaryInterceptor(conf Config) grpc.UnaryServerInterceptor {
	trackPerformance := conf.UsePerformanceTracking()
	reportErrors := conf.ReportErrors
	repanic := conf.Repanic

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Clone the hub for concurrent operations
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		if trackPerformance {
			span := sentry.StartSpan(ctx, "grpc", sentry.TransactionName(info.FullMethod))
			defer span.Finish()
		}

		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTransaction(info.FullMethod)
			scope.SetTag("rpc", "unary")
		})

		defer sentryRecovery(hub, ctx, repanic)
		rep, err := handler(ctx, req)
		if reportErrors && err != nil {
			hub.CaptureException(err)
		}
		return rep, err
	}
}

// StreamInterceptor for gRPC services that are using Sentry. Ensures that Sentry is
// used in a thread-safe manner and that performance, panics, and errors are correctly
// tracked with gRPC method names. Streaming RPCs don't necessarily benefit from Sentry
// performance tracking, but it is helpful to see average durations.
func StreamInterceptor(conf Config) grpc.StreamServerInterceptor {
	trackPerformance := conf.UsePerformanceTracking()
	reportErrors := conf.ReportErrors
	repanic := conf.Repanic

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Clone the hub for concurrent operations
		ctx := stream.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		if trackPerformance {
			span := sentry.StartSpan(ctx, "grpc", sentry.TransactionName(info.FullMethod))
			defer span.Finish()
		}

		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTransaction(info.FullMethod)
			scope.SetTag("rpc", "streaming")
		})

		stream = Stream(stream, ctx)
		defer sentryRecovery(hub, ctx, repanic)

		err = handler(srv, stream)
		if reportErrors && err != nil {
			hub.CaptureException(err)
		}

		return err
	}
}

func sentryRecovery(hub *sentry.Hub, ctx context.Context, repanic bool) {
	if err := recover(); err != nil {
		hub.RecoverWithContext(ctx, err)
		if repanic {
			panic(err)
		}
	}
}
