package events

import (
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
)

// Implements iterator.EventIterator to access to a sequence of events in a topic.
type EventIterator struct {
	ldbiter.Iterator
}

func (i *EventIterator) Event() (*api.EventWrapper, error) {
	return nil, nil
}
