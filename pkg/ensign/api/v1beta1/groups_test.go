package api_test

import (
	"testing"

	"github.com/google/uuid"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/stretchr/testify/require"
	"github.com/twmb/murmur3"
)

func TestGroupsKey(t *testing.T) {
	group := &api.ConsumerGroup{}
	_, err := group.Key()
	require.ErrorIs(t, err, api.ErrNoGroupID)

	uu := uuid.New()

	testCases := []struct {
		Name     string
		ID       []byte
		Expected [16]byte
		Msg      string
	}{
		{
			ID:       uu[:],
			Expected: [16]byte(uu),
			Msg:      "should be able to use a UUID as the ID",
		},
		{
			ID:       []byte{16, 32, 255, 129, 82, 91, 255, 0, 21, 48, 198, 122, 219, 42},
			Expected: [16]uint8{0x8f, 0x8, 0xf5, 0x6c, 0x71, 0x64, 0xc6, 0x55, 0xe7, 0x47, 0x95, 0xb4, 0xa0, 0x52, 0xb7, 0x5},
			Msg:      "should be able to use a short variable length ID hashed as the key",
		},
		{
			ID:       []byte("should be able to use a very long byte string and hash it to a 16 byte key"),
			Expected: [16]uint8{0x38, 0xde, 0x19, 0x43, 0xd9, 0x32, 0xd0, 0xd6, 0xd9, 0xc5, 0x51, 0x5d, 0x9b, 0x46, 0x8f, 0xeb},
			Msg:      "should be able to use a long variable length ID hashed as the key",
		},
		{
			Name:     "testing.test.consumergroup.alpha",
			Expected: [16]uint8{0xf0, 0xa5, 0xb9, 0x50, 0xd, 0x4b, 0x78, 0xe0, 0xdd, 0x7f, 0x15, 0xa, 0x39, 0x15, 0x59, 0x52},
			Msg:      "should be able to use name if ID is not specified",
		},
		{
			ID:       []byte{16, 32, 255, 129, 82, 91, 255, 0, 21, 48, 198, 122, 219, 42},
			Name:     "testing.test.consumergroup.alpha",
			Expected: [16]uint8{0x8f, 0x8, 0xf5, 0x6c, 0x71, 0x64, 0xc6, 0x55, 0xe7, 0x47, 0x95, 0xb4, 0xa0, 0x52, 0xb7, 0x5},
			Msg:      "ID should be used as the key if both are specified",
		},
	}

	for _, tc := range testCases {
		group = &api.ConsumerGroup{Id: tc.ID, Name: tc.Name}
		key, err := group.Key()
		require.NoError(t, err, tc.Msg)
		require.Len(t, key, 16)
		require.Equal(t, tc.Expected, key, tc.Msg)
	}
}

func TestMurmur3Sanity(t *testing.T) {
	// Just making sure we understand how murmur3 works.
	hash := murmur3.New128()
	hash.Write([]byte("this is the song that never ends, it goes on and on my friends, and if you started singing it ..."))
	sum := hash.Sum(nil)
	require.Len(t, sum, 16)
}
