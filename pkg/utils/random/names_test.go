package random_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/random"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	testCases := []struct {
		n          int
		consonants int
		vowels     int
	}{
		{0, 0, 0},
		{1, 1, 0},
		{2, 2, 0},
		{3, 2, 1},
		{4, 3, 1},
		{5, 3, 2},
		{6, 4, 2},
		{7, 4, 3},
		{8, 5, 3},
		{9, 5, 4},
		{100, 51, 49},
	}

	for _, tc := range testCases {
		name := random.Name(tc.n)
		require.Len(t, name, tc.n, "wrong length for n=%d", tc.n)

		var vowels, consonants int
		for _, c := range name {
			switch c {
			case 'a', 'e', 'i', 'o', 'u':
				vowels++
			default:
				consonants++
			}
		}
		require.Equal(t, tc.consonants, consonants, "wrong number of consonants for n=%d", tc.n)
		require.Equal(t, tc.vowels, vowels, "wrong number of vowels for n=%d", tc.n)
	}
}
