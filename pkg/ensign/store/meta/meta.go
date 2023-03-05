package meta

import (
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func Open(conf config.StorageConfig) (store *Store, err error) {
	store = &Store{
		readonly: conf.ReadOnly,
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

func (s *Store) Get(key []byte) (value []byte, err error) {
	if value, err = s.db.Get(key, nil); err != nil {
		return nil, errors.Wrap(err)
	}
	return value, nil
}

func (s *Store) Put(key, value []byte) error {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if err := s.db.Put(key, value, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (s *Store) Create(key, value []byte) error {
	if s.readonly {
		return errors.ErrReadOnly
	}

	s.db.OpenTransaction()
	return nil
}

func (s *Store) Delete(key []byte) error {
	if s.readonly {
		return errors.ErrReadOnly
	}

	if err := s.db.Delete(key, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (s *Store) Count(slice *util.Range) (count uint64, err error) {
	iter := s.db.NewIterator(slice, &opt.ReadOptions{DontFillCache: true})
	defer iter.Release()

	for iter.Next() {
		count++
	}
	return count, iter.Error()
}
