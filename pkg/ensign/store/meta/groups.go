package meta

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Implements the iterator.GroupIterator to provide access to a list of groups.
type GroupIterator struct {
	ldbiter.Iterator
}

// Group unmarshals the next ConsumerGroup in the iterator.
func (i *GroupIterator) Group() (*api.ConsumerGroup, error) {
	group := &api.ConsumerGroup{}
	if err := proto.Unmarshal(i.Value(), group); err != nil {
		return nil, err
	}
	return group, nil
}

// List all consumer groups associated with the specified projectID. Returns a
// GroupIterator that can be used to retrieve each individual consumer group.
func (s *Store) ListGroups(projectID ulid.ULID) iterator.GroupIterator {
	// Iterate over all consumer groups prefixed by project ID and group segment
	prefix := make([]byte, 18)
	projectID.MarshalBinaryTo(prefix[:16])
	copy(prefix[16:18], GroupSegment[:])

	slice := util.BytesPrefix(prefix)
	iter := s.db.NewIterator(slice, nil)
	return &GroupIterator{iter}
}

// Create a consumer group in the database. If the group already exists or if the topic
// is not valid an error is returned. This method uses the keymu lock to avoid
// concurrency issues for multiple writers.
func (s *Store) CreateGroup(group *api.ConsumerGroup) (err error) {
	// Validate the group to ensure it can be saved
	if err = ValidateGroup(group, true); err != nil {
		return err
	}

	// Set the created and modified timestamps
	group.Created = timestamppb.Now()
	group.Modified = group.Created

	// Marshal the group and store it
	var data []byte
	if data, err = proto.Marshal(group); err != nil {
		return errors.Wrap(err)
	}

	if err = s.Create(GroupKey(group), data); err != nil {
		return err
	}
	return nil
}

// GroupKey is a 34 byte value that is the concatenated projectID followed by the
// group segment and then the murmur3 hashed key of the group (unless a 16 byte ID is
// specified by the user).
func GroupKey(group *api.ConsumerGroup) ObjectKey {
	// If the key errors then panic - it is the responsibility of the caller to validate
	// this group before they call this function.
	gkey, err := group.Key()
	if err != nil {
		panic(err)
	}

	var key [34]byte
	copy(key[0:16], group.ProjectId)
	copy(key[16:18], GroupSegment[:])
	copy(key[18:], gkey[:])

	return ObjectKey(key)
}

func ValidateGroup(group *api.ConsumerGroup, partial bool) error {
	if len(group.ProjectId) == 0 {
		return errors.ErrGroupMissingProjectId
	}

	if _, err := ulids.Parse(group.ProjectId); err != nil {
		return errors.ErrGroupInvalidProjectId
	}

	if len(group.Id) == 0 && group.Name == "" {
		return errors.ErrGroupMissingKeyField
	}

	if !partial {
		if !IsValidTimestamp(group.Created) {
			return errors.ErrGroupInvalidCreated
		}

		if !IsValidTimestamp(group.Modified) {
			return errors.ErrGroupInvalidModified
		}
	}
	return nil
}
