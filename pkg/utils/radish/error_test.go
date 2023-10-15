package radish_test

import (
	"errors"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rotationalio/ensign/pkg/utils/radish"
	"github.com/stretchr/testify/require"
)

func TestTaskErrors(t *testing.T) {
	werr := errors.New("significant badness happened")
	err := radish.Errorw(werr)

	require.ErrorIs(t, errors.Unwrap(err), werr, "expected to be able to unwrap an error")
	require.ErrorIs(t, err, werr, "expected the error to wrap an error")
	require.EqualError(t, err, "after 0 attempts: significant badness happened")

	// Append some errors
	err.Append(errors.New("could not reach database"))
	err.Append(errors.New("could not reach database"))
	err.Append(errors.New("could not reach database"))
	err.Append(errors.New("failed precondition"))
	err.Append(errors.New("maximum backoff limit reached"))

	err.Since(time.Now().Add(-10 * time.Second))
	require.EqualError(t, err, "after 5 attempts: significant badness happened")

	err.Dict()
	err.Capture(sentry.CurrentHub().Clone())
}

func TestNilTaskError(t *testing.T) {
	err := &radish.Error{}

	require.Nil(t, errors.Unwrap(err))
	require.EqualError(t, err, "task failed after 0 attempts")

	// Append some errors
	err.Append(errors.New("could not reach database"))
	err.Append(errors.New("could not reach database"))
	err.Append(errors.New("could not reach database"))
	err.Append(errors.New("failed precondition"))
	err.Append(errors.New("maximum backoff limit reached"))

	err.Since(time.Now().Add(-10 * time.Second))
	require.EqualError(t, err, "task failed after 5 attempts")

	err.Dict()
	err.Capture(sentry.CurrentHub().Clone())
}
