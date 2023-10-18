package api_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	. "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	region "github.com/rotationalio/ensign/pkg/ensign/region/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDuplicates(t *testing.T) {
	// For the exact same event, all policies should return duplicate; for completely
	// different events, all policies should return not duplicate.
	alpha := createRandomEvent(mimetype.ApplicationBSON, "RandomData v1.2.3", "FirstKey:rand,SecondKey:rand,ThirdKey:rand,foo:bar,color:blue")
	bravo := createRandomEvent(mimetype.ApplicationJSON, "RandomJSON v2.1.0", "FirstKey:rand,SecondKey:rand,ThirdKey:rand")

	policies := []*Deduplication{
		{Strategy: Deduplication_STRICT},
		{Strategy: Deduplication_DATAGRAM},
		{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"FirstKey", "SecondKey"}},
		{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"FirstKey", "SecondKey"}},
	}

	for _, policy := range policies {
		duplicate, err := alpha.Duplicates(alpha, policy)
		require.NoError(t, err, "could not compute duplicate")
		require.True(t, duplicate, "expected alpha to be a duplicate of itself")

		duplicate, err = alpha.Duplicates(bravo, policy)
		require.NoError(t, err, "could not compute duplicate")
		require.False(t, duplicate, "expected alpha and bravo to not be duplicates")
	}
}

func TestDuplicatesStrict(t *testing.T) {
	testCases := []struct {
		name   string
		alpha  *EventWrapper
		bravo  *EventWrapper
		assert require.BoolAssertionFunc
	}{
		{
			"identical events",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			require.True,
		},
		{
			"different data",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("WdTo+rXJ9+oWa/Bh0Sy8bU5f6yB6DCiJN5j/jFmlK606ym7KliheJ3IS", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			require.False,
		},
		{
			"different metadata",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "baz:zap,color:blue", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			require.False,
		},
		{
			"different mimetype",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationJSON, "TestEvent v1.2.3", ""),
			require.False,
		},
		{
			"different event type",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "MockEvent v1.2.3", ""),
			require.False,
		},
		{
			"different event version",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v2.1.0", ""),
			require.False,
		},
	}

	for i, tc := range testCases {
		duplicates, err := tc.alpha.DuplicatesStrict(tc.bravo)
		require.NoError(t, err, "could not compute duplicate strict")
		tc.assert(t, duplicates, "test case %s (case %d) failed", tc.name, i)
	}
}

func TestDuplicatesDatagram(t *testing.T) {
	testCases := []struct {
		name   string
		alpha  *EventWrapper
		bravo  *EventWrapper
		assert require.BoolAssertionFunc
	}{
		{
			"identical events",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			require.True,
		},
		{
			"different data",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("WdTo+rXJ9+oWa/Bh0Sy8bU5f6yB6DCiJN5j/jFmlK606ym7KliheJ3IS", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			require.False,
		},
		{
			"different metadata",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "baz:zap,color:blue", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			require.True,
		},
		{
			"different mimetype",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationJSON, "TestEvent v1.2.3", ""),
			require.True,
		},
		{
			"different event type",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "MockEvent v1.2.3", ""),
			require.True,
		},
		{
			"different prefix",
			mkwevt("rOfa3JglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v2.1.0", ""),
			require.False,
		},
		{
			"different suffix",
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
			mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqp3nI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v2.1.0", ""),
			require.False,
		},
	}

	for i, tc := range testCases {
		duplicates, err := tc.alpha.DuplicatesDatagram(tc.bravo)
		require.NoError(t, err, "could not compute duplicate datagram")
		tc.assert(t, duplicates, "test case %s (case %d) failed", tc.name, i)
	}
}

