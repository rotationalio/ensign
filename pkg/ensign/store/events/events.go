package events

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
)

func Open(conf config.StorageConfig) (store *Store, err error) {
	store = &Store{
		readonly: conf.ReadOnly,
	}

	var path string
	if path, err = conf.EventPath(); err != nil {
		return nil, err
	}

	if store.db, err = leveldb.OpenFile(path, &opt.Options{ReadOnly: conf.ReadOnly}); err != nil {
		return nil, err
	}
	return store, nil
}

type Store struct {
	db       *leveldb.DB
	readonly bool
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) ReadOnly() bool {
	return s.readonly
}

// Insert an event with the event segment into the database. If the event doesn't have
// an ID or a TopicID, an error is returned. This method also ensures that the localID
// is not stored and is nil. No other validation is performed by the database as this
// method is designed to write as quickly as possible.
func (s *Store) Insert(event *api.EventWrapper) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// The localID should not be stored in the database
	event.LocalId = nil

	var key Key
	if key, err = EventKey(event); err != nil {
		return err
	}

	var value []byte
	if value, err = proto.Marshal(event); err != nil {
		return errors.Wrap(err)
	}

	// Write to the database with fsync to avoid data loss
	if err = s.db.Put(key[:], value, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Returns an iterator of events in the specified topic. If an offset RLID is specified
func (s *Store) List(topicID ulid.ULID) iterator.EventIterator {
	if ulids.IsZero(topicID) {
		return &EventErrorIterator{ErrorIterator: errors.NewIter(errors.ErrKeyNull)}
	}

	// Iterate over all of the events prefixed by the topicID and the event segment
	prefix := make([]byte, 18)
	topicID.MarshalBinaryTo(prefix[:16])
	copy(prefix[16:18], EventSegment[:])
	slice := util.BytesPrefix(prefix)

	iter := s.db.NewIterator(slice, nil)
	return &EventIterator{Iterator: iter, topicID: topicID}
}

// Retrieve a specific event from the database by topic and eventID.
func (s *Store) Retrieve(topicId ulid.ULID, eventID rlid.RLID) (event *api.EventWrapper, err error) {
	var key Key
	if key, err = CreateKey(topicId, eventID, EventSegment); err != nil {
		return nil, err
	}

	var data []byte
	if data, err = s.db.Get(key[:], nil); err != nil {
		return nil, err
	}

	event = &api.EventWrapper{}
	if err = proto.Unmarshal(data, event); err != nil {
		return nil, err
	}

	return event, nil
}

// Count the number of objects that match the specified range by iterating through all
// of the keys and counting them. This is primarily used for testing.
func (s *Store) Count(slice *util.Range) (count uint64, err error) {
	iter := s.db.NewIterator(slice, &opt.ReadOptions{DontFillCache: true})
	defer iter.Release()

	for iter.Next() {
		count++
	}
	return count, iter.Error()
}
