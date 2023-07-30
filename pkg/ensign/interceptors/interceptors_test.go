package interceptors_test

import (
	"os"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/logger"
)

func TestMain(m *testing.M) {
	logger.Discard()
	exitVal := m.Run()
	logger.ResetLogger()
	os.Exit(exitVal)
}