func TestDuplicatesKeyGrouped(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := []struct {
			name   string
			alpha  *EventWrapper
			bravo  *EventWrapper
			keys   []string
			assert require.BoolAssertionFunc
		}{
			{
				"identical events - single key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo"},
				require.True,
			},
			{
				"identical events - multi key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo", "color"},
				require.True,
			},
			{
				"identical data - different keys",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:blue", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"color"},
				require.False,
			},
			{
				"identical data - different multi key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:blue", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo", "color"},
				require.False,
			},
			{
				"mismatch data - same key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("3keqNCPHXchJhH5bPpVW0R4dh7YAJkrciUlzfm8Mr+1fUVepnwp8Ps9J", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"color"},
				require.False,
			},
			{
				"mismatch data - same multi key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("3keqNCPHXchJhH5bPpVW0R4dh7YAJkrciUlzfm8Mr+1fUVepnwp8Ps9J", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"color", "foo"},
				require.False,
			},
		}

		for i, tc := range testCases {
			duplicates, err := tc.alpha.DuplicatesKeyGrouped(tc.bravo, tc.keys)
			require.NoError(t, err, "could not compute duplicate key grouped")
			tc.assert(t, duplicates, "test case %s (case %d) failed", tc.name, i)
		}
	})

	t.Run("KeysRequired", func(t *testing.T) {
		alpha := mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", "")
		bravo := mkwevt("UP3uSUGCpNIIP4EFtOtOPv7d4wbYuvZATwN/o35czZUlCi3cLJdJEiXT", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", "")

		_, err := alpha.DuplicatesKeyGrouped(bravo, nil)
		require.ErrorIs(t, err, ErrNoKeys)

		_, err = alpha.DuplicatesKeyGrouped(bravo, []string{})
		require.ErrorIs(t, err, ErrNoKeys)
	})
}

func TestDuplicatesUniqueKey(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := []struct {
			name   string
			alpha  *EventWrapper
			bravo  *EventWrapper
			keys   []string
			assert require.BoolAssertionFunc
		}{
			{
				"identical events - single key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo"},
				require.True,
			},
			{
				"identical events - multi key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo", "color"},
				require.True,
			},
			{
				"mismatched data - single key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("3keqNCPHXchJhH5bPpVW0R4dh7YAJkrciUlzfm8Mr+1fUVepnwp8Ps9J", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo"},
				require.True,
			},
			{
				"mismatched data - multi key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("3keqNCPHXchJhH5bPpVW0R4dh7YAJkrciUlzfm8Mr+1fUVepnwp8Ps9J", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo", "color"},
				require.True,
			},
			{
				"mismatched key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:baz,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo"},
				require.False,
			},
			{
				"mismatched multi-key",
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:blue", mimetype.ApplicationBSON, "TestEvent v1.2.3", ""),
				[]string{"foo", "color"},
				require.False,
			},
		}

		for i, tc := range testCases {
			duplicates, err := tc.alpha.DuplicatesUniqueKey(tc.bravo, tc.keys)
			require.NoError(t, err, "could not compute duplicate unique key")
			tc.assert(t, duplicates, "test case %s (case %d) failed", tc.name, i)
		}
	})

	t.Run("KeysRequired", func(t *testing.T) {
		alpha := mkwevt("rOfawJglnmlvXWAKhw7aNUdXFlqpZnI", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", "")
		bravo := mkwevt("UP3uSUGCpNIIP4EFtOtOPv7d4wbYuvZATwN/o35czZUlCi3cLJdJEiXT", "foo:bar,color:red", mimetype.ApplicationBSON, "TestEvent v1.2.3", "")

		_, err := alpha.DuplicatesUniqueKey(bravo, nil)
		require.ErrorIs(t, err, ErrNoKeys)

		_, err = alpha.DuplicatesUniqueKey(bravo, []string{})
		require.ErrorIs(t, err, ErrNoKeys)
	})
}

func TestDuplicatesUniqueField(t *testing.T) {
	testCases := []struct {
		name   string
		alpha  *EventWrapper
		bravo  *EventWrapper
		fields []string
		assert require.BoolAssertionFunc
	}{}

	for i, tc := range testCases {
		duplicates, err := tc.alpha.DuplicatesUniqueField(tc.bravo, tc.fields)
		require.NoError(t, err, "could not compute duplicate unique field")
		tc.assert(t, duplicates, "test case %s (case %d) failed", tc.name, i)
	}
}

func TestHash(t *testing.T) {
	// For the same event, but different policies, each hashing method should return
	// a different hash signature of the fixture.
	hashes := make(map[string]struct{})
	event := createRandomEvent(
		mimetype.ApplicationMsgPack,
		"RandomData v1.2.3",
		"FirstKey:rand,SecondKey:rand,ThirdKey:rand,foo:bar,color:blue",
	)

	policies := []*Deduplication{
		{Strategy: Deduplication_STRICT},
		{Strategy: Deduplication_DATAGRAM},
		{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"FirstKey", "SecondKey"}},
		{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"FirstKey", "SecondKey"}},
	}

	for _, policy := range policies {
		sig, err := event.Hash(policy)
		require.NoError(t, err, "could not create a hash of the event")
		hashes[base64.RawStdEncoding.EncodeToString(sig)] = struct{}{}
	}

	require.Len(t, hashes, len(policies), "expected a unique hash for each policy, one of the hashes is duplicated")
}

