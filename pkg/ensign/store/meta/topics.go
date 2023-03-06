package meta

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Implements iterator.TopicIterator to provide access to a list of topics.
type TopicIterator struct {
	ldbiter.Iterator
}

// Topic unmarshals the next topic in the iterator.
func (i *TopicIterator) Topic() (*api.Topic, error) {
	topic := &api.Topic{}
	if err := proto.Unmarshal(i.Value(), topic); err != nil {
		return nil, err
	}
	return topic, nil
}

// TODO: implement pagination in the store!
func (i *TopicIterator) NextPage(in *api.PageInfo) (*api.TopicsPage, error) {
	return nil, errors.ErrNotImplemented
}

// List all of the topics associated with the specified projectID. Returns a
// TopicIterator that can be used to retrieve each individual topic or to create a page
// of topics for paginated requests.
func (s *Store) ListTopics(projectID ulid.ULID) iterator.TopicIterator {
	// Iterate over all of the topics prefixed by the projectID
	slice := util.BytesPrefix(projectID.Bytes())
	iter := s.db.NewIterator(slice, nil)
	topics := &TopicIterator{iter}

	return topics
}

// Create a topic in the database; if the topic already exists or if the topic is not
// valid an error is returned. This method uses the keymu lock to avoid concurrency
// issues for multiple writers.
func (s *Store) CreateTopic(topic *api.Topic) (err error) {
	// Validate the partial topic
	if err = ValidateTopic(topic, true); err != nil {
		return err
	}

	// Create an ID for the topic if one is not set on it
	if len(topic.Id) == 0 {
		topic.Id = ulids.New().Bytes()
	}

	// Set the created and modified timestamps
	topic.Created = timestamppb.Now()
	topic.Modified = topic.Created

	// Marshal the topic and store it
	var data []byte
	if data, err = proto.Marshal(topic); err != nil {
		return errors.Wrap(err)
	}

	if err = s.Create(TopicKey(topic), data); err != nil {
		return err
	}
	return nil
}

// Retrieve a topic from the database.
func (s *Store) RetrieveTopic(topicID ulid.ULID) (topic *api.Topic, err error) {
	var data []byte
	if data, err = s.Retrieve(IndexKey(topicID)); err != nil {
		return nil, err
	}

	topic = &api.Topic{}
	if err = proto.Unmarshal(data, topic); err != nil {
		return nil, errors.Wrap(err)
	}
	return topic, nil
}

// Update a topic by putting the specified topic into the database. This method uses
// the keymu lock to avoid concurrency issues and returns an error if the specified
// topic does not exist or is not valid.
func (s *Store) UpdateTopic(topic *api.Topic) (err error) {
	// Validate the complete topic
	if err = ValidateTopic(topic, false); err != nil {
		return err
	}

	// Update the modified timestamp on the topic.
	topic.Modified = timestamppb.Now()

	// Marshal the topic and store it
	var data []byte
	if data, err = proto.Marshal(topic); err != nil {
		return errors.Wrap(err)
	}

	if err = s.Update(TopicKey(topic), data); err != nil {
		return err
	}
	return nil
}

// Delete a topic from the database. If the topic does not exist, no error is returned.
// This method uses the keymu lock to avoid concurrency issues and also cleans up any
// indices associated with the topic.
func (s *Store) DeleteTopic(topicID ulid.ULID) error {
	if err := s.Destroy(IndexKey(topicID)); err != nil {
		return err
	}
	return nil
}

// TopicKey is a 32 byte value that is the concatenated projectID followed by the topicID.
func TopicKey(topic *api.Topic) ObjectKey {
	var key [32]byte
	copy(key[0:16], topic.ProjectId)
	copy(key[16:], topic.Id)
	return ObjectKey(key)
}

// Validate a topic is ready for storage in the database. If partial is true, then the
// fields that may be set by this package (e.g. ID, created, modified) are not checked,
// otherwise the entire struct is validated.
func ValidateTopic(topic *api.Topic, partial bool) error {
	if len(topic.ProjectId) == 0 {
		return errors.ErrTopicMissingProjectId
	}

	if _, err := ulids.Parse(topic.ProjectId); err != nil {
		return errors.ErrTopicInvalidProjectId
	}

	if len(topic.Id) > 0 {
		if _, err := ulids.Parse(topic.Id); err != nil {
			return errors.ErrTopicInvalidId
		}
	}

	if topic.Name == "" {
		return errors.ErrTopicMissingName
	}

	if !partial {
		if len(topic.Id) == 0 {
			return errors.ErrTopicMissingId
		}

		if !IsValidTimestamp(topic.Created) {
			return errors.ErrTopicInvalidCreated
		}

		if !IsValidTimestamp(topic.Modified) {
			return errors.ErrTopicInvalidModified
		}
	}

	return nil
}

func IsValidTimestamp(s *timestamppb.Timestamp) bool {
	return s.IsValid() && !s.AsTime().IsZero()
}
