package sentry_test

import (
	"fmt"
	"testing"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	. "github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/stretchr/testify/require"
)

func TestZerologCompatibility(t *testing.T) {
	godotenv.Load()

	config := &Config{}
	err := envconfig.Process("sentry", config)
	require.NoError(t, err, "could not process sentry config")

	if !config.UseSentry() {
		t.Skip("to enable this test set the $SENTRY_DSN environment variable")
	}

	config.Environment = "testing"
	err = Init(*config)
	require.NoError(t, err, "could not initialize sentry")

	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelFatal)
		scope.SetContext("error", map[string]interface{}{
			"message": "panic in the disco",
		})
	})
	hub.CaptureException(fmt.Errorf("we're in trouble"))
	hub.Flush(10 * time.Second)
}
