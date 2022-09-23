package logger_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	lastLog map[string]interface{}
	levels  map[zerolog.Level]uint16
}

func (w *testWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if w.levels == nil {
		w.levels = make(map[zerolog.Level]uint16)
	}
	w.levels[level]++
	return w.Write(p)
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	if err = json.Unmarshal(p, &w.lastLog); err != nil {
		return 0, err
	}
	return len(p), nil
}

func TestSeverityHook(t *testing.T) {
	// Initialize zerolog with GCP logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Test writer
	tw := &testWriter{}

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(tw).Hook(gcpHook).With().Timestamp().Logger()

	log.Trace().Msg("just a trace")
	require.Equal(t, uint16(1), tw.levels[zerolog.TraceLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "DEBUG", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "just a trace", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])

	log.Debug().Msg("is it on?")
	require.Equal(t, uint16(1), tw.levels[zerolog.DebugLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "DEBUG", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "is it on?", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])

	log.Info().Str("extra", "foo").Msg("my name is bob")
	require.Equal(t, uint16(1), tw.levels[zerolog.InfoLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "INFO", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "my name is bob", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, "extra")
	require.Equal(t, "foo", tw.lastLog["extra"])

	log.Warn().Msg("don't run with scissors")
	require.Equal(t, uint16(1), tw.levels[zerolog.WarnLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "WARNING", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "don't run with scissors", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])

	log.Error().Err(errors.New("bad things")).Msg("oops")
	require.Equal(t, uint16(1), tw.levels[zerolog.ErrorLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "ERROR", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "oops", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, "error")
	require.Equal(t, "bad things", tw.lastLog["error"])

	// Must use WithLevel or the program will exit and the test will fail.
	log.WithLevel(zerolog.FatalLevel).Err(errors.New("murder")).Msg("dying")
	require.Equal(t, uint16(1), tw.levels[zerolog.FatalLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "CRITICAL", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "dying", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, "error")
	require.Equal(t, "murder", tw.lastLog["error"])

	require.Panics(t, func() {
		log.Panic().Err(errors.New("run away!")).Msg("squeeeee!!!")
	})
	require.Equal(t, uint16(1), tw.levels[zerolog.PanicLevel])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeySeverity)
	require.Equal(t, "ALERT", tw.lastLog[logger.GCPFieldKeySeverity])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyMsg)
	require.Equal(t, "squeeeee!!!", tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, logger.GCPFieldKeyTime)
	require.NotEmpty(t, tw.lastLog[logger.GCPFieldKeyMsg])
	require.Contains(t, tw.lastLog, "error")
	require.Equal(t, "run away!", tw.lastLog["error"])
}
