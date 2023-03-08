package meta

import (
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/utils/keymu"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Open a new metadata store at the metapath supplied by the configuration.
func Open(conf config.StorageConfig) (store *Store, err error) {
	store = &Store{
		readonly: conf.ReadOnly,
		keymu:    keymu.New(),
	}

	var path string
	if path, err = conf.MetaPath(); err != nil {
		return nil, err
	}

	if store.db, err = leveldb.OpenFile(path, &opt.Options{ReadOnly: conf.ReadOnly}); err != nil {
		return nil, err
	}
	return store, nil
}

// Store implements the store.MetaStore interface for interacting with topics and other
// persistent data in the ensign database. Store can be readonly which is enforced both
// by the underlying disk reads to leveldb and the readonly flag at the top level.
type Store struct {
	db       *leveldb.DB
	readonly bool
	keymu    *keymu.Mutex
}

// Close the underlying leveldb gracefully to avoid database corruption.
func (s *Store) Close() error {
	return s.db.Close()
}

// Returns the readonly state of the meta store.
func (s *Store) ReadOnly() bool {
	return s.readonly
}

// Gets a value for the specified key, wrapping any leveldb errors in an errors.Error.
// NOTE: if getting objects it is preferred to use Retrieve to avoid concurrency issues.
func (s *Store) Get(key []byte) (value []byte, err error) {
	if value, err = s.db.Get(key, nil); err != nil {
		return nil, errors.Wrap(err)
	}
	return value, nil
}

// Put a value for the specified pair, wrapping any leveldb errors in an errors.Error.
// NOTE: if putting objects it is preferred to use Create or Update to avoid
// concurrency issues by using key-specific transactions.
func (s *Store) Put(key, value []byte) error {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if err := s.db.Put(key, value, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Delete the specified key, wrapping any leveldb errors in an errors.Error.
// NOTE: if deleting objects it is preferred to use Destroy to avoid concurrency issues.
func (s *Store) Delete(key []byte) error {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if err := s.db.Delete(key, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Has checks the existence of a key in the database, wrapping any leveldb errors in an
// errors.Error for easy error checking and verification.
func (s *Store) Has(key []byte) (ret bool, err error) {
	if ret, err = s.db.Has(key, nil); err != nil {
		return ret, errors.Wrap(err)
	}
	return ret, nil
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

// Create an object in the database with the specified key by saving both the object key
// in the database as well as the object mapped to the key value. If the object already
// exists, then an error is returned. To prevent concurrency issues, Create locks the
// object key before access and performs all writes in batch.
// TODO: extend tests for this method to validate concurrency.
func (s *Store) Create(key ObjectKey, value []byte) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// Acquire a lock on the object key to avoid concurrency issues
	indexKey := key.Key()
	mu := s.keymu.Lock(indexKey)
	defer mu.Unlock()

	// Check if the object already exists in the database
	var exists bool
	if exists, err = s.db.Has(indexKey[:], nil); err != nil {
		return errors.Wrap(err)
	}

	if exists {
		return errors.ErrAlreadyExists
	}

	// Create the batch write to put the object key and the value to the database
	batch := &leveldb.Batch{}
	batch.Put(indexKey[:], key[:])
	batch.Put(key[:], value)

	if err = s.db.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Retrieve an object from the database by first looking up its object key with the
// specified objectID then returning the value at the retrieved object key.
// TODO: implement read locks to avoid concurrency issues.
func (s *Store) Retrieve(key IndexKey) (value []byte, err error) {
	// Fetch the object key from the index
	var data []byte
	if data, err = s.db.Get(key[:], nil); err != nil {
		return nil, errors.Wrap(err)
	}

	var objectKey ObjectKey
	if err = objectKey.UnmarshalValue(data); err != nil {
		return nil, errors.Wrap(err)
	}

	if value, err = s.db.Get(objectKey[:], nil); err != nil {
		return nil, errors.Wrap(err)
	}
	return value, nil
}

// Update an object in the database with the specified object key. If the object does
// not exist, an error is returned (unlike normal Put semantics). To avoid concurrency
// issues, update locks the specified key before performing any writes.
func (s *Store) Update(key ObjectKey, value []byte) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// Acquire a lock on the object key to avoid concurrency issues
	indexKey := key.Key()
	mu := s.keymu.Lock(indexKey)
	defer mu.Unlock()

	// Check if the object already exists in the database
	var data []byte
	if data, err = s.db.Get(indexKey[:], nil); err != nil {
		return errors.Wrap(err)
	}

	var objectKey ObjectKey
	if err = objectKey.UnmarshalValue(data); err != nil {
		return errors.Wrap(err)
	}

	// Update the object value in the database
	if err = s.db.Put(objectKey[:], value, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Destroy an object in the database by first looking its object key with the specified
// objectID then deleting the value at the retrieved object key. To avoid concurrency
// issues, destroy locks the key and deletes both the object key and the object in a
// batch write. If the object does not exist, no error is returned.
func (s *Store) Destroy(key IndexKey) (err error) {
	if s.readonly {
		return errors.ErrReadOnly
	}

	// Acquire a lock on the object key to avoid concurrency issues
	mu := s.keymu.Lock(key)
	defer mu.Unlock()

	// Fetch the object key from the index
	var data []byte
	if data, err = s.db.Get(key[:], nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil
		}
		return errors.Wrap(err)
	}

	var objectKey ObjectKey
	if err = objectKey.UnmarshalValue(data); err != nil {
		return errors.Wrap(err)
	}

	batch := &leveldb.Batch{}
	batch.Delete(objectKey[:])
	batch.Delete(key[:])
	if err = s.db.Write(batch, &opt.WriteOptions{Sync: false}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
