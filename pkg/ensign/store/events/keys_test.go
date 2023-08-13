package events_test

import (
	"bytes"
	"math"
	"math/rand"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/events"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	makeSegmentTest := func(segment events.Segment) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()

			testCases := []struct {
				topicID ulid.ULID
				eventID rlid.RLID
				err     error
			}{
				{ulids.Null, rlid.Null, errors.ErrKeyNull},
				{ulids.Null, rlid.Make(42), errors.ErrKeyNull},
				{ulid.MustParse("01H7QTE9ACE52PN22MTVSF70W0"), rlid.Null, errors.ErrKeyNull},
				{ulid.MustParse("01H7QTE9ACE52PN22MTVSF70W0"), rlid.Make(0), nil},
				{ulid.MustParse("01H7QTE9ACE52PN22MTVSF70W0"), rlid.Make(uint32(rand.Int31n(math.MaxInt32) + 1)), nil},
				{ulids.New(), rlid.Make(uint32(rand.Int31n(math.MaxInt32) + 1)), nil},
				{ulids.New(), rlid.Make(uint32(rand.Int31n(math.MaxInt32) + 1)), nil},
				{ulids.New(), rlid.Make(uint32(rand.Int31n(math.MaxInt32) + 1)), nil},
				{ulids.New(), rlid.Make(uint32(rand.Int31n(math.MaxInt32) + 1)), nil},
			}

			for i, tc := range testCases {
				key, err := events.CreateKey(tc.topicID, tc.eventID, segment)
				if tc.err != nil {
					require.ErrorIs(t, err, tc.err, "test case %d failed", i)
				} else {
					require.NoError(t, err, "test case %d failed", i)
					require.Equal(t, tc.topicID, key.TopicID(), "test case %d failed", i)
					require.Equal(t, tc.eventID, key.EventID(), "test case %d failed", i)
					require.Equal(t, segment, key.Segment(), "test case %d failed", i)
				}
			}
		}
	}

	t.Run("Event", makeSegmentTest(events.EventSegment))
	t.Run("MetaEvent", makeSegmentTest(events.MetaEventSegment))
}

func TestMakeKey(t *testing.T) {
	topicID := ulids.New()
	rlid0 := rlid.Make(0)
	rlid42 := rlid.Make(42)

	testCases := []struct {
		event   *api.EventWrapper
		topicID ulid.ULID
		eventID rlid.RLID
		err     error
	}{
		{&api.EventWrapper{}, ulids.Null, rlid.Null, errors.ErrInvalidKey},
		{&api.EventWrapper{TopicId: ulids.New().Bytes()}, ulids.Null, rlid.Null, errors.ErrInvalidKey},
		{&api.EventWrapper{Id: rlid.Make(42).Bytes()}, ulids.Null, rlid.Null, errors.ErrInvalidKey},
		{&api.EventWrapper{Id: rlid.Null[:], TopicId: ulids.Null[:]}, ulids.Null, rlid.Null, errors.ErrKeyNull},
		{&api.EventWrapper{Id: rlid0[:], TopicId: topicID[:]}, topicID, rlid0, nil},
		{&api.EventWrapper{Id: rlid42[:], TopicId: topicID[:]}, topicID, rlid42, nil},
		{&api.EventWrapper{TopicId: rlid42[:], Id: topicID[:]}, topicID, rlid42, errors.ErrInvalidKey},
	}

	for i, tc := range testCases {
		if tc.err != nil {
			_, err := events.EventKey(tc.event)
			require.ErrorIs(t, err, tc.err, "test case %d failed", i)

			_, err = events.MetaEventKey(tc.event)
			require.ErrorIs(t, err, tc.err, "test case %d failed", i)

			continue
		}

		key, err := events.EventKey(tc.event)
		require.NoError(t, err, "test case %d failed", i)
		require.Equal(t, tc.topicID, key.TopicID(), "test case %d failed", i)
		require.Equal(t, tc.eventID, key.EventID(), "test case %d failed", i)
		require.Equal(t, events.EventSegment, key.Segment(), "test case %d failed", i)

		key, err = events.MetaEventKey(tc.event)
		require.NoError(t, err, "test case %d failed", i)
		require.Equal(t, tc.topicID, key.TopicID(), "test case %d failed", i)
		require.Equal(t, tc.eventID, key.EventID(), "test case %d failed", i)
		require.Equal(t, events.MetaEventSegment, key.Segment(), "test case %d failed", i)

	}
}

func TestEmptyKeyNoPanic(t *testing.T) {
	var key events.Key
	require.Equal(t, bytes.Repeat([]byte{0x00}, 28), key[:])
	require.Equal(t, ulids.Null, key.TopicID())
	require.Equal(t, rlid.Null, key.EventID())
	require.Equal(t, "unknown", key.Segment().String())
}

func TestSegment(t *testing.T) {
	require.Equal(t, []byte("߷"), events.EventSegment[:])
	require.Equal(t, []byte("֍"), events.MetaEventSegment[:])
	require.Equal(t, "event", events.EventSegment.String())
	require.Equal(t, "metaevent", events.MetaEventSegment.String())
	require.Equal(t, "unknown", events.Segment([2]byte{0x00, 0x42}).String())
}
