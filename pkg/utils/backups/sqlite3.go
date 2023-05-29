package backups

import (
	"database/sql"
	"path/filepath"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/rotationalio/ensign/pkg/utils/sqlite"
	"github.com/rs/zerolog/log"
)

const (
	SQLite3Main   = "main"
	PagesPerStep  = 5
	StepSleep     = 50 * time.Millisecond
	DefaultDBName = "backup.db"
)

// SQLite3 implements a single sqlite3 backup that uses the online backup of a running
// database mechanism to copy the db over to a second sqlite3 database at the temporary
// directory location.
type SQLite3 struct {
	DB   *sqlite.Conn
	Name string
}

var _ Backup = &SQLite3{}

// Backup executes the sqlite3 backup strategy.
func (s *SQLite3) Backup(tmpdir string) (err error) {
	var (
		dstDB   *sql.DB
		dstConn *sqlite.Conn
	)

	// Open a second sqlite3 database at the backup location.
	if dstDB, dstConn, err = s.OpenDestDB(tmpdir); err != nil {
		return err
	}

	// Ensure the database connection is closed when the backup is complete; this will
	// also close the underlying sqlite3 connection.
	defer dstDB.Close()

	// Create the backup manager into the destination db from the src connection.
	// NOTE: backup.Finish() MUST be called to prevent panics.
	var backup *sqlite3.SQLiteBackup
	if backup, err = dstConn.Backup(SQLite3Main, s.DB, SQLite3Main); err != nil {
		return err
	}

	// Execute the backup copying the specified number of pages at each step then
	// sleeping to allow concurrent transactions to acquire write locks. This will
	// increase the amount of backup time but preserve normal operations. This means
	// that backups will be most successful during low-volume times.
	var isDone bool
	for !isDone {
		// Backing up a smaller number of pages per step is the most effective way of
		// doing online backups and also allow write transactions to make progress.
		if isDone, err = backup.Step(PagesPerStep); err != nil {
			backup.Finish()
			return err
		}

		log.Debug().
			Int("remaining", backup.Remaining()).
			Int("page_count", backup.PageCount()).
			Msg("sqlite3 backup step")

		// This sleep allows other transactions to write during backups.
		time.Sleep(StepSleep)
	}
	return backup.Finish()
}

func (s *SQLite3) OpenDestDB(tmpdir string) (db *sql.DB, conn *sqlite.Conn, err error) {
	var path string
	if s.Name != "" {
		path = filepath.Join(tmpdir, s.Name)
	} else {
		path = filepath.Join(tmpdir, DefaultDBName)
	}

	if db, err = sql.Open(sqlite.DriverName, path); err != nil {
		return nil, nil, err
	}

	// Ping the database in order to establish a connection
	if err = db.Ping(); err != nil {
		return nil, nil, err
	}

	// Get the last database connection from the sqlite package to establish a backup
	var ok bool
	if conn, ok = sqlite.GetLastConn(); !ok {
		return nil, nil, ErrNilSQLite3Conn
	}

	return db, conn, nil
}
