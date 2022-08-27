package logger_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLevelDecoder(t *testing.T) {
	testTable := []struct {
		value    string
		expected zerolog.Level
	}{
		{
			"panic", zerolog.PanicLevel,
		},
		{
			"FATAL", zerolog.FatalLevel,
		},
		{
			"Error", zerolog.ErrorLevel,
		},
		{
			"   warn   ", zerolog.WarnLevel,
		},
		{
			"iNFo", zerolog.InfoLevel,
		},
		{
			"debug", zerolog.DebugLevel,
		},
		{
			"trace", zerolog.TraceLevel,
		},
	}

	// Test valid cases
	for _, testCase := range testTable {
		var level logger.LevelDecoder
		err := level.Decode(testCase.value)
		require.NoError(t, err)
		require.Equal(t, testCase.expected, zerolog.Level(level))
	}

	// Test error case
	var level logger.LevelDecoder
	err := level.Decode("notalevel")
	require.EqualError(t, err, `unknown log level "notalevel"`)
}

func TestUnmarshaler(t *testing.T) {
	type Config struct {
		Level logger.LevelDecoder
	}

	var yamlConf Config
	err := yaml.Unmarshal([]byte(`level: "warn"`), &yamlConf)
	require.NoError(t, err, "could not unmarshal level decoder in yaml file")
	require.Equal(t, zerolog.WarnLevel, zerolog.Level(yamlConf.Level))

	var jsonConf Config
	err = json.Unmarshal([]byte(`{"level": "panic"}`), &jsonConf)
	require.NoError(t, err, "could not unmarshal level decoder in json file")
	require.Equal(t, zerolog.PanicLevel, zerolog.Level(jsonConf.Level))
}

func TestMarshaler(t *testing.T) {
	confs := []struct {
		Level logger.LevelDecoder `yaml:"level" json:"level"`
	}{
		{logger.LevelDecoder(zerolog.PanicLevel)},
		{logger.LevelDecoder(zerolog.FatalLevel)},
		{logger.LevelDecoder(zerolog.ErrorLevel)},
		{logger.LevelDecoder(zerolog.WarnLevel)},
		{logger.LevelDecoder(zerolog.InfoLevel)},
		{logger.LevelDecoder(zerolog.DebugLevel)},
		{logger.LevelDecoder(zerolog.TraceLevel)},
	}

	for _, conf := range confs {
		data, err := yaml.Marshal(conf)
		require.NoError(t, err, "could not marshal data into yaml")
		require.Equal(t, []byte(fmt.Sprintf("level: %s\u000A", &conf.Level)), data)

		data, err = json.Marshal(conf)
		require.NoError(t, err, "could not marshal data into json")
		require.Equal(t, []byte(fmt.Sprintf(`{"level":%q}`, &conf.Level)), data)
	}

}
