package config

import "github.com/rs/zerolog"

type Config struct {
	Maintenance bool   `split_words:"true" default:"false"`
	BindAddr    string `split_words:"true" default:"8088"`
	Mode        string `split_words:"true" default:"release"`
	LogLevel    string `split_words:"true" default:"info"`
	ConsoleLog  bool   `split_words:"true" default:"false"`
}

func New() (Config, error) {
	return Config{}, nil
}

func (c Config) IsZero() bool {
	// TODO: implement
	return false
}

func (c Config) GetLogLevel() zerolog.Level {
	// TODO: implement
	return zerolog.InfoLevel
}
