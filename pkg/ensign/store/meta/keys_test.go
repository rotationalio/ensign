package meta_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
)

func TestObjectKeys(t *testing.T) {
	parentID := ulids.New()
	objectID := ulids.New()

	key, err := meta.CreateKey(parentID, objectID)
	require.NoError(t, err, "could not create object key")

	parsedParentID, err := key.ParentID()
	require.NoError(t, err, "could not parse parent ID")
	require.Equal(t, parentID, parsedParentID)

	parsedObjectID, err := key.ObjectID()
	require.NoError(t, err, "could not parse object ID")
	require.Equal(t, objectID, parsedObjectID)

	parsedKey := meta.ObjectKey{}
	err = parsedKey.UnmarshalValue(key[:])
	require.NoError(t, err, "could not unmarshal key")
	require.Equal(t, key, parsedKey)

	index := key.Key()
	require.Len(t, index, 16)
	require.Equal(t, [16]byte(objectID), [16]byte(index))

	// Should not be able to create a key with a null ulid
	_, err = meta.CreateKey(ulids.Null, objectID)
	require.ErrorIs(t, err, errors.ErrKeyNull)

	_, err = meta.CreateKey(parentID, ulids.Null)
	require.ErrorIs(t, err, errors.ErrKeyNull)

	_, err = meta.CreateKey(ulids.Null, ulids.Null)
	require.ErrorIs(t, err, errors.ErrKeyNull)

	// Cannot unmarshal incorrect data
	err = parsedKey.UnmarshalValue([]byte{42})
	require.ErrorIs(t, err, errors.ErrKeyWrongSize)
}
