package config

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load required settings from the environment, parses
// and validates them, and loads defaults where necessary in preparation for running
// the Tenant API service. This is the top-level config, any sub configurations
// will need to be defined as properties of this Config.
type Config struct {
	Maintenance  bool                `default:"false"`                                    // $TENANT_MAINTENANCE
	BindAddr     string              `split_words:"true" default:":8088"`                 // $TENANT_BIND_ADDR
	Mode         string              `default:"release"`                                  // $TENANT_MODE
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info"`                  // $TENANT_LOG_LEVEL
	ConsoleLog   bool                `split_words:"true" default:"false"`                 // $TENANT_CONSOLE_LOG
	AllowOrigins []string            `split_words:"true" default:"http://localhost:3000"` // $TENANT_ALLOW_ORIGINS
	Database     DatabaseConfig      `split_words:"true"`
	SendGrid     SendGridConfig      `split_words:"false"`
	Sentry       sentry.Config
	processed    bool // set when the config is properly procesesed from the environment
}

// Configures the connection to trtl for replicated data storage.
type DatabaseConfig struct {
	URL      string `default:"trtl://localhost:4436"`
	Insecure bool   `default:"true"`
	CertPath string `split_words:"true"`
	PoolPath string `split_words:"true"`
	Testing  bool   `default:"false"`
}

// Configures the email and marketing contact APIs for use with the Tenant server.
type SendGridConfig struct {
	APIKey       string `split_words:"true" required:"false"`
	FromEmail    string `split_words:"true" default:"ensign@rotational.io"`
	AdminEmail   string `split_words:"true" default:"admins@rotational.io"`
	EnsignListID string `split_words:"true" required:"false"`
}

// New loads and parses the config from the environment and validates it, marking it as
// processed so that external users can determine if the config is ready for use. This
// should be the only way Config objects are created for use in the application.
func New() (conf Config, err error) {
	if err = envconfig.Process("tenant", &conf); err != nil {
		return Config{}, err
	}

	// Ensure the Sentry release is named correctly
	if conf.Sentry.Release == "" {
		conf.Sentry.Release = fmt.Sprintf("tenant@%s", pkg.Version())
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

	if err = c.Database.Validate(); err != nil {
		return err
	}

	if err = c.SendGrid.Validate(); err != nil {
		return err
	}

	if err = c.Sentry.Validate(); err != nil {
		return err
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

// If not insecure, the cert and pool paths are required.
func (c DatabaseConfig) Validate() (err error) {
	// If in testing mode, configuration is valid
	if c.Testing {
		return nil
	}

	// Ensure that the URL connects to trtl
	var u *url.URL
	if u, err = url.Parse(c.URL); err != nil {
		return errors.New("invalid configuration: could not parse database url")
	}

	if u.Scheme != "trtl" {
		return errors.New("invalid configuration: tenant can only connect to trtl databases")
	}

	if !c.Insecure {
		if c.CertPath == "" || c.PoolPath == "" {
			return errors.New("invalid configuration: connecting to trtl via mTLS requires certs and a pool")
		}
	}
	return nil
}

func (c DatabaseConfig) Endpoint() (_ string, err error) {
	var u *url.URL
	if u, err = url.Parse(c.URL); err != nil {
		return "", err
	}
	return u.Host, nil
}

// The from and admin emails are required if the SendGrid API is enabled.
func (c SendGridConfig) Validate() error {
	if c.Enabled() {
		if c.AdminEmail == "" || c.FromEmail == "" {
			return errors.New("invalid configuration: admin and from emails are required if sendgrid is enabled")
		}
	}
	return nil
}

// Returns true if there is a SendGrid API key available
func (c SendGridConfig) Enabled() bool {
	return c.APIKey != ""
}
