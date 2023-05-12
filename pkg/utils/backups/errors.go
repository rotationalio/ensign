package backups

import "errors"

var (
	ErrNotEnabled = errors.New("the backup manager is not enabled")
)
