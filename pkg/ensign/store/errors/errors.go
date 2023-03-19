/*
Package errors implements standard database read/write errors for the store package.
*/
package errors

import (
	"errors"
	"fmt"

	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrReadOnly         = &Error{"database is readonly: cannot perform operation", leveldb.ErrReadOnly}
	ErrNotFound         = &Error{"object not found", leveldb.ErrNotFound}
	ErrClosed           = &Error{"database is closed: cannot perform operation", leveldb.ErrClosed}
	ErrIterReleased     = &Error{"iterator released", leveldb.ErrIterReleased}
	ErrSnapshotReleased = &Error{"snapshot released", leveldb.ErrSnapshotReleased}
	ErrAlreadyExists    = errors.New("object with specified key already exists")
	ErrNotImplemented   = errors.New("this method has not been implemented yet")

	ErrInvalidTopic          = errors.New("invalid topic")
	ErrTopicMissingProjectId = &Error{"missing project_id field", ErrInvalidTopic}
	ErrTopicInvalidProjectId = &Error{"cannot parse project_id field", ErrInvalidTopic}
	ErrTopicMissingName      = &Error{"missing name field", ErrInvalidTopic}
	ErrTopicMissingId        = &Error{"missing id field", ErrInvalidTopic}
	ErrTopicInvalidId        = &Error{"cannot parse id field", ErrInvalidTopic}
	ErrTopicInvalidCreated   = &Error{"invalid created field", ErrInvalidTopic}
	ErrTopicInvalidModified  = &Error{"invalid modified field", ErrInvalidTopic}

	ErrInvalidGroup          = errors.New("invalid group")
	ErrGroupMissingProjectId = &Error{"missing project_id field", ErrInvalidGroup}
	ErrGroupInvalidProjectId = &Error{"cannot parse project_id field", ErrInvalidGroup}
	ErrGroupMissingId        = &Error{"missing id field", ErrInvalidGroup}
	ErrGroupMissingKeyField  = &Error{"missing one of id or name fields", ErrInvalidGroup}
	ErrGroupInvalidCreated   = &Error{"invalid created field", ErrInvalidGroup}
	ErrGroupInvalidModified  = &Error{"invalid modified field", ErrInvalidGroup}

	ErrInvalidKey   = errors.New("invalid object key")
	ErrKeyWrongSize = &Error{"incorrect key size", ErrInvalidKey}
	ErrKeyNull      = &Error{"no part of the key can be zero-valued", ErrInvalidKey}

	ErrInvalidPage      = errors.New("invalid page token")
	ErrPageTokenExpired = &Error{"next page token has expired", ErrInvalidPage}
	ErrInvalidPageToken = &Error{"invalid next page token field", ErrInvalidPage}
)

func Wrap(err error) error {
	switch {
	case errors.Is(err, leveldb.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, leveldb.ErrReadOnly):
		return ErrReadOnly
	case errors.Is(err, leveldb.ErrClosed):
		return ErrClosed
	case errors.Is(err, leveldb.ErrIterReleased):
		return ErrIterReleased
	case errors.Is(err, leveldb.ErrSnapshotReleased):
		return ErrSnapshotReleased
	case errors.Is(err, pagination.ErrCursorExpired):
		return ErrPageTokenExpired
	case errors.Is(err, pagination.ErrUnparsableToken), errors.Is(err, pagination.ErrTokenQueryMismatch), errors.Is(err, pagination.ErrMissingExpiration):
		return ErrInvalidPageToken
	}

	return &Error{fmt.Sprintf("unhandled store exception occurred: %s", err), err}
}

type Error struct {
	msg  string
	ldbe error
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Is(target error) bool {
	return errors.Is(e.ldbe, target)
}

func (e *Error) Unwrap() error {
	return e.ldbe
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
