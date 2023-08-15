package meta_test

import (
	"bytes"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestObjectKeys(t *testing.T) {
	makeSegmentTest := func(segment meta.Segment) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()

			parentID := ulids.New()
			objectID := ulids.New()

			key, err := meta.CreateKey(parentID, objectID, segment)
			require.NoError(t, err, "could not create object key")

			parsedParentID, err := key.ParentID()
			require.NoError(t, err, "could not parse parent ID")
			require.Equal(t, parentID, parsedParentID)

			parsedObjectID, err := key.ObjectID()
			require.NoError(t, err, "could not parse object ID")
			require.Equal(t, objectID, parsedObjectID)

			parsedSegment, err := key.Segment()
			require.NoError(t, err, "could not parse segment")
			require.Equal(t, segment, parsedSegment)

			parsedKey := meta.ObjectKey{}
			err = parsedKey.UnmarshalValue(key[:])
			require.NoError(t, err, "could not unmarshal key")
			require.Equal(t, key, parsedKey)

			index := key.Key()
			require.Len(t, index, 16)
			require.Equal(t, [16]byte(objectID), [16]byte(index))

			// Should not be able to create a key with a null ulid
			_, err = meta.CreateKey(ulids.Null, objectID, segment)
			require.ErrorIs(t, err, errors.ErrKeyNull)

			_, err = meta.CreateKey(parentID, ulids.Null, segment)
			require.ErrorIs(t, err, errors.ErrKeyNull)

			_, err = meta.CreateKey(ulids.Null, ulids.Null, segment)
			require.ErrorIs(t, err, errors.ErrKeyNull)

			// Cannot unmarshal incorrect data
			err = parsedKey.UnmarshalValue([]byte{42})
			require.ErrorIs(t, err, errors.ErrKeyWrongSize)
		}
	}

	t.Run("TopicSegment", makeSegmentTest(meta.TopicSegment))
	t.Run("TopicNamesSegment", makeSegmentTest(meta.TopicNamesSegment))
	t.Run("TopicInfoSegment", makeSegmentTest(meta.TopicInfoSegment))
	t.Run("GroupSegment", makeSegmentTest(meta.GroupSegment))
}

func TestObjectKeyConvert(t *testing.T) {
	segments := []meta.Segment{
		meta.TopicSegment,
		meta.TopicNamesSegment,
		meta.TopicInfoSegment,
		meta.GroupSegment,
	}

	parentID := ulid.MustParse("01H7V2HDHM6QH6CZ0KATPSQMF1")
	objectID := ulid.MustParse("01H7V2HMSR47TQVSFCNTD4D5EE")

	for _, orig := range segments {
		for _, convert := range segments {
			key, err := meta.CreateKey(parentID, objectID, orig)
			require.NoError(t, err)

			actual, err := key.Segment()
			require.NoError(t, err)
			require.Equal(t, orig, actual)

			key.Convert(convert)

			actual, err = key.Segment()
			require.NoError(t, err)
			require.Equal(t, convert, actual)

		}
	}

}

func TestIndexKey(t *testing.T) {
	_, err := meta.CreateIndex(ulid.ULID{})
	require.ErrorIs(t, err, errors.ErrKeyNull, "cannot create an index for zero-valued key")

	objectID := ulid.MustParse("01H7V28RZBEBZW0W04Q7DBYXJY")
	key, err := meta.CreateIndex(objectID)
	require.NoError(t, err, "could not create index key")
	require.True(t, bytes.Equal(objectID[:], key[:]))
}

func TestSegment(t *testing.T) {
	// Test Segment IDs
	require.Equal(t, []byte("tp"), meta.TopicSegment[:])
	require.Equal(t, []byte("Tn"), meta.TopicNamesSegment[:])
	require.Equal(t, []byte("Ti"), meta.TopicInfoSegment[:])
	require.Equal(t, []byte("GP"), meta.GroupSegment[:])

	// Test Strings
	require.Equal(t, "topic", meta.TopicSegment.String())
	require.Equal(t, "topic_name", meta.TopicNamesSegment.String())
	require.Equal(t, "topic_info", meta.TopicInfoSegment.String())
	require.Equal(t, "group", meta.GroupSegment.String())
	require.Equal(t, "unknown", meta.Segment([2]byte{0x00, 0x42}).String())
}
