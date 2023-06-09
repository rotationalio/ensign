package sentry

import (
	"fmt"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

const HeaderRequestID = "X-Request-ID"

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
		request := TransactionName(c)
		span := sentry.StartSpan(c.Request.Context(), "api", sentry.WithTransactionName(request))
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
			// Set a unique request-ID either from the header or generated
			var requestID string
			if requestID = c.Request.Header.Get(HeaderRequestID); requestID == "" {
				requestID = ulid.Make().String()
			}
			c.Set("request_id", requestID)

			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTags(tags)
				scope.SetTag("method", c.Request.Method)
				scope.SetTag("path", c.Request.URL.Path)
				scope.SetTag("request_id", requestID)
			})

		}
		c.Next()
	}
}

// Gin middleware to capture errors set on the gin context.
func ReportErrors(conf Config) gin.HandlerFunc {
	if !conf.UseSentry() || !conf.ReportErrors {
		return nil
	}

	return func(c *gin.Context) {
		// Handle errors after the request is complete
		c.Next()

		// If there are errors send them to Sentry
		if len(c.Errors) > 0 {
			if hub := sentrygin.GetHubFromContext(c); hub != nil {
				status := c.Writer.Status()
				hub.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetTag("status", strconv.Itoa(status))
				})

				for _, err := range c.Errors {
					hub.CaptureException(err)
				}
			}
		}
	}
}

func TransactionName(c *gin.Context) string {
	return fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
}
