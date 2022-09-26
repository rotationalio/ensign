package logger

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/grpclog"
)

// DisableGRPCLog sets the grpclog V2 logger to the discard logger. This must be called
// before any grpc calls are made because this method is not mutex protected.
func DisableGRPCLog() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(io.Discard, io.Discard, io.Discard))
}

// Implements the grpclog.LoggerV2 interface to pass logging calls to zerolog.
// To enable grpclog with zerolog (e.g. for grpc logging to GCP) then set this logger
// before any grpc calls are made:
//
//	grpclog.SetLoggerV2(&logger.ZeroGRPCV2{})
//
// This logger should respect the zerolog global log level from grpclog calls.
type ZeroGRPCV2 struct{}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (g *ZeroGRPCV2) Info(args ...interface{}) {
	log.Info().Msg(fmt.Sprint(args...))
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (g *ZeroGRPCV2) Infoln(args ...interface{}) {
	log.Info().Msg(fmt.Sprint(args...))
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (g *ZeroGRPCV2) Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (g *ZeroGRPCV2) Warning(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (g *ZeroGRPCV2) Warningln(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (g *ZeroGRPCV2) Warningf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (g *ZeroGRPCV2) Error(args ...interface{}) {
	log.Error().Msg(fmt.Sprint(args...))
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (g *ZeroGRPCV2) Errorln(args ...interface{}) {
	log.Error().Msg(fmt.Sprint(args...))
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (g *ZeroGRPCV2) Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *ZeroGRPCV2) Fatal(args ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(args...))
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *ZeroGRPCV2) Fatalln(args ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(args...))
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (g *ZeroGRPCV2) Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (g *ZeroGRPCV2) V(l int) bool {
	switch zerolog.GlobalLevel() {
	case zerolog.InfoLevel:
		return l <= 0
	case zerolog.WarnLevel:
		return l <= 1
	case zerolog.ErrorLevel:
		return l <= 2
	case zerolog.FatalLevel, zerolog.PanicLevel:
		return l <= 3
	default:
		return false
	}
}
