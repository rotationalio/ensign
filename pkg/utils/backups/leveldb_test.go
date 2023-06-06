package backups_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/backups"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func TestLevelDBBackup(t *testing.T) {
	err := checkLevelDBFixture()
	require.NoError(t, err, "could not create required fixtures")

	// Extract the src database to a temporary directory
	srcPath := t.TempDir()
	_, err = extract(leveldbFixture, srcPath, false)
	require.NoError(t, err, "could not extract %s to %s", leveldbFixture, srcPath)

	// Open a connection to the source database
	srcDB, err := leveldb.OpenFile(srcPath, nil)
	require.NoError(t, err, "could not open leveldb at %s", srcPath)
	defer srcDB.Close()

	// Create the backup system to the src database
	backup := &backups.LevelDB{DB: srcDB}

	// Execute the backup
	dstPath := t.TempDir()
	err = backup.Backup(dstPath)
	require.NoError(t, err, "could not execute leveldb backup")

	// Assert that the backup exists
	require.FileExists(t, filepath.Join(dstPath, "MANIFEST-000000"), "expected log file to exist")

	// Unfortunately hashing an archive doesn't work as a comparison method for leveldb
	// databases, so instead of using the file hash method, we have to compare the data
	// directly from the src database to the destination database.
	dstDB, err := leveldb.OpenFile(dstPath, nil)
	require.NoError(t, err, "could not open leveldb at %s", dstPath)
	defer dstDB.Close()

	// Count the number of rows in the dst db
	nDstRows, err := countRows(dstDB)
	require.NoError(t, err, "could not count dstDB rows")
	require.Equal(t, uint64(MaxBackupRecords), nDstRows, "number of rows in dst db unexpected")

	// Compare the key/val pairs in the dstdb with the srcdb
	compareLevelDB(t, srcDB, dstDB)
}

const leveldbFixture = "testdata/leveldb.tgz"

var keySizes = []int{16, 32, 38, 64}

func checkLevelDBFixture() error {
	return checkFixture(leveldbFixture, func(path string) (err error) {
		// Create a temporary directory to write the fixture into.
		var tmpdir string
		if tmpdir, err = os.MkdirTemp("", "leveldb-fixture-*"); err != nil {
			return err
		}
		defer os.RemoveAll(tmpdir)

		var db *leveldb.DB
		if db, err = leveldb.OpenFile(tmpdir, nil); err != nil {
			return err
		}
		defer db.Close()

		for i := 0; i < MaxBackupRecords; i++ {
			value := make([]byte, 192)
			if _, err = rand.Read(value); err != nil {
				return err
			}

			key := make([]byte, keySizes[rand.Intn(len(keySizes))])
			if _, err = rand.Read(key); err != nil {
				return err
			}

			if err = db.Put(key, value, &opt.WriteOptions{Sync: true}); err != nil {
				return err
			}
		}

		// Compact the database
		if err = db.CompactRange(util.Range{}); err != nil {
			return err
		}
		db.Close()

		// Archive the fixture to the specified path
		return archive(tmpdir, path)
	})
}

// Helper to count the number of rows in a leveldb database
func countRows(db *leveldb.DB) (nrows uint64, err error) {
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		nrows++
	}
	return nrows, iter.Error()
}

// Helper to compare the rows from the dst database to the src database
func compareLevelDB(t *testing.T, target, source *leveldb.DB) {
	iter := source.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		hasKey, err := target.Has(key, nil)
		require.NoError(t, err, "could not check if target db has key")
		require.True(t, hasKey, "target db does not have key from source")

		sval := iter.Value()
		tval, err := target.Get(key, nil)

		require.NoError(t, err, "could not get value for key from target")
		require.Equal(t, sval, tval, "target value does not match source value")
	}
	require.NoError(t, iter.Error(), "iterating source database failed")
}
