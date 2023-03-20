package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func CloneContext(c *gin.Context) context.Context {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		return sentry.SetHubOnContext(context.Background(), hub.Clone())
	}
	return context.Background()
}
