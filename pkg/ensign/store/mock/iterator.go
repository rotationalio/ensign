package mock

import (
	"strconv"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/syndtr/goleveldb/leveldb"
)

type TopicIterator struct {
	index  int
	err    error
	topics []*api.Topic
}

func NewTopicIterator(topics []*api.Topic) *TopicIterator {
	return &TopicIterator{topics: topics}
}

func NewTopicErrorIterator(err error) *TopicIterator {
	return &TopicIterator{index: -1, err: err}
}

func (t *TopicIterator) Key() []byte {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return nil
	}
	key := meta.TopicKey(t.topics[t.index])
	return key[:]
}

func (t *TopicIterator) Next() bool {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return false
	}

	if t.index+1 < len(t.topics) {
		t.index++
		return true
	}
	return false
}

func (t *TopicIterator) Prev() bool {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return false
	}

	if t.index-1 > 0 {
		t.index--
		return true
	}
	return false
}

func (t *TopicIterator) Error() error {
	return t.err
}

func (t *TopicIterator) Release() {
	t.topics = nil
	t.index = -1
}

func (t *TopicIterator) Topic() (*api.Topic, error) {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return nil, nil
	}
	return t.topics[t.index], nil
}

func (t *TopicIterator) NextPage(in *api.PageInfo) (out *api.TopicsPage, err error) {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return &api.TopicsPage{}, nil
	}

	if len(t.topics) == 0 {
		return &api.TopicsPage{}, nil
	}

	if in.PageSize == 0 {
		in.PageSize = 100
	}

	idx := 0
	if in.NextPageToken != "" {
		if idx, err = strconv.Atoi(in.NextPageToken); err != nil {
			return nil, err
		}
	}

	jdx := idx + int(in.PageSize)
	if jdx >= len(t.topics) {
		jdx = len(t.topics) - 1
	}

	out = &api.TopicsPage{
		Topics: t.topics[idx:jdx],
	}

	if jdx < len(t.topics)-1 {
		out.NextPageToken = strconv.Itoa(jdx + 1)
	}
	return out, nil
}

type EventIterator struct {
	index  int
	err    error
	events []*api.EventWrapper
}

func NewEventIterator(events []*api.EventWrapper) *EventIterator {
	return &EventIterator{events: events}
}

func NewEventErrorIterator(err error) *EventIterator {
	return &EventIterator{index: -1, err: err}
}

func (t *EventIterator) Key() []byte {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return nil
	}

	// TODO: return actual event key
	return t.events[t.index].Id
}

func (t *EventIterator) Next() bool {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return false
	}

	if t.index+1 < len(t.events) {
		t.index++
		return true
	}
	return false
}

func (t *EventIterator) Prev() bool {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return false
	}

	if t.index-1 > 0 {
		t.index--
		return true
	}
	return false
}

func (t *EventIterator) Error() error {
	return t.err
}

func (t *EventIterator) Release() {
	t.events = nil
	t.index = -1
}

func (t *EventIterator) Event() (*api.EventWrapper, error) {
	if t.index < 0 {
		if t.err == nil {
			t.err = leveldb.ErrIterReleased
		}
		return nil, nil
	}
	return t.events[t.index], nil
}
