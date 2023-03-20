package iterator

import api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"

// Iterators allow memory safe list operations from the Store.
type Iterator interface {
	Key() []byte
	Next() bool
	Prev() bool
	Error() error
	Release()
}

// TopicIterator allows access to Topic models in the database
type TopicIterator interface {
	Iterator
	Topic() (*api.Topic, error)
	NextPage(in *api.PageInfo) (*api.TopicsPage, error)
}

// GroupIterator allows access to ConsumerGroup models in the datbase
type GroupIterator interface {
	Iterator
	Group() (*api.ConsumerGroup, error)
}
