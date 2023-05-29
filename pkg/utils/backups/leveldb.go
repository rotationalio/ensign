package backups

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDB implements a single leveldb backup that takes a snapshot of the database and
// writes it to a second leveldb database at the temporary directory location.
type LevelDB struct {
	DB *leveldb.DB
}

var _ Backup = &LevelDB{}

// Backup executes the leveldb backup strategy.
func (l *LevelDB) Backup(tmpdir string) (err error) {
	// Open a second leveldb database at the backup location
	var arcdb *leveldb.DB
	if arcdb, err = leveldb.OpenFile(tmpdir, nil); err != nil {
		return fmt.Errorf("could not open archive database: %w", err)
	}

	// Copy all recrods to the archive database
	var narchived uint64
	if narchived, err = CopyLevelDB(l.DB, arcdb); err != nil {
		arcdb.Close()
		return fmt.Errorf("could not write all records to archive database, wrote %d records: %s", narchived, err)
	}
	log.Debug().Uint64("records", narchived).Msg("leveldb archive complete")

	// Close the archived database
	if err = arcdb.Close(); err != nil {
		return fmt.Errorf("could not close archive db: %w", err)
	}

	return nil
}

func CopyLevelDB(src, dst *leveldb.DB) (ncopied uint64, err error) {
	// Create a new batch write to the destination database, writing every 100 records
	// as we iterate over all of the data in the source database.
	var nrows uint64
	batch := new(leveldb.Batch)
	iter := src.NewIterator(nil, nil)
	for iter.Next() {
		nrows++
		batch.Put(iter.Key(), iter.Value())

		if nrows%100 == 0 {
			if err = dst.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
				return ncopied, fmt.Errorf("could not write next 100 rows after %d rows: %s", ncopied, err)
			}
			batch.Reset()
			ncopied += 100
		}
	}

	// Release the iterator and check for errors, just in case we didn't write anything
	iter.Release()
	if err = iter.Error(); err != nil {
		return ncopied, fmt.Errorf("could not iterate over GDS store: %s", err)
	}

	// Write final rows to the database
	if err = dst.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
		return ncopied, fmt.Errorf("could not write final %d rows after %d rows: %s", nrows-ncopied, ncopied, err)
	}
	batch.Reset()
	ncopied += (nrows - ncopied)
	return ncopied, nil
}
