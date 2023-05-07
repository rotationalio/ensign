package api_test

import (
	"crypto/rand"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
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
}
