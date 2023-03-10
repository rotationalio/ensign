package mock

import (
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type topicsIterator struct {
	index  int
	err    error
	topics []*api.Topic
}

func (t *topicsIterator) Next() bool {
	if t.index < 0 {
		t.err = leveldb.ErrIterReleased
		return false
	}

	if t.index+1 < len(t.topics) {
		t.index++
		return true
	}
	return false
}

func (t *topicsIterator) Prev() bool {
	if t.index < 0 {
		t.err = leveldb.ErrIterReleased
		return false
	}

	if t.index-1 > 0 {
		t.index--
		return true
	}
	return false
}

func (t *topicsIterator) Error() error {
	return t.err
}

func (t *topicsIterator) Release() {
	t.topics = nil
	t.index = -1
}

func (t *topicsIterator) Topic() (*api.Topic, error) {
	if t.index < 0 {
		t.err = leveldb.ErrIterReleased
		return nil, nil
	}
	return t.topics[t.index], nil
}

func (t *topicsIterator) NextPage(in *api.PageInfo) (*api.TopicsPage, error) {
	if t.index < 0 {
		t.err = leveldb.ErrIterReleased
		return nil, nil
	}
	return nil, errors.ErrNotImplemented
}
