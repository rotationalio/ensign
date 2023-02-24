package db

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Transaction struct {
	leveldb.Transaction
}

func (t *Transaction) Get(key []byte, m Model) (err error) {
	data, err := t.Transaction.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, leveldb.ErrClosed) {
			return ErrNotConnected
		}
	}
	return m.Unmarshal(data)
}

func (t *Transaction) Put(m Model) (err error) {
	// Marshal the model
	var key, value []byte
	if key, err = m.Key(); err != nil {
		return err
	}

	if value, err = m.Marshal(); err != nil {
		return err
	}

	if err = t.Transaction.Put(key, value, nil); err != nil {
		if errors.Is(err, leveldb.ErrClosed) {
			return ErrNotConnected
		}
	}
	return err
}

func (t *Transaction) Delete(m Model) (err error) {
	var key []byte
	if key, err = m.Key(); err != nil {
		return err
	}

	if err = t.Transaction.Delete(key, nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return ErrNotFound
		}

		if errors.Is(err, leveldb.ErrClosed) {
			return ErrNotConnected
		}
	}
	return err
}

func (t *Transaction) NewIterator(slice *util.Range) iterator.Iterator {
	iter := t.Transaction.NewIterator(slice, nil)
	return iter
}
