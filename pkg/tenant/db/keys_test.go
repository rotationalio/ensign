package db_test

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	parentID := ulid.MustParse("01F1ZQZJXQZJXQZJXQZJXQZJXQ")
	objectID := ulid.MustParse("02ABCQZJXQZJXQZJXQZJXQZJXD")

	// Should be able to create a key with no parent ID
	key, err := db.NewKey(ulid.ULID{}, objectID)
	require.NoError(t, err, "expected no error when creating key")

	// Should be able to marshal and unmarshal the key.
	key, err = db.NewKey(parentID, objectID)
	require.NoError(t, err, "expected no error when creating key")
	data, err := key.MarshalValue()
	require.NoError(t, err, "expected no error when marshaling key")

	// ID of the key is the object ID.
	id, err := key.Key()
	require.NoError(t, err, "expected no error when reading key ID")
	require.Equal(t, objectID[:], id, "expected key ID to be the object ID")

	k := &db.Key{}
	require.NoError(t, k.UnmarshalValue(data), "expected no error when unmarshaling key")
	require.Equal(t, key, k, "expected key to be equal")

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
	key = &db.Key{}
	_, err = key.Key()
	require.ErrorIs(t, err, db.ErrKeyNoID, "expected error when ID is empty")
}
