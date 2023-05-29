package backups

import "errors"

var (
	ErrNotEnabled        = errors.New("the backup manager is not enabled")
	ErrTmpDirUnavailable = errors.New("cannot create temporary directories for backup")
	ErrInvalidStorageDSN = errors.New("could not parse storage dsn, specify scheme:///relative/path/")
	ErrNotADirectory     = errors.New("incorrectly configured: backup storage is not a directory")
	ErrNilSQLite3Conn    = errors.New("could not fetch underlying sqlite3 connection for backup")
)