func TestHashStrict(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			&EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			fixtures[0], []byte{0xb, 0xa6, 0x75, 0xb8, 0x91, 0xd7, 0xee, 0xc8, 0xa4, 0xe1, 0xf, 0x84, 0x7d, 0x8a, 0x14, 0x9b}, nil,
		},
		{
			fixtures[1], []byte{0xe6, 0xab, 0x9c, 0x80, 0xd, 0xec, 0xa8, 0x2a, 0x5d, 0xaf, 0x5a, 0xbd, 0x99, 0x16, 0xb3, 0x66}, nil,
		},
		{
			fixtures[2], []byte{0x32, 0x1e, 0xf3, 0xfb, 0x10, 0x56, 0xb0, 0x12, 0x2d, 0x87, 0x1b, 0x1b, 0xc3, 0x65, 0xf1, 0x5e}, nil,
		},
		{
			fixtures[3], []byte{0xc8, 0x33, 0x8f, 0xf9, 0x7b, 0x87, 0x1, 0x30, 0xc7, 0xc5, 0x77, 0xf4, 0x1b, 0x8a, 0x6c, 0x95}, nil,
		},
		{
			fixtures[4], []byte{0xf9, 0xdd, 0xed, 0xfc, 0xd8, 0x7e, 0xdc, 0x47, 0x31, 0x58, 0x80, 0x41, 0x34, 0x82, 0x22, 0x39}, nil,
		},
		{
			fixtures[5], []byte{0x9e, 0xd2, 0xca, 0xe, 0x35, 0xd2, 0x11, 0xc8, 0xe2, 0x2a, 0x46, 0xf1, 0x46, 0xb0, 0xf7, 0x71}, nil,
		},
		{
			fixtures[6], []byte{0x7f, 0x6a, 0xad, 0xbb, 0xec, 0x1a, 0xcc, 0x9, 0x9f, 0xb0, 0x6e, 0x9e, 0x64, 0x67, 0xf0, 0x5c}, nil,
		},
		{
			fixtures[7], []byte{0xbe, 0x68, 0x6a, 0xd5, 0x70, 0x5, 0xef, 0x12, 0xc9, 0xd6, 0x20, 0x37, 0xc2, 0x3b, 0xd1, 0x89}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashStrict()
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashDatagram(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			&EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			fixtures[0], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[1], []byte{0x57, 0xad, 0xd1, 0x29, 0x33, 0xd1, 0xfb, 0xda, 0xec, 0x1f, 0x5d, 0x8d, 0x8, 0x22, 0x2, 0xc2}, nil,
		},
		{
			fixtures[2], []byte{0x42, 0x60, 0x33, 0x63, 0xdb, 0xff, 0x5, 0x2c, 0x82, 0x92, 0xbd, 0x6a, 0x35, 0x42, 0x7f, 0x72}, nil,
		},
		{
			fixtures[3], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[4], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[5], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[6], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[7], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashDatagram()
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashKeyGroup(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		keys     []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			[]string{"foo", "bar"}, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, fixtures[0], nil, ErrNoKeys,
		},
		{
			[]string{"alpha"}, fixtures[0], []byte{0xb6, 0x8c, 0x41, 0x51, 0x48, 0x79, 0x3d, 0xed, 0x22, 0xd3, 0xc2, 0xd9, 0x8a, 0xba, 0x97, 0xf1}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[0], []byte{0x1c, 0x7c, 0x5f, 0x13, 0xd3, 0xe1, 0xa7, 0x3d, 0xb4, 0x73, 0xbc, 0x29, 0x6c, 0x7c, 0xba, 0xd0}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[0], []byte{0xec, 0x41, 0xdc, 0xed, 0x2, 0x6c, 0xfe, 0xd3, 0xf6, 0xbc, 0xab, 0xb3, 0x32, 0x6a, 0xff, 0x36}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[0], []byte{0xc6, 0x94, 0x6f, 0xf2, 0x9a, 0xcd, 0xbb, 0x7a, 0x7b, 0x3d, 0x66, 0x7c, 0xba, 0x6b, 0xd0, 0xf8}, nil,
		},
		{
			[]string{"alpha"}, fixtures[1], []byte{0x37, 0x13, 0x18, 0xee, 0x86, 0xcd, 0xad, 0xef, 0x44, 0xf9, 0x7d, 0xb4, 0xe2, 0x9, 0x4b, 0x65}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[1], []byte{0x85, 0xe7, 0x60, 0xa5, 0x23, 0x1f, 0x7f, 0x6b, 0x61, 0x2b, 0x43, 0x6e, 0x3c, 0x76, 0x92, 0x50}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[1], []byte{0xb1, 0xf1, 0xcd, 0xa4, 0x27, 0x5f, 0x9c, 0xfd, 0x88, 0x6b, 0x16, 0x3f, 0x3b, 0x9d, 0x75, 0xaf}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[1], []byte{0xc6, 0x7b, 0x7f, 0xcf, 0x2f, 0x8, 0x82, 0x32, 0xb0, 0xea, 0x32, 0x3f, 0x29, 0x37, 0x4a, 0x28}, nil,
		},
		{
			[]string{"alpha"}, fixtures[3], []byte{0xe0, 0xab, 0x20, 0x5c, 0xca, 0xe6, 0x97, 0xa1, 0xae, 0x1f, 0x76, 0xda, 0xb4, 0x97, 0x1b, 0x54}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[3], []byte{0x47, 0x80, 0x8, 0x9, 0xad, 0xe1, 0x7d, 0x42, 0xc, 0xd4, 0x76, 0xe0, 0x9d, 0xee, 0xeb, 0x3a}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[3], []byte{0xb9, 0x5c, 0x2, 0x4f, 0xc0, 0xda, 0xd8, 0xf4, 0xe4, 0xe6, 0x6e, 0x76, 0x4b, 0x83, 0xcb, 0x37}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[3], []byte{0xc6, 0x60, 0x2, 0xf, 0x9c, 0xa0, 0x12, 0x8b, 0x77, 0x83, 0x58, 0xb8, 0x3a, 0x1b, 0x21, 0x1b}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashKeyGrouped(tc.keys)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashUniqueKey(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		keys     []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			[]string{"foo", "bar"}, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, fixtures[0], nil, ErrNoKeys,
		},
		{
			[]string{"alpha"}, fixtures[0], []byte{0x32, 0x59, 0xfa, 0xa7, 0xed, 0x90, 0xe9, 0xae, 0xab, 0xbc, 0x75, 0xac, 0xa1, 0xf7, 0xac, 0x9a}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[0], []byte{0x5a, 0xbb, 0x4a, 0x32, 0xa9, 0x3d, 0x54, 0x8e, 0xeb, 0xc8, 0xb6, 0x1e, 0xc6, 0x8, 0xef, 0xb4}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[0], []byte{0x85, 0x53, 0xb2, 0x1a, 0x1e, 0x2f, 0x63, 0x29, 0x5b, 0x1, 0x6, 0x1f, 0x64, 0x7a, 0x45, 0x87}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[0], []byte{0x4b, 0xa4, 0x42, 0x74, 0x80, 0xbf, 0x8f, 0x56, 0x19, 0x88, 0x77, 0x0, 0xa4, 0x8, 0x63, 0x9c}, nil,
		},
		{
			[]string{"alpha"}, fixtures[1], []byte{0x32, 0x59, 0xfa, 0xa7, 0xed, 0x90, 0xe9, 0xae, 0xab, 0xbc, 0x75, 0xac, 0xa1, 0xf7, 0xac, 0x9a}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[1], []byte{0x5a, 0xbb, 0x4a, 0x32, 0xa9, 0x3d, 0x54, 0x8e, 0xeb, 0xc8, 0xb6, 0x1e, 0xc6, 0x8, 0xef, 0xb4}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[1], []byte{0x85, 0x53, 0xb2, 0x1a, 0x1e, 0x2f, 0x63, 0x29, 0x5b, 0x1, 0x6, 0x1f, 0x64, 0x7a, 0x45, 0x87}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[1], []byte{0x4b, 0xa4, 0x42, 0x74, 0x80, 0xbf, 0x8f, 0x56, 0x19, 0x88, 0x77, 0x0, 0xa4, 0x8, 0x63, 0x9c}, nil,
		},
		{
			[]string{"alpha"}, fixtures[3], []byte{0xfe, 0x48, 0xa9, 0xff, 0x5b, 0x55, 0x53, 0xb0, 0x8c, 0x78, 0x8b, 0x65, 0x2f, 0x38, 0x52, 0x30}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[3], []byte{0x65, 0x36, 0xf5, 0x79, 0x3f, 0x74, 0xfa, 0x3e, 0x94, 0xd, 0xe3, 0xed, 0x71, 0x95, 0x61, 0x63}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[3], []byte{0xf1, 0x80, 0x3f, 0xaa, 0x3a, 0xb7, 0xfe, 0x2b, 0xe9, 0x63, 0xdb, 0xc8, 0x15, 0x62, 0x32, 0x9}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[3], []byte{0x1c, 0xdc, 0x82, 0x74, 0x32, 0x4b, 0x84, 0x11, 0xbd, 0xf0, 0xe3, 0x9, 0x78, 0x9c, 0x95, 0x11}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashUniqueKey(tc.keys)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashUniqueField(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		fields   []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			[]string{"foo", "bar"}, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, fixtures[0], nil, ErrNoFields,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashUniqueField(tc.fields)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestDuplicateReferencing(t *testing.T) {
	setPublisherBravo := func(w *EventWrapper) {
		w.Region = region.Region_LKE_EU_WEST_1A
		w.Publisher = &Publisher{
			PublisherId: "01HD1M0AVDVHPA4WA73MAA7NH7",
			Ipaddr:      "192.148.21.133",
			ClientId:    "data-ingestor-bravo",
			UserAgent:   "Go-Ensign v0.11.0",
		}
	}

	clone := func(t *testing.T, w *EventWrapper) *EventWrapper {
		data, err := proto.Marshal(w)
		require.NoError(t, err, "could not marshal event")

		o := &EventWrapper{}
		err = proto.Unmarshal(data, o)
		require.NoError(t, err, "could not unmarshal event")
		return o
	}

	execTest := func(t *testing.T, evt, dup *EventWrapper, policy *Deduplication) {
		// Clone the duplicate for comparison purposes.
		org := clone(t, dup)

		// Mark the dup as a duplicate of evt
		err := dup.DuplicateOf(evt, policy)
		require.NoError(t, err, "could not take duplicate of event")
		require.True(t, dup.IsDuplicate)
		require.Equal(t, evt.Id, dup.DuplicateId)
		require.NotNil(t, dup.Event, "expected event data remaining")
		require.Less(t, len(dup.Event), len(evt.Event), "expected the duplicate event size to be reduced")
		require.NotEqual(t, evt.Committed, dup.Committed, "expected committed timestamp to still differ")
		require.False(t, proto.Equal(org, dup), "expected the duplicate event to be modified")

		// Restore the duplicate from the original event
		err = dup.DuplicateFrom(evt)
		require.NoError(t, err, "could restore duplicate from event")
		require.True(t, org.Equals(dup), "expected event to be equal after duplicate resoration")
	}

	t.Run("Strict", func(t *testing.T) {
		evt := mkwevt("tgRkJ4Q2otHC/wJFpePDd+YLwkBQPefGPzfldeIzZKz2wvzfFvNuUrifU7E7ritvFKUcmxW0aR8uGOVTNq1jBA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T09:41:19-05:00")
		dup := mkwevt("tgRkJ4Q2otHC/wJFpePDd+YLwkBQPefGPzfldeIzZKz2wvzfFvNuUrifU7E7ritvFKUcmxW0aR8uGOVTNq1jBA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T16:32:21-05:00")
		setPublisherBravo(dup)

		// Clone the duplicate for comparison purposes.
		org := clone(t, dup)

		// Mark the dup as a duplicate of evt
		err := dup.DuplicateOf(evt, &Deduplication{Strategy: Deduplication_STRICT})
		require.NoError(t, err, "could not take duplicate of event")
		require.True(t, dup.IsDuplicate)
		require.Equal(t, evt.Id, dup.DuplicateId)
		require.Nil(t, dup.Event, "expected event data to be nil")
		require.Nil(t, dup.Encryption, "expected encryption to be nil")
		require.Nil(t, dup.Compression, "expected compression to be nil")
		require.NotEqual(t, evt.Committed, dup.Committed, "expected committed timestamp to still differ")
		require.False(t, proto.Equal(org, dup), "expeced the duplicate event to be modified")

		// Restore the duplicate from the original event
		err = dup.DuplicateFrom(evt)
		require.NoError(t, err, "could restore duplicate from event")
		require.True(t, org.Equals(dup), "expected event to be equal after duplicate resoration")
	})

	t.Run("Datagram", func(t *testing.T) {
		evt := mkwevt("tgRkJ4Q2otHC/wJFpePDd+YLwkBQPefGPzfldeIzZKz2wvzfFvNuUrifU7E7ritvFKUcmxW0aR8uGOVTNq1jBA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T09:41:19-05:00")
		dup := mkwevt("tgRkJ4Q2otHC/wJFpePDd+YLwkBQPefGPzfldeIzZKz2wvzfFvNuUrifU7E7ritvFKUcmxW0aR8uGOVTNq1jBA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T16:32:21-05:00")
		setPublisherBravo(dup)
		execTest(t, evt, dup, &Deduplication{Strategy: Deduplication_DATAGRAM})
	})

	t.Run("Datagram_Metadata", func(t *testing.T) {
		evt := mkwevt("tgRkJ4Q2otHC/wJFpePDd+YLwkBQPefGPzfldeIzZKz2wvzfFvNuUrifU7E7ritvFKUcmxW0aR8uGOVTNq1jBA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T09:41:19-05:00")
		dup := mkwevt("tgRkJ4Q2otHC/wJFpePDd+YLwkBQPefGPzfldeIzZKz2wvzfFvNuUrifU7E7ritvFKUcmxW0aR8uGOVTNq1jBA==", "foo:bar,color:blue,jack:jones", mimetype.MIME_APPLICATION_JSON, "TestVersion 1.4.4", "2023-10-18T16:32:21-05:00")
		setPublisherBravo(dup)
		execTest(t, evt, dup, &Deduplication{Strategy: Deduplication_DATAGRAM})
	})

	t.Run("KeyGrouped", func(t *testing.T) {
		// TODO: implement this test!
		t.Skip("not implemented yet")
	})

	t.Run("UniqueKey", func(t *testing.T) {
		evt := mkwevt("VzcrPlcPKWG/dapHfLnfiA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T09:41:19-05:00")
		dup := mkwevt("VzcrPlcPKWG/dapHfLnfiA==", "foo:bar,color:blue,jack:jones", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T16:32:21-05:00")
		setPublisherBravo(dup)
		execTest(t, evt, dup, &Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"foo"}})
	})

	t.Run("UniqueKey_Data", func(t *testing.T) {
		evt := mkwevt("VzcrPlcPKWG/dapHfLnfiA==", "foo:bar,color:red", mimetype.MIME_APPLICATION_BSON, "TestVersion 1.2.3", "2023-10-18T09:41:19-05:00")
		dup := mkwevt("64azNXY/sTehGqN3Fk7kbw==", "foo:bar,color:red,jack:jones", mimetype.MIME_APPLICATION_JSON, "TestVersion 1.4.4", "2023-10-18T16:32:21-05:00")
		setPublisherBravo(dup)
		execTest(t, evt, dup, &Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"foo"}})
	})

	t.Run("UniqueField", func(t *testing.T) {
		// TODO: implement this test!
		t.Skip("not implemented yet")
	})
}

func TestDeduplicationEquals(t *testing.T) {
	testCases := []struct {
		alpha  *Deduplication
		bravo  *Deduplication
		assert require.BoolAssertionFunc
	}{
		{
			&Deduplication{}, &Deduplication{}, require.True,
		},
		{
			&Deduplication{Strategy: Deduplication_STRICT},
			&Deduplication{Strategy: Deduplication_STRICT},
			require.True,
		},
		{
			&Deduplication{Strategy: Deduplication_DATAGRAM, Offset: Deduplication_OFFSET_LATEST},
			&Deduplication{Strategy: Deduplication_DATAGRAM, Offset: Deduplication_OFFSET_LATEST},
			require.True,
		},
		{
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"alpha", "bravo"}},
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"bravo", "alpha"}},
			require.True,
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"alpha", "bravo"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"bravo", "alpha", "alpha"}},
			require.True,
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"alpha", "bravo", "bravo"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"bravo", "alpha", "alpha"}},
			require.True,
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"alpha", "bravo", "bravo"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"bravo", "alpha", "alpha", "charlie"}},
			require.False,
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"alpha", "bravo", "bravo", "delta"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"bravo", "alpha", "alpha", "charlie"}},
			require.False,
		},
		{
			&Deduplication{Strategy: Deduplication_STRICT},
			&Deduplication{Strategy: Deduplication_DATAGRAM},
			require.False,
		},
		{
			&Deduplication{Strategy: Deduplication_DATAGRAM},
			&Deduplication{Strategy: Deduplication_DATAGRAM, Offset: Deduplication_OFFSET_LATEST},
			require.False,
		},
		{
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"alpha", "bravo", "charlie"}},
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"bravo", "alpha"}},
			require.False,
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"alpha", "bravo", "delta"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"bravo", "alpha", "alpha", "charlie"}},
			require.False,
		},
	}

	for i, tc := range testCases {
		tc.assert(t, tc.alpha.Equals(tc.bravo), "test case %d failed", i)
	}
}

