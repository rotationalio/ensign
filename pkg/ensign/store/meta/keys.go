package meta

import (
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

type IndexKey [16]byte
type ObjectKey [32]byte

func CreateIndex(objectID ulid.ULID) (key IndexKey, err error) {
	if ulids.IsZero(objectID) {
		return key, errors.ErrKeyNull
	}
	return IndexKey(objectID), nil
}

func CreateKey(parentID, objectID ulid.ULID) (key ObjectKey, err error) {
	if ulids.IsZero(parentID) || ulids.IsZero(objectID) {
		return key, errors.ErrKeyNull
	}

	if err = parentID.MarshalBinaryTo(key[:16]); err != nil {
		return key, err
	}

	if err = objectID.MarshalBinaryTo(key[16:]); err != nil {
		return key, err
	}
	return key, nil
}

func (k *ObjectKey) Key() IndexKey {
	return IndexKey(*(*[16]byte)(k[16:]))
}

func (k *ObjectKey) UnmarshalValue(data []byte) error {
	if len(data) != 32 {
		return errors.ErrKeyWrongSize
	}
	copy(k[:], data)
	return nil
}

func (k *ObjectKey) ParentID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[:16])
	return id, err
}

func (k *ObjectKey) ObjectID() (id ulid.ULID, err error) {
	err = id.UnmarshalBinary(k[16:])
	return id, err
}
