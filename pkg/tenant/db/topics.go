package db

import (
	"context"
	"strings"
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

func (t *Topic) Validate() error {
	if t.ProjectID.Compare(ulid.ULID{}) == 0 {
		return ErrMissingID
	}

	topicName := t.Name

	if topicName == "" {
		return ErrMissingTopicName
	}

	if strings.ContainsAny(string(topicName[0]), "0123456789") {
		return ErrNumberFirstCharacter
	}

	if strings.ContainsAny(topicName, " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		return ErrSpecialCharacters
	}
	return nil
}

// CreateTopic adds a new topic to the database.
func CreateTopic(ctx context.Context, topic *Topic) (err error) {
	// Validate topic data.
	if err = topic.Validate(); err != nil {
		return err
	}

	// TODO: Use crypto rand and monotonic entropy with ulid.New
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

		// Marshal and unmarshal the data with msgPack.
		topic.MarshalData()
		topic.UnmarshalData(data)
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

	// TODO: Use crypto rand and monotonic entropy with ulid.New
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

// Marshals data with msgPack.
func (p *Topic) MarshalData() (data []byte, err error) {
	if data, err = p.MarshalValue(); err != nil {
		return nil, err
	}
	return data, nil
}

// Unmarshals data with msgPack.
func (p *Topic) UnmarshalData(data []byte) (err error) {
	if err := p.UnmarshalValue(data); err != nil {
		return err
	}
	return nil
}
