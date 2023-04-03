package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

const (
	APIStatusEndpoint    = "GET /v1/status"
	EnsignStatusEndpoint = "/ensign.v1beta1.Ensign/Status"
	StatusSampleRate     = 0.005
)

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

func NewStatusSampler(defaultSampleRate float64) sentry.TracesSampler {
	sampler := NewSampler(defaultSampleRate)
	sampler.AddRoute(APIStatusEndpoint, StatusSampleRate)
	sampler.AddRoute(EnsignStatusEndpoint, StatusSampleRate)

	return sampler.TracesSampler()
}

func (s *Sampler) Sample(ctx sentry.SamplingContext) float64 {
	// Inherit decision from parent
	if ctx.Parent != nil && ctx.Parent.Sampled != sentry.SampledUndefined {
		if ctx.Parent.Sampled == sentry.SampledTrue {
			log.Trace().Float64("sample", 1.0).Str("op", ctx.Span.Op).Msg("parent sampled true")
			return 1.0
		}
		log.Trace().Float64("sample", 0.0).Str("op", ctx.Span.Op).Msg("parent sampled false")
		return 0.0
	}

	// Determine if specific route should be sampled at a different rate
	if sample, ok := s.routes[ctx.Span.Op]; ok {
		log.Trace().Float64("sample", sample).Str("op", ctx.Span.Op).Msg("txn sampler")
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
