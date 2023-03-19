package meta

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
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

// GetOrCreate a consumer group in the database. This method is the primary mechanism of
// retrieving a consumer group object from the database. The user should only specify
// the minimum required fields on the group and if the group exists the object will be
// updated in place with other values. Otherwise, the group will be created with default
// values set ready service a new group. If the group was created this function will
// return true.
func (s *Store) GetOrCreateGroup(group *api.ConsumerGroup) (created bool, err error) {
	// Validate the group to ensure it can be saved
	if err = ValidateGroup(group, true); err != nil {
		return false, err
	}

	// Acquire a lock on the object key to avoid concurrency issues
	key := GroupKey(group)
	mu := s.keymu.Lock(key)
	defer mu.Unlock()

	// Check if the group already exists in the database
	var exists bool
	if exists, err = s.db.Has(key[:], nil); err != nil {
		return false, errors.Wrap(err)
	}

	if exists {
		// Retrieve the group and return!
		var data []byte
		if data, err = s.db.Get(key[:], nil); err != nil {
			return false, errors.Wrap(err)
		}

		if err = proto.Unmarshal(data, group); err != nil {
			return false, errors.Wrap(err)
		}
		return false, nil
	}

	// Otherwise create the group in the database
	if s.readonly {
		return false, errors.ErrReadOnly
	}

	// If the ID is not set compute it from the name for faster future retrieval.
	if len(group.Id) == 0 {
		var key [16]byte
		if key, err = group.Key(); err != nil {
			return false, err
		}
		group.Id = key[:]
	}

	// Set the created and modified timestamps
	group.Created = timestamppb.Now()
	group.Modified = group.Created

	// Marshal the group and store it
	var data []byte
	if data, err = proto.Marshal(group); err != nil {
		return false, errors.Wrap(err)
	}

	if err = s.db.Put(key[:], data, &opt.WriteOptions{Sync: true}); err != nil {
		return false, errors.Wrap(err)
	}
	return true, nil
}

// Create a consumer group in the database. If the group already exists or if the topic
// is not valid an error is returned. The entire object key for the group is locked in
// the keymu to prevent concurrent writes to the group.
func (s *Store) CreateGroup(group *api.ConsumerGroup) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// Validate the group to ensure it can be saved
	if err = ValidateGroup(group, true); err != nil {
		return err
	}

	// If the ID is not set compute it from the name for faster future retrieval.
	if len(group.Id) == 0 {
		var key [16]byte
		if key, err = group.Key(); err != nil {
			return err
		}
		group.Id = key[:]
	}

	// Set the created and modified timestamps
	group.Created = timestamppb.Now()
	group.Modified = group.Created

	// Acquire a lock on the object key to avoid concurrency issues
	// NOTE: this method cannot use s.Create because it is not guaranteed that the user
	// will specify a unique index ID. The uniqueness constraint for groups is the
	// project ID + the user-supplied key (e.g. name or ID).
	key := GroupKey(group)
	mu := s.keymu.Lock(key)
	defer mu.Unlock()

	// Check if the group already exists in the database
	var exists bool
	if exists, err = s.db.Has(key[:], nil); err != nil {
		return errors.Wrap(err)
	}

	if exists {
		return errors.ErrAlreadyExists
	}

	// Marshal the group and store it
	var data []byte
	if data, err = proto.Marshal(group); err != nil {
		return errors.Wrap(err)
	}

	if err = s.db.Put(key[:], data, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Retrieve updates the group pointer in place with the current values in the database.
func (s *Store) RetrieveGroup(group *api.ConsumerGroup) (err error) {
	if err = ValidateGroup(group, true); err != nil {
		return err
	}

	// Fetch the data from the database by he group key.
	var data []byte
	key := GroupKey(group)
	if data, err = s.db.Get(key[:], nil); err != nil {
		return errors.Wrap(err)
	}

	if err = proto.Unmarshal(data, group); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Update a group by putting the data from the input group into the database. If the
// group does not exist, an error is returned (unlike normal Put semantics). To avoid
// concurrency issues, this method locks the object key before performing writes.
func (s *Store) UpdateGroup(group *api.ConsumerGroup) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// Validate the complete group - if the group does not have a project ID or a group
	// key field (e.g. id or name) then an error is returned preventing panics when the
	// group key is created for storage.
	if err = ValidateGroup(group, false); err != nil {
		return err
	}

	// Update the modified timestamp on the group.
	group.Modified = timestamppb.Now()

	// Acquire a lock on the object key to avoid concurrency issues
	// NOTE: this method cannot use s.Update because it is not guaranteed that the user
	// will specify a unique index ID. The uniqueness constraint for groups is the
	// project ID + the user-supplied key (e.g. name or ID).
	key := GroupKey(group)
	mu := s.keymu.Lock(key)
	defer mu.Unlock()

	// Check if the group already exists in the database
	var exists bool
	if exists, err = s.db.Has(key[:], nil); err != nil {
		return errors.Wrap(err)
	}

	if !exists {
		return errors.ErrNotFound
	}

	// Marshal the group and store it
	var data []byte
	if data, err = proto.Marshal(group); err != nil {
		return errors.Wrap(err)
	}

	if err = s.db.Put(key[:], data, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// Delete a group from the database. If the group does not exist, no error is returned.
func (s *Store) DeleteGroup(group *api.ConsumerGroup) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if err = ValidateGroup(group, true); err != nil {
		return err
	}

	key := GroupKey(group)
	if err = s.db.Delete(key[:], &opt.WriteOptions{Sync: false}); err != nil {
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
	if group == nil {
		return errors.ErrGroupMissingKeyField
	}

	if len(group.Id) == 0 && group.Name == "" {
		return errors.ErrGroupMissingKeyField
	}

	if len(group.ProjectId) == 0 {
		return errors.ErrGroupMissingProjectId
	}

	if _, err := ulids.Parse(group.ProjectId); err != nil {
		return errors.ErrGroupInvalidProjectId
	}

	if !partial {
		if len(group.Id) == 0 {
			return errors.ErrGroupMissingId
		}

		if !IsValidTimestamp(group.Created) {
			return errors.ErrGroupInvalidCreated
		}

		if !IsValidTimestamp(group.Modified) {
			return errors.ErrGroupInvalidModified
		}
	}
	return nil
}
