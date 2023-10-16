package radish

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Options configure the task beyond the input context allowing for retries or backoff
// delays in task processing when there are failures or other task-specific handling.
type Option func(*TaskHandler)

// Specify the number of times to retry a task when it returns an error (default 0).
func WithRetries(retries int) Option {
	return func(o *TaskHandler) {
		o.retries = retries
	}
}

// Backoff strategy to use when retrying (default exponential backoff).
func WithBackoff(backoff backoff.BackOff) Option {
	return func(o *TaskHandler) {
		o.backoff = backoff
	}
}

// Specify a base context to be used as the parent context when the task is executed
// and on all subsequent retries.
//
// NOTE: it is recommended that this context does not contain a deadline, otherwise the
// deadline may expire before the specified number of retries. Use WithTimeout instead.
func WithContext(ctx context.Context) Option {
	return func(o *TaskHandler) {
		o.ctx = ctx
	}
}

// Specify a timeout to add to the context before passing it into the task function.
func WithTimeout(timeout time.Duration) Option {
	return func(o *TaskHandler) {
		o.timeout = timeout
	}
}

// Log a specific error if all retries failed under the provided context. This error
// will be bundled with the errors that caused the retry failure and reported in a
// single error log message.
func WithError(err error) Option {
	return func(o *TaskHandler) {
		o.err = Errorw(err)
	}
}

// Log a specific error as WithError but using fmt.Errorf semantics to create the err.
func WithErrorf(format string, a ...any) Option {
	return func(o *TaskHandler) {
		o.err = Errorf(format, a...)
	}
}
