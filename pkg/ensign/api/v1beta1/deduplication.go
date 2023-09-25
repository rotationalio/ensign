package api

import (
	"fmt"

	"github.com/twmb/murmur3"
	"google.golang.org/protobuf/proto"
)

//===========================================================================
// Event Hashing Methods
//===========================================================================

func (w *EventWrapper) Hash(policy *Deduplication) ([]byte, error) {
	switch policy.Strategy {
	case Deduplication_NONE:
		return nil, nil
	case Deduplication_STRICT:
		return w.HashStrict()
	case Deduplication_DATAGRAM:
		return w.HashDatagram()
	case Deduplication_KEY_GROUPED:
		return w.HashKeyGrouped(policy.Keys)
	case Deduplication_UNIQUE_KEY:
		return w.HashUniqueKey(policy.Keys)
	case Deduplication_UNIQUE_FIELD:
		return w.HashUniqueField(policy.Keys)
	default:
		return nil, fmt.Errorf("unknown deduplication strategy %s", policy.Strategy)
	}
}

// Strict hashing is used to detect duplicates where two events have identical metadata,
// data, mimetype, and type.
func (w *EventWrapper) HashStrict() (_ []byte, err error) {
	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return nil, err
	}

	// Set any field that should not be in the hash to nil
	event.Created = nil

	// If there is no type, add the unspecified type
	if event.Type == nil || event.Type.IsZero() {
		event.Type = UnspecifiedType
	}

	var data []byte
	if data, err = proto.Marshal(event); err != nil {
		return nil, err
	}

	hash := murmur3.New128()
	if _, err = hash.Write(data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// Datagram hashing is used to detect duplicates in data only, ignoring metadata,
// mimetype, and type as in strict hashing. This method returns a murmur3 hash of the
// data field of the event only.
func (w *EventWrapper) HashDatagram() (_ []byte, err error) {
	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return nil, err
	}

	hash := murmur3.New128()
	if _, err = hash.Write(event.Data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func (w *EventWrapper) HashKeyGrouped(keys []string) ([]byte, error) {
	return nil, nil
}

func (w *EventWrapper) HashUniqueKey(keys []string) ([]byte, error) {
	return nil, nil
}

func (w *EventWrapper) HashUniqueField(fields []string) ([]byte, error) {
	return nil, nil
}