func TestDeduplicationNormalize(t *testing.T) {
	testCases := []struct {
		in       *Deduplication
		expected *Deduplication
	}{
		{
			&Deduplication{},
			&Deduplication{Strategy: Deduplication_NONE, Offset: Deduplication_OFFSET_EARLIEST},
		},
		{
			&Deduplication{Strategy: Deduplication_DATAGRAM, Offset: Deduplication_OFFSET_LATEST},
			&Deduplication{Strategy: Deduplication_DATAGRAM, Offset: Deduplication_OFFSET_LATEST},
		},
		{
			&Deduplication{Strategy: Deduplication_STRICT, Keys: []string{"foo", "bar"}, Fields: []string{"alpha", "bravo"}},
			&Deduplication{Strategy: Deduplication_STRICT, Offset: Deduplication_OFFSET_EARLIEST},
		},
		{
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST},
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST, Keys: []string{}},
		},
		{
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST, Fields: []string{"foo", "bar"}},
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST, Keys: []string{}},
		},
		{
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST, Keys: []string{"alpha"}},
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST, Keys: []string{"alpha"}},
		},
		{
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"foo", "bar", "alpha", "foo", "alpha", "foo", "bar", "foo"}},
			&Deduplication{Strategy: Deduplication_KEY_GROUPED, Offset: Deduplication_OFFSET_EARLIEST, Keys: []string{"alpha", "bar", "foo"}},
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"FOO", "bar", "Alpha", "foo", "alpha", "foo", "bar", "foo", "Alpha", "FOO"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_KEY, Offset: Deduplication_OFFSET_EARLIEST, Keys: []string{"Alpha", "FOO", "alpha", "bar", "foo"}},
		},
		{
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Fields: []string{"alpha", "bravo", "charlie", "alpha", "bravo", "charlie"}},
			&Deduplication{Strategy: Deduplication_UNIQUE_FIELD, Offset: Deduplication_OFFSET_EARLIEST, Fields: []string{"alpha", "bravo", "charlie"}},
		},
	}

	for i, tc := range testCases {
		tc.in.Normalize()
		require.Equal(t, tc.expected, tc.in, "test case %d failed", i)

		// Default Invariants
		require.NotEqual(t, tc.in.Offset, Deduplication_OFFSET_UNKNOWN, "expected offset to not be unknown after normalization in test case %d", i)
		require.NotEqual(t, tc.in.Strategy, Deduplication_UNKNOWN, "expected strategy to not be unknown after normalization in test case %d", i)

		// Strategy Based Invariants
		switch tc.in.Strategy {
		case Deduplication_NONE, Deduplication_STRICT, Deduplication_DATAGRAM:
			require.Nil(t, tc.in.Keys, "expected keys to be nil for %s strategy (test case %d)", tc.in.Strategy, i)
			require.Nil(t, tc.in.Fields, "expected fields to be nil for %s strategy (test case %d)", tc.in.Strategy, i)
		case Deduplication_KEY_GROUPED, Deduplication_UNIQUE_KEY:
			require.NotNil(t, tc.in.Keys, "expected keys to not be nil for %s strategy (test case %d)", tc.in.Strategy, i)
			require.Nil(t, tc.in.Fields, "expected fields to be nil for %s strategy (test case %d)", tc.in.Strategy, i)
		case Deduplication_UNIQUE_FIELD:
			require.Nil(t, tc.in.Keys, "expected keys to be nil for %s strategy (test case %d)", tc.in.Strategy, i)
			require.NotNil(t, tc.in.Fields, "expected fields to not be nil for %s strategy (test case %d)", tc.in.Strategy, i)
		}

	}
}

