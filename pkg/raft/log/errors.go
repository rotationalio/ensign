package log

import "errors"

var (
	ErrSyncRequired = errors.New("cannot load log from disk without a sync")
)
