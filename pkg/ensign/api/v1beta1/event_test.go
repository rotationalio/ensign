package api_test

import (
	"crypto/rand"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestEventWrapper(t *testing.T) {
	// Create an event wrapper and random event data
	wrap := &api.EventWrapper{
		Id:      ulid.Make().Bytes(),
		TopicId: ulid.Make().Bytes(),
		Offset:  421,
		Epoch:   23,
	}

	evt := &api.Event{
		Data: make([]byte, 128),
	}
	rand.Read(evt.Data)

	err := wrap.Wrap(evt)
	require.NoError(t, err, "should be able to wrap an event in an event wrapper")

	cmp, err := wrap.Unwrap()
	require.NoError(t, err, "should be able to unwrap an event in an event wrapper")
	require.NotNil(t, cmp, "the unwrapped event should not be nil")
	require.True(t, proto.Equal(evt, cmp), "the unwrapped event should match the original")
	require.NotSame(t, evt, cmp, "a pointer to the same event should not be returned")

	wrap.Event = nil
	empty, err := wrap.Unwrap()
	require.EqualError(t, err, "event wrapper contains no event")
	require.Empty(t, empty, "no data event should be zero-valued")

	wrap.Event = []byte("foo")
	_, err = wrap.Unwrap()
	require.Error(t, err, "should not be able to unwrap non-protobuf data")
}

func TestEventWrapperIDParsing(t *testing.T) {
	testCases := []struct {
		eventID  []byte
		expected rlid.RLID
		err      error
	}{
		{nil, rlid.RLID{}, rlid.ErrDataSize},
		{[]byte{}, rlid.RLID{}, rlid.ErrDataSize},
		{[]byte("foo"), rlid.RLID{}, rlid.ErrDataSize},
		{rlid.MustParse("064zwbj8vg00000n").Bytes(), rlid.MustParse("064zwbj8vg00000n"), nil},
	}

	for i, tc := range testCases {
		event := &api.EventWrapper{Id: tc.eventID}
		eventID, err := event.ParseEventID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, rlid.Null, eventID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, eventID, "test case %d failed", i)
		}
	}
}

func TestEventWrapperTopicIDParsing(t *testing.T) {
	testCases := []struct {
		topicID  []byte
		expected ulid.ULID
		err      error
	}{
		{nil, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte{}, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte("foo"), ulid.ULID{}, ulid.ErrDataSize},
		{ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ").Bytes(), ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ"), nil},
	}

	for i, tc := range testCases {
		event := &api.EventWrapper{TopicId: tc.topicID}
		topicID, err := event.ParseTopicID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, ulid.ULID{}, topicID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, topicID, "test case %d failed", i)
		}
	}
}

func TestResolveType(t *testing.T) {
	testCases := []struct {
		event    *api.Event
		expected *api.Type
	}{
		{&api.Event{}, api.UnspecifiedType},
		{&api.Event{Type: nil}, api.UnspecifiedType},
		{&api.Event{Type: &api.Type{}}, api.UnspecifiedType},
		{&api.Event{Type: &api.Type{Name: "TestType", MajorVersion: 1}}, &api.Type{Name: "TestType", MajorVersion: 1}},
	}

	for i, tc := range testCases {
		actual := tc.event.ResolveType()
		require.Equal(t, tc.expected, actual, "test case %d failed", i)
	}
}

func TestResolveClientID(t *testing.T) {
	testCases := []struct {
		pub      *api.Publisher
		expected string
	}{
		{&api.Publisher{}, ""},
		{&api.Publisher{PublisherId: "testpub1"}, "testpub1"},
		{&api.Publisher{ClientId: "testclient1"}, "testclient1"},
		{&api.Publisher{PublisherId: "testpub1", ClientId: "testclient1"}, "testclient1"},
	}

	for i, tc := range testCases {
		actual := tc.pub.ResolveClientID()
		require.Equal(t, tc.expected, actual, "test case %d failed", i)
	}
}

func TestTypeEquality(t *testing.T) {
	testCases := []struct {
		in      *api.Type
		require require.BoolAssertionFunc
	}{
		{nil, require.False},
		{&api.Type{}, require.False},
		{&api.Type{Name: "TESTTYPE", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "testtype", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "testType", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 4, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 1, MinorVersion: 7, PatchVersion: 3}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 1, MinorVersion: 2, PatchVersion: 14}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.True},
	}

	etype := &api.Type{
		Name:         "TestType",
		MajorVersion: 1,
		MinorVersion: 2,
		PatchVersion: 3,
	}

	for i, tc := range testCases {
		tc.require(t, etype.Equals(tc.in), "test case %d failed", i)
	}
}

func TestTypeIsZero(t *testing.T) {
	testCases := []struct {
		in      *api.Type
		require require.BoolAssertionFunc
	}{
		{&api.Type{}, require.True},
		{api.UnspecifiedType, require.True},
		{&api.Type{Name: "TestType", MajorVersion: 0, MinorVersion: 0, PatchVersion: 0}, require.False},
		{&api.Type{Name: "", MajorVersion: 1, MinorVersion: 0, PatchVersion: 0}, require.False},
		{&api.Type{Name: "", MajorVersion: 0, MinorVersion: 2, PatchVersion: 0}, require.False},
		{&api.Type{Name: "", MajorVersion: 0, MinorVersion: 0, PatchVersion: 3}, require.False},
	}

	for i, tc := range testCases {
		tc.require(t, tc.in.IsZero(), "test case %d failed", i)
	}
}
