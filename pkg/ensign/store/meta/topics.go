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

type TopicIterator struct {
	ldbiter.Iterator
}

func (i *TopicIterator) Topic() (*api.Topic, error) {
	topic := &api.Topic{}
	if err := proto.Unmarshal(i.Value(), topic); err != nil {
		return nil, err
	}
	return topic, nil
}

func (i *TopicIterator) NextPage(in *api.PageInfo) (*api.TopicsPage, error) {
	return nil, nil
}

func (s *Store) ListTopics(orgID, projectID ulid.ULID) iterator.TopicIterator {
	// Iterate over all of the topics prefixed by the projectID
	slice := util.BytesPrefix(projectID.Bytes())
	iter := s.db.NewIterator(slice, nil)
	topics := &TopicIterator{iter}

	return topics
}

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

	if err = s.Put(TopicKey(topic), data); err != nil {
		return err
	}
	return nil
}

func (s *Store) RetrieveTopic(topicID ulid.ULID) (*api.Topic, error) {
	return nil, nil
}

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

	if err = s.Put(TopicKey(topic), data); err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteTopic(topicID ulid.ULID) error {
	if s.readonly {
		return errors.ErrReadOnly
	}
	return nil
}

// TopicKey is a 32 byte value that is the concatenated projectID followed by the topicID.
func TopicKey(topic *api.Topic) []byte {
	key := make([]byte, 32)
	copy(key[0:16], topic.ProjectId)
	copy(key[16:], topic.Id)
	return key
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
