package api

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/twmb/murmur3"
	"golang.org/x/exp/slices"
)

//===========================================================================
// Deduplication Methods
//===========================================================================

// Duplicates uses a non-hash equality method to determine if the input event is a
// duplicate of the current event using the deduplication policy. Hashing shoud be used
// to determine duplication candidates, but the duplicates method should be used to
// confirm if two events are duplicates or not.
func (w *EventWrapper) Duplicates(o *EventWrapper, policy *Deduplication) (bool, error) {
	switch policy.Strategy {
	case Deduplication_NONE:
		return false, nil
	case Deduplication_STRICT:
		return w.DuplicatesStrict(o)
	case Deduplication_DATAGRAM:
		return w.DuplicatesDatagram(o)
	case Deduplication_KEY_GROUPED:
		return w.DuplicatesKeyGrouped(o, policy.Keys)
	case Deduplication_UNIQUE_KEY:
		return w.DuplicatesUniqueKey(o, policy.Keys)
	case Deduplication_UNIQUE_FIELD:
		return w.DuplicatesUniqueField(o, policy.Fields)
	default:
		return false, fmt.Errorf("unknown deduplication strategy %s", policy.Strategy)
	}
}

// Strict deduplication requires that the events data, metadata, mimetype, and type are
// all identical in order for an event to be marked a duplicate. This method uses the
// Event.Equals method for comparing the wrapped events in the source and target.
func (w *EventWrapper) DuplicatesStrict(o *EventWrapper) (_ bool, err error) {
	var we *Event
	if we, err = w.Unwrap(); err != nil {
		return false, err
	}

	var oe *Event
	if oe, err = o.Unwrap(); err != nil {
		return false, err
	}

	return we.Equals(oe), nil
}

// Datagram duplicates only compare the event's data to determine duplication, ignoring
// the metadata, mimetype, and type fields. This method uses Event.DataEquals.
func (w *EventWrapper) DuplicatesDatagram(o *EventWrapper) (_ bool, err error) {
	var we *Event
	if we, err = w.Unwrap(); err != nil {
		return false, err
	}

	var oe *Event
	if oe, err = o.Unwrap(); err != nil {
		return false, err
	}

	return we.DataEquals(oe), nil
}

// Key grouped duplicates must have identical values for the specified keys (if not
// then the events are not considered duplicates, even if the data is the same in both
// events), then the events must have identical data. This method uses the
// Event.MetaEquals first, then Event.DataEquals second.
func (w *EventWrapper) DuplicatesKeyGrouped(o *EventWrapper, keys []string) (_ bool, err error) {
	if len(keys) == 0 {
		return false, ErrNoKeys
	}

	var we *Event
	if we, err = w.Unwrap(); err != nil {
		return false, err
	}

	var oe *Event
	if oe, err = o.Unwrap(); err != nil {
		return false, err
	}

	// First check that the keys match (otherwise they're not in the same group)
	if !we.MetaEquals(oe, keys...) {
		return false, nil
	}

	// If these events are in the same group; e.g. the keys match, then check the data.
	return we.DataEquals(oe), nil
}

// Unique key duplication only checks that the events have the same values for the keys
// specified in the policy, ignoring other keys, data, mimetype, and type information.
// This method uses Event.MetaEquals to perform the comparison.
func (w *EventWrapper) DuplicatesUniqueKey(o *EventWrapper, keys []string) (_ bool, err error) {
	if len(keys) == 0 {
		return false, ErrNoKeys
	}

	var we *Event
	if we, err = w.Unwrap(); err != nil {
		return false, err
	}

	var oe *Event
	if oe, err = o.Unwrap(); err != nil {
		return false, err
	}

	return we.MetaEquals(oe, keys...), nil
}

// Unique field duplication focuses on data duplication but rather than checking the
// entire datagram, parses the data and only compares specified fields. This requires
// Ensign to be able to parse the data and unparseable mimetypes (such as protocol
// buffers) will return an error.
//
// BUG: this is currently unimplemented
func (w *EventWrapper) DuplicatesUniqueField(o *EventWrapper, fields []string) (_ bool, err error) {
	if len(fields) == 0 {
		return false, ErrNoFields
	}

	if _, err = w.Unwrap(); err != nil {
		return false, err
	}

	if _, err = o.Unwrap(); err != nil {
		return false, err
	}

	return false, errors.New("hash unique field is not implemented")
}

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
	etype := event.ResolveType()
	if _, err = hash.Write([]byte(etype.Name)); err != nil {
		return nil, fmt.Errorf("could not write type name to hash: %w", err)
	}

	// Note: ensure all integers are hashed with the same byte order.
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

