package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load the required settings from the environment, parse and
// validate them, loading defaults where necessary in preparation for running the
// Quarterdeck API service. This is the top-level config, all sub configurations need
// to be defined as properties of this Config.
type Config struct {
	Maintenance bool                `default:"false"`                    // $QUARTERDECK_MAINTENANCE
	BindAddr    string              `split_words:"true" default:":8088"` // $QUARTERDECK_BIND_ADDR
	Mode        string              `default:"release"`                  // $QUARTERDECK_MODE
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info"`  // $QUARTERDECK_LOG_LEVEL
	ConsoleLog  bool                `split_words:"true" default:"false"` // $QUARTERDECK_CONSOLE_LOG
	processed   bool                // set when the config is properly procesesed from the environment
}

// New loads and parses the config from the environment and validates it, marking it as
// processed so that external users can determine if the config is ready for use. This
// should be the only way Config objects are created for use in the application.
func New() (conf Config, err error) {
	if err = envconfig.Process("quarterdeck", &conf); err != nil {
		return Config{}, err
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

	return nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}
