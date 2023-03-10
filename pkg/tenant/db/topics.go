package db

import (
	"context"
	"regexp"
	"time"

	"github.com/oklog/ulid/v2"
	pb "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/vmihailenco/msgpack/v5"
)

const TopicNamespace = "topics"

// Topic names must be URL safe and begin with a letter.
var TopicNameRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9.-_]*$")

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
	// ProjectID and TopicID are required
	if ulids.IsZero(t.ProjectID) {
		return nil, ErrMissingProjectID
	}

	if ulids.IsZero(t.ID) {
		return nil, ErrMissingID
	}

	var k Key
	if k, err = CreateKey(t.ProjectID, t.ID); err != nil {
		return nil, err
	}

	return k.MarshalValue()
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
	if err = PutObjectKey(ctx, topic); err != nil {
		return err
	}
	return nil
}

// RetrieveTopic gets a topic from the database by the given project ID and topic ID.
func RetrieveTopic(ctx context.Context, topicID ulid.ULID) (topic *Topic, err error) {
	// Lookup the topic key in the database
	var key []byte
	if key, err = GetObjectKey(ctx, topicID); err != nil {
		return nil, err
	}

	// Use the key to lookup the topic
	var data []byte
	if data, err = getRequest(ctx, TopicNamespace, key); err != nil {
		return nil, err
	}

	// Unmarshal the data into the topic
	topic = &Topic{}
	if err = topic.UnmarshalValue(data); err != nil {
		return nil, err
	}

	return topic, nil
}

// ListTopics retrieves all topics assigned to a project.
func ListTopics(ctx context.Context, projectID, topicID ulid.ULID, prev *pg.Cursor) (topics []*Topic, next *pg.Cursor, err error) {
	// Store the project ID as the prefix.
	var prefix []byte
	if projectID.Compare(ulid.ULID{}) != 0 {
		prefix = projectID[:]
	}

	var seekKey []byte
	if topicID.Compare(ulid.ULID{}) != 0 {
		seekKey = topicID[:]
	}

	// Check to see if a default cursor exists and create one if it does not.
	if prev == nil {
		prev = pg.New("", "", 0)
	}

	if prev.PageSize <= 0 {
		return nil, nil, ErrMissingPageSize
	}

	topics = make([]*Topic, 0)
	onListItem := func(item *trtl.KVPair) error {
		topic := &Topic{}
		if err = topic.UnmarshalValue(item.Value); err != nil {
			return err
		}
		topics = append(topics, topic)
		return nil
	}

	if next, err = List(ctx, prefix, seekKey, TopicNamespace, onListItem, prev); err != nil {
		return nil, nil, err
	}

	return topics, next, nil
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

	// Retrieve the topic key to update the topic.
	// Note: There is a possible concurrency issue here if the topic is deleted between
	// Get and Put.
	var key []byte
	if key, err = GetObjectKey(ctx, topic.ID); err != nil {
		return err
	}

	topic.Modified = time.Now()
	if topic.Created.IsZero() {
		topic.Created = topic.Modified
	}

	var data []byte
	if data, err = topic.MarshalValue(); err != nil {
		return err
	}

	if err = putRequest(ctx, TopicNamespace, key, data); err != nil {
		return err
	}

	return nil
}

// DeleteTopic deletes a topic by the given project ID and topic ID.
func DeleteTopic(ctx context.Context, topicID ulid.ULID) (err error) {
	// Retrieve the topic key to delete the topic.
	var key []byte
	if key, err = GetObjectKey(ctx, topicID); err != nil {
		return err
	}

	// Delete the project and its key from the database.
	if err = deleteRequest(ctx, TopicNamespace, key); err != nil {
		return err
	}

	if err = DeleteObjectKey(ctx, key); err != nil {
		return err
	}
	return nil
}
