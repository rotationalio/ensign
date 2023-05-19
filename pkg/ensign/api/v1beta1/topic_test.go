package api_test

import (
	"bytes"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestTopicNameHash(t *testing.T) {
	testCases := []string{
		"",
		"testing",
		"testing.topic.foo",
		"_testing1234",
		"areallylongtopicnamethatismuchlongerthan16bytes",
	}

	for _, tc := range testCases {
		// Create the topic and the name hashes
		topic := &api.Topic{Name: tc}
		topicHash := topic.NameHash()
		nameHash := api.TopicNameHash(tc)

		require.Len(t, topicHash, 16, "expect the topic hash to be 16 bytes for key storage in leveldb")
		require.True(t, bytes.Equal(topicHash, nameHash), "expected topic hashing and name hashing to be identical")
	}
}
