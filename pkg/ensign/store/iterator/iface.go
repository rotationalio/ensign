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

// EventIterator allows access to events in the database
type EventIterator interface {
	Iterator
	Event() (*api.EventWrapper, error)
}

// TopicIterator allows access to Topic models in the database
type TopicIterator interface {
	Iterator
	Topic() (*api.Topic, error)
	NextPage(in *api.PageInfo) (*api.TopicsPage, error)
}

// TopicIterator allows access to Topic names index in the database
type TopicNamesIterator interface {
	Iterator
	TopicName() (*api.TopicName, error)
	NextPage(in *api.PageInfo) (*api.TopicNamesPage, error)
}

// GroupIterator allows access to ConsumerGroup models in the datbase
type GroupIterator interface {
	Iterator
	Group() (*api.ConsumerGroup, error)
}

// Paginator iterators allow the fetching of multiple items at a time. Used primarily
// for testing paginated interfaces, the NextPage() methods are used in production.
type Paginator interface {
	Page(*api.PageInfo) ([]interface{}, string, error)
}

// Valuer interfaces fetch the item at the cursor as an interface. Used primarily for
// testing iterators, the type-specific methods are used in production.
type Valuer interface {
	Value() (interface{}, error)
}
