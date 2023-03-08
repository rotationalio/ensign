/*
Ensign maintains two separate storage locations on disk: the event store which is
intended to be an append-only fast disk write for incoming events and a meta store which
is used to persist operational metadata such as topic and placement information.
*/
package store

import (
	"io"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/events"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
)

func Open(conf config.StorageConfig) (data EventStore, meta MetaStore, err error) {
	if data, err = OpenEvents(conf); err != nil {
		return nil, nil, err
	}

	if meta, err = OpenMeta(conf); err != nil {
		data.Close()
		return nil, nil, err
	}

	return data, meta, nil
}

func OpenEvents(conf config.StorageConfig) (*events.Store, error) {
	return events.Open(conf)
}

func OpenMeta(conf config.StorageConfig) (*meta.Store, error) {
	return meta.Open(conf)
}

type Store interface {
	io.Closer
	ReadOnly() bool
}

type EventStore interface {
	Store
}

type MetaStore interface {
	Store
	TopicStore
}

type TopicStore interface {
	ListTopics(projectID ulid.ULID) iterator.TopicIterator
	CreateTopic(*api.Topic) error
	RetrieveTopic(topicID ulid.ULID) (*api.Topic, error)
	UpdateTopic(*api.Topic) error
	DeleteTopic(topicID ulid.ULID) error
}
