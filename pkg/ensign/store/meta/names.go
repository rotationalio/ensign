package meta

import (
	"encoding/base64"
	"fmt"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Implements iterator.TopicNamesIterator to provide access to the topic names index.
type TopicNamesIterator struct {
	ldbiter.Iterator
}

func (i *TopicNamesIterator) Error() (err error) {
	if err = i.Iterator.Error(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// TopicName unmarshals the next topic name in the iterator
func (i *TopicNamesIterator) TopicName() (_ *api.TopicName, err error) {
	// Right now because the name is the murmur3 hash in the index we cannot retrieve
	// the name directly from the index; so we have to retrieve the object; but in the
	// future we should include the name as part of the value in the index so it's
	// retrievable directly from the index.
	name := &api.TopicName{}

	// Parse the key which is projectid:segment:namehash
	var nameKey ObjectKey
	if err = nameKey.UnmarshalValue(i.Key()); err != nil {
		return nil, err
	}

	// Parse the value which is projectid:segement:topicid
	var objectKey ObjectKey
	if err = objectKey.UnmarshalValue(i.Value()); err != nil {
		return nil, err
	}

	// HACK: right now we're returning the base64 encoded murmur3 hash but it would be
	// better to return the name so the SDKs don't have to do so much work.
	name.Name = base64.RawURLEncoding.EncodeToString(nameKey[18:])

	var projectID ulid.ULID
	if projectID, err = ulids.Parse(nameKey[:16]); err != nil {
		return nil, err
	}

	var topicID ulid.ULID
	if topicID, err = ulids.Parse(objectKey[18:]); err != nil {
		return nil, err
	}

	name.ProjectId = projectID.String()
	name.TopicId = topicID.String()
	return name, nil
}

func (i *TopicNamesIterator) NextPage(in *api.PageInfo) (page *api.TopicNamesPage, err error) {
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
			return &api.TopicNamesPage{}, nil
		}
	}

	hasNextPage := false
	page = &api.TopicNamesPage{
		TopicNames: make([]*api.TopicName, 0, in.PageSize),
	}

	for i.Next() {
		// Check if we're done iterating; if we have a full page, then we've gone one
		// item over the page size, so we can create the next page cursor.
		if len(page.TopicNames) == int(in.PageSize) {
			hasNextPage = true
			break
		}

		// Append the current topic to the page
		var topic *api.TopicName
		if topic, err = i.TopicName(); err != nil {
			sentry.Error(nil).Err(err).Bytes("topic_name_key", i.Key()).Msg("could not parse topic stored in database")
			continue
		}

		page.TopicNames = append(page.TopicNames, topic)
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

func (s *Store) ListTopicNames(projectID ulid.ULID) iterator.TopicNamesIterator {
	// Iterate over the topic names index, prefixed by the projectID
	prefix := make([]byte, 18)
	projectID.MarshalBinaryTo(prefix[:16])
	copy(prefix[16:18], TopicNamesSegment[:])

	slice := util.BytesPrefix(prefix)
	iter := s.db.NewIterator(slice, nil)

	return &TopicNamesIterator{iter}
}

func (s *Store) TopicExists(in *api.TopicName) (_ *api.TopicExistsInfo, err error) {
	var projectID ulid.ULID
	if projectID, err = ulids.Parse(in.ProjectId); err != nil || ulids.IsZero(projectID) {
		return nil, errors.ErrTopicInvalidProjectId
	}

	var (
		topicExists bool
		nameExists  bool
	)

	if in.Name != "" {
		key := TopicNameKey(&api.Topic{ProjectId: projectID.Bytes(), Name: in.Name})
		if nameExists, err = s.db.Has(key[:], nil); err != nil {
			return nil, errors.Wrap(err)
		}
	}

	if in.TopicId != "" {
		var topicID ulid.ULID
		if topicID, err = ulids.Parse(in.TopicId); err != nil || ulids.IsZero(topicID) {
			return nil, errors.ErrTopicInvalidId
		}

		key := TopicKey(&api.Topic{ProjectId: projectID.Bytes(), Id: topicID.Bytes()})
		if topicExists, err = s.db.Has(key[:], nil); err != nil {
			return nil, errors.Wrap(err)
		}
	}

	info := &api.TopicExistsInfo{}

	switch {
	case in.Name != "" && in.TopicId != "":
		info.Query = fmt.Sprintf("name=%q and topic=%q", in.Name, in.TopicId)
		info.Exists = nameExists && topicExists
	case in.Name != "":
		info.Query = fmt.Sprintf("name=%q", in.Name)
		info.Exists = nameExists
	case in.TopicId != "":
		info.Query = fmt.Sprintf("topic=%q", in.TopicId)
		info.Exists = topicExists
	default:
		return nil, errors.ErrTopicMissingName
	}

	return info, nil
}

// TopicName returns the name of the topic for the specified topicID.
func (s *Store) TopicName(topicID ulid.ULID) (_ string, err error) {
	// TODO: use names index rather than the topic to get the name.
	var topic *api.Topic
	if topic, err = s.RetrieveTopic(topicID); err != nil {
		return "", err
	}
	return topic.Name, nil
}

// LookupTopicName returns a topicID for the specified topic name and project.
func (s *Store) LookupTopicName(name string, projectID ulid.ULID) (topicID ulid.ULID, err error) {
	if name == "" || ulids.IsZero(projectID) {
		return ulids.Null, errors.ErrNotFound
	}

	// Create a stub topic to get the index key from.
	key := TopicNameKey(&api.Topic{ProjectId: projectID[:], Name: name})

	// Get the object key from the database and parse it to fetch the topicID
	var data []byte
	if data, err = s.db.Get(key[:], nil); err != nil {
		return ulids.Null, errors.Wrap(err)
	}

	var objectKey ObjectKey
	if err = objectKey.UnmarshalValue(data); err != nil {
		return ulids.Null, errors.Wrap(err)
	}

	// The last section of bytes is expected to be the topic ID
	var topicExists bool
	if topicExists, err = s.db.Has(objectKey[:], nil); err != nil {
		return ulids.Null, errors.Wrap(err)
	}

	if !topicExists {
		return ulids.Null, errors.ErrNotFound
	}

	if topicID, err = ulids.Parse(objectKey[18:]); err != nil {
		return ulids.Null, errors.Wrap(err)
	}

	return topicID, nil
}

// TopicNameKey is a 34 byte value that is the concatenated projectID followed by the
// topic segment and then the murmur3 hashed topic name. This allows us to ensure that
// topic names are unique to the project.
func TopicNameKey(topic *api.Topic) ObjectKey {
	var key [34]byte
	copy(key[0:16], topic.ProjectId)
	copy(key[16:18], TopicNamesSegment[:])
	copy(key[18:], topic.NameHash())

	return ObjectKey(key)
}
