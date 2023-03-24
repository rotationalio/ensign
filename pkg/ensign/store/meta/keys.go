package meta

import (
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

// An Index Key is a ULID that must be unique across all objects; e.g. Topics must have
// unique keys across projects. This can be guaranteed by using ulids.New for creating
// IDs. If a key can be created by a user e.g. for Groups, it must not use the object
// key accessor since we cannot guarantee that users will create unique keys.
type IndexKey [16]byte

// ObjectKey is the 16 byte project ID followed by a 2 byte object segment then a 16
// byte unique object ID. The IndexKey maps to the object key to allow for easy lookups.
type ObjectKey [34]byte

// Segments ensure that different objects are stored contiguously in the database
// ordered by their project then their ID to make it easy to scan for objects.
type Segment [2]byte

// Segments currently in use by Ensign
var (
	TopicSegment      = Segment{0x74, 0x70}
	TopicNamesSegment = Segment{0x54, 0x6e}
	GroupSegment      = Segment{0x47, 0x50}
)

func CreateIndex(objectID ulid.ULID) (key IndexKey, err error) {
	if ulids.IsZero(objectID) {
		return key, errors.ErrKeyNull
	}
	return IndexKey(objectID), nil
}

func CreateKey(parentID, objectID ulid.ULID, segment Segment) (key ObjectKey, err error) {
	if ulids.IsZero(parentID) || ulids.IsZero(objectID) {
		return key, errors.ErrKeyNull
	}

	if err = parentID.MarshalBinaryTo(key[:16]); err != nil {
		return key, err
	}

	copy(key[16:18], segment[:])

	if err = objectID.MarshalBinaryTo(key[18:]); err != nil {
		return key, err
	}
	return key, nil
}

func (k *ObjectKey) Key() IndexKey {
	return IndexKey(*(*[16]byte)(k[18:]))
}

func (k *ObjectKey) UnmarshalValue(data []byte) error {
	if len(data) != 34 {
		return errors.ErrKeyWrongSize
	}
	copy(k[:], data)
	return nil
}

func (k *ObjectKey) ParentID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[:16])
	return id, err
}

func (k *ObjectKey) ObjectID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[18:])
	return id, err
}

func (k *ObjectKey) Segment() (Segment, error) {
	return Segment(*(*[2]byte)(k[16:18])), nil
}
