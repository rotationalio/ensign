package sentry

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Initialize the Sentry SDK with the configuration; must be called before servers are started
func Init(conf Config) (err error) {
	if err = sentry.Init(conf.ClientOptions()); err != nil {
		return fmt.Errorf("could not initialize sentry: %w", err)
	}

	log.Debug().
		Bool("track_performance", conf.TrackPerformance).
		Float64("sample_rate", conf.SampleRate).
		Msg("sentry tracing is enabled")
	return nil
}

// Flush the Sentry log, usually called before shutting down servers
func Flush(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}

// Gin middleware that tracks HTTP request performance with Sentry
func TrackPerformance(tags map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Setup span performance prior to request:
		request := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		span := sentry.StartSpan(c.Request.Context(), "api", sentry.TransactionName(request))
		for k, v := range tags {
			span.SetTag(k, v)
		}

		// Execute request and compute the performance of the request:
		c.Next()
		span.Finish()
	}
}

// Gin middleware that adds request-level tags to the Sentry scope.
func UseTags(tags map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			for k, v := range tags {
				hub.Scope().SetTag(k, v)
			}

			hub.Scope().SetTag("path", c.Request.URL.Path)
			hub.Scope().SetTag("method", c.Request.Method)
		}
		c.Next()
	}
}