//===========================================================================
// Duplicate Modification
//===========================================================================

// DuplicateOf marks the original event (w) as a duplicate of the original event (o).
// In other words, the original event (w) becomes a duplicate reference to the original.
// The duplicate is updated in place to mark the wrapper as a duplicate of the original
// and to reduce the data storage depending on the policy. For example, in strict mode,
// the event data is nilified and only the wrapper metadata is kept, whereas in unique
// keys mode, the data may not be removed depending on the policy.
//
// NOTE: this method does not check if the events are duplicates! Use the Duplicates()
// method for verification that the two events are duplicates of each other.
func (w *EventWrapper) DuplicateOf(o *EventWrapper, policy *Deduplication) (err error) {
	if policy.Strategy == Deduplication_UNKNOWN || policy.Strategy == Deduplication_NONE {
		return ErrDuplicatesNotAllowed
	}

	// Mark the current event as a duplicate of the original event
	w.IsDuplicate = true
	w.DuplicateId = o.Id

	// Remove fields as necessary to reduce the storage load of the database.
	// In STRICT mode - simply set the event to nil so we're not storing that data.
	// Since there will be no encryption or compression, set those to nil as well.
	// NOTE: this strategy will cause us to lose the created timestamp from the event,
	// though the committed timestamp on the wrapper will still be kept.
	if policy.Strategy == Deduplication_STRICT {
		w.Event = nil
		w.Encryption = nil
		w.Compression = nil
		return nil
	}

	// For all other modes we'll need to update the event record itself.
	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return err
	}

	var orig *Event
	if orig, err = o.Unwrap(); err != nil {
		return err
	}

	// If the policy is datagram or key grouped, then we know the data is identical so
	// set it to nil. Otherwise, check to make sure the data is identical before
	// deleting the data from the event
	if policy.Strategy == Deduplication_DATAGRAM || policy.Strategy == Deduplication_KEY_GROUPED {
		event.Data = nil
	} else if bytes.Equal(event.Data, orig.Data) {
		event.Data = nil
	}

	// Deduplicate the metadata, keeping only the keys that differ from the original
	var meta map[string]string
	for key, val := range event.Metadata {
		if val != orig.Metadata[key] {
			if meta == nil {
				meta = make(map[string]string)
			}
			meta[key] = val
		}
	}
	event.Metadata = meta

	// Deduplicate the mimetype if they match
	if event.Mimetype == orig.Mimetype {
		event.Mimetype = 0
	}

	// Deduplicate the event type if they match
	if event.Type.Equals(orig.Type) {
		event.Type = nil
	}

	// Rewrap the event into the wrapper
	return w.Wrap(event)
}

// DuplicateFrom is the inverse of DuplicateOf: it modifies the event w, populating it
// with the duplicated data from the original event. It still keeps the event w marked
// as a duplicate (this has to be undone manually), but it allows the duplicate to be
// returned to the user with any unique information it may have contained.
func (w *EventWrapper) DuplicateFrom(o *EventWrapper) (err error) {
	// If the event data is nil, simply copy over the original event, encryption, and
	// compression and return w. This is the case in strict mode.
	if len(w.Event) == 0 {
		w.Event = o.Event
		w.Encryption = o.Encryption
		w.Compression = o.Compression
		return nil
	}

	// Otherwise unwrap the events to merge the data together
	var event *Event
	if event, err = w.Unwrap(); err != nil {
		return err
	}

	var orig *Event
	if orig, err = o.Unwrap(); err != nil {
		return err
	}

	// If there is no event data, merge it from the original event.
	if len(event.Data) == 0 {
		event.Data = orig.Data
	}

	// If there is no event metadata, copy it, otherwise only copy the keys that are
	// in original but not in the event.
	//
	// NOTE: this has the unfortunate side-effect that keys that were in the original
	// event but were not in the original duplicate event will be copied over, causing
	// a small data inconsistency issue.
	if len(event.Metadata) == 0 {
		event.Metadata = orig.Metadata
	} else {
		for key, val := range orig.Metadata {
			if _, ok := event.Metadata[key]; !ok {
				event.Metadata[key] = val
			}
		}
	}

	if event.Mimetype == 0 {
		event.Mimetype = orig.Mimetype
	}

	if event.Type == nil {
		event.Type = orig.Type
	}

	return w.Wrap(event)
}

