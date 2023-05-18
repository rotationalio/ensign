package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"ENSIGN_MAINTENANCE":               "true",
	"ENSIGN_LOG_LEVEL":                 "debug",
	"ENSIGN_CONSOLE_LOG":               "true",
	"ENSIGN_BIND_ADDR":                 ":8888",
	"ENSIGN_META_TOPIC_ENABLED":        "true",
	"ENSIGN_META_TOPIC_TOPIC_NAME":     "ensign.testing",
	"ENSIGN_META_TOPIC_CLIENT_ID":      "test1234",
	"ENSIGN_META_TOPIC_CLIENT_SECRET":  "supersecret",
	"ENSIGN_META_TOPIC_ENDPOINT":       "ensign.ninja:443",
	"ENSIGN_META_TOPIC_AUTH_URL":       "https://auth.ensign.world",
	"ENSIGN_MONITORING_ENABLED":        "true",
	"ENSIGN_MONITORING_BIND_ADDR":      ":8889",
	"ENSIGN_MONITORING_NODE_ID":        "test1234",
	"ENSIGN_STORAGE_READ_ONLY":         "true",
	"ENSIGN_STORAGE_DATA_PATH":         "/data/db",
	"ENSIGN_AUTH_KEYS_URL":             "http://localhost:8088/.well-known/jwks.json",
	"ENSIGN_AUTH_AUDIENCE":             "http://localhost:3000",
	"ENSIGN_AUTH_ISSUER":               "http://localhost:8088",
	"ENSIGN_AUTH_MIN_REFRESH_INTERVAL": "10m",
	"ENSIGN_SENTRY_DSN":                "http://testing.sentry.test/1234",
	"ENSIGN_SENTRY_SERVER_NAME":        "test1234",
	"ENSIGN_SENTRY_ENVIRONMENT":        "testing",
	"ENSIGN_SENTRY_RELEASE":            "", // This should always be empty!
	"ENSIGN_SENTRY_TRACK_PERFORMANCE":  "true",
	"ENSIGN_SENTRY_SAMPLE_RATE":        "0.95",
	"ENSIGN_SENTRY_DEBUG":              "true",
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
	require.True(t, conf.MetaTopic.Enabled)
	require.Equal(t, testEnv["ENSIGN_META_TOPIC_TOPIC_NAME"], conf.MetaTopic.TopicName)
	require.Equal(t, testEnv["ENSIGN_META_TOPIC_CLIENT_ID"], conf.MetaTopic.ClientID)
	require.Equal(t, testEnv["ENSIGN_META_TOPIC_CLIENT_SECRET"], conf.MetaTopic.ClientSecret)
	require.Equal(t, testEnv["ENSIGN_META_TOPIC_ENDPOINT"], conf.MetaTopic.Endpoint)
	require.Equal(t, testEnv["ENSIGN_META_TOPIC_AUTH_URL"], conf.MetaTopic.AuthURL)
	require.True(t, conf.Monitoring.Enabled)
	require.Equal(t, testEnv["ENSIGN_MONITORING_BIND_ADDR"], conf.Monitoring.BindAddr)
	require.Equal(t, testEnv["ENSIGN_MONITORING_NODE_ID"], conf.Monitoring.NodeID)
	require.False(t, conf.Storage.Testing)
	require.True(t, conf.Storage.ReadOnly)
	require.Equal(t, testEnv["ENSIGN_STORAGE_DATA_PATH"], conf.Storage.DataPath)
	require.Equal(t, testEnv["ENSIGN_AUTH_KEYS_URL"], conf.Auth.KeysURL)
	require.Equal(t, testEnv["ENSIGN_AUTH_AUDIENCE"], conf.Auth.Audience)
	require.Equal(t, testEnv["ENSIGN_AUTH_ISSUER"], conf.Auth.Issuer)
	require.Equal(t, 10*time.Minute, conf.Auth.MinRefreshInterval)
	require.Equal(t, testEnv["ENSIGN_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["ENSIGN_SENTRY_SERVER_NAME"], conf.Sentry.ServerName)
	require.Equal(t, testEnv["ENSIGN_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.True(t, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.95, conf.Sentry.SampleRate)
	require.True(t, conf.Sentry.Debug)

	// Ensure the sentry release is correctly set
	require.True(t, strings.HasPrefix(conf.Sentry.GetRelease(), "ensign@"))
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
	setEnv("ENSIGN_MAINTENANCE", "ENSIGN_BIND_ADDR", "ENSIGN_STORAGE_DATA_PATH")

	conf, err := config.Load("testdata/partial.yaml")
	require.NoError(t, err, "could not load configuration from file")

	require.True(t, conf.Maintenance)
	require.Equal(t, zerolog.InfoLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Equal(t, testEnv["ENSIGN_BIND_ADDR"], conf.BindAddr)
}

func TestValidateMetaTopicConfig(t *testing.T) {
	conf := config.MetaTopicConfig{
		Enabled: false,
	}

	err := conf.Validate()
	require.NoError(t, err, "disabled config should be valid")

	conf.Enabled = true
	require.EqualError(t, conf.Validate(), "invalid meta topic config: missing topic name")

	conf.TopicName = "foo"
	require.EqualError(t, conf.Validate(), "invalid meta topic config: missing client id or secret")

	conf.ClientID = "foo"
	require.EqualError(t, conf.Validate(), "invalid meta topic config: missing client id or secret")

	conf.ClientSecret = "foo"
	require.NoError(t, conf.Validate(), "topic name, client id, client secret are all that's required")
}

func TestStoragePaths(t *testing.T) {
	dir := t.TempDir()
	conf := config.StorageConfig{
		ReadOnly: false,
		DataPath: dir,
	}

	expectedMetaPath := filepath.Join(dir, "metadata")
	expectedEventPath := filepath.Join(dir, "events")

	require.NoDirExists(t, expectedMetaPath, "expected no metadata directory to exist")
	require.NoDirExists(t, expectedEventPath, "expected no event log directory to exist")

	actualMetaPath, err := conf.MetaPath()
	require.NoError(t, err, "could not get the metadata path")
	require.Equal(t, expectedMetaPath, actualMetaPath)

	actualEventPath, err := conf.EventPath()
	require.NoError(t, err, "could not get the event log path")
	require.Equal(t, expectedEventPath, actualEventPath)

	require.DirExists(t, expectedMetaPath, "expected metadata directory to exist")
	require.DirExists(t, expectedEventPath, "expected event log directory to exist")

	// Should be able to get the paths again without recreating the dirs
	_, err = conf.MetaPath()
	require.NoError(t, err, "could not get meta path")

	_, err = conf.EventPath()
	require.NoError(t, err, "could not get event path")

	// Require the paths to be directories not files
	require.NoError(t, os.Remove(expectedMetaPath))
	require.NoError(t, os.Remove(expectedEventPath))

	_, err = os.Create(expectedMetaPath)
	require.NoError(t, err, "could not create meta file")
	_, err = os.Create(expectedEventPath)
	require.NoError(t, err, "could not create event file")

	_, err = conf.MetaPath()
	require.Error(t, err, "expected is not a directory error")

	_, err = conf.EventPath()
	require.Error(t, err, "expected is not a directory error")
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
