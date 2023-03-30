package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load the required settings from the environment, parse and
// validate them, loading defaults where necessary in preparation for running the
// Quarterdeck API service. This is the top-level config, all sub configurations need
// to be defined as properties of this Config.
type Config struct {
	Maintenance  bool                `default:"false"`                                    // $QUARTERDECK_MAINTENANCE
	BindAddr     string              `split_words:"true" default:":8088"`                 // $QUARTERDECK_BIND_ADDR
	Mode         string              `default:"release"`                                  // $QUARTERDECK_MODE
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info"`                  // $QUARTERDECK_LOG_LEVEL
	ConsoleLog   bool                `split_words:"true" default:"false"`                 // $QUARTERDECK_CONSOLE_LOG
	AllowOrigins []string            `split_words:"true" default:"http://localhost:3000"` // $QUARTERDECK_ALLOW_ORIGINS
	EmailURL     URLConfig           `split_words:"true"`
	SendGrid     emails.Config       `split_words:"false"`
	RateLimit    RateLimitConfig     `split_words:"true"`
	Database     DatabaseConfig
	Token        TokenConfig
	Sentry       sentry.Config
	processed    bool // set when the config is properly processed from the environment
}

type URLConfig struct {
	Base   string `split_words:"true" default:"https://rotational.app"`
	Verify string `split_words:"true" default:"/verify"`
	Invite string `split_words:"true" default:"/invite"`
}

type DatabaseConfig struct {
	URL      string `default:"sqlite3:////data/db/quarterdeck.db"` // $QUARTERDECK_DATABASE_URL
	ReadOnly bool   `split_words:"true" default:"false"`           // $QUARTERDECK_DATABASE_READ_ONLY
}

type TokenConfig struct {
	Keys            map[string]string `required:"false"`                      // $QUARTERDECK_TOKEN_KEYS
	Audience        string            `default:"https://rotational.app"`      // $QUARTERDECK_TOKEN_AUDIENCE
	RefreshAudience string            `required:"false"`                      // $QUARTERDECK_TOKEN_REFRESH_AUDIENCE
	Issuer          string            `default:"https://auth.rotational.app"` // $QUARTERDECK_TOKEN_ISSUER
	AccessDuration  time.Duration     `split_words:"true" default:"1h"`       // $QUARTERDECK_TOKEN_ACCESS_DURATION
	RefreshDuration time.Duration     `split_words:"true" default:"2h"`       // $QUARTERDECK_TOKEN_REFRESH_DURATION
	RefreshOverlap  time.Duration     `split_words:"true" default:"-15m"`     // $QUARTERDECK_TOKEN_REFRESH_OVERLAP
}

// Used by the Rate Limiter middleware
// Limit: represents the number of tokens that can be added to the token bucket per second
// Burst: maximum number of tokens/requests in a "token bucket" which is initially full
// empty token bucket results in failed requests
// TTL: //number of minutes before an IP is removed from the ratelimiter map
// NOTE: If Burst is not included the config, then all requests will be rejected!
// The Validate() method checks to see if the all required values for the RateLimiter middleware
// are populated and will fail startup if they are not populated
type RateLimitConfig struct {
	PerSecond float64       `default:"10" split_words:"true"`
	Burst     int           `default:"30"`
	TTL       time.Duration `default:"5m"`
}

// New loads and parses the config from the environment and validates it, marking it as
// processed so that external users can determine if the config is ready for use. This
// should be the only way Config objects are created for use in the application.
func New() (conf Config, err error) {
	if err = envconfig.Process("quarterdeck", &conf); err != nil {
		return Config{}, err
	}

	// Ensure the Sentry release is named correctly
	if conf.Sentry.Release == "" {
		conf.Sentry.Release = fmt.Sprintf("quarterdeck@%s", pkg.Version())
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Returns true if the config has not been correctly processed from the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Mark a manually constructed config as processed as long as it is valid.
func (c Config) Mark() (_ Config, err error) {
	if err = c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

// Custom validations are added here, particularly validations that require one or more
// fields to be processed before the validation occurs.
// NOTE: ensure that all nested config validation methods are called here.
func (c Config) Validate() (err error) {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("invalid configuration: %q is not a valid gin mode", c.Mode)
	}

	if err = c.EmailURL.Validate(); err != nil {
		return err
	}

	if err = c.SendGrid.Validate(); err != nil {
		return err
	}

	if err = c.Sentry.Validate(); err != nil {
		return err
	}

	if (RateLimitConfig{}) == c.RateLimit {
		return fmt.Errorf("invalid configuration: RateLimitConfig needs to be populated")
	}

	if c.RateLimit.PerSecond == 0.00 {
		return fmt.Errorf("invalid configuration: RateLimitConfig.PerSecond needs to be populated and must be a nonzero value")
	}

	if c.RateLimit.Burst == 0 {
		return fmt.Errorf("invalid configuration: RateLimitConfig.Burst needs to be populated and must be a nonzero value")
	}

	if c.RateLimit.TTL == 0*time.Second {
		return fmt.Errorf("invalid configuration: RateLimitConfig.TTL needs to be populated and must be a nonzero value")
	}

	return nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

// Returns true if the allow origins slice contains one entry that is a "*"
func (c Config) AllowAllOrigins() bool {
	if len(c.AllowOrigins) == 1 && c.AllowOrigins[0] == "*" {
		return true
	}
	return false
}

func (c URLConfig) Validate() (err error) {
	if c.Base == "" {
		return fmt.Errorf("invalid email url configuration: base URL is required")
	}

	if c.Invite == "" {
		return fmt.Errorf("invalid email url configuration: invite path is required")
	}

	if c.Verify == "" {
		return fmt.Errorf("invalid email url configuration: verify path is required")
	}

	return nil
}

// Construct an invite URL from the token.
func (c URLConfig) InviteURL(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	base, _ := url.Parse(c.Base)
	url := base.ResolveReference(&url.URL{Path: c.Invite, RawQuery: url.Values{"token": []string{token}}.Encode()})
	return url.String(), nil
}

// Construct a verify URL from the token.
func (c URLConfig) VerifyURL(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	base, _ := url.Parse(c.Base)
	url := base.ResolveReference(&url.URL{Path: c.Verify, RawQuery: url.Values{"token": []string{token}}.Encode()})
	return url.String(), nil
}
