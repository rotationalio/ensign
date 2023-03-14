package health

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/ensign/pkg/uptime/db"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EnsignMonitor wraps a gRPC Ensign client to perform health checks.
type EnsignMonitor struct {
	cc       *grpc.ClientConn
	client   api.EnsignClient
	checked  time.Time
	attempts uint32
}

var _ Monitor = &EnsignMonitor{}

func NewEnsignMonitor(endpoint string, opts ...MonitorOption) (mon *EnsignMonitor, err error) {
	var conf *Options
	if conf, err = NewOptions(opts...); err != nil {
		return nil, err
	}

	dialer := make([]grpc.DialOption, 0, 1)
	if conf.EnsignInsecure {
		dialer = append(dialer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialer = append(dialer, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	mon = &EnsignMonitor{}
	if mon.cc, err = grpc.Dial(endpoint, dialer...); err != nil {
		return nil, err
	}

	mon.client = api.NewEnsignClient(mon.cc)
	return mon, nil
}

func (h *EnsignMonitor) Status(ctx context.Context) (_ ServiceStatus, err error) {
	out := &api.HealthCheck{
		Attempts:      h.attempts,
		LastCheckedAt: timestamppb.New(h.checked),
	}

	state := &EnsignStatus{
		BaseStatus: BaseStatus{
			Timestamp: time.Now(),
		},
	}

	var in *api.ServiceState
	if in, err = h.client.Status(ctx, out); err != nil {
		state.Error = err.Error()

		if serr, ok := status.FromError(err); ok {
			state.Code = serr.Code()
			state.Error = serr.Message()

			switch state.Code {
			case codes.Canceled, codes.DeadlineExceeded, codes.ResourceExhausted:
				state.ErrorType = Degraded
			case codes.Aborted, codes.Unavailable, codes.DataLoss:
				state.ErrorType = Unhealthy
			}

		} else if errors.Is(err, context.DeadlineExceeded) {
			state.ErrorType = Degraded
		}

		return state, nil
	}

	state.Code = codes.OK
	state.EnsignState = in.Status
	state.EnsignUptime = in.Uptime.AsDuration()
	state.EnsignVersion = in.Version
	return state, nil
}

// EnsignStatus determines the health of the status by its status response.
type EnsignStatus struct {
	BaseStatus
	Code          codes.Code              `msgpack:"code"`
	Error         string                  `msgpack:"error"`
	ErrorType     Status                  `msgpack:"error_type"`
	EnsignState   api.ServiceState_Status `msgpack:"ensign_state"`
	EnsignUptime  time.Duration           `msgpack:"uptime"`
	EnsignVersion string                  `msgpack:"version"`
}

var _ ServiceStatus = &EnsignStatus{}

// Unmarshal from msgpack binary data.
func (h *EnsignStatus) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, h)
}

// Marshal to msgpack binary data for storage.
func (h *EnsignStatus) Marshal() ([]byte, error) {
	return msgpack.Marshal(h)
}

// Return the previous EnsignStatus
func (h *EnsignStatus) Prev() (_ ServiceStatus, err error) {
	var sid uuid.UUID
	if sid, err = h.GetServiceID(); err != nil {
		return nil, err
	}

	prev := &EnsignStatus{}
	if err = db.LastServiceStatus(sid, prev); err != nil {
		return nil, err
	}
	return prev, nil
}

// Hashes the status code, error, and error type for comparison purposes.
func (h *EnsignStatus) Hash() []byte {
	if h.hash == nil {
		sig := fnv.New128()

		// Write the Status code
		code := make([]byte, 4)
		binary.LittleEndian.PutUint32(code, uint32(h.Code))
		sig.Write(code)

		// Write the parsed status
		state := h.Status()
		binstate := make([]byte, 2)
		binary.LittleEndian.PutUint16(binstate, uint16(state))
		sig.Write(binstate)

		// Write the version
		sig.Write([]byte(h.Version()))

		// Write the error
		sig.Write([]byte(h.Error))

		// Write the error type
		etype := make([]byte, 2)
		binary.LittleEndian.PutUint16(etype, uint16(h.ErrorType))
		sig.Write(etype)

		buf := make([]byte, 0, 16)
		h.hash = sig.Sum(buf)
	}
	return h.hash
}

func (h *EnsignStatus) Status() Status {
	// If we have neither a status code nor an error then return unknown.
	if h.Code == codes.Unknown && h.Error == "" {
		return Unknown
	}

	// If we have an error then determine the error from the error type otherwise simply
	// report offline since we were unable to connect and make a request to the server.
	if h.Error != "" {
		if h.ErrorType > 0 {
			return h.ErrorType
		}
		return Offline
	}

	switch h.EnsignState {
	case api.ServiceState_HEALTHY:
		return Online
	case api.ServiceState_MAINTENANCE:
		return Maintenance
	case api.ServiceState_UNHEALTHY:
		return Unhealthy
	case api.ServiceState_DANGER:
		return Offline
	case api.ServiceState_OFFLINE:
		return Offline
	default:
		return Unknown
	}
}

// Version attempts to get version information from the API response.
func (h *EnsignStatus) Version() string {
	return h.EnsignVersion
}
