package mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"google.golang.org/protobuf/encoding/protojson"
)

// Constants are used to reference store methods in mock code
const (
	Close         = "Close"
	ReadOnly      = "ReadOnly"
	ListTopics    = "ListTopics"
	CreateTopic   = "CreateTopic"
	RetrieveTopic = "RetrieveTopic"
	UpdateTopic   = "UpdateTopic"
	DeleteTopic   = "DeleteTopic"
)

// Implements both a store.EventStore and a store.MetaStore for testing purposes.
type Store struct {
	readonly        bool
	calls           map[string]int
	OnClose         func() error
	OnReadOnly      func() bool
	OnListTopics    func(ulid.ULID) iterator.TopicIterator
	OnCreateTopic   func(*api.Topic) error
	OnRetrieveTopic func(topicID ulid.ULID) (*api.Topic, error)
	OnUpdateTopic   func(*api.Topic) error
	OnDeleteTopic   func(topicID ulid.ULID) error
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
	s.OnListTopics = nil
	s.OnCreateTopic = nil
	s.OnRetrieveTopic = nil
	s.OnUpdateTopic = nil
	s.OnDeleteTopic = nil
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
	case ListTopics:
		items := make([]interface{}, 0)
		if err = json.Unmarshal(data, &items); err != nil {
			return fmt.Errorf("could not json unmarshal fixture: %v", err)
		}

		out := make([]*api.Topic, 0, len(items))
		for _, item := range items {
			var buf []byte
			if buf, err = json.Marshal(item); err != nil {
				return err
			}

			topic := &api.Topic{}
			if err = jsonpb.Unmarshal(buf, topic); err != nil {
				return err
			}
			out = append(out, topic)
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

func (s *Store) incrCalls(call string) {
	if s.calls == nil {
		s.calls = make(map[string]int)
	}
	s.calls[call]++
}
