/*
Package health implements health check clients to return server system statuses.
*/
package health

import (
	"fmt"
	"strings"
)

type Status uint

const (
	Unknown Status = iota
	Online
	Maintenance
	Stopping
	Degraded
	Unhealthy
	Offline
	Outage
)

func (s Status) String() string {
	switch s {
	case Online:
		return "online"
	case Maintenance:
		return "maintenance"
	case Stopping:
		return "stopping"
	case Degraded:
		return "degraded"
	case Unhealthy:
		return "unhealthy"
	case Offline:
		return "offline"
	case Outage:
		return "outage"
	default:
		return "unknown"
	}
}

// Parse Status attempts to parse a status message from a string.
func ParseStatus(s string) (Status, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "ok", "online", "healthy", "ready":
		return Online, nil
	case "maintenance":
		return Maintenance, nil
	case "stopping", "stopped":
		return Stopping, nil
	case "unhealthy":
		return Unhealthy, nil
	case "offline":
		return Offline, nil
	case "outage":
		return Outage, nil
	case "degraded":
		return Degraded, nil
	case "unknown":
		return Unknown, nil
	}
	return Unknown, fmt.Errorf("could not parse status %q", s)
}

type StatusChange uint

const (
	NoChange StatusChange = iota
	BackOnline
	Deescalating
	NoLongerHealthy
	Escalating
)

func CompareStatus(previous, current Status) StatusChange {
	// Check if the incident is escalating, e.g. going from a less concerning state to
	// a more concerning state (e.g. unhealthy to outage).
	if current > previous {
		// Are we starting from healthy or from another state?
		if previous == Online {
			return NoLongerHealthy
		}
		return Escalating
	}

	// Check if the incident is deescalating, e.g. going from a more concerning state to
	// a less concerning state (e.g. outage to online).
	if current < previous {
		// Are we going to a healthy state?
		if current == Online {
			return BackOnline
		}
		return Deescalating
	}

	return NoChange
}
