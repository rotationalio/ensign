package api

import (
	"bytes"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"google.golang.org/protobuf/proto"
)

// Unspecified is the type name of the unspecified type.
const Unspecified = "Unspecified"

// UnspecifiedType is returned when the event has no type.
var UnspecifiedType = &Type{Name: Unspecified}

//===========================================================================
// Event Wrapper Helper Methods
//===========================================================================

// Wrap an event inside of the event wrapper, marshaling the event into bytes and
// storing it in its raw form so that it doesn't have to be parsed during wrapper
// unmarshaling (the Broker uses the event wrapper metadata not the event itself).
func (w *EventWrapper) Wrap(e *Event) (err error) {
	if w.Event, err = proto.Marshal(e); err != nil {
		return err
	}
	return nil
}

// Unwrap an event from the event wrapper, marshaling the event bytes into an event
// protocol buffer for event-specific processing.
func (w *EventWrapper) Unwrap() (e *Event, err error) {
	if len(w.Event) == 0 {
		return nil, ErrNoEvent
	}

	e = &Event{}
	if err = proto.Unmarshal(w.Event, e); err != nil {
		return nil, err
	}
	return e, nil
}

// Parse the topicID on the event wrapper as a ULID.
func (w *EventWrapper) ParseTopicID() (topicID ulid.ULID, err error) {
	topicID = ulid.ULID{}
	if err = topicID.UnmarshalBinary(w.TopicId); err != nil {
		return topicID, err
	}
	return topicID, nil
}

// Parwse the eventID on the event wrapper as an RLID.
func (w *EventWrapper) ParseEventID() (eventID rlid.RLID, err error) {
	eventID = rlid.RLID{}
	if err = eventID.UnmarshalBinary(w.Id); err != nil {
		return eventID, err
	}
	return eventID, nil
}

//===========================================================================
// Event Helper Methods
//===========================================================================

// ResolveType returns the event's type if it has one, otherwise if the event's type is
// nil or empty, returns the "Unspecified" type, which is the default type for typeless
// events. It is important to have a named unspecified type for type checking and
// downstream event logging (such a logging in tenant).
func (e *Event) ResolveType() *Type {
	if e.Type == nil || e.Type.IsZero() {
		return UnspecifiedType
	}
	return e.Type
}

// Equals returns strict equality of an event. The event's mimetype and type must match
// and the data must be equal. Finally, the events must have identical metadata - e.g.
// the same keys and values (without omission). Note that the created timestamp is not
// included in the equality check.
func (e *Event) Equals(o *Event) bool {
	// Handle nil event comparison
	if (o == nil) != (e == nil) {
		return false
	} else {
		if o == nil && e == nil {
			return true
		}
	}

	if e.Mimetype != o.Mimetype {
		return false
	}

	// Type checks should resolve the type first to ensure that nil and zero-valued
	// types are compared the same way (e.g. as unspecified types).
	if !e.ResolveType().Equals(o.ResolveType()) {
		return false
	}

	if !bytes.Equal(e.Data, o.Data) {
		return false
	}

	// If the number of keys doesn't match, then it's impossible for the events to have
	// identical metadata; this check allows us to loop over the keys from one of the
	// metadata without having to worry about the intersection of keys. See the source
	// code from reflect.DeepEqual (the Map case) as a comparison.
	//
	// NOTE: nil metadata will be equal to empty metadata in this case.
	if len(e.Metadata) != len(o.Metadata) {
		return false
	}

	for key, vala := range e.Metadata {
		if valb, ok := o.Metadata[key]; !ok || vala != valb {
			return false
		}
	}

	return true
}

//===========================================================================
// Publisher Helper Methods
//===========================================================================

// Returns the user-specified client ID if set, otherwise returns the publisher ID.
func (p *Publisher) ResolveClientID() string {
	if p.ClientId != "" {
		return p.ClientId
	}
	return p.PublisherId
}

//===========================================================================
// Type Helper Methods
//===========================================================================

// Type equality checking, the names must match (currently case-sensitive) and the
// major, minor, and patch versions, must also match. Two zero valued types will be
// equal with one another.
func (t *Type) Equals(o *Type) bool {
	// Handle nil type comparison
	if (o == nil) != (t == nil) {
		return false
	} else {
		if o == nil && t == nil {
			return true
		}
	}

	return t.Name == o.Name && t.MajorVersion == o.MajorVersion && t.MinorVersion == o.MinorVersion && t.PatchVersion == o.PatchVersion
}

// IsZero returns true if the name is empty or unspecified and the major, minor, and
// patch versions are equal to zero.
func (t *Type) IsZero() bool {
	return (t.Name == "" || t.Name == Unspecified) && t.MajorVersion == 0 && t.MinorVersion == 0 && t.PatchVersion == 0
}
