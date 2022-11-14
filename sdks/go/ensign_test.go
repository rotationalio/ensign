package ensign_test

import (
	"testing"

	ensign "github.com/rotationalio/ensign/sdks/go"
	"github.com/stretchr/testify/require"
)

func TestNilWilNilOpts(t *testing.T) {
	_, err := ensign.New(nil)
	require.NoError(t, err, "could not pass nil into ensign")
}
