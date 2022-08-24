package logger_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
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
