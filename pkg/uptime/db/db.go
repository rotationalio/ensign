package db

import (
	"errors"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
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
	return db.OpenTransaction()
}
