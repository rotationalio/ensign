package config_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// Test environment for all config tests that is manipulated by curEnv and setEnv
var testEnv = map[string]string{
	"TENANT_MAINTENANCE":              "false",
	"TENANT_BIND_ADDR":                ":3636",
	"TENANT_MODE":                     gin.TestMode,
	"TENANT_LOG_LEVEL":                "error",
	"TENANT_CONSOLE_LOG":              "true",
	"TENANT_ALLOW_ORIGINS":            "http://localhost:8888,http://localhost:8080",
	"TENANT_AUTH_KEYS_URL":            "http://localhost:8080/.well-known/jwks.json",
	"TENANT_AUTH_AUDIENCE":            "audience",
	"TENANT_AUTH_ISSUER":              "issuer",
	"TENANT_AUTH_COOKIE_DOMAIN":       "localhost",
	"TENANT_DATABASE_URL":             "trtl://localhost:4436",
	"TENANT_DATABASE_INSECURE":        "true",
	"TENANT_DATABASE_CERT_PATH":       "path/to/certs.pem",
	"TENANT_DATABASE_POOL_PATH":       "path/to/pool.pem",
	"TENANT_QUARTERDECK_URL":          "https://localhost:8080",
	"TENANT_ENSIGN_ENDPOINT":          "localhost:5356",
	"TENANT_ENSIGN_AUTH_URL":          "http://localhost:8080",
	"TENANT_ENSIGN_CLIENT_ID":         "client-id",
	"TENANT_ENSIGN_CLIENT_SECRET":     "client-secret",
	"TENANT_ENSIGN_TOPIC_NAME":        "topic-id",
	"TENANT_ENSIGN_INSECURE":          "true",
	"TENANT_ENSIGN_NO_AUTHENTICATION": "true",
	"TENANT_ENSIGN_WAIT_FOR_READY":    "10m",
	"TENANT_META_TOPIC_ENABLED":       "true",
	"TENANT_META_TOPIC_TOPIC_NAME":    "ensign.metatopic.topics",
	"TENANT_META_TOPIC_CLIENT_ID":     "client-id",
	"TENANT_META_TOPIC_CLIENT_SECRET": "client-secret",
	"TENANT_META_TOPIC_ENDPOINT":      "ensign.rotational.app:443",
	"TENANT_META_TOPIC_AUTH_URL":      "https://auth.rotational.app",
	"TENANT_SENDGRID_API_KEY":         "SG.testing.123-331-test",
	"TENANT_SENDGRID_FROM_EMAIL":      "test@example.com",
	"TENANT_SENDGRID_ADMIN_EMAIL":     "admin@example.com",
	"TENANT_SENDGRID_ENSIGN_LIST_ID":  "cb385e60-b43c-4db2-89ad-436ec277eacb",
	"TENANT_SENTRY_DSN":               "http://testing.sentry.test/1234",
	"TENANT_SENTRY_SERVER_NAME":       "tnode",
	"TENANT_SENTRY_ENVIRONMENT":       "testing",
	"TENANT_SENTRY_RELEASE":           "", // This should always be empty
	"TENANT_SENTRY_TRACK_PERFORMANCE": "true",
	"TENANT_SENTRY_SAMPLE_RATE":       "0.95",
	"TENANT_SENTRY_DEBUG":             "true",
}

