/*
Ensign maintains two separate storage locations on disk: the event store which is
intended to be an append-only fast disk write for incoming events and a meta store which
is used to persist operational metadata such as topic and placement information.
*/
package store

import (
	"io"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/events"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
)

func Open(conf config.StorageConfig) (data EventStore, meta MetaStore, err error) {
	// If in testing mode return a mock store for both data and meta.
	if conf.Testing {
		var mockStore *mock.Store
		mockStore, err = mock.Open(conf)
		return mockStore, mockStore, err
	}

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
	Insert(*api.EventWrapper) error
	List(topicID ulid.ULID) iterator.EventIterator
	Retrieve(topicID ulid.ULID, eventID rlid.RLID) (*api.EventWrapper, error)
}

type MetaStore interface {
	Store
	TopicStore
	TopicNamesStore
	TopicInfoStore
}

type TopicStore interface {
	AllowedTopics(projectID ulid.ULID) ([]ulid.ULID, error)
	ListTopics(projectID ulid.ULID) iterator.TopicIterator
	CreateTopic(*api.Topic) error
	RetrieveTopic(topicID ulid.ULID) (*api.Topic, error)
	UpdateTopic(*api.Topic) error
	DeleteTopic(topicID ulid.ULID) error
}

type TopicNamesStore interface {
	ListTopicNames(projectID ulid.ULID) iterator.TopicNamesIterator
	TopicExists(in *api.TopicName) (*api.TopicExistsInfo, error)
	TopicName(topicID ulid.ULID) (string, error)
	LookupTopicID(name string, projectID ulid.ULID) (topicID ulid.ULID, err error)
}

type TopicInfoStore interface {
	TopicInfo(topicID ulid.ULID) (*api.TopicInfo, error)
	UpdateTopicInfo(*api.TopicInfo) error
}

type GroupStore interface {
	ListGroups(projectID ulid.ULID) iterator.GroupIterator
	GetOrCreateGroup(*api.ConsumerGroup) (bool, error)
	UpdateGroup(*api.ConsumerGroup) error
	DeleteGroup(*api.ConsumerGroup) error
}
