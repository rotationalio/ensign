package incident

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/ensign/pkg/uptime/db"
	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/rotationalio/ensign/pkg/uptime/services"
	"github.com/rs/zerolog/log"
)

// Incidents are grouped by day and stored together in the database to make it easier
// to fetch incidents for a specific date rather than have to serialize each incident
// individually. To update an incident group ensure a transaction is used.
type Group struct {
	Date      time.Time   `msgpack:"date"`
	Incidents []*Incident `msgpack:"incidents"`
}

type Incident struct {
	ServiceID      uuid.UUID     `msgpack:"service_id"`
	ServiceName    string        `msgpack:"service_name"`
	Description    string        `msgpack:"description"`
	StartTime      time.Time     `msgpack:"start_time"`
	EndTime        time.Time     `msgpack:"end_time"`
	Status         health.Status `msgpack:"status"`
	PreviousStatus health.Status `msgpack:"previous_status"`
}

// New creates an incident based on the previous and current statuses.
func New(previous, current health.ServiceStatus, service *services.Service) (err error) {
	// Determine if there is a change in the service status
	if !health.Equal(previous, current) {
		// Determine if a status change has occurred
		if previous.Status() != current.Status() {
			if err = NewStatusChange(previous, current, service); err != nil {
				return err
			}
		}

		// Determine if a version change has occurred
		if health.VersionChanged(previous, current) {
			if err = NewVersionChange(previous, current, service); err != nil {
				return err
			}
		}
		return nil
	}

	// Determine if an outage has escalated
	relative := health.RelativeStatus(previous, current)
	if relative != service.Status {
		return NewRelativeStatusChange(previous, current, service)
	}
	return nil
}

func NewStatusChange(previous, current health.ServiceStatus, service *services.Service) (err error) {
	incident := &Incident{
		ServiceID:      service.ID,
		ServiceName:    service.Title,
		StartTime:      current.CheckedAt(),
		Status:         current.Status(),
		PreviousStatus: previous.Status(),
	}

	incident.Description = incident.DescriptionFromStatus()

	switch health.CompareStatus(incident.PreviousStatus, incident.Status) {
	case health.NoLongerHealthy:
		// Simplest case, simply create a new incident
		if err = Create(incident); err != nil {
			return err
		}
		log.Info().
			Str("service_id", incident.ServiceID.String()).
			Str("from_status", incident.PreviousStatus.String()).
			Str("to_status", incident.Status.String()).
			Msg("status change incident created")
	case health.BackOnline:
		// Find the previous incident and update its endtime
		if err = Conclude(incident.ServiceID, previous.CheckedAt(), current.CheckedAt()); err != nil {
			log.Warn().Err(err).Str("service_id", incident.ServiceID.String()).Time("start_time", previous.CheckedAt()).Msg("could not conclude previous incident")
		}
		log.Info().
			Str("service_id", incident.ServiceID.String()).
			Str("from_status", incident.PreviousStatus.String()).
			Str("to_status", incident.Status.String()).
			Msg("status change incident concluded")

	case health.Escalating, health.Deescalating:
		// Find the previous incident, update its endtime and create a new incident
		if err = Conclude(incident.ServiceID, previous.CheckedAt(), current.CheckedAt()); err != nil {
			log.Warn().Err(err).Str("service_id", incident.ServiceID.String()).Time("start_time", previous.CheckedAt()).Msg("could not conclude previous incident")
		}
		if err = Create(incident); err != nil {
			return err
		}
		log.Info().
			Str("service_id", incident.ServiceID.String()).
			Str("from_status", incident.PreviousStatus.String()).
			Str("to_status", incident.Status.String()).
			Msg("status change incident continuing")
	default:
		return errors.New("unhandled compare status value")
	}
	return nil
}

