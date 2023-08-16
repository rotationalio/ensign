package api_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
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
		uid, err := tc.topic.ParseTopicID()
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

func TestTopicInfoParseTopicID(t *testing.T) {
	testCases := []struct {
		topicID  []byte
		expected ulid.ULID
		err      error
	}{
		{nil, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte{}, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte("foo"), ulid.ULID{}, ulid.ErrDataSize},
		{ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ").Bytes(), ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ"), nil},
	}

	for i, tc := range testCases {
		info := &api.TopicInfo{TopicId: tc.topicID}
		topicID, err := info.ParseTopicID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, ulid.ULID{}, topicID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, topicID, "test case %d failed", i)
		}
	}
}

func TestTopicInfoParseProjectID(t *testing.T) {
	testCases := []struct {
		projectID []byte
		expected  ulid.ULID
		err       error
	}{
		{nil, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte{}, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte("foo"), ulid.ULID{}, ulid.ErrDataSize},
		{ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ").Bytes(), ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ"), nil},
	}

	for i, tc := range testCases {
		info := &api.TopicInfo{ProjectId: tc.projectID}
		projectID, err := info.ParseProjectID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, ulid.ULID{}, projectID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, projectID, "test case %d failed", i)
		}
	}
}

func TestTopicInfoParseEventOffsetID(t *testing.T) {
	testCases := []struct {
		eventID  []byte
		expected rlid.RLID
		err      error
	}{
		{nil, rlid.RLID{}, rlid.ErrDataSize},
		{[]byte{}, rlid.RLID{}, rlid.ErrDataSize},
		{[]byte("foo"), rlid.RLID{}, rlid.ErrDataSize},
		{rlid.MustParse("064zwbj8vg00000n").Bytes(), rlid.MustParse("064zwbj8vg00000n"), nil},
	}

	for i, tc := range testCases {
		info := &api.TopicInfo{EventOffsetId: tc.eventID}
		eventID, err := info.ParseEventOffsetID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, rlid.RLID{}, eventID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, eventID, "test case %d failed", i)
		}
	}
}

func TestFindEventTypeInfo(t *testing.T) {
	info := &api.TopicInfo{}
	require.Len(t, info.Types, 0)

	etype := info.FindEventTypeInfo(api.UnspecifiedType)
	require.Len(t, info.Types, 1)
	etype.Events = 42
	etype.Duplicates = 1
	etype.DataSizeBytes = 4096

	compat := info.FindEventTypeInfo(api.UnspecifiedType)
	require.Len(t, info.Types, 1)
	require.Equal(t, uint64(42), compat.Events)
	require.Equal(t, uint64(1), compat.Duplicates)
	require.Equal(t, uint64(4096), compat.DataSizeBytes)

	// Manually add a type info to the info
	info.Types = append(info.Types, &api.EventTypeInfo{Type: &api.Type{Name: "TestType", MajorVersion: 1}, Events: 100, DataSizeBytes: 1.24e6})
	require.Len(t, info.Types, 2)

	etype = info.FindEventTypeInfo(&api.Type{Name: "TestType", MajorVersion: 1})
	require.Len(t, info.Types, 2)
	require.Equal(t, uint64(100), etype.Events)
	require.Equal(t, uint64(0), etype.Duplicates)
	require.Equal(t, uint64(1.24e6), etype.DataSizeBytes)
}
