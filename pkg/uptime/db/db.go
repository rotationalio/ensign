package db

import (
	"errors"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	db      *leveldb.DB  // connection to the database managed by the package
	connmu  sync.RWMutex // synchronize connect and close
	connect sync.Once    // ensure the database is only connected to once
)

var (
	ErrNotConnected = errors.New("not connected to the database")
	ErrNotFound     = errors.New("record not found")
)

// Connect to the leveldb, opening the database at the specified path.
func Connect(path string, readonly bool) (err error) {
	connmu.Lock()
	defer connmu.Unlock()

	connect.Do(func() {
		if db, err = leveldb.OpenFile(path, &opt.Options{ReadOnly: readonly}); err != nil {
			return
		}
	})

	return err
}

// Close the connection to the leveldb and allow reconnect.
func Close() (err error) {
	connmu.Lock()
	if db != nil {
		err = db.Close()
		db = nil
		connect = sync.Once{}
	}
	connmu.Unlock()
	return err
}

// OpenTransaction opens an atomic DB transaction. Only one transaction can be opened at
// a time. Subsequent call to Write and OpenTransaction will be blocked until in-flight
// transaction is committed or discarded. The returned transaction handle is safe for
// concurrent use.
func BeginTx() (*leveldb.Transaction, error) {
	connmu.RLock()
	defer connmu.RUnlock()
	return db.OpenTransaction()
}

// Get a value from the database by the key and unmarshal it into the specified model.
func Get(key []byte, m Model) (err error) {
	connmu.RLock()
	if db == nil {
		connmu.RUnlock()
		return ErrNotConnected
	}

	// Execute the Get request and immediately unlock
	data, err := db.Get(key, nil)
	connmu.RUnlock()

	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, leveldb.ErrClosed) {
			return ErrNotConnected
		}

		return err
	}

	// Unmarshal the data into the model
	return m.Unmarshal(data)
}

// Put a value into the database by marshaling it from the specified model.
func Put(m Model) (err error) {
	// Marshal the model
	var key, value []byte
	if key, err = m.Key(); err != nil {
		return err
	}

	if value, err = m.Marshal(); err != nil {
		return err
	}

	connmu.RLock()
	defer connmu.RUnlock()
	if db == nil {
		return ErrNotConnected
	}

	if err = db.Put(key, value, nil); err != nil {
		if errors.Is(err, leveldb.ErrClosed) {
			return ErrNotConnected
		}
	}
	return nil
}

// Delete a model from the database.
func Delete(m Model) (err error) {
	var key []byte
	if key, err = m.Key(); err != nil {
		return err
	}

	connmu.RLock()
	defer connmu.RUnlock()
	if db == nil {
		return ErrNotConnected
	}

	if err = db.Delete(key, nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, leveldb.ErrClosed) {
			return ErrNotConnected
		}

		return err
	}
	return nil
}

// Create a new iterator - note that it is possible to close the DB connection during
// iteration, which will cause the iterator to return a leveldb.ErrClosed error.
func NewIterator(slice *util.Range) iterator.Iterator {
	connmu.RLock()
	defer connmu.RUnlock()
	if db == nil {
		return iterator.NewEmptyIterator(ErrNotConnected)
	}

	iter := db.NewIterator(slice, nil)
	return iter
}
