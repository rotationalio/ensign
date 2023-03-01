package db_test

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
)

func TestNullKeys(t *testing.T) {
	require.Len(t, db.NullID, 16, "wrong null ID length")
	for _, b := range db.NullID {
		require.Zero(t, b, "null ID should be zero valued")
	}
}

func TestKey(t *testing.T) {
	parentID := ulid.MustParse("01F1ZQZJXQZJXQZJXQZJXQZJXQ")
	objectID := ulid.MustParse("02ABCQZJXQZJXQZJXQZJXQZJXD")

	// Should be able to create a key from the parent ID and object ID.
	key, err := db.CreateKey(parentID, objectID)
	require.NoError(t, err, "expected no error when creating key")

	// Should be able to marshal the key.
	data, err := key.MarshalValue()
	require.NoError(t, err, "expected no error when marshaling key")

	// ID of the key is the object ID.
	id, err := key.Key()
	require.NoError(t, err, "expected no error when reading key ID")
	require.Equal(t, objectID[:], id, "expected key ID to be the object ID")

	// Should be able to unmarshal the key.
	k := &db.Key{}
	require.NoError(t, k.UnmarshalValue(data), "expected no error when unmarshaling key")
	require.Equal(t, key, *k, "expected key to be equal")

	// Should fail to unmarshal wrong length data.
	require.ErrorIs(t, k.UnmarshalValue([]byte{0, 0, 0, 0}), db.ErrKeyWrongSize, "expected error when unmarshaling wrong length data")

	// Should be able to read the parent ID and object ID.
	parent, err := key.ParentID()
	require.NoError(t, err, "expected no error when reading parent ID")
	require.Equal(t, parentID, parent, "expected parent ID to be equal")

	object, err := key.ObjectID()
	require.NoError(t, err, "expected no error when reading object ID")
	require.Equal(t, objectID, object, "expected object ID to be equal")

	// Empty key should not have a key
	empty := &db.Key{}
	_, err = empty.Key()
	require.ErrorIs(t, err, db.ErrKeyNoID, "expected error when ID is empty")
}
