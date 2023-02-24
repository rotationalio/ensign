package uptime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rotationalio/ensign/pkg/uptime/db"
	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/rotationalio/ensign/pkg/uptime/incident"
	"github.com/rotationalio/ensign/pkg/uptime/services"
	"github.com/rs/zerolog/log"
)

const defaultTimeout = 5 * time.Second

type signal struct{}

// Monitor wraps several service monitors to routinely conduct health checks and save
// the status (if changed) to the database.
type Monitor struct {
	sync.Mutex
	monitors []health.Monitor
	services *services.Info
	interval time.Duration
	stop     chan<- signal
	done     <-chan signal
	running  bool
}

// NewMonitor loads the services definition from the path on disk and creates monitors
// for each of them to perform heartbeat checks on the specified interval.
func NewMonitor(interval time.Duration, infoPath string) (mon *Monitor, err error) {
	mon = &Monitor{
		interval: interval,
	}

	// Load the service info from disk
	if mon.services, err = services.Load(infoPath); err != nil {
		return nil, err
	}

	// Create monitors from each of the services
	mon.monitors = make([]health.Monitor, 0, mon.services.Len())
	for _, info := range mon.services.Services() {
		var monitor health.Monitor
		switch info.Type {
		case services.APIServiceType:
			if monitor, err = health.NewAPIMonitor(info.Endpoint); err != nil {
				return nil, err
			}
		case services.HTTPServiceType:
			if monitor, err = health.NewHTTPMonitor(info.Endpoint); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown monitor type %q", info.Type)
		}
		mon.monitors = append(mon.monitors, monitor)
	}
	return mon, nil
}

// Start the monitor background routine.
func (m *Monitor) Start() {
	m.Lock()
	defer m.Unlock()

	if m.running {
		return
	}

	stop := make(chan signal, 1)
	done := make(chan signal, 1)

	go m.Run(stop, done)
	m.stop = stop
	m.done = done
	m.running = true
}

// Stop the monitor background routine (blocks until shutdown is complete).
func (m *Monitor) Stop(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	if !m.running {
		return nil
	}

	m.stop <- signal{}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.done:
	}

	m.stop = nil
	m.done = nil
	m.running = false
	log.Debug().Msg("uptime monitor stopped")
	return nil
}

// The background routine that executes every interval; it has self contained signaling.
func (m *Monitor) Run(stop <-chan signal, done chan<- signal) {
	// Ensure the done signal is sent when this go routine exits.
	defer func() {
		done <- signal{}
	}()

	log.Info().Dur("interval", m.interval).Int("monitors", len(m.monitors)).Msg("uptime monitor started")

	for {
		// Use a wait channel rather than a ticker so that monitoring delays and time
		// don't cause back pressure and reduce the deterministic nature of pings.
		wait := time.After(m.interval)

		// Wait for interval or until a stop signal is received
		select {
		case <-stop:
			return
		case <-wait:
		}

		// All errors and handling should happen in the RunChecks method.
		if err := m.RunChecks(); err != nil {
			log.Error().Err(err).Msg("could not run uptime monitor checks")
		}
	}
}

func (m *Monitor) RunChecks() error {
	if len(m.monitors) == 0 {
		return errors.New("no monitors have been configured")
	}

	nerrors := 0
	for i, service := range m.services.Services() {
		monitor := m.monitors[i]
		if err := m.CheckStatus(monitor, service); err != nil {
			log.Error().Err(err).Str("service", service.Title).Msg("could not check status for service")
			nerrors++
		}
	}

	// Update the current service status on disk
	if err := db.Put(m.services); err != nil {
		log.Error().Err(err).Msg("could not update current status on disk")
	}

	log.Debug().Int("monitors", len(m.monitors)).Int("errors", nerrors).Msg("uptime monitor checks complete")
	return nil
}

// Check status executes the service status for the specified monitor, saves the status
// to disk if it has changed and creates any incidents if required. It then updates the
// service pointer to save the current status after all monitor checks have happened.
func (m *Monitor) CheckStatus(monitor health.Monitor, service *services.Service) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Get current status
	var status health.ServiceStatus
	if status, err = monitor.Status(ctx); err != nil {
		return err
	}

	// Set the service ID on the status
	status.SetServiceID(service.ID)

	// Update the service with the status so that the current status is saved
	service.Status = status.Status()
	service.LastUpdate = status.CheckedAt()

	// Load previous status from the database if it exists
	var prev health.ServiceStatus
	if prev, err = status.Prev(); err != nil {
		// If there is no previous status save the current status to the database
		if errors.Is(err, db.ErrNotFound) {
			log.Debug().Str("service_id", service.ID.String()).Msg("saving first status for service")
			if err = db.Put(status); err != nil {
				return err
			}

			if err = incident.NewVersionDetected(status, service); err != nil {
				return err
			}

			return nil
		}
		return err
	}

	// Compare the statuses, if there is a change in the status save the new status to disk
	if !health.Equal(prev, status) {
		log.Debug().Str("service_id", service.ID.String()).Msg("saving service status change")
		if err = db.Put(status); err != nil {
			return err
		}
	}

	// Handle Incidents; note that the incident requires the service in its previous state
	// for comparison and escalation purposes!
	if !health.Equal(prev, status) || service.Status != health.RelativeStatus(prev, status) {
		if err = incident.New(prev, status, service); err != nil {
			return err
		}
	}

	// Update the service status to the relative status (must happen last).
	service.Status = health.RelativeStatus(prev, status)
	log.Debug().
		Str("service_id", service.ID.String()).
		Str("service", service.Title).
		Str("status", service.Status.String()).
		Msg("status check complete")
	return nil
}
