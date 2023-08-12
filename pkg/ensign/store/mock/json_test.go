package mock_test

import (
	"os"
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalEventList(t *testing.T) {
	data, err := os.ReadFile("testdata/events.pb.json")
	require.NoError(t, err)

	events, err := mock.UnmarshalEventList(data)
	require.NoError(t, err)
	require.Len(t, events, 0)
}

func TestUnmarshalTopicList(t *testing.T) {
	data, err := os.ReadFile("testdata/topics.pb.json")
	require.NoError(t, err)

	topics, err := mock.UnmarshalTopicList(data)
	require.NoError(t, err)
	require.Len(t, topics, 7)
}

func TestUnmarshalTopicNamesList(t *testing.T) {
	data, err := os.ReadFile("testdata/topicnames.pb.json")
	require.NoError(t, err)

	names, err := mock.UnmarshalTopicNamesList(data)
	require.NoError(t, err)
	require.Len(t, names, 7)
}
