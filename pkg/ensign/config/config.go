package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
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
	Maintenance bool                `default:"false" yaml:"maintenance"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info" yaml:"log_level"`
	ConsoleLog  bool                `split_words:"true" default:"false" yaml:"console_log"`
	BindAddr    string              `split_words:"true" default:":7777" yaml:"bind_addr"`
	Sentry      sentry.Config
	Monitoring  MonitoringConfig
	processed   bool
	file        string
}

// MonitoringConfig maintains the parameters for the o11y server that the Prometheus
// scraper will fetch the configured observability metrics from.
type MonitoringConfig struct {
	Enabled  bool   `default:"true" yaml:"enabled"`
	BindAddr string `split_words:"true" default:":9090" yaml:"bind_addr"`
	NodeID   string `split_words:"true" required:"false" yaml:"node"`
}

// New creates and processes a Config from the environment ready for use. If the
// configuration is invalid or it cannot be processed an error is returned.
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

// Load and process a Config from the specified configuration file, then process the
// configuration from the environment. If the configuration is invalid or cannot be
// processed either from the file or the environment, then an error is returned.
// The configuration file is processed based on its file extension. YAML files with a
// .yaml or .yml extension are preferred, but JSON (.json) and TOML (.toml) files will
// also be processed. If the path has an unrecognized extension an error is returned.
// HACK: this is a beta function right now and is not fully tested; use with care!
func Load(path string) (conf Config, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return Config{}, err
	}
	defer f.Close()

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		if err = yaml.NewDecoder(f).Decode(&conf); err != nil {
			return Config{}, err
		}
	case ".json":
		if err = json.NewDecoder(f).Decode(&conf); err != nil {
			return Config{}, err
		}
	case ".toml":
		if _, err = toml.NewDecoder(f).Decode(&conf); err != nil {
			return Config{}, err
		}
	default:
		return Config{}, fmt.Errorf("unrecognized file extension %q", filepath.Ext(path))
	}

	// Load the configuration from the environment in order to merge it with the file
	// based configuration ensuring that the values from the environment take priority
	// and that default values are populated where not set by configuration file.
	// NOTE: this next section relies on the fact that the envconfig gets its values
	// from the environment otherwise sets a default from the struct tags. The merge
	// rules populate the conf from the envconf in two cases: when the conf field is
	// zero valued or when an environment variable exists for the field. This code is
	// somewhat fragile because we don't have a method to export the actual environment
	// variable names from envconfig and would have to port code to that. We may want to
	// consider looking into other libraries or porting the code so we can modify it.
	// BUG: if a value is required this will error even if specified in the conf file.
	var envconf Config
	if err = envconfig.Process(prefix, &envconf); err != nil {
		return Config{}, err
	}

	if err = mergenv(&conf, &envconf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return conf, err
	}

	conf.file = path
	conf.processed = true
	return conf, nil
}

// Parse and return the zerolog log level for configuring global logging.
func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

// A Config is zero-valued if it hasn't been processed by a file or the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Mark a manually constructed config as processed as long as its valid.
func (c Config) Mark() (Config, error) {
	if err := c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

// Validates the config is ready for use in the application and that configuration
// semantics such as requiring multiple required configuration parameters are enforced.
func (c Config) Validate() (err error) {
	if err = c.Sentry.Validate(); err != nil {
		return err
	}

	return nil
}

// The path to the configuration file on disk when it was loaded. Returns empty string
// if the config was not loaded from a configuration file.
func (c Config) Path() string {
	return c.file
}
