package radish_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/utils/radish"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		conf Config
		err  error
	}{
		{Config{}, ErrNoWorkers},
		{Config{Workers: 4}, ErrNoServerName},
		{Config{Workers: 4, ServerName: "radish"}, nil},
	}

	for i, tc := range testCases {
		err := tc.conf.Validate()
		require.ErrorIs(t, err, tc.err, "test case %d failed", i)
	}
}
