package api

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/twmb/murmur3"
	"golang.org/x/exp/slices"
)

var (
	ErrNoKeys   = errors.New("no keys specified for key based hashing")
	ErrNoFields = errors.New("no fields specified for field based hashing")
)

//===========================================================================
// Event Hashing Methods
//===========================================================================

// Hash uses the deduplication policy to determine the hash signature of the event
// wrapped by the event wrapper and returns the appropriate signature that should be
// used to detect duplicates in the event stream.
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
// data, mimetype, and type. This method works by setting any non-hash fields to zero
// values then marshaling the protocol buffers of the event and computing the murmur3
// hash on the serialized data.
func (w *EventWrapper) HashStrict() (_ []byte, err error) {
	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return nil, err
	}

	// Create the hash and write the event data to the hash.
	hash := murmur3.New128()
	if _, err = hash.Write(event.Data); err != nil {
		return nil, fmt.Errorf("could not write event data to hash: %w", err)
	}

	// Sort the keys in the metadata to write them in lexicographic order
	keys := make([]string, 0, len(event.Metadata))
	for key := range event.Metadata {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	// Write the metadata to the to the hash
	for _, key := range keys {
		val := event.Metadata[key]
		if _, err = hash.Write([]byte(key + val)); err != nil {
			return nil, fmt.Errorf("could not write metadata key value pair to hash: %w", err)
		}
	}

	// Write the mimetype to the hash
	if err = binary.Write(hash, binary.LittleEndian, event.Mimetype); err != nil {
		return nil, fmt.Errorf("could not hash mimetype: %w", err)
	}

	// Write the type to the hash
	var etype *Type
	if event.Type != nil && !event.Type.IsZero() {
		etype = event.Type
	} else {
		etype = UnspecifiedType
	}

	if _, err = hash.Write([]byte(etype.Name)); err != nil {
		return nil, fmt.Errorf("could not write type name to hash: %w", err)
	}

	if err = binary.Write(hash, binary.LittleEndian, etype.MajorVersion); err != nil {
		return nil, fmt.Errorf("could not hash type major version: %w", err)
	}

	if err = binary.Write(hash, binary.LittleEndian, etype.MinorVersion); err != nil {
		return nil, fmt.Errorf("could not hash type minor version: %w", err)
	}

	if err = binary.Write(hash, binary.LittleEndian, etype.PatchVersion); err != nil {
		return nil, fmt.Errorf("could not hash type patch version: %w", err)
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

// Key grouped hashing returns the murmur3 hash of the data of the event prefixed with
// the metadata values of the the specified keys. E.g. if the data is foobar and the
// hash is grouped by the key month - then for two events with month jan and month feb
// will have different hashes: murmur3(janfoobar) and murmur3(febfoobar).
//
// NOTE: this method does not take into account mimetype or type but in the future we
// may have "reserved keys" to factor in these elements to the hash.
func (w *EventWrapper) HashKeyGrouped(keys []string) (_ []byte, err error) {
	if len(keys) == 0 {
		return nil, ErrNoKeys
	}

	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return nil, err
	}

	hash := murmur3.New128()
	for _, key := range keys {
		val := event.Metadata[key]
		if _, err = hash.Write([]byte(key + val)); err != nil {
			return nil, err
		}
	}

	if _, err = hash.Write(event.Data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// Unique key hashes determine duplicates not from the event data but from the keys
// specified in the metadata (useful for creating lookup indexes). The hash is the
// murmur3 hash of the concatenated key values for the specified keys.
func (w *EventWrapper) HashUniqueKey(keys []string) (_ []byte, err error) {
	if len(keys) == 0 {
		return nil, ErrNoKeys
	}

	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return nil, err
	}

	hash := murmur3.New128()
	for _, key := range keys {
		val := event.Metadata[key]
		if _, err = hash.Write([]byte(key + val)); err != nil {
			return nil, err
		}
	}

	return hash.Sum(nil), nil
}

// Unique field hashing determines duplicates not from the entire datagram, but rather
// from specified fields in the datagram. This requires Ensign to be able to parse the
// data, and unparsable mimetypes (such as protocol buffers) will return an error.
//
// BUG: this is currently unimplemented
func (w *EventWrapper) HashUniqueField(fields []string) (_ []byte, err error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}

	if _, err = w.Unwrap(); err != nil {
		return nil, err
	}

	return nil, errors.New("hash unique field is not implemented")
}
