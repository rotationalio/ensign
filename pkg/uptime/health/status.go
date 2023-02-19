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
	Maintenance
	Stopping
	Online
	Degraded
	Unhealthy
	Offline
	Outage
)

func (s Status) String() string {
	switch s {
	case Maintenance:
		return "maintenance"
	case Stopping:
		return "stopping"
	case Online:
		return "online"
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
	case "ok", "online":
		return Online, nil
	case "maintenance":
		return Maintenance, nil
	case "stopping":
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
