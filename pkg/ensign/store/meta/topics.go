package meta

import (
	"encoding/base64"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const MaxTopicNameLength = 512

// Implements iterator.TopicIterator to provide access to a list of topics.
type TopicIterator struct {
	ldbiter.Iterator
}

func (i *TopicIterator) Error() (err error) {
	if err = i.Iterator.Error(); err != nil {
		return errors.Wrap(err)
	}
	return nil
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

// AllowedTopics returns the topic ids for all of the topics specified in the project.
// This method is routinely used for validating topics available to the user.
func (s *Store) AllowedTopics(projectID ulid.ULID) (_ []ulid.ULID, err error) {
	allowed := make([]ulid.ULID, 0)
	topics := s.ListTopics(projectID)
	defer topics.Release()

	for topics.Next() {
		var key ObjectKey
		copy(key[:], topics.Key())

		// Parse the ULID
		var topicID ulid.ULID
		if topicID, err = key.ObjectID(); err != nil {
			return nil, errors.ErrInvalidKey
		}

		allowed = append(allowed, topicID)
	}

	if err = topics.Error(); err != nil {
		return nil, err
	}
	return allowed, nil
}

// List all of the topics associated with the specified projectID. Returns a
// TopicIterator that can be used to retrieve each individual topic or to create a page
// of topics for paginated requests.
func (s *Store) ListTopics(projectID ulid.ULID) iterator.TopicIterator {
	// Iterate over all of the topics prefixed by the projectID
	prefix := make([]byte, 18)
	projectID.MarshalBinaryTo(prefix[:16])
	copy(prefix[16:18], TopicSegment[:])

	slice := util.BytesPrefix(prefix)
	iter := s.db.NewIterator(slice, nil)
	return &TopicIterator{iter}
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

	// Create the topic name key for uniqueness constraint reasons
	uniqueName := TopicNameKey(topic)

	// Marshal the topic and store it
	var data []byte
	if data, err = proto.Marshal(topic); err != nil {
		return errors.Wrap(err)
	}

	if err = s.Create(TopicKey(topic), data, uniqueName); err != nil {
		if errors.Is(err, errors.ErrUniqueConstraint) {
			return errors.ErrUniqueTopicName
		}
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
//
// NOTE: We must return an error if the topic name has changed otherwise we will also
// have to modify the uniqueness constraints on topic name and check them as well.
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

	if err = s.Update(TopicKey(topic), data, TopicNameKey(topic)); err != nil {
		if errors.Is(err, errors.ErrUniqueConstraintChanged) {
			return errors.ErrTopicNameChanged
		}
		return err
	}
	return nil
}

// Delete a topic from the database. If the topic does not exist, no error is returned.
// This method uses the keymu lock to avoid concurrency issues and also cleans up any
// indices associated with the topic.
func (s *Store) DeleteTopic(topicID ulid.ULID) (err error) {
	// Lookup the topic name to get the unique constraints to also delete
	var topic *api.Topic
	if topic, err = s.RetrieveTopic(topicID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil
		}
		return err
	}

	// Delete the topic as well as the topic name index.
	// NOTE: because no error is returned if the topic exists, there shouldn't be a
	// concurrency issue between the retrieve above and the delete below.
	if err := s.Destroy(IndexKey(topicID), TopicNameKey(topic)); err != nil {
		return err
	}
	return nil
}

// TopicKey is a 34 byte value that is the concatenated projectID followed by the topic
// segment and then the topicID. We expect that topicIDs are unique in the database.
func TopicKey(topic *api.Topic) ObjectKey {
	var key [34]byte
	copy(key[0:16], topic.ProjectId)
	copy(key[16:18], TopicSegment[:])
	copy(key[18:], topic.Id)
	return ObjectKey(key)
}

// Validate a topic is ready for storage in the database. If partial is true, then the
// fields that may be set by this package (e.g. ID, created, modified) are not checked,
// otherwise the entire struct is validated.
func ValidateTopic(topic *api.Topic, partial bool) error {
	if topic == nil {
		return errors.ErrTopicInvalidId
	}

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

	if len(topic.Name) > MaxTopicNameLength {
		return errors.ErrTopicNameTooLong
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
