package meta_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/stretchr/testify/require"
)

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
