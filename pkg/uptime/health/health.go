package health

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/ensign/pkg/uptime/db"
	"github.com/rotationalio/ensign/pkg/utils/ulid"
)

// If the time between relative statuses is greater than this threshold then the status
// should be escalated; e.g. degraded to unhealthy or offline to outage.
const (
	DegradedThreshold  = 1 * time.Minute
	UnhealthyThreshold = 2 * time.Minute
	OfflineThreshold   = 5 * time.Minute
)

// Monitor is meant to wrap a service endpoint and return a service status. These
// checks can be http or gRPC based or some other network connection. They can also
// generically check if a service is online or parse a specific service's expected
// response to generate incident reports.
type Monitor interface {
	Status(ctx context.Context) (ServiceStatus, error)
}

type ServiceStatus interface {
	db.Model

	// In order to support quick comparisons between service statuses of the same type,
	// the data should be hashed uniquely into bytes so that it can be compared with
	// another hash in a deterministic fashion. This is used by the
	// ComparableServiceStatus embedding to implement the Equal() method. The data to
	// be hashed should only include data that users are interested in changes for, e.g.
	// the readiness state, the version number, etc. It should exclude variable or
	// continuous data such as uptime, number of requests, etc.
	//
	// If computing the hash errors then the method should panic or return a nil hash
	// depending on the severity of possible errors.
	Hash() []byte

	// Status should return the current service status as determined by the health check.
	Status() Status

	// CheckedAt should return the time that the check was conducted.
	CheckedAt() time.Time

	// SetServiceID allows external users to specify an ID prefix for the Key.
	SetServiceID(sid uuid.UUID)

	// Prev returns the previous status from the database
	Prev() (ServiceStatus, error)
}

// Relative status computes the status between two consecutive service statuses. For
// example, if the previous status was unhealthy and the current status is also
// unhealthy but the time between statuses is greater than 5 minutes then return outage.
func RelativeStatus(previous, current ServiceStatus) Status {
	// Relative status checks if the two checks are the same
	status := current.Status()
	if status == previous.Status() {
		delta := current.CheckedAt().Sub(previous.CheckedAt())
		if status == Degraded && delta > DegradedThreshold {
			status = Unhealthy
		}

		if status == Unhealthy && delta > UnhealthyThreshold {
			status = Offline
		}

		if status == Offline && delta > OfflineThreshold {
			status = Outage
		}
	}

	// If the above checks fail simply return the current status
	return status
}

// Equal compares the statuses to determine if they're the same without taking into
// account variable data between checks such as timestamps.
func Equal(first, second ServiceStatus) bool {
	return bytes.Equal(first.Hash(), second.Hash())
}

type BaseStatus struct {
	ID        []byte    `msgpack:"id"`
	Timestamp time.Time `msgpack:"timestamp"`
	sid       uuid.UUID
	hash      []byte
}

// Key returns the service status ID as the key and error if no key is set.
func (h *BaseStatus) Key() ([]byte, error) {
	if len(h.ID) == 0 {
		if h.sid == uuid.Nil {
			return nil, ErrNoServiceID
		}

		if h.Timestamp.IsZero() {
			return nil, ErrNoTimestamp
		}

		// Compose the ID as a composite of the service id and a ULID with the timestamp
		h.ID = make([]byte, 32)
		copy(h.ID[0:16], h.sid[:])

		tsid := ulid.FromTime(h.Timestamp)
		copy(h.ID[16:], tsid[:])
	}
	return h.ID, nil
}

func (h *BaseStatus) CheckedAt() time.Time {
	return h.Timestamp
}

func (h *BaseStatus) SetServiceID(sid uuid.UUID) {
	h.sid = sid
}

func (h *BaseStatus) GetServiceID() (uuid.UUID, error) {
	if h.sid == uuid.Nil {
		if len(h.ID) == 0 {
			return uuid.Nil, ErrNoServiceID
		}
		h.sid = uuid.UUID{}
		if err := h.sid.UnmarshalBinary(h.ID[0:16]); err != nil {
			return uuid.Nil, err
		}
	}
	return h.sid, nil
}
