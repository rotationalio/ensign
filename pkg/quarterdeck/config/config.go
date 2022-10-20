package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/ensign/pkg"
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
	Database     DatabaseConfig
	Token        TokenConfig
	Sentry       sentry.Config
	processed    bool // set when the config is properly procesesed from the environment
}

type DatabaseConfig struct {
	URL      string `default:"sqlite3:////data/db/quarterdeck.db"` // $QUARTERDECK_DATABASE_URL
	ReadOnly bool   `split_words:"true" default:"false"`           // $QUARTERDECK_DATABASE_READ_ONLY
}

type TokenConfig struct {
	Keys     map[string]string `required:"false"`                      // $QUARTERDECK_TOKEN_KEYS
	Audience string            `default:"ensign.rotational.app:443"`   // $QUARTERDECK_TOKEN_AUDIENCE
	Issuer   string            `default:"https://auth.rotational.app"` // $QUARTERDECK_TOKEN_ISSUER
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
