package log

import "errors"

var (
	ErrAppendEarlierTerm     = errors.New("cannot append entry in earlier term")
	ErrAppendSmallerIndex    = errors.New("cannot append entry with smaller index")
	ErrAppendSkipIndex       = errors.New("cannot skip index")
	ErrCommitInvalidIndex    = errors.New("cannot commit invalid index")
	ErrIndexAlreadyCommitted = errors.New("index already committed")
	ErrTruncInvalidIndex     = errors.New("cannot truncate invalid index")
	ErrTruncCommittedIndex   = errors.New("cannot truncate already committed index")
	ErrTruncTermMismatch     = errors.New("the first entry being truncated must match expected term")
	ErrSyncRequired          = errors.New("cannot load log from disk without a sync")
)
