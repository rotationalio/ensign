package mock

import (
	"strconv"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/syndtr/goleveldb/leveldb"
)

type MockIterator struct {
	index  int
	err    error
	keys   [][]byte
	values []interface{}
}

func (i *MockIterator) Key() []byte {
	if i.index < -1 {
		if i.err == nil {
			i.err = leveldb.ErrIterReleased
		}
		return nil
	}
	return i.keys[i.index]
}

func (i *MockIterator) Value() (interface{}, error) {
	if i.index < -1 {
		if i.err == nil {
			i.err = leveldb.ErrIterReleased
		}
		return nil, i.err
	}
	return i.values[i.index], nil
}

func (i *MockIterator) Next() bool {
	if i.index < -1 {
		if i.err == nil {
			i.err = leveldb.ErrIterReleased
		}
		return false
	}

	if i.index+1 < len(i.values) {
		i.index++
		return true
	}
	return false
}

func (i *MockIterator) Prev() bool {
	if i.index < -1 {
		if i.err == nil {
			i.err = leveldb.ErrIterReleased
		}
		return false
	}

	if i.index-1 > -1 {
		i.index--
		return true
	}
	return false
}

func (i *MockIterator) Error() error {
	return i.err
}

func (i *MockIterator) Release() {
	i.values = nil
	i.index = -2
}

func (i *MockIterator) Page(in *api.PageInfo) (out []interface{}, token string, err error) {
	if i.index < -1 {
		if i.err == nil {
			i.err = leveldb.ErrIterReleased
		}
		return nil, "", i.err
	}

	if len(i.values) == 0 {
		return i.values, "", nil
	}

	if in.PageSize == 0 {
		in.PageSize = 100
	}

	idx := 0
	if in.NextPageToken != "" {
		if idx, err = strconv.Atoi(in.NextPageToken); err != nil {
			return nil, "", err
		}
	}

	jdx := idx + int(in.PageSize)
	if jdx >= len(i.values) {
		jdx = len(i.values)
	}

	out = i.values[idx:jdx]
	if jdx < len(i.values) {
		token = strconv.Itoa(jdx)
	}

	return out, token, nil
}

type EventIterator struct {
	MockIterator
}

func NewEventIterator(events []*api.EventWrapper) *EventIterator {
	keys := make([][]byte, 0, len(events))
	values := make([]interface{}, 0, len(events))

	for _, event := range events {
		// TODO: append the correct key for the event iterator
		keys = append(keys, event.Id)
		values = append(values, event)
	}

	return &EventIterator{MockIterator{keys: keys, values: values, index: -1}}
}

func NewEventErrorIterator(err error) *EventIterator {
	return &EventIterator{MockIterator{index: -2, err: err}}
}

func (t *EventIterator) Event() (*api.EventWrapper, error) {
	value, err := t.Value()
	if err != nil {
		return nil, err
	}
	return value.(*api.EventWrapper), nil
}

type TopicIterator struct {
	MockIterator
}

func NewTopicIterator(topics []*api.Topic) *TopicIterator {
	keys := make([][]byte, 0, len(topics))
	values := make([]interface{}, 0, len(topics))

	for _, topic := range topics {
		key := meta.TopicKey(topic)
		keys = append(keys, key[:])
		values = append(values, topic)
	}

	return &TopicIterator{MockIterator{keys: keys, values: values, index: -1}}
}

func NewTopicErrorIterator(err error) *TopicIterator {
	return &TopicIterator{MockIterator{index: -2, err: err}}
}

func (t *TopicIterator) Topic() (*api.Topic, error) {
	value, err := t.Value()
	if err != nil {
		return nil, err
	}
	return value.(*api.Topic), nil
}

func (t *TopicIterator) NextPage(in *api.PageInfo) (out *api.TopicsPage, err error) {
	var values []interface{}
	out = &api.TopicsPage{}
	if values, out.NextPageToken, err = t.Page(in); err != nil {
		return out, err
	}

	out.Topics = make([]*api.Topic, 0, len(values))
	for _, value := range values {
		out.Topics = append(out.Topics, value.(*api.Topic))
	}
	return out, nil
}

type TopicNamesIterator struct {
	MockIterator
}

func NewTopicNamesIterator(names []*api.TopicName) *TopicNamesIterator {
	keys := make([][]byte, 0, len(names))
	values := make([]interface{}, 0, len(names))

	for _, name := range names {
		projectID := ulid.MustParse(name.ProjectId)
		key := meta.TopicNameKey(&api.Topic{ProjectId: projectID[:], Name: name.Name})
		keys = append(keys, key[:])
		values = append(values, name)
	}

	return &TopicNamesIterator{MockIterator{keys: keys, values: values, index: -1}}
}

func NewTopicNamesErrorIterator(err error) *TopicNamesIterator {
	return &TopicNamesIterator{MockIterator{index: -2, err: err}}
}

func (t *TopicNamesIterator) TopicName() (*api.TopicName, error) {
	value, err := t.Value()
	if err != nil {
		return nil, err
	}
	return value.(*api.TopicName), nil
}

func (t *TopicNamesIterator) NextPage(in *api.PageInfo) (out *api.TopicNamesPage, err error) {
	var values []interface{}
	out = &api.TopicNamesPage{}
	if values, out.NextPageToken, err = t.Page(in); err != nil {
		return out, err
	}

	out.TopicNames = make([]*api.TopicName, 0, len(values))
	for _, value := range values {
		out.TopicNames = append(out.TopicNames, value.(*api.TopicName))
	}
	return out, nil
}
