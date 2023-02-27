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

const (
	TopicNamespace     = "topics"
	TopicKeysNamespace = "topic_keys"
)

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

// TopicKey stores the components of the topic key to enable direct lookup from the
// topic ID.
type TopicKey struct {
	ProjectID ulid.ULID `msgpack:"project_id"`
	ID        ulid.ULID `msgpack:"id"`
}

var _ Model = &Topic{}
var _ Model = &TopicKey{}

// Key is a 32 composite key combining the project ID and the topic ID.
func (t *Topic) Key() (key []byte, err error) {
	// ProjectID and TopicID are required
	if ulids.IsZero(t.ProjectID) {
		return nil, ErrMissingProjectID
	}

	if ulids.IsZero(t.ID) {
		return nil, ErrMissingID
	}

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
		return ErrMissingProjectID
	}

	if t.Name == "" {
		return ErrMissingTopicName
	}

	if !alphaNum.MatchString(t.Name) {
		return ErrInvalidTopicName
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

func (t *TopicKey) Key() (key []byte, err error) {
	if ulids.IsZero(t.ID) {
		return nil, ErrMissingID
	}

	return t.ID[:], nil
}

func (t *TopicKey) Namespace() string {
	return TopicKeysNamespace
}

func (t *TopicKey) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(t)
}

func (t *TopicKey) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, t)
}

// CreateTopic adds a new topic to the database.
func CreateTopic(ctx context.Context, topic *Topic) (err error) {
	if ulids.IsZero(topic.ID) {
		topic.ID = ulids.New()
	}

	// Validate topic data.
	if err = topic.Validate(); err != nil {
		return err
	}

	topic.Created = time.Now()
	topic.Modified = topic.Created

	if err = Put(ctx, topic); err != nil {
		return err
	}

	// Store the topic key in the database to allow direct lookups by topic ID.
	topicKey := &TopicKey{
		ProjectID: topic.ProjectID,
		ID:        topic.ID,
	}
	if err = Put(ctx, topicKey); err != nil {
		return err
	}
	return nil
}

// RetrieveTopic gets a topic from the database by the given project ID and topic ID.
func RetrieveTopic(ctx context.Context, topicID ulid.ULID) (topic *Topic, err error) {
	// Lookup the topic key in the database
	key := &TopicKey{
		ID: topicID,
	}

	if err = Get(ctx, key); err != nil {
		return nil, err
	}

	// Use the key to lookup the topic
	topic = &Topic{
		ProjectID: key.ProjectID,
		ID:        key.ID,
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
	if ulids.IsZero(topic.ID) {
		return ErrMissingID
	}

	// Validate topic data.
	if err = topic.Validate(); err != nil {
		return err
	}

	// Retrieve the topic key to update the project.
	// Note: There is a possible concurrency issue here if the topic is deleted between
	// Get and Put.
	key := &TopicKey{
		ID: topic.ID,
	}
	if err = Get(ctx, key); err != nil {
		return err
	}

	topic.Modified = time.Now()
	if topic.Created.IsZero() {
		topic.Created = topic.Modified
	}

	topic.ProjectID = key.ProjectID
	if err = Put(ctx, topic); err != nil {
		return err
	}

	return nil
}

// DeleteTopic deletes a topic by the given project ID and topic ID.
func DeleteTopic(ctx context.Context, topicID ulid.ULID) (err error) {
	topic := &Topic{
		ID: topicID,
	}

	// Retrieve the topic key to delete the project.
	key := &TopicKey{
		ID: topicID,
	}
	if err = Get(ctx, key); err != nil {
		return err
	}

	// Delete the project and its key from the database.
	topic.ProjectID = key.ProjectID
	if err = Delete(ctx, topic); err != nil {
		return err
	}

	if err = Delete(ctx, key); err != nil {
		return err
	}
	return nil
}
