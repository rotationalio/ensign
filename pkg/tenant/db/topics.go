package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const TopicNamespace = "topics"

type Topic struct {
	ProjectID ulid.ULID `msgpack:"project_id"`
	ID        ulid.ULID `msgpack:"id"`
	Name      string    `msgpack:"name"`
	Created   time.Time `msgpack:"created"`
	Modified  time.Time `msgpack:"modified"`
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

// CreateTopic adds a new topic to the database.
func CreateTopic(ctx context.Context, topic *Topic) (err error) {
	if topic.ID.Compare(ulid.ULID{}) == 0 {
		topic.ID = ulid.Make()
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

	// Parse the topics from the data.
	topics = make([]*Topic, 0, len(values))
	for _, data := range values {
		topic := &Topic{}
		if data, err = topic.MarshalValue(); err != nil {
			return nil, err
		}
		if err = topic.UnmarshalValue(data); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

// UpdateTopic updates the record of a topic by a given ID.
func UpdateTopic(ctx context.Context, topic *Topic) (err error) {
	if topic.ID.Compare(ulid.ULID{}) == 0 {
		return ErrMissingID
	}

	topic.Modified = time.Now()

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
