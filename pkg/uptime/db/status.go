package db

import (
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var KeyCurrentStatus = []byte{0, 0, 0, 0, 0, 0, 114, 111, 116, 97, 116, 105, 111, 110, 97, 108}

// Because service status keys are ordered by time using ULIDs, it is possible to fetch
// the last service status for a service ID by using the ID as a prefix to iterate over
// the keys and return the last item; unmarshaling it into the specified model.
func LastServiceStatus(id uuid.UUID, m Model) (err error) {
	iter := NewIterator(util.BytesPrefix(id[:]))
	defer iter.Release()

	if !iter.Last() {
		return ErrNotFound
	}

	data := iter.Value()
	if err = iter.Error(); err != nil {
		return err
	}
	return m.Unmarshal(data)
}
