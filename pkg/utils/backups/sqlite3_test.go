package backups_test

import (
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rotationalio/ensign/pkg/utils/backups"
	"github.com/rotationalio/ensign/pkg/utils/sqlite"
	"github.com/stretchr/testify/require"
)

func TestSQLite3Backup(t *testing.T) {
	err := checkSQLite3Fixture()
	require.NoError(t, err, "could not create required fixtures")

	// Take a hash to the source database for comparison purposes
	srcsig, err := fileHash(sqlite3Fixture)
	require.NoError(t, err, "could not get file hash of %s", sqlite3Fixture)

	// Open a connection to the source database
	srcDB, err := sql.Open("ensign_sqlite3", sqlite3Fixture)
	require.NoError(t, err, "could not open sqlite3 database at %s", sqlite3Fixture)
	require.NoError(t, srcDB.Ping(), "could not ping src database")
	defer srcDB.Close()

	// Create the backup system to the src database
	backup := &backups.SQLite3{Name: "example.db"}
	backup.DB, _ = sqlite.GetLastConn()

	// Execute the backup
	tmpdir := t.TempDir()
	err = backup.Backup(tmpdir)
	require.NoError(t, err, "could not execute sqlite3 backup")

	// Assert that the backup exists
	dstPath := filepath.Join(tmpdir, "example.db")
	require.FileExists(t, dstPath, "expected the backup database to exist")

	// Take a hash of the backup database to compare to the src sig
	dstsig, err := fileHash(dstPath)
	require.NoError(t, err, "could not get file hash of %s", dstPath)
	require.Equal(t, srcsig, dstsig, "expected the SHA512 signature of the backup to match the source")
}

const (
	sqlite3Fixture = "testdata/sqlite.db"
	sqlite3Schema  = `CREATE TABLE IF NOT EXISTS entries (
							id       INTEGER PRIMARY KEY,
							name     TEXT NOT NULL,
							blob     BLOB NOT NULL,
							created  TEXT NOT NULL,
							modified TEXT NOT NULL
						);`
)

func checkSQLite3Fixture() (err error) {
	return checkFixture(sqlite3Fixture, func(path string) (err error) {
		var db *sql.DB
		if db, err = sql.Open("sqlite3", path); err != nil {
			return err
		}
		defer db.Close()

		var tx *sql.Tx
		if tx, err = db.Begin(); err != nil {
			return err
		}
		defer tx.Rollback()

		// Execute the schema SQL
		if _, err = tx.Exec(sqlite3Schema); err != nil {
			return err
		}

		now := time.Now()
		for i := 0; i < MaxBackupRecords; i++ {

			data := make([]byte, 192)
			if _, err = rand.Read(data); err != nil {
				return err
			}

			created := now.Add(time.Duration(i*750) * time.Millisecond)
			modified := created.Add(time.Duration(rand.Int63n(20000)) * time.Millisecond)

			params := []interface{}{
				sql.Named("name", fmt.Sprintf("%04x", rand.Int63())),
				sql.Named("data", data),
				sql.Named("created", created.Format(time.RFC3339Nano)),
				sql.Named("modified", modified.Format(time.RFC3339Nano)),
			}

			if _, err = tx.Exec("INSERT INTO entries (name, blob, created, modified) VALUES (:name, :data, :created, :modified)", params...); err != nil {
				return err
			}
		}
		return tx.Commit()
	})
}
