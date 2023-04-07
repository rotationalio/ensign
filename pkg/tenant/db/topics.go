package db

import (
	"context"
	"regexp"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	pb "github.com/rotationalio/go-ensign/api/v1beta1"
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

	// Store the topic ID as a key and the topic org ID as a value in the database for org verification.
	if err = PutOrgIndex(ctx, topic.ID, topic.OrgID); err != nil {
		return err
	}

	// Store the topic key in the database to allow direct lookups by topic ID.
	if err = PutObjectKey(ctx, topic); err != nil {
		return err
	}
	return nil
}

// RetrieveTopic gets a topic from the database by the given topic ID.
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

// ListTopics retrieves a paginated list of topics.
func ListTopics(ctx context.Context, projectID ulid.ULID, c *pg.Cursor) (topics []*Topic, cursor *pg.Cursor, err error) {
	// Store the project ID as the prefix.
	var prefix []byte
	if projectID.Compare(ulid.ULID{}) != 0 {
		prefix = projectID[:]
	}

	var seekKey []byte
	if c.EndIndex != "" {
		var start ulid.ULID
		if start, err = ulid.Parse(c.EndIndex); err != nil {
			return nil, nil, err
		}
		seekKey = start[:]
	}

	// Check to see if a default cursor exists and create one if it does not.
	if c == nil {
		c = pg.New("", "", 0)
	}

	if c.PageSize <= 0 {
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

	if cursor, err = List(ctx, prefix, seekKey, TopicNamespace, onListItem, c); err != nil {
		return nil, nil, err
	}

	return topics, cursor, nil
}

// UpdateTopic updates the record of a topic from its database model.
func UpdateTopic(ctx context.Context, topic *Topic) (err error) {
	if ulids.IsZero(topic.ID) {
		return ErrMissingID
	}

	// Validate topic data.
	if err = topic.Validate(); err != nil {
		return err
	}

	topic.Modified = time.Now()
	if topic.Created.IsZero() {
		topic.Created = topic.Modified
	}

	return Put(ctx, topic)
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
