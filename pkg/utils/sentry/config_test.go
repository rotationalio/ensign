package sentry_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/stretchr/testify/require"
)

func TestConfigValidation(t *testing.T) {
	conf := sentry.Config{
		DSN:         "",
		Environment: "",
		Release:     "1.4",
		Debug:       true,
	}

	// If DSN is empty, then Sentry is not enabled
	err := conf.Validate()
	require.NoError(t, err, "expected no validation error when sentry is not enabled")

	// If Sentry is enabled, then the environment is required
	conf.DSN = "https://something.ingest.sentry.io"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: environment must be configured when Sentry is enabled")

	conf.Environment = "test"
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration")
}
