package ulid_test

import (
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
	ulidlib "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
)

func TestIsZero(t *testing.T) {
	testCases := []struct {
		input  ulid.ULID
		assert require.BoolAssertionFunc
	}{
		{ulid.ULID{}, require.True},
		{ulid.ULID{0x00}, require.True},
		{ulid.ULID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, require.True},
		{ulid.Make(), require.False},
	}

	for _, tc := range testCases {
		tc.assert(t, ulidlib.IsZero(tc.input))
	}
}

func TestNew(t *testing.T) {
	// Should be able to concurrently create 100 ulids.
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			uid := ulidlib.New()
			require.False(t, ulidlib.IsZero(uid))
		}()
	}
	wg.Wait()
}
