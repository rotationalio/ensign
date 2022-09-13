package config_test

import (
	"os"
	"testing"

	"github.com/rotationalio/ensign/pkg/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"ENSIGN_MAINTENANCE":          "true",
	"ENSIGN_LOG_LEVEL":            "debug",
	"ENSIGN_CONSOLE_LOG":          "true",
	"ENSIGN_BIND_ADDR":            ":8888",
	"ENSIGN_MONITORING_ENABLED":   "true",
	"ENSIGN_MONITORING_BIND_ADDR": ":8889",
	"ENSIGN_MONITORING_NODE_ID":   "test1234",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after the test is complete.
	t.Cleanup(cleanupEnv())
	setEnv()

	conf, err := config.New()
	require.NoError(t, err, "could not process configuration from the environment")
	require.False(t, conf.IsZero(), "processed config should not be zero valued")
	require.Empty(t, conf.Path(), "no path should be available when loaded from environment")

	// Ensure configuration is correctly set from the environment
	require.True(t, conf.Maintenance)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Equal(t, testEnv["ENSIGN_BIND_ADDR"], conf.BindAddr)
	require.True(t, conf.Monitoring.Enabled)
	require.Equal(t, testEnv["ENSIGN_MONITORING_BIND_ADDR"], conf.Monitoring.BindAddr)
	require.Equal(t, testEnv["ENSIGN_MONITORING_NODE_ID"], conf.Monitoring.NodeID)
}

func TestLoadConfig(t *testing.T) {
	config.ResetLocalEnviron()
	conf, err := config.Load("testdata/config.yaml")
	require.NoError(t, err, "could not process configuration from file")
	require.False(t, conf.IsZero(), "processed config should not be zero valued")
	require.Equal(t, "testdata/config.yaml", conf.Path(), "path should match config file path")

	require.True(t, conf.Maintenance)
	require.Equal(t, zerolog.WarnLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Equal(t, "127.0.0.1:7778", conf.BindAddr)

	// TODO: Test JSON and TOML files by serializing the YAML config to those formats
	// in a temporary directory then reading them using config.Load and verifying them.

	// Ensure that an error is returned if the file cannot be opened or has an invalid ext
	_, err = config.Load("testdata/missing.json")
	require.Error(t, err, "should not be able to load from a missing config file")

	// TODO: Test a file that has an invalid extension by writing a temporary file.
}

func TestLoadConfigPriorities(t *testing.T) {
	// Data should be loaded from the configuration file if it is specified in the file,
	// from the environment or from a default if it doesn't and an enviornment variable
	// should take precedence if both are set.
	t.Cleanup(cleanupEnv())
	config.ResetLocalEnviron()

	// This test depends on partial_config.yaml containing the console log and bind addr
	// configurations and the environment containing the maintenance and bind addr. The
	// log level is omitted from both and should be set to the default.
	setEnv("ENSIGN_MAINTENANCE", "ENSIGN_BIND_ADDR")

	conf, err := config.Load("testdata/partial.yaml")
	require.NoError(t, err, "could not load configuration from file")

	require.True(t, conf.Maintenance)
	require.Equal(t, zerolog.InfoLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Equal(t, testEnv["ENSIGN_BIND_ADDR"], conf.BindAddr)
}

// Returns the current environment for the specified keys, or if no keys are specified
// then it returns the current environment for all keys in the testEnv variable.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv variable. If no keys are specified,
// then this function sets all environment variables from the testEnv.
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

// Cleanup helper function that can be run when the tests are complete to reset the
// environment back to its previous state before the test was run.
func cleanupEnv(keys ...string) func() {
	prevEnv := curEnv(keys...)
	return func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
