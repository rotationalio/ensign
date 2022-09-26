package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// The test environment for all config tests, manipulated using curEnv and setEnv
var testEnv = map[string]string{
	"QUARTERDECK_MAINTENANCE":              "false",
	"QUARTERDECK_BIND_ADDR":                ":3636",
	"QUARTERDECK_MODE":                     gin.TestMode,
	"QUARTERDECK_LOG_LEVEL":                "error",
	"QUARTERDECK_CONSOLE_LOG":              "true",
	"QUARTERDECK_ALLOW_ORIGINS":            "http://localhost:8888,http://localhost:8080",
	"QUARTERDECK_SENTRY_DSN":               "http://testing.sentry.test/1234",
	"QUARTERDECK_SENTRY_SERVER_NAME":       "tnode",
	"QUARTERDECK_SENTRY_ENVIRONMENT":       "testing",
	"QUARTERDECK_SENTRY_RELEASE":           "", // This should always be empty!
	"QUARTERDECK_SENTRY_TRACK_PERFORMANCE": "true",
	"QUARTERDECK_SENTRY_SAMPLE_RATE":       "0.95",
	"QUARTERDECK_SENTRY_DEBUG":             "true",
}

func TestConfig(t *testing.T) {
	// Set the required environment variables and cleanup after.
	prevEnv := curEnv()
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})
	setEnv()

	// At this point in the test, the environment should contain testEnv
	conf, err := config.New()
	require.NoError(t, err, "could not create a default config")
	require.False(t, conf.IsZero(), "default config should be processed")

	// Test the configuration
	require.False(t, conf.Maintenance)
	require.Equal(t, testEnv["QUARTERDECK_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, testEnv["QUARTERDECK_MODE"], conf.Mode)
	require.Equal(t, zerolog.ErrorLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Len(t, conf.AllowOrigins, 2)
	require.Equal(t, testEnv["QUARTERDECK_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["QUARTERDECK_SENTRY_SERVER_NAME"], conf.Sentry.ServerName)
	require.Equal(t, testEnv["QUARTERDECK_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.True(t, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.95, conf.Sentry.SampleRate)
	require.True(t, conf.Sentry.Debug)

	// Ensure the sentry release is configured correctly
	require.True(t, strings.HasPrefix(conf.Sentry.GetRelease(), "quarterdeck@"))
}

func TestValidation(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err, "could not create default config")

	modes := []string{gin.ReleaseMode, gin.DebugMode, gin.TestMode}
	for _, mode := range modes {
		conf.Mode = mode
		require.NoError(t, conf.Validate(), "expected config to be valid in %q mode", mode)
	}

	// Ensure conf is invalid on wrong mode
	conf.Mode = "invalid"
	require.EqualError(t, conf.Validate(), `invalid configuration: "invalid" is not a valid gin mode`, "expected gin mode validation error")
}

func TestIsZero(t *testing.T) {
	// An empty config should always return IsZero
	require.True(t, config.Config{}.IsZero(), "an empty config should always be zero valued")

	// A processed config should not be zero valued
	conf, err := config.New()
	require.NoError(t, err, "should have been able to load the config")
	require.False(t, conf.IsZero(), "expected a processed config to be non-zero valued")

	// Custom config not processed
	conf = config.Config{
		Maintenance: false,
		BindAddr:    "127.0.0.1:0",
		LogLevel:    logger.LevelDecoder(zerolog.TraceLevel),
		Mode:        "invalid",
	}
	require.True(t, config.Config{}.IsZero(), "a non-empty config that isn't marked will be zero valued")

	// Should not be able to mark a custom config that is invalid
	conf, err = conf.Mark()
	require.EqualError(t, err, `invalid configuration: "invalid" is not a valid gin mode`, "expected gin mode validation error")

	// Should be able to mark a custom config that is valid as processed
	conf.Mode = gin.ReleaseMode
	conf, err = conf.Mark()
	require.NoError(t, err, "should be able to mark a valid config")
	require.False(t, conf.IsZero(), "a marked config should not be zero-valued")
}

func TestAllowAllOrigins(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err, "could not create default configuration")
	require.Equal(t, []string{"http://localhost:3000"}, conf.AllowOrigins, "allow origins should be localhost by default")
	require.False(t, conf.AllowAllOrigins(), "expected allow all origins to be false by default")

	conf.AllowOrigins = []string{"https://ensign.rotational.dev", "https://ensign.io"}
	require.False(t, conf.AllowAllOrigins(), "expected allow all origins to be false when allow origins is set")

	conf.AllowOrigins = []string{}
	require.False(t, conf.AllowAllOrigins(), "expected allow all origins to be false when allow origins is empty")

	conf.AllowOrigins = []string{"*"}
	require.True(t, conf.AllowAllOrigins(), "expect allow all origins to be true when * is set")
}

// Returns the current environment for the specified keys, or if no keys are specified
// then returns the current environment for all keys in testEnv.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)

	if len(keys) > 0 {
		// Process the keys passed in by the user
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		// Process all the keys in testEnv
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv, if no keys are specified then sets
// all environment variables that are specified in the testEnv.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}
