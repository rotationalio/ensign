package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rotationalio/confire"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
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
	Maintenance bool                `default:"false" yaml:"maintenance"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info" yaml:"log_level"`
	ConsoleLog  bool                `split_words:"true" default:"false" yaml:"console_log"`
	BindAddr    string              `split_words:"true" default:":5356" yaml:"bind_addr"`
	MetaTopic   MetaTopicConfig     `split_words:"true"`
	Monitoring  MonitoringConfig
	Storage     StorageConfig
	Auth        AuthConfig
	Sentry      sentry.Config
	processed   bool
	file        string
}

// MetaTopicConfig defines the topics and events that the Ensign node publishes along
// with the credentials and connection endpoints to connect to Ensign on.
type MetaTopicConfig struct {
	Enabled      bool   `default:"true" yaml:"enabled"`
	TopicName    string `split_words:"true" default:"ensign.metatopic.topics"`
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Endpoint     string `default:"ensign.rotational.app:443"`
	AuthURL      string `split_words:"true" default:"https://auth.rotational.app"`
}

// MonitoringConfig maintains the parameters for the o11y server that the Prometheus
// scraper will fetch the configured observability metrics from.
type MonitoringConfig struct {
	Enabled  bool   `default:"true" yaml:"enabled"`
	BindAddr string `split_words:"true" default:":1205" yaml:"bind_addr"`
	NodeID   string `split_words:"true" required:"false" yaml:"node"`
}

// StorageConfig defines on disk where Ensign keeps its data. Users must specify the
// DataPath directory where Ensign will store it's data.
type StorageConfig struct {
	ReadOnly bool   `default:"false" split_words:"true" yaml:"read_only"`
	DataPath string `split_words:"true" yaml:"data_path"`
	Testing  bool   `default:"false" yaml:"testing"`
}

// AuthConfig defines how Ensign connects to Quarterdeck in order to authorize requests.
type AuthConfig struct {
	KeysURL            string        `split_words:"true" default:"https://auth.rotational.app/.well-known/jwks.json"`
	Audience           string        `default:"https://ensign.rotational.app"`
	Issuer             string        `default:"https://auth.rotational.app"`
	MinRefreshInterval time.Duration `split_words:"true" default:"5m"`
}

// New creates and processes a Config from the environment ready for use. If the
// configuration is invalid or it cannot be processed an error is returned.
func New() (conf Config, err error) {
	if err = confire.Process(prefix, &conf); err != nil {
		return conf, err
	}

	// Ensure the Sentry release is set to ensign.
	if conf.Sentry.Release == "" {
		conf.Sentry.Release = fmt.Sprintf("ensign@%s", pkg.Version())
	}

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
	if err = c.Storage.Validate(); err != nil {
		return err
	}

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

func (c MetaTopicConfig) Validate() error {
	if c.Enabled {
		if c.TopicName == "" {
			return errors.New("invalid meta topic config: missing topic name")
		}

		if c.ClientID == "" || c.ClientSecret == "" {
			return errors.New("invalid meta topic config: missing client id or secret")
		}
	}
	return nil
}

func (c StorageConfig) Validate() (err error) {
	if c.DataPath == "" {
		return errors.New("invalid storage config: missing data path")
	}
	return nil
}

// MetaPath returns the path to the metadata store for Ensign, checking to make sure
// that the directory exists and that it is a directory. If it doesn't exist, the
// directory is created; an error is returned if the path is invalid or cannot be
// created.
func (c StorageConfig) MetaPath() (path string, err error) {
	path = filepath.Join(c.DataPath, "metadata")
	if err = c.checkPath(path); err != nil {
		return "", err
	}
	return path, nil
}

// EventPath returns the path to the event data store for Ensign, checking to make sure
// that the directory exists and that it is a directory. If it doesn't exist, the
// directory is created; an error is returned if the path is invalid or cannot be
// created.
func (c StorageConfig) EventPath() (path string, err error) {
	path = filepath.Join(c.DataPath, "events")
	if err = c.checkPath(path); err != nil {
		return "", err
	}
	return path, nil
}

func (c StorageConfig) checkPath(path string) (err error) {
	var info os.FileInfo
	if info, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Attempt to create the directory if it doesn't exist
			if err = os.MkdirAll(path, 0744); err != nil {
				return err
			}
			return nil
		}

		// Return any permissions error or other os errors.
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("invalid configuration: %s is not a directory", path)
	}
	return nil
}

func (c AuthConfig) AuthOptions() []middleware.AuthOption {
	return []middleware.AuthOption{
		middleware.WithAuthOptions(middleware.AuthOptions{
			KeysURL:            c.KeysURL,
			Audience:           c.Audience,
			Issuer:             c.Issuer,
			MinRefreshInterval: c.MinRefreshInterval,
		}),
	}
}
