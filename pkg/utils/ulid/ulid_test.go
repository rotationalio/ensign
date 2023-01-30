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

func TestParse(t *testing.T) {
	example := ulidlib.New()

	testCases := []struct {
		input    any
		expected ulid.ULID
		err      error
	}{
		{example.String(), example, nil},
		{example.Bytes(), example, nil},
		{example, example, nil},
		{[16]byte(example), example, nil},
		{"", ulidlib.Null, nil},
		{uint64(14), ulidlib.Null, ulidlib.ErrUnknownType},
		{"foo", ulidlib.Null, ulid.ErrDataSize},
		{[]byte{0x14, 0x21}, ulidlib.Null, ulid.ErrDataSize},
		{ulidlib.Null.String(), ulidlib.Null, nil},
	}

	for i, tc := range testCases {
		actual, err := ulidlib.Parse(tc.input)
		require.ErrorIs(t, err, tc.err, "could not compare error on test case %d", i)
		require.Equal(t, tc.expected, actual, "expected result not returned")
	}
}
