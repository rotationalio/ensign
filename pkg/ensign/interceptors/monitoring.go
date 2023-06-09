package interceptors

import (
	"context"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Monitoring does triple duty, handling Sentry tracking, Prometheus metrics, and
// logging with zerolog. These are piled into the same interceptor so that the
// monitoring uses the same latency and tagging constructs and so that this interceptor
// can be the outermost interceptor for unary calls.
func UnaryMonitoring(conf config.Config) grpc.UnaryServerInterceptor {
	// TODO: chain sentry interceptors rather than integrating them into monitoring.
	useSentry := conf.Sentry.UsePerformanceTracking()
	usePrometheus := conf.Monitoring.Enabled
	version := pkg.Version()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Parse the method for tags
		service, method := ParseMethod(info.FullMethod)

		// Trace the sentry transaction span
		if useSentry {
			span := sentry.StartSpan(ctx, "grpc", sentry.WithTransactionName(info.FullMethod))
			defer span.Finish()
		}

		// Monitor how many RPCs have been started
		if usePrometheus {
			o11y.RPCStarted.WithLabelValues("unary", service, method).Inc()
		}

		// Handle the request and trace how long the request takes.
		start := time.Now()
		rep, err := handler(ctx, req)
		duration := time.Since(start)
		code := status.Code(err)

		// Monitor how many RPCs have been completed
		if usePrometheus {
			o11y.RPCHandled.WithLabelValues("unary", service, method, code.String()).Inc()
			o11y.RPCDuration.WithLabelValues("unary", service, method).Observe(duration.Seconds())
		}

		// Prepare log context for logging
		logctx := log.With().
			Str("type", "unary").
			Str("service", service).
			Str("method", method).
			Str("version", version).
			Bool("use_sentry", useSentry).
			Bool("use_prometheus", usePrometheus).
			Uint32("code", uint32(code)).
			Dur("duration", duration).
			Logger()

		// Log based on the error code if it is not unknown
		switch code {
		case codes.OK:
			logctx.Info().Uint32("code", uint32(code)).Dur("duration", duration).Msg(info.FullMethod)
		case codes.Unknown:
			logctx.Error().Err(err).Msgf("unknown error handling %s", info.FullMethod)
		case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
			logctx.Info().Err(err).Msg(info.FullMethod)
		case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unavailable:
			logctx.Warn().Err(err).Msg(info.FullMethod)
		case codes.Unimplemented, codes.Internal, codes.DataLoss:
			logctx.Error().Err(err).Str("full_method", info.FullMethod).Msg(err.Error())
		default:
			logctx.Error().Err(err).Msgf("unhandled error code %s: %s", code, info.FullMethod)
		}
		return rep, err
	}
}

// Monitoring does double duty, handling Prometheus metrics, and logging with zerolog.
// These are piled into the same interceptor so that the monitoring uses the same
// latency and tagging constructs and so that this interceptor can be the outermost
// interceptor for stream calls.
// NOTE: Sentry is excluded from stream monitoring because we do not work to minimize
// the duration of stream processing but rather to maximize it in Ensign.
func StreamMonitoring(conf config.Config) grpc.StreamServerInterceptor {
	usePrometheus := conf.Monitoring.Enabled
	version := pkg.Version()

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Parse the method for tags
		service, method := ParseMethod(info.FullMethod)

		// Monitor how many streams have been started
		if usePrometheus {
			o11y.RPCStarted.WithLabelValues("stream", service, method).Inc()
			stream = &MonitoredStream{stream, service, method}
		}

		// Handle the request and trace how long the request takes.
		start := time.Now()
		err = handler(srv, stream)
		duration := time.Since(start)
		code := status.Code(err)

		// Monitor how many RPCs have been completed
		if usePrometheus {
			o11y.RPCHandled.WithLabelValues("stream", service, method, code.String()).Inc()
			o11y.RPCDuration.WithLabelValues("stream", service, method).Observe(duration.Seconds())
		}

		// Prepare log context for logging
		logctx := log.With().
			Str("type", "stream").
			Str("service", service).
			Str("method", method).
			Str("version", version).
			Bool("use_prometheus", usePrometheus).
			Uint32("code", uint32(code)).
			Dur("duration", duration).
			Logger()

		switch code {
		case codes.OK:
			logctx.Info().Uint32("code", uint32(code)).Dur("duration", duration).Msg(info.FullMethod)
		case codes.Unknown:
			logctx.Error().Err(err).Msgf("unknown error handling %s", info.FullMethod)
		case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
			logctx.Info().Err(err).Msg(info.FullMethod)
		case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unavailable:
			logctx.Warn().Err(err).Msg(info.FullMethod)
		case codes.Unimplemented, codes.Internal, codes.DataLoss:
			logctx.Error().Err(err).Str("full_method", info.FullMethod).Msg(err.Error())
		default:
			logctx.Error().Err(err).Msgf("unhandled error code %s: %s", code, info.FullMethod)
		}
		return err
	}
}

func ParseMethod(method string) (string, string) {
	method = strings.TrimPrefix(method, "/") // remove leading slash
	if i := strings.Index(method, "/"); i >= 0 {
		return method[:i], method[i+1:]
	}
	return "unknown", "unknown"
}

func StreamType(info *grpc.MethodInfo) string {
	if !info.IsClientStream && !info.IsServerStream {
		return "unary"
	}
	if info.IsClientStream && !info.IsServerStream {
		return "client_stream"
	}
	if !info.IsClientStream && info.IsServerStream {
		return "server_stream"
	}
	return "bidirectional"
}

// MonitoredStream wraps a grpc.ServerStream allowing it to increment Sent and Recv
// message counters when they are called by the application.
type MonitoredStream struct {
	grpc.ServerStream
	service string
	method  string
}

// Increment the number of sent messages if there is no error on Send.
func (s *MonitoredStream) SendMsg(m interface{}) (err error) {
	if err = s.ServerStream.SendMsg(m); err == nil {
		o11y.StreamMsgSent.WithLabelValues("stream", s.service, s.method).Inc()
	}
	return err
}

// Increment the number of received messages if there is no error on Recv.
func (s *MonitoredStream) RecvMsg(m interface{}) (err error) {
	if err = s.ServerStream.RecvMsg(m); err == nil {
		o11y.StreamMsgRecv.WithLabelValues("stream", s.service, s.method).Inc()
	}
	return err
}
