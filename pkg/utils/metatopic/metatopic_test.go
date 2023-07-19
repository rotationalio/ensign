package metatopic_test

import (
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/metatopic"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	update := &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateCreated,
		ProjectID:  ulids.New(),
		TopicID:    ulids.New(),
	}

	// OrgID, ProjectID, and TopicID are required for all updates
	require.ErrorIs(t, update.Validate(), metatopic.ErrMissingOrgID, "expected error for missing org ID")

	update.OrgID = ulids.New()
	update.ProjectID = ulids.Null
	require.ErrorIs(t, update.Validate(), metatopic.ErrMissingProjectID, "expected error for missing project ID")

	update.ProjectID = ulids.New()
	update.TopicID = ulids.Null
	require.ErrorIs(t, update.Validate(), metatopic.ErrMissingTopicID, "expected error for missing topic ID")

	// Topic is required for created and modified updates
	for _, updateType := range []metatopic.TopicUpdateType{metatopic.TopicUpdateCreated, metatopic.TopicUpdateModified} {
		update.UpdateType = updateType
		update.TopicID = ulids.New()
		update.Topic = nil

		require.ErrorIs(t, update.Validate(), metatopic.ErrMissingTopic, "expected error for missing topic for update type %s", updateType)

		update.Topic = &metatopic.Topic{
			Publishers:  &metatopic.Activity{},
			Subscribers: &metatopic.Activity{},
			Created:     time.Now(),
			Modified:    time.Now(),
		}

		require.ErrorIs(t, update.Validate(), metatopic.ErrMissingName, "expected error for missing topic name for update type %s", updateType)

		update.Topic.Name = "testing"
		update.Topic.Events = -1
		require.ErrorIs(t, update.Validate(), metatopic.ErrInvalidEvents, "expected error for invalid events for update type %s", updateType)

		update.Topic.Events = 0
		update.Topic.Storage = -1.0
		require.ErrorIs(t, update.Validate(), metatopic.ErrInvalidStorage, "expected error for invalid storage for update type %s", updateType)

		update.Topic.Storage = 0.0
		update.Topic.Publishers = nil
		require.ErrorIs(t, update.Validate(), metatopic.ErrMissingPublishers, "expected error for missing publishers for update type %s", updateType)

		update.Topic.Publishers = &metatopic.Activity{}
		update.Topic.Subscribers = nil
		require.ErrorIs(t, update.Validate(), metatopic.ErrMissingSubscribers, "expected error for missing subscribers for update type %s", updateType)

		update.Topic.Subscribers = &metatopic.Activity{}
		update.Topic.Created = time.Time{}
		require.ErrorIs(t, update.Validate(), metatopic.ErrMissingCreated, "expected error for missing created time for update type %s", updateType)

		update.Topic.Created = time.Now()
		update.Topic.Modified = time.Time{}
		require.ErrorIs(t, update.Validate(), metatopic.ErrMissingModified, "expected error for missing modified time for update type %s", updateType)
	}

	// Unknown update types are not allowed
	update.UpdateType = metatopic.TopicUpdateUnknown
	require.ErrorIs(t, update.Validate(), metatopic.ErrUnknownUpdateType, "expected error for unknown update type")

	// State change and deleted updates do not require a topic
	for _, updateType := range []metatopic.TopicUpdateType{metatopic.TopicUpdateStateChange, metatopic.TopicUpdateDeleted} {
		update.UpdateType = updateType
		update.Topic = nil
		require.NoError(t, update.Validate(), "expected no error for update type %s", updateType)
	}
}

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
			Events:    1000000,
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