func NewRelativeStatusChange(previous, current health.ServiceStatus, service *services.Service) (err error) {
	relativeStatus := health.RelativeStatus(previous, current)
	switch health.CompareStatus(previous.Status(), relativeStatus) {
	case health.Escalating, health.Deescalating:
		// Find the previous incident, update its description and status
		if err = Update(service.ID, previous.CheckedAt(), func(i *Incident) error {
			i.Status = relativeStatus
			i.Description = i.DescriptionFromStatus()
			return nil
		}); err != nil {
			// If it's a not found error, then create the incident
			if errors.Is(err, db.ErrNotFound) {
				incident := &Incident{
					ServiceID:      service.ID,
					ServiceName:    service.Title,
					StartTime:      previous.CheckedAt(),
					Status:         relativeStatus,
					PreviousStatus: previous.Status(),
				}
				incident.Description = incident.DescriptionFromStatus()

				if err = Create(incident); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	default:
		// Note that health.NoLongerHealthy and health.BackOnline should not occur in a
		// relative status change since these changes are time based/heuristic.
		return errors.New("unhandled compare status value")
	}

	log.Info().
		Str("service_id", service.ID.String()).
		Str("status", health.RelativeStatus(previous, current).String()).
		Msg("relative status change detected")
	return nil
}

func NewVersionChange(previous, current health.ServiceStatus, service *services.Service) (err error) {
	incident := &Incident{
		ServiceID:   service.ID,
		ServiceName: service.Title,

		StartTime: current.CheckedAt(),
		Status:    health.Maintenance,
	}

	var fromVersion string
	if versioned, ok := previous.(health.Versioned); ok {
		fromVersion = versioned.Version()
		if fromVersion != "" && !strings.HasPrefix(fromVersion, "v") {
			fromVersion = "v" + fromVersion
		}
	}

	var toVersion string
	if versioned, ok := current.(health.Versioned); ok {
		toVersion = versioned.Version()
		if toVersion != "" && !strings.HasPrefix(toVersion, "v") {
			toVersion = "v" + toVersion
		}
	}

	if fromVersion != "" && toVersion != "" {
		incident.Description = fmt.Sprintf("Detected a version change for the %s, from %s to %s", incident.ServiceName, fromVersion, toVersion)
	} else if toVersion != "" {
		incident.Description = fmt.Sprintf("Detected a version change for the %s to %s", incident.ServiceName, toVersion)
	} else if fromVersion != "" {
		incident.Description = fmt.Sprintf("Detected a version change for the %s, from %s to an unknown version", incident.ServiceName, fromVersion)
	} else {
		incident.Description = fmt.Sprintf("Detected a version change for the %s", incident.ServiceName)
	}

	if err = Create(incident); err != nil {
		return err
	}

	log.Info().
		Str("service_id", incident.ServiceID.String()).
		Str("from_version", fromVersion).
		Str("to_version", toVersion).
		Msg("version change incident created")
	return nil
}

func NewVersionDetected(current health.ServiceStatus, service *services.Service) (err error) {
	incident := &Incident{
		ServiceID:   service.ID,
		ServiceName: service.Title,

		StartTime: current.CheckedAt(),
		Status:    health.Maintenance,
	}

	var toVersion string
	if versioned, ok := current.(health.Versioned); ok {
		toVersion = versioned.Version()
		if toVersion != "" && !strings.HasPrefix(toVersion, "v") {
			toVersion = "v" + toVersion
		}
	}

	if toVersion != "" {
		incident.Description = fmt.Sprintf("Detected a version change for the %s to %s", incident.ServiceName, toVersion)
	} else {
		incident.Description = fmt.Sprintf("Detected a version change for the %s", incident.ServiceName)
	}

	if err = Create(incident); err != nil {
		return err
	}

	log.Info().
		Str("service_id", incident.ServiceID.String()).
		Str("to_version", toVersion).
		Msg("version change incident created")
	return nil
}

func (i *Incident) DescriptionFromStatus() string {
	switch i.Status {
	case health.Online:
		return fmt.Sprintf("%s is online and healthy", i.ServiceName)
	case health.Maintenance:
		return fmt.Sprintf("Currently conducting maintenance on the %s; the service should be back online soon", i.ServiceName)
	case health.Stopping:
		return fmt.Sprintf("Detected servers stopping or rebooting (%s)", i.ServiceName)
	case health.Degraded:
		return fmt.Sprintf("The %s has slowed down or is experiencing degraded performance", i.ServiceName)
	case health.Unhealthy:
		return fmt.Sprintf("Partial %s service outage or unhealthy state detected", i.ServiceName)
	case health.Offline:
		return fmt.Sprintf("The %s is offline and cannot currently be accessed", i.ServiceName)
	case health.Outage:
		return fmt.Sprintf("A major outage has been detected for the %s", i.ServiceName)
	default:
		return fmt.Sprintf("The %s has gone into an unknown state, it is under investigation.", i.ServiceName)
	}
}
