package iterator

import api "github.com/rotationalio/go-ensign/api/v1beta1"

// Iterators allow memory safe list operations from the Store.
type Iterator interface {
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
