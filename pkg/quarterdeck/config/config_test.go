package config_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// The test environment for all config tests, manipulated using curEnv and setEnv
var testEnv = map[string]string{
	"QUARTERDECK_MAINTENANCE":                "false",
	"QUARTERDECK_BIND_ADDR":                  ":3636",
	"QUARTERDECK_MODE":                       gin.TestMode,
	"QUARTERDECK_LOG_LEVEL":                  "error",
	"QUARTERDECK_CONSOLE_LOG":                "true",
	"QUARTERDECK_ALLOW_ORIGINS":              "http://localhost:8888,http://localhost:8080",
	"QUARTERDECK_EMAIL_URL_BASE":             "http://localhost:8888",
	"QUARTERDECK_EMAIL_URL_INVITE":           "/invite",
	"QUARTERDECK_EMAIL_URL_VERIFY":           "/verify",
	"QUARTERDECK_SENDGRID_API_KEY":           "SG.1234",
	"QUARTERDECK_SENDGRID_FROM_EMAIL":        "test@example.com",
	"QUARTERDECK_SENDGRID_ADMIN_EMAIL":       "admin@example.com",
	"QUARTERDECK_SENDGRID_ENSIGN_LIST_ID":    "1234",
	"QUARTERDECK_RATE_LIMIT_ENABLED":         "false",
	"QUARTERDECK_RATE_LIMIT_PER_SECOND":      "20",
	"QUARTERDECK_RATE_LIMIT_BURST":           "100",
	"QUARTERDECK_RATE_LIMIT_TTL":             "1h",
	"QUARTERDECK_REPORTING_ENABLE_DAILY_PLG": "true",
	"QUARTERDECK_REPORTING_DOMAIN":           "ensign.world",
	"QUARTERDECK_REPORTING_DASHBOARD_URL":    "https://grafana.rotational.dev",
	"QUARTERDECK_DATABASE_URL":               "sqlite3:///test.db",
	"QUARTERDECK_DATABASE_READ_ONLY":         "true",
	"QUARTERDECK_TOKEN_KEYS":                 "01GECSDK5WJ7XWASQ0PMH6K41K:testdata/01GECSDK5WJ7XWASQ0PMH6K41K.pem,01GECSJGDCDN368D0EENX23C7R:testdata/01GECSJGDCDN368D0EENX23C7R.pem",
	"QUARTERDECK_TOKEN_AUDIENCE":             "http://localhost:8888",
	"QUARTERDECK_TOKEN_ISSUER":               "http://localhost:1025",
	"QUARTERDECK_TOKEN_ACCESS_DURATION":      "5m",
	"QUARTERDECK_TOKEN_REFRESH_DURATION":     "10m",
	"QUARTERDECK_TOKEN_REFRESH_OVERLAP":      "-2m",
	"QUARTERDECK_SENTRY_DSN":                 "http://testing.sentry.test/1234",
	"QUARTERDECK_SENTRY_SERVER_NAME":         "tnode",
	"QUARTERDECK_SENTRY_ENVIRONMENT":         "testing",
	"QUARTERDECK_SENTRY_RELEASE":             "", // This should always be empty!
	"QUARTERDECK_SENTRY_TRACK_PERFORMANCE":   "true",
	"QUARTERDECK_SENTRY_SAMPLE_RATE":         "0.95",
	"QUARTERDECK_SENTRY_DEBUG":               "true",
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
	require.Equal(t, testEnv["QUARTERDECK_EMAIL_URL_BASE"], conf.EmailURL.Base)
	require.Equal(t, testEnv["QUARTERDECK_EMAIL_URL_INVITE"], conf.EmailURL.Invite)
	require.Equal(t, testEnv["QUARTERDECK_EMAIL_URL_VERIFY"], conf.EmailURL.Verify)
	require.Equal(t, testEnv["QUARTERDECK_SENDGRID_API_KEY"], conf.SendGrid.APIKey)
	require.Equal(t, testEnv["QUARTERDECK_SENDGRID_FROM_EMAIL"], conf.SendGrid.FromEmail)
	require.Equal(t, testEnv["QUARTERDECK_SENDGRID_ADMIN_EMAIL"], conf.SendGrid.AdminEmail)
	require.Equal(t, testEnv["QUARTERDECK_SENDGRID_ENSIGN_LIST_ID"], conf.SendGrid.EnsignListID)
	require.True(t, conf.Reporting.EnableDailyPLG)
	require.Equal(t, testEnv["QUARTERDECK_REPORTING_DOMAIN"], conf.Reporting.Domain)
	require.Equal(t, testEnv["QUARTERDECK_REPORTING_DASHBOARD_URL"], conf.Reporting.DashboardURL)
	require.Equal(t, testEnv["QUARTERDECK_DATABASE_URL"], conf.Database.URL)
	require.True(t, conf.Database.ReadOnly)
	require.Len(t, conf.Token.Keys, 2)
	require.Equal(t, testEnv["QUARTERDECK_TOKEN_AUDIENCE"], conf.Token.Audience)
	require.Equal(t, testEnv["QUARTERDECK_TOKEN_ISSUER"], conf.Token.Issuer)
	require.Equal(t, 5*time.Minute, conf.Token.AccessDuration)
	require.Equal(t, 10*time.Minute, conf.Token.RefreshDuration)
	require.Equal(t, -2*time.Minute, conf.Token.RefreshOverlap)
	require.Equal(t, testEnv["QUARTERDECK_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["QUARTERDECK_SENTRY_SERVER_NAME"], conf.Sentry.ServerName)
	require.Equal(t, testEnv["QUARTERDECK_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.True(t, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.95, conf.Sentry.SampleRate)
	require.True(t, conf.Sentry.Debug)
	require.False(t, conf.RateLimit.Enabled)
	require.Equal(t, 20.00, conf.RateLimit.PerSecond)
	require.Equal(t, 100, conf.RateLimit.Burst)
	require.Equal(t, 60*time.Minute, conf.RateLimit.TTL)

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

	// Should not be able to mark a config if the email URL values are not set
	conf.Mode = gin.ReleaseMode
	conf.EmailURL.Invite = "/invite"
	conf.EmailURL.Verify = "/verify"
	conf, err = conf.Mark()
	require.EqualError(t, err, "invalid email url configuration: base URL is required", "expected EmailURL validation error")

	// Should not be able to mark a config that does not contain required values for the RateLimiter middleware
	conf.EmailURL.Base = "https://localhost:8080"
	conf.RateLimit.Enabled = true

	conf, err = conf.Mark()
	require.EqualError(t, err, "invalid configuration: RateLimitConfig.PerSecond needs to be populated and must be a nonzero value")

	conf.RateLimit.Burst = 120
	conf.RateLimit.PerSecond = 20.00
	conf.RateLimit.TTL = 5 * time.Minute

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

func TestEmailURL(t *testing.T) {
	conf := &config.URLConfig{
		Base:   "https://auth.rotational.app",
		Invite: "/invite",
		Verify: "/verify",
	}
	require.NoError(t, conf.Validate(), "expected config to be valid")

	// Token is required for invite URL
	_, err := conf.InviteURL("")
	require.EqualError(t, err, "token is required", "expected error when token is empty")

	// Constructing a valid invite URL
	inviteURL, err := conf.InviteURL("1234")
	require.NoError(t, err, "expected no error when constructing invite URL")
	require.Equal(t, "https://auth.rotational.app/invite?token=1234", inviteURL, "expected invite URL to be constructed")

	// Token is required for verify URL
	_, err = conf.VerifyURL("")
	require.EqualError(t, err, "token is required", "expected error when token is empty")

	// Constructing a valid verify URL
	verifyURL, err := conf.VerifyURL("1234")
	require.NoError(t, err, "expected no error when constructing verify URL")
	require.Equal(t, "https://auth.rotational.app/verify?token=1234", verifyURL, "expected verify URL to be constructed")

	// Base URL is required
	conf.Base = ""
	require.EqualError(t, conf.Validate(), "invalid email url configuration: base URL is required", "expected base URL validation error")

	// Invite URL is required
	conf.Base = "https://auth.rotational.app"
	conf.Invite = ""
	require.EqualError(t, conf.Validate(), "invalid email url configuration: invite path is required", "expected invite URL validation error")

	// Verify URL is required
	conf.Invite = "/invite"
	conf.Verify = ""
	require.EqualError(t, conf.Validate(), "invalid email url configuration: verify path is required", "expected verify URL validation error")
}

func TestRateLimitConfigValidate(t *testing.T) {
	conf := config.RateLimitConfig{Enabled: false}
	require.NoError(t, conf.Validate(), "disabled config should be valid")

	// Burst must be greater than 0
	conf = config.RateLimitConfig{
		Enabled:   true,
		Burst:     0,
		PerSecond: 20.00,
		TTL:       5 * time.Minute,
	}
	require.EqualError(t, conf.Validate(), "invalid configuration: RateLimitConfig.Burst needs to be populated and must be a nonzero value", "burst must be greater than 0")

	// PerSecond must be greater than 0
	conf = config.RateLimitConfig{
		Enabled:   true,
		Burst:     120,
		PerSecond: 0.00,
		TTL:       5 * time.Minute,
	}
	require.EqualError(t, conf.Validate(), "invalid configuration: RateLimitConfig.PerSecond needs to be populated and must be a nonzero value", "per-second must be greater than 0")

	// TTL must be greater than 0
	conf = config.RateLimitConfig{
		Enabled:   true,
		Burst:     120,
		PerSecond: 20.00,
		TTL:       0,
	}
	require.EqualError(t, conf.Validate(), "invalid configuration: RateLimitConfig.TTL needs to be populated and must be a nonzero value", "ttl must be greater than 0")

	// Valid configuration
	conf = config.RateLimitConfig{
		Enabled:   true,
		Burst:     120,
		PerSecond: 20.00,
		TTL:       5 * time.Minute,
	}
	require.NoError(t, conf.Validate(), "expected valid configuration")
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
