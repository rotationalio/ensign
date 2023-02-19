package health

import (
	"bytes"
	"context"
	"time"
)

// Monitor is meant to wrap a service endpoint and return a service status. These
// checks can be http or gRPC based or some other network connection. They can also
// generically check if a service is online or parse a specific service's expected
// response to generate incident reports.
type Monitor interface {
	Status(ctx context.Context) (ServiceStatus, error)
}

type ServiceStatus interface {
	// Marshal and Unmarshal the service status to save it to disk.
	Marshal() ([]byte, error)
	Unmarshal([]byte) error

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
}

func Equal(first, second ServiceStatus) bool {
	return bytes.Equal(first.Hash(), second.Hash())
}