const fixturePath = "testdata/events.json"

func loadFixtures() (_ []*EventWrapper, err error) {
	if _, err = os.Stat(fixturePath); os.IsNotExist(err) {
		if err = generateFixtures(); err != nil {
			return nil, err
		}
	}

	var f *os.File
	if f, err = os.Open(fixturePath); err != nil {
		return nil, err
	}
	defer f.Close()

	events := make([]*EventWrapper, 0)
	if err = json.NewDecoder(f).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}

func generateFixtures() (err error) {
	block1 := rands(1024)
	block2 := rands(1024)
	block3 := rands(1892)

	key1 := rands(92)
	key2 := rands(112)
	key3 := rands(32)

	events := []*EventWrapper{
		mkwevt(
			block1,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key1, key2, key3),
			mimetype.ApplicationOctetStream,
			"RandomData v2.11.3",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block2,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key1, key2, key3),
			mimetype.ApplicationOctetStream,
			"RandomData v2.11.3",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block3,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key3, key1, key2),
			mimetype.ApplicationOctetStream,
			"RandomData v2.11.3",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block1,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key2, key3, key1),
			mimetype.ApplicationOctetStream,
			"RandomData v2.11.3",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block1,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key1, key2, key3),
			mimetype.ApplicationAvro,
			"RandomData v2.11.3",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block1,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key1, key2, key3),
			mimetype.ApplicationOctetStream,
			"",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block1,
			fmt.Sprintf("alpha:%s,bravo:%s,charlie:%s", key1, key2, key3),
			mimetype.ApplicationOctetStream,
			"RandomData v2.12.0",
			time.Now().Format(time.RFC3339),
		),
		mkwevt(
			block1,
			"",
			mimetype.ApplicationOctetStream,
			"RandomData v2.11.3",
			time.Now().Format(time.RFC3339),
		),
	}

	var f *os.File
	if f, err = os.Create(fixturePath); err != nil {
		return err
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(events); err != nil {
		return err
	}
	return nil
}

