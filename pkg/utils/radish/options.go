package radish

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Options configure the task beyond the input context allowing for retries or backoff
// delays in task processing when there are failures or other task-specific handling.
type Option func(*options)

// Specify the number of times to retry a task when it returns an error (default 0).
func WithRetries(retries int) Option {
	return func(o *options) {
		o.retries = retries
	}
}

// Backoff strategy to use when retrying (default exponential backoff).
func WithBackoff(backoff backoff.BackOff) Option {
	return func(o *options) {
		o.backoff = backoff
	}
}

// Specify a timeout to add to the context before passing it into the task function.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// Log a specific error if all retries failed under the provided context. This error
// will be bundled with the errors that caused the retry failure and reported in a
// single error log message.
func WithError(err error) Option {
	return func(o *options) {
		o.err = Errorw(err)
	}
}

// Log a specific error as WithError but using fmt.Errorf semantics to create the err.
func WithErrorf(format string, a ...any) Option {
	return func(o *options) {
		o.err = Errorf(format, a...)
	}
}

// Default options are 0 retries, exponential backoff, and a radish.Error.
type options struct {
	retries int
	backoff backoff.BackOff
	timeout time.Duration
	err     *Error
}

// Helper function to create internal options from defaults and variadic options.
func makeOptions(opts ...Option) *options {
	// Create default options
	o := &options{
		retries: 0,
		err:     &Error{},
	}

	// Override options with user preferences
	for _, opt := range opts {
		opt(o)
	}

	// If retries is greater than 0 but no backoff is set, set default backoff.
	if o.retries > 0 && o.backoff == nil {
		o.backoff = backoff.NewExponentialBackOff()
	}
	return o
}
