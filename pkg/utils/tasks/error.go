package tasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

var (
	ErrTaskManagerStopped = errors.New("the task manager is not running")
)

type Error struct {
	err      error          // user supplied errors
	retries  int            // number of retries attempted
	taskerrs map[string]int // the errors that occurred in the task mapped to how many times they ocurred.
	duration time.Duration  // the amount of time the task was tried before failure
}

func NewError(err error) *Error {
	return &Error{err: err}
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("after %d retries: %s", e.retries, e.err.Error())
	}
	return fmt.Sprintf("task failed with %d types of errors after %d retries", len(e.taskerrs), e.retries)
}

func (e *Error) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Append(err error) {
	if e.taskerrs == nil {
		e.taskerrs = make(map[string]int)
	}
	e.retries++
	e.taskerrs[err.Error()]++
}

func (e *Error) Since(started time.Time) {
	e.duration = time.Since(started)
}

func (e *Error) Log(log zerolog.Logger) *zerolog.Event {
	retryErrors := make([]error, 0, len(e.taskerrs))
	for err, count := range e.taskerrs {
		retryErrors = append(retryErrors, fmt.Errorf("%q occurred %d times", err, count))
	}

	return log.Error().
		Err(e).
		Errs("retry_errors", retryErrors).
		Dur("retry_duration", e.duration).
		Int("retries", e.retries)
}

func (e *Error) Capture(hub *sentry.Hub) {
	if hub != nil {
		hub.ConfigureScope(func(scope *sentry.Scope) {
			info := map[string]interface{}{
				"retries":        e.retries,
				"retry_duration": e.duration,
				"retry_errors":   e.taskerrs,
			}
			scope.SetContext("error", info)
			scope.SetLevel(sentry.LevelError)
		})
		hub.CaptureException(e)
	}
}
