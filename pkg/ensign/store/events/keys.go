package events

import (
	"bytes"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

// Key is is the 16 byte topic ID followed by a 2 byte segment that splits events from
// meta events then the 10 byte RLID for an event ID.
type Key [28]byte

// Segments ensure that a topic's events and meta events are stored contiguously.
type Segment [2]byte

// Segments in use by the Ensign event store
var (
	EventSegment     = Segment{0xdf, 0xb7}
	MetaEventSegment = Segment{0xd6, 0x8d}
)

func CreateKey(topicID ulid.ULID, eventID rlid.RLID, segment Segment) (key Key, err error) {
	if ulids.IsZero(topicID) || rlid.IsZero(eventID) {
		return key, errors.ErrKeyNull
	}

	copy(key[:16], topicID[:])
	copy(key[16:18], segment[:])
	copy(key[18:], eventID[:])

	return key, nil
}

func EventKey(event *api.EventWrapper) (key Key, err error) {
	return makeKey(event, EventSegment)
}

func MetaEventKey(event *api.EventWrapper) (key Key, err error) {
	return makeKey(event, MetaEventSegment)
}

func makeKey(event *api.EventWrapper, segment Segment) (key Key, err error) {
	// Validate the event key
	switch {
	case len(event.Id) == 0:
		return key, errors.ErrEventMissingId
	case len(event.TopicId) == 0:
		return key, errors.ErrEventMissingTopicId
	case len(event.Id) != 10 || bytes.Equal(event.Id, rlid.Null[:]):
		return key, errors.ErrEventInvalidId
	case len(event.TopicId) != 16 || bytes.Equal(event.TopicId, ulids.Null[:]):
		return key, errors.ErrEventInvalidTopicId
	}

	copy(key[:16], event.TopicId)
	copy(key[16:18], segment[:])
	copy(key[18:], event.Id)

	return key, nil
}

// TopicID parses and returns the topicID ULID from the key.
func (k *Key) TopicID() (topicID ulid.ULID) {
	if err := topicID.UnmarshalBinary(k[:16]); err != nil {
		// a panic should not happen since the key is fixed length.
		panic(err)
	}
	return topicID
}

// EventID parses and returns the eventID RLID from the key.
func (k *Key) EventID() (eventID rlid.RLID) {
	if err := eventID.UnmarshalBinary(k[18:]); err != nil {
		// a panic should not happen since the key is fixed length.
		panic(err)
	}
	return eventID
}

func (k *Key) Segment() Segment {
	return Segment(*(*[2]byte)(k[16:18]))
}

func (s Segment) String() string {
	switch s {
	case EventSegment:
		return "event"
	case MetaEventSegment:
		return "metaevent"
	default:
		return "unknown"
	}
}
