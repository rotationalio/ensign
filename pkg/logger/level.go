package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// LogLevelDecoder deserializes the log level from a config string.
type LevelDecoder zerolog.Level

// Decode implements envconfig.Decoder
func (ll *LevelDecoder) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "panic":
		*ll = LevelDecoder(zerolog.PanicLevel)
	case "fatal":
		*ll = LevelDecoder(zerolog.FatalLevel)
	case "error":
		*ll = LevelDecoder(zerolog.ErrorLevel)
	case "warn":
		*ll = LevelDecoder(zerolog.WarnLevel)
	case "info":
		*ll = LevelDecoder(zerolog.InfoLevel)
	case "debug":
		*ll = LevelDecoder(zerolog.DebugLevel)
	case "trace":
		*ll = LevelDecoder(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown log level %q", value)
	}
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (ll *LevelDecoder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ls string
	if err := unmarshal(&ls); err != nil {
		return err
	}

	return ll.Decode(ls)
}

// UnmarshalJSON implements json.Unmarshaler
func (ll *LevelDecoder) UnmarshalJSON(data []byte) error {
	var ls string
	if err := json.Unmarshal(data, &ls); err != nil {
		return err
	}
	return ll.Decode(ls)
}
