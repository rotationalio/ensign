package events

import (
	"errors"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
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

func (s *Store) Insert(event *api.EventWrapper) error {
	return errors.New("not implemented yet")
}

func (s *Store) List(topicID ulid.ULID, eventID rlid.RLID) iterator.EventIterator {
	return nil
}

func (s *Store) Retrieve(topicId ulid.ULID, eventID rlid.RLID) (*api.EventWrapper, error) {
	return nil, errors.New("not implemented yet")
}
