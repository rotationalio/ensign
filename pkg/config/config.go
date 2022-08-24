package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/ensign/pkg/logger"
	"github.com/rs/zerolog"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be ENSIGN_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const prefix = "ensign"

// Config contains all of the configuration parameters for an Ensign server and is
// loaded from the environment or a configuration file with reasonable defaults for
// values that are omitted. The Config should be validated in preparation for running
// the Ensign server to ensure that all server operations work as expected.
type Config struct {
	Maintenance bool                `split_words:"true" default:"false"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool                `split_words:"true" default:"false"`
	BindAddr    string              `split_words:"true" default:"7777"`
	processed   bool
	file        string
}

func New() (conf Config, err error) {
	if err = envconfig.Process(prefix, &conf); err != nil {
		return conf, err
	}

	if err = conf.Validate(); err != nil {
		return conf, err
	}

	conf.processed = true
	return conf, nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c Config) IsZero() bool {
	return !c.processed
}

func (c Config) Validate() error {
	return nil
}

func (c Config) Path() string {
	return c.file
}
