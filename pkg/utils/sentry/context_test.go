package sentry_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	sentryutils "github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/stretchr/testify/require"
)

func TestCloneContext(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost", nil)
	require.NoError(t, err)

	// Create a gin context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	// If the hub is not set, CloneContext should return a background context
	ctx := sentryutils.CloneContext(c)
	require.Equal(t, context.Background(), ctx, "expected background context when sentry hub is not set")

	// Invoke the sentrygin middleware
	handler := sentrygin.New(sentrygin.Options{})
	handler(c)

	// CloneContext should return a context with the sentry hub
	clone := sentryutils.CloneContext(c)
	require.NotNil(t, sentry.GetHubFromContext(clone), "expected sentry hub on context")
}
