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

// Config uses envconfig to load required settings from the environment, parses
// and validates them, and loads defaults where necessary in preparation for running
// the Tenant API service. This is the top-level config, any sub configurations
// will need to be defined as properties of this Config.
type Config struct {
	Maintenance  bool                `default:"false"`                                    // $TENANT_MAINTENANCE
	BindAddr     string              `split_words:"true" default:":8080"`                 // $TENANT_BIND_ADDR
	Mode         string              `default:"release"`                                  // $TENANT_MODE
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info"`                  // $TENANT_LOG_LEVEL
	ConsoleLog   bool                `split_words:"true" default:"false"`                 // $TENANT_CONSOLE_LOG
	AllowOrigins []string            `split_words:"true" default:"http://localhostL3000"` // $TENANT_ALLOW_ORIGINS
	Sentry       sentry.Config
	processed    bool //is set when the config is properly processed from the environment
}

// New loads and parses the config from the environment and validates it. It is then
// marked as processed so that external users can determine if the config is ready for use.
// This should be the only way Config objects are created for use in the application.
func New() (conf Config, err error) {
	if err = envconfig.Process("tenant", &conf); err != nil {
		return Config{}, err
	}

	// Ensures the Sentry release is named correctly
	if conf.Sentry.Release == "" {
		conf.Sentry.Release = fmt.Sprintf("tenant@%s", pkg.Version())
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Returns true if the config has not been processed correctly from the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Marks a manually constructed config as processed if it is valid.
func (c Config) Mark() (_ Config, err error) {
	if err = c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

// Custom validations are added here, especially validations that require one or more
// fields to be processed before validation occurs. All nested config validation methods
// are called here.
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

func (c Config) AllowAllOrigins() bool {
	if len(c.AllowOrigins) == 1 && c.AllowOrigins[0] == "*" {
		return true
	}
	return false
}
