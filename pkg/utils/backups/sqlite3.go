package backups

import (
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

const (
	SQLite3Main  = "main"
	PagesPerStep = 5
)

// SQLite3 implements a single sqlite3 backup that uses the online backup of a running
// database mechanism to copy the db over to a second sqlite3 database at the temporary
// directory location.
type SQLite3 struct {
	DB *sqlite3.SQLiteConn
}

var _ Backup = &SQLite3{}

// Backup executes the sqlite3 backup strategy.
func (s *SQLite3) Backup(tmpdir string) (err error) {
	var backup *sqlite3.SQLiteBackup
	if backup, err = s.DB.Backup(SQLite3Main, s.DB, SQLite3Main); err != nil {
		return err
	}
	defer backup.Close()

	var remaining int
	for remaining > 0 {
		var done bool
		if done, err = backup.Step(5); err != nil {
			return err
		}

		if done {
			return
		}

		remaining = backup.Remaining()
		time.Sleep(250 * time.Millisecond)
	}

	return nil
}
