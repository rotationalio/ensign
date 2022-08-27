package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// LogLevelDecoder deserializes the log level from a config string.
type LevelDecoder zerolog.Level

// Names of log levels for use in encoding/decoding from strings.
const (
	llPanic = "panic"
	llFatal = "fatal"
	llError = "error"
	llWarn  = "warn"
	llInfo  = "info"
	llDebug = "debug"
	llTrace = "trace"
)

// Decode implements envconfig.Decoder
func (ll *LevelDecoder) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case llPanic:
		*ll = LevelDecoder(zerolog.PanicLevel)
	case llFatal:
		*ll = LevelDecoder(zerolog.FatalLevel)
	case llError:
		*ll = LevelDecoder(zerolog.ErrorLevel)
	case llWarn:
		*ll = LevelDecoder(zerolog.WarnLevel)
	case llInfo:
		*ll = LevelDecoder(zerolog.InfoLevel)
	case llDebug:
		*ll = LevelDecoder(zerolog.DebugLevel)
	case llTrace:
		*ll = LevelDecoder(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown log level %q", value)
	}
	return nil
}

// Encode converts the loglevel into a string for use in YAML and JSON
func (ll *LevelDecoder) Encode() (string, error) {
	switch zerolog.Level(*ll) {
	case zerolog.PanicLevel:
		return llPanic, nil
	case zerolog.FatalLevel:
		return llFatal, nil
	case zerolog.ErrorLevel:
		return llError, nil
	case zerolog.WarnLevel:
		return llWarn, nil
	case zerolog.InfoLevel:
		return llInfo, nil
	case zerolog.DebugLevel:
		return llDebug, nil
	case zerolog.TraceLevel:
		return llTrace, nil
	default:
		return "", fmt.Errorf("unknown log level %d", ll)
	}
}

func (ll LevelDecoder) String() string {
	ls, _ := ll.Encode()
	return ls
}

// UnmarshalYAML implements yaml.Unmarshaler
func (ll *LevelDecoder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ls string
	if err := unmarshal(&ls); err != nil {
		return err
	}

	return ll.Decode(ls)
}

// MarshalYAML implements yaml.Marshaler
func (ll LevelDecoder) MarshalYAML() (interface{}, error) {
	return ll.Encode()
}

// UnmarshalJSON implements json.Unmarshaler
func (ll *LevelDecoder) UnmarshalJSON(data []byte) error {
	var ls string
	if err := json.Unmarshal(data, &ls); err != nil {
		return err
	}
	return ll.Decode(ls)
}

// MarshalJSON implements json.Marshaler
func (ll LevelDecoder) MarshalJSON() ([]byte, error) {
	ls, err := ll.Encode()
	if err != nil {
		return nil, err
	}
	return json.Marshal(ls)
}
