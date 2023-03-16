package sentry

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var sentryToZerologLevel = map[sentry.Level]zerolog.Level{
	sentry.LevelDebug:   zerolog.DebugLevel,
	sentry.LevelInfo:    zerolog.InfoLevel,
	sentry.LevelWarning: zerolog.WarnLevel,
	sentry.LevelError:   zerolog.ErrorLevel,
	sentry.LevelFatal:   zerolog.FatalLevel,
}

// Handle is a top level function for dealing with errors in a robust manner. It logs
// the error using zerolog at the specified level, sends an error to Sentry if the hub
// is available, adds the error to the gin context if it's available and performs other
// tasks related to monitoring and alerting of errors in the Ensign project.
//
// This should only be used if the error needs to generate an alert; otherwise use
// zerolog directly rather than using this method.
//
// The sentry level is mapped to the zerolog level, which means the zerolog.TraceLevel
// and zerolog.PanicLevel are not available in this method.
//
// The ctx should be either a gin.Context or a context.Context; the hub is extracted
// from the context if it was set by middleware or interceptors.
//
// The error can be specified or it can be nil; if nil a Sentry CaptureMessage is used,
// otherwise CaptureException is used. The error will be added to the zerolog message if
// it is not nil.
//
// The extra data is used to set the context on the Sentry message as well as on the
// zerolog event; it is purposefully generic to enable multiple error handling.
//
// The message is optional and is formatted for zerolog if format args are specified. It
// is also used to format the error that is sent to Sentry so use with care!
func Handle(level sentry.Level, ctx interface{}, err error, extra map[string]interface{}, msg string, args ...interface{}) {
	// Attempt to fetch the hub from the context
	var hub *sentry.Hub
	switch c := ctx.(type) {
	case *gin.Context:
		hub = sentrygin.GetHubFromContext(c)
		if err != nil {
			c.Error(err)
		}
	case context.Context:
		hub = sentry.GetHubFromContext(c)
	}

	if hub != nil {
		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetContext("error", extra)
			scope.SetLevel(level)
		})

		if err != nil {
			err = &ServiceError{err: err, msg: msg, args: args}
			hub.CaptureException(err)
		} else {
			hub.CaptureMessage(fmt.Sprintf(msg, args...))
		}
	}

	// Prepare the log context
	// NOTE: this is absolutely the least performant method of zerologging ...
	dict := zerolog.Dict()
	for key, val := range extra {
		dict.Interface(key, val)
	}

	// Send the zerolog log message
	log.WithLevel(sentryToZerologLevel[level]).Err(err).Dict("extra", dict).Msgf(msg, args...)
}

// Reports a debug level event to Sentry and logs a debug message. Use this method when
// the debug message should produce an alert that the team can take action on (which
// should happen only very rarely in code). Most of the time you should use zerolog.Debug
// directly unless this is at the top level of the stack.
func Debug(ctx interface{}, err error, extra map[string]interface{}, msg string, args ...interface{}) {
	Handle(sentry.LevelDebug, ctx, err, extra, msg, args...)
}

// Reports an info level event to Sentry and logs an info message. Use this method when
// the info message should produce an alert that the team can take action on (which
// should happen very rarely in code and is probably related to a third party service
// such as rate limits or usage thresholds). Most of the time you should use zerolog.Info
// directly unless this is at the top level of the stack.
func Info(ctx interface{}, err error, extra map[string]interface{}, msg string, args ...interface{}) {
	Handle(sentry.LevelInfo, ctx, err, extra, msg, args...)
}

// Report a warning level event to Sentry and logs a warning messages. Use this method
// on top level service handlers to produce alerts that something is going wrong in the
// code such as bad requests or not found errors. The team will likely not take action
// on these errors but will get a general sense of what is going on in code. When not
// in a service handler it is better to use zerolog.Warn directly.
func Warn(ctx interface{}, err error, extra map[string]interface{}, msg string, args ...interface{}) {
	Handle(sentry.LevelWarning, ctx, err, extra, msg, args...)
}

// Report an error to Sentry and log an error message. This is the most commonly used
// method for Sentry on top level service handlers and is intended to produce alerts
// that something is going wrong and that the team needs to handle it. When not in a
// service handler, feel free to use zerolog.Error but probably zerolog.Warn is more
// appropriate for most cases.
func Error(ctx interface{}, err error, extra map[string]interface{}, msg string, args ...interface{}) {
	Handle(sentry.LevelError, ctx, err, extra, msg, args...)
}

// Report a critical error to Sentry and log a fatal error message. While this method
// will not cause the process to exit, it should create a serious alert that will cause
// on call personnel to immediately act. Use with care!
func Fatal(ctx interface{}, err error, extra map[string]interface{}, msg string, args ...interface{}) {
	Handle(sentry.LevelFatal, ctx, err, extra, msg, args...)
}

// A standardized error type for fingerprinting inside of Sentry.
type ServiceError struct {
	msg  string
	args []interface{}
	err  error
}

func (e *ServiceError) Error() string {
	if e.msg == "" {
		return e.err.Error()
	}

	msg := fmt.Sprintf(e.msg, e.args...)
	return fmt.Sprintf("%s: %s", msg, e.err)
}
