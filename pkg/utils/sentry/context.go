package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func CloneContext(c *gin.Context) context.Context {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		return sentry.SetHubOnContext(context.Background(), hub.Clone())
	}
	return context.Background()
}

// Wrap a grpc.ServerStream handler with the sentry context.
func Stream(s grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return &stream{s, ctx}
}

type stream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *stream) Context() context.Context {
	return s.ctx
}
