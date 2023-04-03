package sentry_test

import (
	"testing"

	sentrylib "github.com/getsentry/sentry-go"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/stretchr/testify/require"
)

func TestSampler(t *testing.T) {
	sampler := sentry.NewSampler(0.2)
	traces := sampler.TracesSampler()

	// Should return the default sample if no routes have been added
	ctx := sentrylib.SamplingContext{
		Span: &sentrylib.Span{Op: "GET /v1/status"},
	}
	require.Equal(t, 0.2, traces.Sample(ctx))

	// If a route is added should return that route's sample rate
	sampler.AddRoute("GET /v1/status", 0.001)
	require.Equal(t, 0.001, traces.Sample(ctx))
}
