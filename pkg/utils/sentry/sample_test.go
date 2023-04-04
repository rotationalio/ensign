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

func TestStatusSampler(t *testing.T) {
	testCases := []struct {
		Op       string
		Expected float64
	}{
		{sentry.APIStatusEndpoint, sentry.StatusSampleRate},
		{sentry.EnsignStatusEndpoint, sentry.StatusSampleRate},
		{"POST /v1/projects", 0.2},
		{"/ensign.v1beta1.Ensign/Publish", 0.2},
	}

	traces := sentry.NewStatusSampler(0.2)
	for _, tc := range testCases {
		ctx := sentrylib.SamplingContext{
			Span: &sentrylib.Span{Op: tc.Op},
		}
		require.Equal(t, tc.Expected, traces.Sample(ctx))
	}
}