var seq rlid.Sequence

// Helper to quickly make event wrappers with a topicID and committed timestamp, using a
// similar methodology to mkevt to generate the underlying event.
func mkwevt(data, kvs string, mime mimetype.MIME, etype, created string) *EventWrapper {
	event := mkevt(data, kvs, mime, etype, created)
	eventID := seq.Next()

	wrap := &EventWrapper{
		Id:        eventID.Bytes(),
		TopicId:   ulid.MustParse("01HBETJKP2ES10XXMK27M651GA").Bytes(),
		Committed: timestamppb.Now(),
		Offset:    uint64(eventID.Sequence()),
		Epoch:     uint64(0xe1),
		Region:    region.Region_LKE_US_EAST_1A,
		Publisher: &Publisher{
			PublisherId: "01HD1KY309F3SHSV1M89GSQBDF",
			Ipaddr:      "192.168.1.1",
			ClientId:    "data-ingestor-alpha",
			UserAgent:   "PyEnsign v0.12.0",
		},
		Encryption:  &Encryption{EncryptionAlgorithm: Encryption_PLAINTEXT},
		Compression: &Compression{Algorithm: Compression_NONE},
	}
	wrap.Wrap(event)
	return wrap
}

// Create random event wrapper using the mkevt methodology but with randomly generated
// data (1024 bytes). If any of the values of the metadata are "rand" then a randomly
// generated string is populated for that value.
func createRandomEvent(mime mimetype.MIME, etype, meta string) *EventWrapper {
	event := mkevt("", meta, mime, etype, "")
	event.Data = randb(1024)

	for key, val := range event.Metadata {
		if val == "rand" {
			event.Metadata[key] = rands(96)
		}
	}

	eventID := seq.Next()
	wrap := &EventWrapper{
		Id:        eventID.Bytes(),
		TopicId:   ulid.MustParse("01HBETJKP2ES10XXMK27M651GA").Bytes(),
		Committed: timestamppb.Now(),
		Offset:    uint64(eventID.Sequence()),
		Epoch:     uint64(0xe1),
		Region:    region.Region_LKE_US_EAST_1A,
		Publisher: &Publisher{
			PublisherId: "01HD1KY309F3SHSV1M89GSQBDF",
			Ipaddr:      "192.168.1.1",
			ClientId:    "data-ingestor-alpha",
			UserAgent:   "PyEnsign v0.12.0",
		},
		Encryption:  &Encryption{EncryptionAlgorithm: Encryption_PLAINTEXT},
		Compression: &Compression{Algorithm: Compression_NONE},
	}
	wrap.Wrap(event)
	return wrap
}

// Generate random byte array without error or panic.
func randb(s int) []byte {
	data := make([]byte, s)
	rand.Read(data)
	return data
}

// Generate random base64 encoded string without error or panic.
func rands(s int) string {
	return base64.RawStdEncoding.EncodeToString(randb(s))
}
