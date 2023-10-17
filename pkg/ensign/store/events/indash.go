package events

import (
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Indash stores an index hash of an event for duplicate detection purposes. In the
// underlying database, the indash key is the topicID + indashSegment + hash and the
// value is the eventID. Primarily the hashes are used to construct bloom filters for
// online duplicate detection, but the hash can also be used to lookup the event that
// represents the duplicate for equality comparison.
func (s *Store) Indash(topicID ulid.ULID, hash []byte, eventID rlid.RLID) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// The key is the topicID + indash segment + hash
	key := make([]byte, 18+len(hash))
	if err = topicID.MarshalBinaryTo(key[:16]); err != nil {
		return errors.Wrap(err)
	}
	copy(key[16:18], IndashSegment[:])
	copy(key[18:], hash)

	// The value is the eventID
	var value []byte
	if value, err = eventID.MarshalBinary(); err != nil {
		return errors.Wrap(err)
	}

	// Write to the database with fsync to avoid data loss
	if err = s.db.Put(key, value, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// LoadIndash returns an iterator that exposes all hashes in the database for the
// specified topicID. The iterator will strip off the topicID and segment from the key
// to return the Hash value.
func (s *Store) LoadIndash(topicID ulid.ULID) iterator.IndashIterator {
	if ulids.IsZero(topicID) {
		return &IndashErrorIterator{ErrorIterator: errors.NewIter(errors.ErrKeyNull)}
	}

	// Iterate over all of the hashes prefixed by the topicID and indash segment
	prefix := make([]byte, 18)
	topicID.MarshalBinaryTo(prefix[:16])
	copy(prefix[16:18], IndashSegment[:])
	slice := util.BytesPrefix(prefix)

	iter := s.db.NewIterator(slice, &opt.ReadOptions{DontFillCache: true})
	return &IndashIterator{Iterator: iter}
}

// ClearIndash deletes all of the index hashes for the the specified topic.
func (s *Store) ClearIndash(topicID ulid.ULID) error {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if ulids.IsZero(topicID) {
		return errors.ErrKeyNull
	}

	// Iterate over all of the hashes prefixed by the topicID and indash segment
	prefix := make([]byte, 18)
	topicID.MarshalBinaryTo(prefix[:16])
	copy(prefix[16:18], IndashSegment[:])
	slice := util.BytesPrefix(prefix)

	batch := &leveldb.Batch{}
	iter := s.db.NewIterator(slice, &opt.ReadOptions{DontFillCache: true})
	defer iter.Release()

	for iter.Next() {
		batch.Delete(iter.Key())
	}

	if err := iter.Error(); err != nil {
		return err
	}

	if err := s.db.Write(batch, &opt.WriteOptions{Sync: false, NoWriteMerge: true}); err != nil {
		return err
	}

	return nil
}