//===========================================================================
// Policy Helpers
//===========================================================================

// Equals compares deduplication polices to see if they would be implemented the same.
// It first compares the strategy, and returns false if the strategies are different. If
// the strategies are identical, it then compares keys for the key grouped and unique
// key strategies and fields for the unique fields strategy.
//
// NOTE: This method normalizes both deduplication policy structs, which might change
// the underlying data stored in the pointer.
func (d *Deduplication) Equals(o *Deduplication) bool {
	// The policies must be normalized before comparison.
	d.Normalize()
	o.Normalize()

	if d.Strategy != o.Strategy {
		return false
	}

	if d.Offset != o.Offset {
		return false
	}

	switch d.Strategy {
	case Deduplication_KEY_GROUPED, Deduplication_UNIQUE_KEY:
		// Keys should not be nil after normalization
		if len(d.Keys) != len(o.Keys) {
			return false
		}

		// Keys are expected to be deduplicated and sorted during normalization
		for i, key := range d.Keys {
			if key != o.Keys[i] {
				return false
			}
		}
	case Deduplication_UNIQUE_FIELD:
		// Fields should not b enil after normalization
		if len(d.Fields) != len(o.Fields) {
			return false
		}

		// Fields are expected to be deduplicated and sorted during normalization
		for i, field := range d.Fields {
			if field != o.Fields[i] {
				return false
			}
		}
	}

	return true
}

// Normalize the deduplication policy based on the strategy. If the strategy does not
// require keys or fields, then keys and fields are set to nil (no matter user input),
// if the strategy does require keys or fields then they are sorted and deduplicated.
//
// NOTE: This method also sets the offset to the default if it is unknown.
// NOTE: This method sets the deduplication strategy to None if it is unknown
func (d *Deduplication) Normalize() *Deduplication {
	// Set the offset to the default offset if unknown.
	if d.Offset == Deduplication_OFFSET_UNKNOWN {
		d.Offset = Deduplication_OFFSET_EARLIEST
	}

	// Set the strategy to none if it is unknown
	if d.Strategy == Deduplication_UNKNOWN {
		d.Strategy = Deduplication_NONE
	}

	switch d.Strategy {
	case Deduplication_NONE, Deduplication_STRICT, Deduplication_DATAGRAM:
		d.Keys = nil
		d.Fields = nil
	case Deduplication_KEY_GROUPED, Deduplication_UNIQUE_KEY:
		d.Keys = uniqueSort(d.Keys)
		d.Fields = nil
	case Deduplication_UNIQUE_FIELD:
		d.Keys = nil
		d.Fields = uniqueSort(d.Fields)
	}

	return d
}

// Validates that the deduplication strategy can be implemented after normalization.
func (d *Deduplication) Validate() error {
	switch d.Strategy {
	case Deduplication_UNKNOWN, Deduplication_NONE, Deduplication_STRICT, Deduplication_DATAGRAM:
		if len(d.Keys) > 0 {
			return ErrKeysNotAllowed
		}
		if len(d.Fields) > 0 {
			return ErrFieldsNotAllowed
		}
	case Deduplication_KEY_GROUPED, Deduplication_UNIQUE_KEY:
		if len(d.Keys) == 0 {
			return ErrNoKeys
		}
		if len(d.Fields) > 0 {
			return ErrFieldsNotAllowed
		}
	case Deduplication_UNIQUE_FIELD:
		if len(d.Fields) == 0 {
			return ErrNoFields
		}
		if len(d.Keys) > 0 {
			return ErrKeysNotAllowed
		}
	}

	return nil
}

func uniqueSort(s []string) []string {
	if s == nil {
		return make([]string, 0)
	}

	if len(s) <= 1 {
		return s
	}

	uniques := make(map[string]struct{}, len(s))
	for _, item := range s {
		uniques[item] = struct{}{}
	}

	r := make([]string, 0, len(uniques))
	for item := range uniques {
		r = append(r, item)
	}

	slices.Sort(r)
	return r
}
