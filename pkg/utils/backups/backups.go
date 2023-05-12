package backups

import (
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
		}

		// Begin the backups process
		start := time.Now()
		log.Debug().Msg("starting backup")

		// Perform the backup
		if err := m.backup.Backup(""); err != nil {
			// Do not continue if there was a backup error. This is a critical error
			// since backups are a safety mechanism, therefore log with the fatal level.
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not complete backup")
			continue backups
		}

		// TODO: remove any previous backups that don't meet the keep requirement.

		log.Info().Dur("duration", time.Since(start)).Msg("backup complete")
	}
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
