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
	stop    chan struct{}
	ticker  *time.Ticker
	running bool
}

// Return a new backup manager ready to be run.
// NOTE: the backup manager is not started on new; it must explicitly be run for the
// backup routine to be effective. This follows the startup/shutdown services model.
func New(conf Config) (*Manager, error) {
	return &Manager{
		conf: conf,
		stop: make(chan struct{}),
	}, nil
}

// Run the main backup manager routine which periodically wakes up and creates a backup
// of the specified database. The backup manager can be started and stopped as necessary.
func (m *Manager) Run() error {
	m.Lock()
	if !m.conf.Enabled {
		return ErrNotEnabled
	}

	// Do not restart the go routine if it is already running
	if m.running {
		return nil
	}

	// TODO: get backup storage and validate that the routine can run.

	m.ticker = time.NewTicker(m.conf.Interval)
	m.running = true
	m.Unlock()

	go m.run()
	return nil
}

// Run the backup loop in a go routine.
func (m *Manager) run() {
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
		if err := m.backup(); err != nil {
			// Do not continue if there was a backup error. This is a critical error
			// since backups are a safety mechanism, therefore log with the fatal level.
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not complete backup")
			continue backups
		}

		// TODO: remove any previous backups that don't meet the keep requirement.

		log.Info().Dur("duration", time.Since(start)).Msg("backup complete")
	}
}

func (m *Manager) backup() error {
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

	// Cleanup
	m.running = false
	m.ticker = nil
	return nil
}
