package config

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/confire"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
)

type Config struct {
	BindAddr       string              `split_words:"true" default:":8090"`
	Mode           string              `default:"release"`
	LogLevel       logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog     bool                `split_words:"true" default:"false"`
	AllowOrigins   []string            `split_words:"true" default:"http://localhost:8090"`
	StatusInterval time.Duration       `split_words:"true" default:"15s"`
	DataPath       string              `split_words:"true" default:".uptime"`
	ServiceInfo    string              `split_words:"true" required:"false"`
	processed      bool                // set when the config is properly processed from the environment
}

// New loads and parses the config from the environment and validates it, marking it as
// processed so that external users can determine if the config is ready for use. This
// should be the only way Config objects are created for use in the application.
func New() (conf Config, err error) {
	if err = confire.Process("uptime", &conf); err != nil {
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