func TestConfig(t *testing.T) {
	// Sets the required environment variables and cleanup after.
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

	// The environment should contain the testEnv at this point in the test.
	conf, err := config.New()
	require.NoError(t, err, "could not create a default config")
	require.False(t, conf.IsZero(), "default config should be processed")

	// Tests the configuration
	require.False(t, conf.Maintenance)
	require.Equal(t, testEnv["TENANT_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, testEnv["TENANT_MODE"], conf.Mode)
	require.Equal(t, zerolog.ErrorLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Len(t, conf.AllowOrigins, 2)
	require.Equal(t, testEnv["TENANT_AUTH_KEYS_URL"], conf.Auth.KeysURL)
	require.Equal(t, testEnv["TENANT_AUTH_AUDIENCE"], conf.Auth.Audience)
	require.Equal(t, testEnv["TENANT_AUTH_ISSUER"], conf.Auth.Issuer)
	require.Equal(t, testEnv["TENANT_AUTH_COOKIE_DOMAIN"], conf.Auth.CookieDomain)
	require.Equal(t, testEnv["TENANT_DATABASE_URL"], conf.Database.URL)
	require.True(t, conf.Database.Insecure)
	require.Equal(t, testEnv["TENANT_DATABASE_CERT_PATH"], conf.Database.CertPath)
	require.Equal(t, testEnv["TENANT_DATABASE_POOL_PATH"], conf.Database.PoolPath)
	require.Equal(t, testEnv["TENANT_ENSIGN_ENDPOINT"], conf.Ensign.Endpoint)
	require.True(t, conf.Ensign.Enabled)
	require.Equal(t, testEnv["TENANT_QUARTERDECK_URL"], conf.Quarterdeck.URL)
	require.Equal(t, testEnv["TENANT_ENSIGN_ENDPOINT"], conf.Ensign.Endpoint)
	require.Equal(t, testEnv["TENANT_ENSIGN_CLIENT_ID"], conf.Ensign.ClientID)
	require.Equal(t, testEnv["TENANT_ENSIGN_CLIENT_SECRET"], conf.Ensign.ClientSecret)
	require.True(t, conf.Ensign.Insecure)
	require.True(t, conf.Ensign.NoAuthentication)
	require.Equal(t, 10*time.Minute, conf.Ensign.WaitForReady)
	require.True(t, conf.MetaTopic.Enabled)
	require.Equal(t, testEnv["TENANT_META_TOPIC_TOPIC_NAME"], conf.MetaTopic.TopicName)
	require.Equal(t, testEnv["TENANT_META_TOPIC_CLIENT_ID"], conf.MetaTopic.ClientID)
	require.Equal(t, testEnv["TENANT_META_TOPIC_CLIENT_SECRET"], conf.MetaTopic.ClientSecret)
	require.Equal(t, testEnv["TENANT_META_TOPIC_ENDPOINT"], conf.MetaTopic.Endpoint)
	require.Equal(t, testEnv["TENANT_META_TOPIC_AUTH_URL"], conf.MetaTopic.AuthURL)
	require.Equal(t, testEnv["TENANT_SENDGRID_API_KEY"], conf.SendGrid.APIKey)
	require.Equal(t, testEnv["TENANT_SENDGRID_FROM_EMAIL"], conf.SendGrid.FromEmail)
	require.Equal(t, testEnv["TENANT_SENDGRID_ADMIN_EMAIL"], conf.SendGrid.AdminEmail)
	require.Equal(t, testEnv["TENANT_SENDGRID_ENSIGN_LIST_ID"], conf.SendGrid.EnsignListID)
	require.Equal(t, testEnv["TENANT_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["TENANT_SENTRY_SERVER_NAME"], conf.Sentry.ServerName)
	require.Equal(t, testEnv["TENANT_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.True(t, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.95, conf.Sentry.SampleRate)
	require.True(t, conf.Sentry.Debug)

	// Ensures the Sentry release is cocnfigured correctly
	require.True(t, strings.HasPrefix(conf.Sentry.GetRelease(), "tenant@"))
}

func TestValidation(t *testing.T) {
	setEnv("TESTING")
	conf, err := config.New()
	require.NoError(t, err, "could not create default config")

	modes := []string{gin.ReleaseMode, gin.DebugMode, gin.TestMode}
	for _, mode := range modes {
		conf.Mode = mode
		require.NoError(t, conf.Validate(), "expected config to be valid in %q mode", mode)
	}

	// Ensures that conf is invalid on wrong mode
	conf.Mode = "invalid"
	require.EqualError(t, conf.Validate(), `invalid configuration: "invalid" is not a valid gin mode`, "expected gin mode validation error")
}

func TestIsZero(t *testing.T) {
	// An empty config should always return IsZero
	require.True(t, config.Config{}.IsZero(), "an empty config should always be zero valued")

	// A processed config should not have a zero value
	setEnv("TESTING")
	conf, err := config.New()
	require.NoError(t, err, "should have been able to load the config")
	require.False(t, conf.IsZero(), "expected a processed config to be non-zero valued")

	// Custom config is not processed
	conf = config.Config{
		Maintenance: false,
		BindAddr:    "127.0.0.1:0",
		LogLevel:    logger.LevelDecoder(zerolog.TraceLevel),
		Mode:        "invalid",
		Database: config.DatabaseConfig{
			URL:      "trtl://localhost:4437",
			Insecure: true,
		},
	}

	require.True(t, config.Config{}.IsZero(), "a non-empty config that isn't marked will be zero valued")

	// Should not be able to mark a custom config that is invalid
	conf, err = conf.Mark()
	require.EqualError(t, err, `invalid configuration: "invalid" is not a valid gin mode`, "expected gin mode validation error")

	// Should not be able to mark a custom config that is valid as processed
	conf.Mode = gin.ReleaseMode
	conf, err = conf.Mark()
	require.NoError(t, err, "should be able to mark a valid config")
	require.False(t, conf.IsZero(), "a marked config should not be zero-valued")
}

func TestAllowAllOrigins(t *testing.T) {
	setEnv("TESTING")
	conf, err := config.New()
	require.NoError(t, err, "could not create default configuration")
	require.Equal(t, []string{"http://localhost:3000"}, conf.AllowOrigins, "allow origins should be localhost by default")
	require.False(t, conf.AllowAllOrigins(), "expected allow all origins to be false by default")

	conf.AllowOrigins = []string{"https://ensign.rotational.dev", "https://ensign.io"}
	require.False(t, conf.AllowAllOrigins(), "expected allow all origins to be false when allow origins is set")

	conf.AllowOrigins = []string{}
	require.False(t, conf.AllowAllOrigins(), "expected allow all origins to be false when allow origins is empty")

	conf.AllowOrigins = []string{"*"}
	require.True(t, conf.AllowAllOrigins(), "expected allow all origins to be true when * is set")
}

func TestAuth(t *testing.T) {
	// TODO: test AuthConfig validation
}

func TestDatabase(t *testing.T) {
	// TODO: test DatabaseConfig validation
}

func TestSDK(t *testing.T) {
	conf := &config.SDKConfig{
		Enabled:      true,
		Endpoint:     "ensign.io:5356",
		Insecure:     true,
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		AuthURL:      "auth.rotational.app",
	}

	require.Len(t, conf.ClientOptions(), 3, "config should return client options")

	// Client ID is required when using authentication
	conf.ClientID = ""
	require.EqualError(t, conf.Validate(), "invalid meta topic config: missing client id", "config should be invalid if client id is empty")

	// Client secret is required when using authentication
	conf.ClientID = "client-id"
	conf.ClientSecret = ""
	require.EqualError(t, conf.Validate(), "invalid meta topic config: missing client secret", "config should be invalid if client secret is empty")

	// Auth URL is required when using authentication
	conf.ClientSecret = "client-secret"
	conf.AuthURL = ""
	require.EqualError(t, conf.Validate(), "invalid meta topic config: missing auth url", "config should be invalid if auth url is empty")

	// Valid configuration
	conf.ClientID = "client-id"
	conf.ClientSecret = "client-secret"
	conf.AuthURL = "auth.rotational.app"
	require.NoError(t, conf.Validate(), "config should be valid when topic ID is set")

	// Client ID and secret are not required if no authentication is set
	conf.ClientID = ""
	conf.ClientSecret = ""
	conf.AuthURL = ""
	conf.NoAuthentication = true
	require.NoError(t, conf.Validate(), "client id and secret should not be required if no authentication is set")

	// Client ID and secret are not required in testing mode
	conf.Testing = true
	conf.NoAuthentication = false
	require.NoError(t, conf.Validate(), "client id and secret should not be required in testing mode")

	// Disabled config should be valid
	empty := &config.SDKConfig{Enabled: false}
	require.NoError(t, empty.Validate(), "disabled config should be valid")
	require.Len(t, empty.ClientOptions(), 0, "disabled config should retun no options")
}

func TestQuarterdeck(t *testing.T) {
	conf := &config.QuarterdeckConfig{
		URL: "trtl://localhost:4437",
	}
	require.Error(t, conf.Validate(), "config should be invalid when URL scheme is not http(s)")

	conf.URL = "http://localhost:8088"
	require.NoError(t, conf.Validate(), "config should be valid when URL scheme is http")

	conf.URL = "https://localhost:8088"
	require.NoError(t, conf.Validate(), "config should be valid when URL scheme is https")
}

// Returns the current environment for the specified keys. If no keys are
// specified then returns the current environment for all keys in the testEnv.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)

	if len(keys) > 0 {
		// Processes the keys passed in by the user
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		// Processes all keys in the testEnv
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv. If no keys are specified
// then sets all environment variables that are specified in the testEnv.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			} else {
				if key == "TESTING" {
					os.Setenv("TENANT_ENSIGN_TESTING", "true")
					os.Setenv("TENANT_META_TOPIC_TESTING", "true")
				}
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}
