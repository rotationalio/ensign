package backups

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Configure the backup manager routine to ensure it is backing up the correct databases
// at the correct interval; removing old backups as necessary and storing the backups in
// the correct storage location.
type Config struct {
	// If false, the backup manager will not run.
	Enabled bool `default:"false"`

	// The interval between backups.
	Interval time.Duration `default:"24h"`

	// The path to a local disk directory to store compressed backups, e.g.
	// file:///rel/path/ or a cloud location such as s3://bucket.
	StorageDSN string `split_words:"true" required:"false"`

	// Temporary directory to perform local backup to. If not set, then the OS tmpdir
	// is used. Backups are generated in this folder and then moved to the storage
	// location. This folder needs enough space to contain the complete backup.
	TempDir string `required:"false"`

	// Prefix specifies the filename of the backup, e.g.g prefix-200601021504.tgz
	Prefix string `default:"backup"`

	// The number of previous backup versions to keep.
	Keep int `default:"1"`
}

// Validate the Config
func (c Config) Validate() error {
	if c.Enabled {
		if c.StorageDSN == "" {
			return errors.New("invalid backup configuration: storage dsn is required")
		}
	}
	return nil
}

// Storage returns the storage configuration specified by the storage DSN.
func (c Config) Storage() (_ Storage, err error) {
	// Parse the DSN specified by the user
	var dsn *url.URL
	if dsn, err = url.Parse(c.StorageDSN); err != nil {
		return nil, err
	}

	if dsn.Scheme == "" || dsn.Path == "" {
		return nil, ErrInvalidStorageDSN
	}

	// Normalization
	scheme := strings.ToLower(dsn.Scheme)
	path := strings.TrimPrefix(dsn.Path, "/")

	// Based on the scheme return the Storage adapter
	switch scheme {
	case "file":
		return NewFileStorage(path, c.Prefix)
	case "inmem":
		return &MemoryStorage{root: path}, nil
	default:
		return nil, fmt.Errorf("invalid backup storage dsn: unknown scheme %q", scheme)
	}
}

// ArchiveName returns the name of the next archive using the current timestamp.
func (c Config) ArchiveName() string {
	prefix := c.Prefix
	if prefix == "" {
		prefix = "backup"
	}
	return fmt.Sprintf("%s-%s.tgz", prefix, time.Now().UTC().Format("200601021504"))
}
