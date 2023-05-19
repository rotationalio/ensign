package metatopic_test

import (
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/metatopic"
	"github.com/stretchr/testify/require"
)

func TestTopicUpdateSerialization(t *testing.T) {
	orgID := ulid.Make()
	projectID := ulid.Make()
	topicID := ulid.Make()

	tut := &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateCreated,
		OrgID:      orgID,
		ProjectID:  projectID,
		TopicID:    topicID,
		ClientID:   "testing123",
		Topic: &metatopic.Topic{
			ID:        topicID.Bytes(),
			ProjectID: projectID.Bytes(),
			Name:      "testing",
			ReadOnly:  true,
			Offset:    1331042,
			Shards:    1,
			Storage:   16.00007915496826,
			Publishers: &metatopic.Activity{
				Active:   7,
				Inactive: 28,
			},
			Subscribers: &metatopic.Activity{
				Active:   23,
				Inactive: 4,
			},
			Created:  time.Now().Truncate(1 * time.Microsecond),
			Modified: time.Now().Truncate(1 * time.Microsecond),
		},
	}

	// Marshal the topic update
	data, err := tut.Marshal()
	require.NoError(t, err, "could not marshal topic update")
	require.NotEmpty(t, data, "expected marshaled topic update data back")

	// Unmarshal the topic update
	cmp := &metatopic.TopicUpdate{}
	err = cmp.Unmarshal(data)
	require.NoError(t, err, "could not unmarshal topic update")

	require.Equal(t, tut, cmp, "expected marshaled and unmarshaled structs to match")

}

func TestTopicUpdateType(t *testing.T) {
	testCases := []struct {
		tut      metatopic.TopicUpdateType
		expected string
	}{
		{metatopic.TopicUpdateUnknown, "unknown"},
		{metatopic.TopicUpdateCreated, "created"},
		{metatopic.TopicUpdateModified, "modified"},
		{metatopic.TopicUpdateStateChange, "state_change"},
		{metatopic.TopicUpdateDeleted, "deleted"},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, tc.tut.String())
	}
}

func TestActivity(t *testing.T) {
	things := &metatopic.Activity{
		Active:   227,
		Inactive: 773,
	}

	require.Equal(t, uint64(1000), things.Total())
	require.Equal(t, 1.0, things.PercentActive()+things.PercentInactive())
}

func TestParseVersion(t *testing.T) {
	require.Equal(t, "1.0.0", metatopic.SchemaVersion, "version has changed test needs to be updated")
	major, minor, patch := metatopic.ParseVersion()
	require.Equal(t, uint32(1), major)
	require.Equal(t, uint32(0), minor)
	require.Equal(t, uint32(0), patch)
}
