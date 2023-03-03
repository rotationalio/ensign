package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rs/zerolog/log"
)

// GinLogger returns a new Gin middleware that performs logging for our JSON APIs using
// zerolog rather than the default Gin logger which is a standard HTTP logger. Provide
// the server name (e.g. adminAPI or BFF) to help us parse the logs.
// NOTE: we previously used github.com/dn365/gin-zerolog but wanted more customization.
func GinLogger(server string) gin.HandlerFunc {
	version := pkg.Version()

	// TODO: Initialize prometheus collectors (this function has a sync.Once so it's safe to call more than once)
	// metrics.Serve()

	return func(c *gin.Context) {
		// Before request
		started := time.Now()

		// TODO prometheus metrics - increment active requests
		// switch server {
		// case "quarterdeck":
		// 	metrics.ActiveQDRequests.Inc()
		// case "tenant":
		// 	metrics.ActiveTenantRequests.Inc()
		// }

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		// Handle the request
		c.Next()

		// TODO prometheus metrics - decrement active requests
		// switch server {
		// case "quarterdeck":
		// 	metrics.ActiveQDRequests.Dec()
		// case "tenant":
		// 	metrics.ActiveTenantRequests.Dec()
		// }

		// After request
		status := c.Writer.Status()
		logctx := log.With().
			Str("path", path).
			Str("ser_name", server).
			Str("version", version).
			Str("method", c.Request.Method).
			Dur("resp_time", time.Since(started)).
			Int("resp_bytes", c.Writer.Size()).
			Int("status", status).
			Str("client_ip", c.ClientIP()).
			Logger()

		// This field requires us to append errors to the Gin context before a 500
		msg := c.Errors.String()
		if msg == "" {
			msg = fmt.Sprintf("%s %s %s %d", server, c.Request.Method, c.Request.URL.Path, status)
		}

		// TODO prometheus metrics - log request duration and type
		// duration := time.Since(started)
		// metrics.RequestDuration.WithLabelValues(server, http.StatusText(status), path).Observe(duration.Seconds())
		// metrics.RequestsHandled.WithLabelValues(server, http.StatusText(status), path).Inc()

		switch {
		case status >= 400 && status < 500:
			logctx.Warn().Msg(msg)
		case status >= 500:
			logctx.Error().Msg(msg)
		default:
			logctx.Info().Msg(msg)
		}
	}
}
