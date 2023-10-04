package events

import (
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
)

func (s *Store) Indash(topicID ulid.ULID, hash []byte, eventID rlid.RLID) error {
	return nil
}

func (s *Store) LoadIndash(topicID ulid.ULID) iterator.IndashIterator {
	return &IndashIterator{}
}
