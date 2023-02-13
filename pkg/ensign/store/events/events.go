package events

import (
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
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
