package logger

import (
	"io"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	mu   sync.Mutex
	orig *zerolog.Logger
)

func ResetLogger() {
	mu.Lock()
	defer mu.Unlock()
	if orig != nil {
		log.Logger = *orig
	}
}

func Testing(tb testing.TB) {
	mu.Lock()
	defer mu.Unlock()
	orig = &log.Logger
	log.Logger = log.Output(zerolog.NewTestWriter(tb))
}

func Discard() {
	mu.Lock()
	defer mu.Unlock()
	orig = &log.Logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: io.Discard})
}
