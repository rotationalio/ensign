package db

import (
	"bytes"
	"context"

	"github.com/oklog/ulid/v2"
)

// Keys Namespace maps object IDs to their keys that are prefixed with a parent ID.
const KeysNamespace = "object_keys"

// Key is composed of two concatenated IDs. The first 16 bytes are the of the parent
// and the second 16 bytes are the ID of the object.
type Key [32]byte

var NullKey = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var NullID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var _ Model = &Key{}

// NewKey constructs a Key from the parent and object IDs.
func NewKey(parentID, objectID ulid.ULID) (key *Key, err error) {
	key = &Key{}
	if err := parentID.MarshalBinaryTo(key[:16]); err != nil {
		return nil, err
	}

	if err := objectID.MarshalBinaryTo(key[16:]); err != nil {
		return nil, err
	}
	return key, nil
}

// Keys are stored by object ID. Since object IDs are locked monotonically increasing
// ulids they are guaranteed to be unique.
func (k *Key) Key() ([]byte, error) {
	if bytes.Compare(k[16:], NullID) == 0 {
		return nil, ErrKeyNoID
	}

	return k[16:], nil
}

func (k *Key) Namespace() string {
	return KeysNamespace
}

func (k *Key) MarshalValue() ([]byte, error) {
	return k[:], nil
}

func (k *Key) UnmarshalValue(data []byte) error {
	if len(data) != 32 {
		return ErrKeyWrongSize
	}

	copy(k[:], data)
	return nil
}

func (k *Key) ParentID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[:16])
	return id, err
}

func (k *Key) ObjectID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[16:])
	return id, err
}

// Helper to retrieve an object's key from its ID from the database.
func GetObjectKey(ctx context.Context, objectID ulid.ULID) (_ *Key, err error) {
	// Construct the Key to retrieve from the database.
	var key *Key
	if key, err = NewKey(ulid.ULID{}, objectID); err != nil {
		return nil, err
	}

	// Retrieve the Key from the database.
	if err = Get(ctx, key); err != nil {
		return nil, err
	}

	return key, nil
}

// Helper to store an object's key in the database.
func PutObjectKey(ctx context.Context, object Model) (err error) {
	var keyBytes []byte
	if keyBytes, err = object.Key(); err != nil {
		return err
	}

	// Construct the Key to store in the database.
	key := &Key{}
	if err = key.UnmarshalValue(keyBytes); err != nil {
		return err
	}

	return Put(ctx, key)
}
