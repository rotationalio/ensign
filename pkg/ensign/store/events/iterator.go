package events

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
	"google.golang.org/protobuf/proto"
)

// Implements iterator.EventIterator to access to a sequence of events in a topic.
type EventIterator struct {
	ldbiter.Iterator
	topicID ulid.ULID
}

func (i *EventIterator) Event() (*api.EventWrapper, error) {
	event := &api.EventWrapper{}
	if err := proto.Unmarshal(i.Value(), event); err != nil {
		return nil, err
	}
	return event, nil
}

func (t *EventIterator) Seek(eventID rlid.RLID) bool {
	key, err := CreateKey(t.topicID, eventID, EventSegment)
	if err != nil {
		return false
	}
	return t.Iterator.Seek(key[:])
}

func (i *EventIterator) Error() error {
	if err := i.Iterator.Error(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Implements iterator.EventIterator to return an error.
type EventErrorIterator struct {
	errors.ErrorIterator
}

func (i *EventErrorIterator) Event() (*api.EventWrapper, error) {
	return nil, i.Error()
}

func (i *EventErrorIterator) Seek(eventID rlid.RLID) bool {
	return false
}

// Implements iterator.IndashIterator to access a sequence of hashes for a topic.
type IndashIterator struct {
	ldbiter.Iterator
}

// Hash returns the key with the prefixed topic and segment stripped. If the key is not
// the right length then an error is returned but no other parsing happens.
func (i *IndashIterator) Hash() ([]byte, error) {
	key := i.Key()
	if len(key) < 18 {
		return nil, errors.ErrKeyWrongSize
	}

	// Must make a copy of the data because the key byte array may change in iteration
	hash := make([]byte, len(key)-18)
	copy(hash, key[18:])
	return hash, nil
}

// Error wraps an internal leveldb error so that ensign users can test against it.
func (i *IndashIterator) Error() error {
	if err := i.Iterator.Error(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Implements iterator.IndashIterator to return an error.
type IndashErrorIterator struct {
	errors.ErrorIterator
}

func (i *IndashErrorIterator) Hash() ([]byte, error) {
	return nil, i.Error()
}
