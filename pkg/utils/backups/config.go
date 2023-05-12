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

	// The number of previous backup versions to keep.
	Keep int `default:"1"`
}
