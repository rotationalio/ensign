package api_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestTopicULID(t *testing.T) {
	testCases := []struct {
		topic    *api.Topic
		expected ulid.ULID
		err      error
	}{
		{
			&api.Topic{Id: nil}, ulid.ULID{}, errors.New("ulid: bad data size when unmarshaling"),
		},
		{
			&api.Topic{Id: make([]byte, 0)}, ulid.ULID{}, errors.New("ulid: bad data size when unmarshaling"),
		},
		{
			&api.Topic{Id: make([]byte, 8)}, ulid.ULID{}, errors.New("ulid: bad data size when unmarshaling"),
		},
		{
			&api.Topic{Id: make([]byte, 16)}, ulid.ULID{}, nil,
		},
		{
			&api.Topic{Id: ulid.MustParse("01H0TB6F9MAMMQ8T9DZAZGQ5RH").Bytes()}, ulid.MustParse("01H0TB6F9MAMMQ8T9DZAZGQ5RH"), nil,
		},
	}

	for _, tc := range testCases {
		uid, err := tc.topic.ULID()
		if tc.err != nil {
			require.EqualError(t, err, tc.err.Error())
		} else {
			require.Equal(t, 0, tc.expected.Compare(uid))
		}
	}
}

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

		require.Len(t, topicHash, api.NameHashLength, "expect the topic hash to be 16 bytes for key storage in leveldb")
		require.True(t, bytes.Equal(topicHash, nameHash), "expected topic hashing and name hashing to be identical")
	}
}

func TestValidateTopicName(t *testing.T) {
	valid := []string{
		"topic",
		"snake_case_topic",
		"CamelCaseTopic",
		"dot.separated.topic",
		"dash-separated-topic",
		"topic34",
		"blue_42_blue_42",
		"HutHutHIKE42",
		"dot.dash-underscore_topic",
	}

	invalid := []string{
		"",
		strings.Repeat("ABRACADBRA", 200),
		"no spaces",
		"42noleadnums",
		"-noleadhypen",
		".noleaddot",
		"no$special!chars",
	}

	for _, name := range valid {
		require.True(t, api.ValidTopicName(name), "expected %q to be a valid topic name", name)
	}

	for _, name := range invalid {
		require.False(t, api.ValidTopicName(name), "expected %q to be invalid", name)
	}
}
