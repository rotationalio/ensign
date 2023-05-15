package backups

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Manager runs as an independent service which periodically backs up a database to a
// compressed backup location on disk or to cloud storage. At a specified interval, the
// manager runs a db-specific backup routine that generally clones the database to a tmp
// directory; it then compresses the clone and moves it to a final backup location.
type Manager struct {
	sync.Mutex
	conf    Config
	backup  Backup
	stop    chan struct{}
	ticker  *time.Ticker
	running bool
}

// Backup is an interface that enables different types of database backups, e.g. for
// backuping up a leveldb database (or multiple leveldb databases) vs a sqlite database.
// It is expected that the backup is written to the temporary directory specified by
// the input variable. It is unnecessary to compress the contents of the backup as the
// backup manager will compress the contents before transferring them to the backup
// storage location. Once the backup is complete (error or not) the tmpdir is removed.
type Backup interface {
	Backup(tmpdir string) error
}

// Return a new backup manager ready to be run.
// NOTE: the backup manager is not started on new; it must explicitly be run for the
// backup routine to be effective. This follows the startup/shutdown services model.
func New(conf Config, backup Backup) *Manager {
	return &Manager{
		conf:   conf,
		backup: backup,
	}
}

// Run the main backup manager routine which periodically wakes up and creates a backup
// of the specified database. The backup manager can be started and stopped as necessary.
func (m *Manager) Run() (err error) {
	m.Lock()
	defer m.Unlock()
	if !m.conf.Enabled {
		return ErrNotEnabled
	}

	// Do not restart the go routine if it is already running
	if m.running {
		return nil
	}

	// Get backup storage and validate that the routine can run.
	var storage Storage
	if storage, err = m.conf.Storage(); err != nil {
		return err
	}

	// Check that temporary directories can be created.
	if path, err := m.MkdirTemp(); err != nil {
		return ErrTmpDirUnavailable
	} else {
		os.Remove(path)
	}

	m.ticker = time.NewTicker(m.conf.Interval)
	m.running = true
	m.stop = make(chan struct{})

	go m.run(storage)
	return nil
}

// Run the backup loop in a go routine.
func (m *Manager) run(storage Storage) {
backups:
	for {
		// Wait for next tick or a stop message
		select {
		case <-m.stop:
			log.Debug().Msg("backup manager received stop signal")
			return
		case <-m.ticker.C:
			if err := m.runOnce(storage); err != nil {
				// Backup errors are considered critical since backups are a safety
				// mechanism; log with a fatal level, but do not terminate the loop in
				// case the error is transient (e.g. disk full).
				log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not complete backup")
				continue backups
			}
		}
	}
}

// Run one instance of the backup loop
func (m *Manager) runOnce(storage Storage) (err error) {
	// Begin the backups process
	start := time.Now()
	log.Debug().Msg("starting backup")

	// Create a temporary directory; note that this must be cleaned up when done!
	var tmpdir string
	if tmpdir, err = m.MkdirTemp(); err != nil {
		// Do not continue with backup if we cannot create a tmpdir.
		return err
	}

	// Ensure that the temporary directory is cleaned up when done.
	defer os.RemoveAll(tmpdir)

	// Perform the backup
	if err = m.backup.Backup(tmpdir); err != nil {
		// Do not continue if there was a backup error.
		return err
	}

	// Write the archive to the backup storage directory.
	archive := m.conf.ArchiveName()

	// Open the file in storage to begin writing to
	var w io.WriteCloser
	if w, err = storage.Open(archive); err != nil {
		return err
	}

	// Archive the contents of the tmpdir
	if err = m.archive(tmpdir, w); err != nil {
		return err
	}

	// Remove any previous backups that don't meet the keep requirement.
	if err = m.cleanup(storage); err != nil {
		return err
	}

	log.Info().Dur("duration", time.Since(start)).Msg("backup complete")
	return nil
}

// Writes the contents of the directory at the path dir as gzip to w.
func (m *Manager) archive(dir string, w io.WriteCloser) (err error) {
	// Prepare for gzip compressing the archive
	defer w.Close()
	gw := gzip.NewWriter(w)
	defer gw.Close()

	// Create a tar file.
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Walk the archive and write to the tar file.
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var hdr *tar.Header
		if hdr, err = tar.FileInfoHeader(info, ""); err != nil {
			return err
		}

		hdr.Name = path[len(dir):]
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var tmp *os.File
		if tmp, err = os.Open(path); err != nil {
			return err
		}
		defer tmp.Close()

		if _, err = io.Copy(tw, tmp); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (m *Manager) cleanup(storage Storage) (err error) {
	var archives []string
	if archives, err = storage.ListArchives(); err != nil {
		return err
	}

	if len(archives) > m.conf.Keep {
		var removed int
		defer log.Debug().Int("kept", m.conf.Keep).Int("removed", removed).Msg("backup storage cleaned up")

		for _, archive := range archives[:len(archives)-m.conf.Keep] {
			log.Debug().Str("archive", archive).Msg("deleting archive")
			if err = storage.Remove(archive); err != nil {
				return err
			}
			removed++
		}
	}
	return nil
}

func (m *Manager) Shutdown() error {
	m.Lock()
	defer m.Unlock()

	// Do not shutdown if we're not running
	if !m.running {
		return nil
	}

	// Send stop signals
	m.ticker.Stop()
	m.stop <- struct{}{}
	close(m.stop)

	// Cleanup
	m.running = false
	m.ticker = nil
	m.stop = nil
	return nil
}

// Create a temporary direcory in the configured path. If no configured directory is
// specified then os.MkdirTemp is used. It is the callers responsibility to cleanup the
// directory that was created.
func (m *Manager) MkdirTemp() (path string, err error) {
	return os.MkdirTemp(m.conf.TempDir, "backup-*")
}
