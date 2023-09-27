package db

import (
	"bytes"
	"context"

	"github.com/oklog/ulid/v2"
)

// Keys Namespace maps object IDs to their fully qualified database keys.
const KeysNamespace = "object_keys"

// Key is composed of two concatenated IDs. The first 16 bytes are the ID of parent and
// the second 16 bytes are the ID of the object.
type Key [32]byte

var NullID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var _ Model = &Key{}

// CreateKey creates a new key from a parent ID and object ID so that callers can
// lookup the object ID from its namespace.
func CreateKey(parentID, objectID ulid.ULID) (key Key, err error) {
	if err = parentID.MarshalBinaryTo(key[:16]); err != nil {
		return Key{}, err
	}

	if err = objectID.MarshalBinaryTo(key[16:]); err != nil {
		return Key{}, err
	}
	return key, nil
}

// Keys are stored by object ID. Since object IDs are locked monotonically increasing
// ulids they are guaranteed to be unique.
func (k Key) Key() ([]byte, error) {
	if bytes.Equal(k[16:], NullID) {
		return nil, ErrKeyNoID
	}

	return k[16:], nil
}

func (k Key) Namespace() string {
	return KeysNamespace
}

func (k Key) MarshalValue() ([]byte, error) {
	return k[:], nil
}

func (k *Key) UnmarshalValue(data []byte) error {
	if len(data) != 32 {
		return ErrKeyWrongSize
	}

	copy(k[:], data)
	return nil
}

func (k Key) ParentID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[:16])
	return id, err
}

func (k Key) ObjectID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[16:])
	return id, err
}

// String returns a string representation of the key with the parent ID concatenated
// with the object ID.
func (k Key) String() (_ string, err error) {
	if len(k) != 32 {
		return "", ErrKeyWrongSize
	}

	var parentID, objectID ulid.ULID
	if parentID, err = k.ParentID(); err != nil {
		return "", err
	}

	if objectID, err = k.ObjectID(); err != nil {
		return "", err
	}

	return parentID.String() + objectID.String(), nil
}

// ParseKey parses a string representation of a key into a Key struct.
func ParseKey(s string) (key Key, err error) {
	if len(s) != 52 {
		return Key{}, ErrKeyWrongSize
	}

	var parentID, objectID ulid.ULID
	if parentID, err = ulid.Parse(s[:26]); err != nil {
		return Key{}, err
	}

	if objectID, err = ulid.Parse(s[26:]); err != nil {
		return Key{}, err
	}

	return CreateKey(parentID, objectID)
}

// Helper to retrieve an object's key from its ID from the database.
func GetObjectKey(ctx context.Context, objectID ulid.ULID) (key []byte, err error) {
	return getRequest(ctx, KeysNamespace, objectID[:])
}

// Helper to store an object's key in the database.
func PutObjectKey(ctx context.Context, object Model) (err error) {
	var keyData []byte
	if keyData, err = object.Key(); err != nil {
		return err
	}

	key := &Key{}
	if err = key.UnmarshalValue(keyData); err != nil {
		return err
	}

	return Put(ctx, key)
}

// Helper to delete an object's key from the database.
func DeleteObjectKey(ctx context.Context, key []byte) (err error) {
	k := &Key{}
	if err = k.UnmarshalValue(key); err != nil {
		return err
	}

	return Delete(ctx, k)
}
