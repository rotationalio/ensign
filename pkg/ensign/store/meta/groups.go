package meta

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
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
	return nil
}
