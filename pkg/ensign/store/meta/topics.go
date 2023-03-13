package meta

import (
	"encoding/base64"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/rs/zerolog/log"
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

// NextPage seeks the iterator to the next page of results and returns a page of topics.
func (i *TopicIterator) NextPage(in *api.PageInfo) (page *api.TopicsPage, err error) {
	if in == nil {
		in = &api.PageInfo{}
	}

	if in.PageSize == 0 {
		in.PageSize = uint32(pagination.DefaultPageSize)
	}

	if in.NextPageToken != "" {
		// Parse the next page cursor
		var cursor *pagination.Cursor
		if cursor, err = pagination.Parse(in.NextPageToken); err != nil {
			return nil, errors.Wrap(err)
		}

		// Seek the iterator to the correct page
		var seekKey []byte
		if seekKey, err = base64.RawStdEncoding.DecodeString(cursor.EndIndex); err != nil {
			return nil, errors.ErrInvalidPageToken
		}

		if !i.Seek(seekKey) {
			// Return an empty page if the seek returns empty
			return &api.TopicsPage{}, nil
		}
	}

	hasNextPage := false
	page = &api.TopicsPage{
		Topics: make([]*api.Topic, 0, in.PageSize),
	}

	for i.Next() {
		// Check if we're done iterating; if we have a full page, then we've gone one
		// item over the page size, so we can create the next page cursor.
		if len(page.Topics) == int(in.PageSize) {
			hasNextPage = true
			break
		}

		// Append the current topic to the page
		var topic *api.Topic
		if topic, err = i.Topic(); err != nil {
			log.Error().Err(err).Bytes("topic_key", i.Key()).Msg("could not parse topic stored in database")
			continue
		}

		page.Topics = append(page.Topics, topic)
	}

	if hasNextPage {
		i.Prev()
		endIndex := base64.RawStdEncoding.EncodeToString(i.Key())
		cursor := pagination.New("", endIndex, int32(in.PageSize))

		if page.NextPageToken, err = cursor.NextPageToken(); err != nil {
			return nil, err
		}
	}
	return page, nil
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
