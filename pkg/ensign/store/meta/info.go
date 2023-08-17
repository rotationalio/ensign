package meta

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TopicInfo returns the topic info struct for the given topic by first checking if the
// topic exists in the database, then returning either the info stored in the database
// or an empty topic info if the value has not been stored before. If the topic does not
// exist then a not found error is returned.
func (s *Store) TopicInfo(topicID ulid.ULID) (_ *api.TopicInfo, err error) {
	// Retrieve the projectId from the topicID index key
	var index IndexKey
	if index, err = CreateIndex(topicID); err != nil {
		return nil, err
	}

	var keyData []byte
	if keyData, err = s.db.Get(index[:], nil); err != nil {
		return nil, errors.Wrap(err)
	}

	var topicKey ObjectKey
	if err = topicKey.UnmarshalValue(keyData); err != nil {
		return nil, errors.Wrap(err)
	}

	// Convert the topic key into a topic info key
	topicKey.Convert(TopicInfoSegment)

	var data []byte
	if data, err = s.db.Get(topicKey[:], nil); err != nil {
		// If the error is not found, then return an empty topicInfo struct
		if errors.Is(err, leveldb.ErrNotFound) {
			return &api.TopicInfo{
				TopicId:   topicKey[18:],
				ProjectId: topicKey[:16],
			}, nil
		}

		// Otherwise return the error
		return nil, errors.Wrap(err)
	}

	info := &api.TopicInfo{}
	if err = proto.Unmarshal(data, info); err != nil {
		return nil, errors.Wrap(err)
	}

	return info, nil
}

// Replaces the current topic info value in the database with the one specified,
// updating the modified timestamp as it does. Beware reads than writes to the database
// without some kind of external locking or serialization as it is possible to overwrite
// another routine's write without care.
func (s *Store) UpdateTopicInfo(info *api.TopicInfo) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if err = ValidateTopicInfo(info); err != nil {
		return err
	}

	// Set the modified timestamp on the struct
	info.Modified = timestamppb.Now()
	key := TopicInfoKey(info)

	// Marshal the protocol buffer
	var value []byte
	if value, err = proto.Marshal(info); err != nil {
		return errors.Wrap(err)
	}

	if err = s.db.Put(key[:], value, nil); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func TopicInfoKey(info *api.TopicInfo) ObjectKey {
	var key ObjectKey
	copy(key[0:16], info.ProjectId)
	copy(key[16:18], TopicInfoSegment[:])
	copy(key[18:], info.TopicId)
	return key
}

func ValidateTopicInfo(info *api.TopicInfo) error {
	switch {
	case info == nil:
		return errors.ErrTopicInfoInvalidTopicId
	case len(info.ProjectId) == 0:
		return errors.ErrTopicInfoMissingProjectId
	case len(info.TopicId) == 0:
		return errors.ErrTopicInfoMissingTopicId
	}

	if projectID, err := ulids.Parse(info.ProjectId); err != nil || ulids.IsZero(projectID) {
		return errors.ErrTopicInfoInvalidProjectId
	}

	if topicID, err := ulids.Parse(info.TopicId); err != nil || ulids.IsZero(topicID) {
		return errors.ErrTopicInfoInvalidTopicId
	}

	return nil
}
