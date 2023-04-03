package sentry

import "github.com/getsentry/sentry-go"

type Sampler struct {
	defaultSampleRate float64
	routes            map[string]float64
}

func NewSampler(defaultSampleRate float64) *Sampler {
	return &Sampler{
		defaultSampleRate: defaultSampleRate,
		routes:            make(map[string]float64),
	}
}

func (s *Sampler) Sample(ctx sentry.SamplingContext) float64 {
	// Inherit decision from parent
	if ctx.Parent != nil && ctx.Parent.Sampled != sentry.SampledUndefined {
		return 1.0
	}

	// Determine if specific route should be sampled at a different rate
	if sample, ok := s.routes[ctx.Span.Op]; ok {
		return sample
	}

	// Return the default sample rate if nothing else
	return s.defaultSampleRate
}

func (s *Sampler) TracesSampler() sentry.TracesSampler {
	return sentry.TracesSampler(s.Sample)
}

func (s *Sampler) AddRoute(route string, sample float64) {
	s.routes[route] = sample
}
