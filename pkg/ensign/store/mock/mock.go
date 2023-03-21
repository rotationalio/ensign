package mock

import (
	"errors"
	"fmt"
	"os"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/protobuf/encoding/protojson"
)

// Constants are used to reference store methods in mock code
const (
	Close          = "Close"
	ReadOnly       = "ReadOnly"
	AllowedTopics  = "AllowedTopics"
	ListTopics     = "ListTopics"
	CreateTopic    = "CreateTopic"
	RetrieveTopic  = "RetrieveTopic"
	UpdateTopic    = "UpdateTopic"
	DeleteTopic    = "DeleteTopic"
	ListTopicNames = "ListTopicNames"
	TopicExists    = "TopicExists"
)

// Implements both a store.EventStore and a store.MetaStore for testing purposes.
type Store struct {
	readonly         bool
	calls            map[string]int
	OnClose          func() error
	OnReadOnly       func() bool
	OnAllowedTopics  func(ulid.ULID) ([]ulid.ULID, error)
	OnListTopics     func(ulid.ULID) iterator.TopicIterator
	OnCreateTopic    func(*api.Topic) error
	OnRetrieveTopic  func(topicID ulid.ULID) (*api.Topic, error)
	OnUpdateTopic    func(*api.Topic) error
	OnDeleteTopic    func(topicID ulid.ULID) error
	OnListTopicNames func(ulid.ULID) iterator.TopicNamesIterator
	OnTopicExists    func(*api.TopicName) (*api.TopicExistsInfo, error)
}

func Open(conf config.StorageConfig) (*Store, error) {
	if !conf.Testing {
		return nil, errors.New("invalid configuration: must be in testing mode")
	}
	return &Store{
		readonly: conf.ReadOnly,
		calls:    make(map[string]int),
	}, nil
}

func (s *Store) Reset() {
	for key := range s.calls {
		s.calls[key] = 0
	}

	s.OnClose = nil
	s.OnReadOnly = nil
	s.OnAllowedTopics = nil
	s.OnListTopics = nil
	s.OnCreateTopic = nil
	s.OnRetrieveTopic = nil
	s.OnUpdateTopic = nil
	s.OnDeleteTopic = nil
	s.OnListTopicNames = nil
	s.OnTopicExists = nil
}

func (s *Store) Calls(call string) int {
	if s.calls == nil {
		return 0
	}
	return s.calls[call]
}

func (s *Store) UseFixture(call, path string) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return fmt.Errorf("could not read fixture: %v", err)
	}

	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch call {
	case AllowedTopics:
		var topics []*api.Topic
		if topics, err = UnmarshalTopicList(data, jsonpb); err != nil {
			return err
		}

		out := make([]ulid.ULID, 0, len(topics))
		for _, topic := range topics {
			out = append(out, ulids.MustParse(topic.Id))
		}

		s.OnAllowedTopics = func(projectID ulid.ULID) ([]ulid.ULID, error) {
			return out, nil
		}
	case ListTopics:
		var out []*api.Topic
		if out, err = UnmarshalTopicList(data, jsonpb); err != nil {
			return err
		}
		s.OnListTopics = func(projectID ulid.ULID) iterator.TopicIterator {
			return NewTopicIterator(out)
		}
	case RetrieveTopic:
		out := &api.Topic{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnRetrieveTopic = func(ulid.ULID) (*api.Topic, error) {
			return out, nil
		}
	default:
		return fmt.Errorf("unhandled call %q", call)
	}
	return nil
}

func (s *Store) UseError(call string, err error) error {
	switch call {
	case Close:
		s.OnClose = func() error { return err }
	case AllowedTopics:
		s.OnAllowedTopics = func(ulid.ULID) ([]ulid.ULID, error) {
			return nil, err
		}
	case ListTopics:
		s.OnListTopics = func(ulid.ULID) iterator.TopicIterator {
			return NewErrorIterator(err)
		}
	case CreateTopic:
		s.OnCreateTopic = func(*api.Topic) error { return err }
	case RetrieveTopic:
		s.OnRetrieveTopic = func(ulid.ULID) (*api.Topic, error) { return nil, err }
	case UpdateTopic:
		s.OnUpdateTopic = func(*api.Topic) error { return err }
	case DeleteTopic:
		s.OnDeleteTopic = func(ulid.ULID) error { return err }
	default:
		return fmt.Errorf("unhandled call %q", call)
	}
	return nil
}

func (s *Store) Close() error {
	s.incrCalls(Close)
	if s.OnClose != nil {
		return s.OnClose()
	}
	return nil
}

func (s *Store) ReadOnly() bool {
	s.incrCalls(ReadOnly)
	if s.OnReadOnly != nil {
		return s.OnReadOnly()
	}
	return s.readonly
}

func (s *Store) AllowedTopics(projectID ulid.ULID) ([]ulid.ULID, error) {
	s.incrCalls(AllowedTopics)
	if s.OnAllowedTopics != nil {
		return s.OnAllowedTopics(projectID)
	}
	return nil, nil
}

func (s *Store) ListTopics(projectID ulid.ULID) iterator.TopicIterator {
	s.incrCalls(ListTopics)
	return s.OnListTopics(projectID)
}

func (s *Store) CreateTopic(topic *api.Topic) error {
	s.incrCalls(CreateTopic)
	if s.OnCreateTopic != nil {
		return s.OnCreateTopic(topic)
	}
	return errors.New("mock database cannot create topic")
}

func (s *Store) RetrieveTopic(topicID ulid.ULID) (*api.Topic, error) {
	s.incrCalls(RetrieveTopic)
	if s.OnRetrieveTopic != nil {
		return s.OnRetrieveTopic(topicID)
	}
	return nil, errors.New("mock database cannot retrieve topic")
}

func (s *Store) UpdateTopic(topic *api.Topic) error {
	s.incrCalls(UpdateTopic)
	if s.OnUpdateTopic != nil {
		return s.OnUpdateTopic(topic)
	}
	return errors.New("mock database cannot update topic")
}

func (s *Store) DeleteTopic(topicID ulid.ULID) error {
	s.incrCalls(DeleteTopic)
	if s.OnDeleteTopic != nil {
		return s.OnDeleteTopic(topicID)
	}
	return errors.New("mock database cannot delete topic")
}

func (s *Store) ListTopicNames(projectID ulid.ULID) iterator.TopicNamesIterator {
	s.incrCalls(ListTopicNames)
	return s.OnListTopicNames(projectID)
}

func (s *Store) TopicExists(in *api.TopicName) (*api.TopicExistsInfo, error) {
	s.incrCalls(TopicExists)
	if s.OnTopicExists != nil {
		return s.OnTopicExists(in)
	}
	return nil, errors.New("mock database cannot check if topic exists")
}

func (s *Store) incrCalls(call string) {
	if s.calls == nil {
		s.calls = make(map[string]int)
	}
	s.calls[call]++
}
