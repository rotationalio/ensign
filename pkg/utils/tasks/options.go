package tasks

import (
	"github.com/cenkalti/backoff/v4"
)

// Option allows retries and backoff to be configured for individual tasks.
type Option func(*options)

type options struct {
	retries int
	backoff backoff.BackOff
	err     error
}

func makeOptions(opts ...Option) *options {
	// Create default options
	o := &options{
		retries: 0,
		backoff: backoff.NewExponentialBackOff(),
		err:     &Error{},
	}

	// Override options with user preferences
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Number of retries to attempt before giving up, default 0
func WithRetries(retries int) Option {
	return func(o *options) {
		o.retries = retries
	}
}

// Backoff strategy to use when retrying, default is an exponential backoff
func WithBackoff(backoff backoff.BackOff) Option {
	return func(o *options) {
		o.backoff = backoff
	}
}

// Log a specific error if all retries failed under the provided context. This error
// will be bundled with the errors that caused the retry failure and reported in a
// single error log message.
func WithError(err error) Option {
	return func(o *options) {
		o.err = err
	}
}
