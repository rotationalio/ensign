package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	pb "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

const TopicNamespace = "topics"

type Topic struct {
	OrgID              ulid.ULID                `msgpack:"org_id"`
	ProjectID          ulid.ULID                `msgpack:"project_id"`
	ID                 ulid.ULID                `msgpack:"id"`
	Name               string                   `msgpack:"name"`
	State              pb.TopicTombstone_Status `msgpack:"state"`
	ConfirmDeleteToken string                   `msgpack:"confirm_delete_token"`
	Created            time.Time                `msgpack:"created"`
	Modified           time.Time                `msgpack:"modified"`
}

var _ Model = &Topic{}

// Key is a 32 composite key combining the project ID and the topic ID.
func (t *Topic) Key() (key []byte, err error) {
	// Create a 32 byte array so that the first 16 bytes hold the project ID
	// and the last 16 bytes hold the topic ID.
	key = make([]byte, 32)

	// Marshal the project ID to the first 16 bytes of the key.
	if err = t.ProjectID.MarshalBinaryTo(key[0:16]); err != nil {
		return nil, err
	}

	// Marshal the topic ID to the last 16 bytes of the key.
	if err = t.ID.MarshalBinaryTo(key[16:]); err != nil {
		return nil, err
	}

	return key, err
}

func (t *Topic) Namespace() string {
	return TopicNamespace
}

func (t *Topic) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(t)
}

func (t *Topic) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, t)
}

func (t *Topic) Validate() error {
	if ulids.IsZero(t.OrgID) {
		return ErrMissingOrgID
	}

	if ulids.IsZero(t.ProjectID) {
		return ErrMissingID
	}

	if t.Name == "" {
		return ErrMissingTopicName
	}

	if !alphaNum.MatchString(t.Name) {
		return ValidationError("topic")
	}

	return nil
}

// Convert the model to an API response.
func (t *Topic) ToAPI() *api.Topic {
	return &api.Topic{
		ID:       t.ID.String(),
		Name:     t.Name,
		State:    t.State.String(),
		Created:  TimeToString(t.Created),
		Modified: TimeToString(t.Modified),
	}
}

// CreateTopic adds a new topic to the database.
func CreateTopic(ctx context.Context, topic *Topic) (err error) {
	// Validate topic data.
	if err = topic.Validate(); err != nil {
		return err
	}

	if ulids.IsZero(topic.ID) {
		topic.ID = ulids.New()
	}

	topic.Created = time.Now()
	topic.Modified = topic.Created

	if err = Put(ctx, topic); err != nil {
		return err
	}

	return nil
}

// RetrieveTopic gets a topic from the database by a given ID.
func RetrieveTopic(ctx context.Context, id ulid.ULID) (topic *Topic, err error) {
	topic = &Topic{
		ID: id,
	}

	if err = Get(ctx, topic); err != nil {
		return nil, err
	}

	return topic, nil
}

// ListTopics retrieves all topics assigned to a project.
func ListTopics(ctx context.Context, projectID ulid.ULID) (topics []*Topic, err error) {
	// Store the project ID as the prefix.
	var prefix []byte
	if projectID.Compare(ulid.ULID{}) != 0 {
		prefix = projectID[:]
	}

	var values [][]byte
	if values, err = List(ctx, prefix, TopicNamespace); err != nil {
		return nil, err
	}

	// Parse the topics from the data
	topics = make([]*Topic, 0, len(values))
	for _, data := range values {
		topic := &Topic{}
		if err = topic.UnmarshalValue(data); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

// UpdateTopic updates the record of a topic by a given ID.
func UpdateTopic(ctx context.Context, topic *Topic) (err error) {
	// Validate topic data.
	if err = topic.Validate(); err != nil {
		return err
	}

	if ulids.IsZero(topic.ID) {
		return ErrMissingID
	}

	topic.Modified = time.Now()
	if topic.Created.IsZero() {
		topic.Created = topic.Modified
	}

	if err = Put(ctx, topic); err != nil {
		return err
	}

	return nil
}

// DeleteTopic deletes a topic by a given ID.
func DeleteTopic(ctx context.Context, id ulid.ULID) (err error) {
	topic := &Topic{
		ID: id,
	}

	if err = Delete(ctx, topic); err != nil {
		return err
	}

	return nil
}
