package api

import (
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"google.golang.org/protobuf/proto"
)

// Unspecified is the type name of the unspecified type.
const Unspecified = "Unspecified"

// UnspecifiedType is returned when the event has no type.
var UnspecifiedType = &Type{Name: Unspecified}

func (w *EventWrapper) Wrap(e *Event) (err error) {
	if w.Event, err = proto.Marshal(e); err != nil {
		return err
	}
	return nil
}

func (w *EventWrapper) Unwrap() (e *Event, err error) {
	if len(w.Event) == 0 {
		return nil, errors.New("event wrapper contains no event")
	}

	e = &Event{}
	if err = proto.Unmarshal(w.Event, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (w *EventWrapper) ParseTopicID() (topicID ulid.ULID, err error) {
	topicID = ulid.ULID{}
	if err = topicID.UnmarshalBinary(w.TopicId); err != nil {
		return topicID, err
	}
	return topicID, nil
}

func (w *EventWrapper) ParseEventID() (eventID rlid.RLID, err error) {
	eventID = rlid.RLID{}
	if err = eventID.UnmarshalBinary(w.Id); err != nil {
		return eventID, err
	}
	return eventID, nil
}

// ResolveType returns the event's type if it has one, otherwise if the event's type is
// nil or empty, returns the "Unspecified" type, which is the default type for typeless
// events.
func (w *Event) ResolveType() *Type {
	if w.Type == nil || w.Type.IsZero() {
		return UnspecifiedType
	}
	return w.Type
}

// Returns the user-specified client ID if set, otherwise returns the publisher ID.
func (p *Publisher) ResolveClientID() string {
	if p.ClientId != "" {
		return p.ClientId
	}
	return p.PublisherId
}

func (t *Type) Equals(o *Type) bool {
	if o == nil {
		return false
	}
	return t.Name == o.Name && t.MajorVersion == o.MajorVersion && t.MinorVersion == o.MinorVersion && t.PatchVersion == o.PatchVersion
}

func (t *Type) IsZero() bool {
	return (t.Name == "" || t.Name == Unspecified) && t.MajorVersion == 0 && t.MinorVersion == 0 && t.PatchVersion == 0
}
