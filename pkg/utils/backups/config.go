package backups

import "time"

// Configure the backup manager routine to ensure it is backing up the correct databases
// at the correct interval; removing old backups as necessary and storing the backups in
// the correct storage location.
type Config struct {
	// If false, the backup manager will not run.
	Enabled bool `default:"true"`

	// The interval between backups.
	Interval time.Duration `default:"24h"`

	// The path to a local disk directory to store compressed backups, e.g.
	// file:///rel/path/ or a cloud location such as s3://bucket.
	StorageDSN string `split_words:"true" required:"true"`

	// Temporary directory to perform local backup to. If not set, then the OS tmpdir
	// is used. Backups are generated in this folder and then moved to the storage
	// location. This folder needs enough space to contain the complete backup.
	TempDir string `required:"false"`

	// Prefix specifies the filename of the backup, e.g.g prefix-200601021504.tgz
	Prefix string `default:"backup"`

	// The number of previous backup versions to keep.
	Keep int `default:"1"`
}

// Storage returns the storage configuration specified by the storage DSN.
func (c Config) Storage() (Storage, error) {
	return nil, nil
}
