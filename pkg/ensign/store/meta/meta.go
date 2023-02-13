package meta

import (
	"errors"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/iterator"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
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

func (s *Store) ListTopics(orgID, projectID ulid.ULID) iterator.TopicIterator {
	return nil
}

func (s *Store) CreateTopic(*api.Topic) error {
	return errors.New("not implemented")
}

func (s *Store) RetrieveTopic(topicID ulid.ULID) (*api.Topic, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) UpdateTopic(*api.Topic) error {
	return errors.New("not implemented")
}

func (s *Store) DeleteTopic(topicID ulid.ULID) error {
	return errors.New("not implemented")
}
