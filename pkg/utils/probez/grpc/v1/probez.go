package health

import (
	context "context"
	"errors"
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	DefaultService   = "default"
	StatusServing    = HealthCheckResponse_SERVING
	StatusNotServing = HealthCheckResponse_NOT_SERVING
)

// ProbeServer implements the grpc.health.v1.Health service to provide status
// information on different services that the gRPC server may be running. The intended
// use case is to embed the ProbeServer into a gRPC server implementation and to set the
// status of the probes services using the state method.
type ProbeServer struct {
	sync.RWMutex
	UnimplementedHealthServer
	services map[string]HealthCheckResponse_ServingStatus
	watchers map[string]map[string]chan<- HealthCheckResponse_ServingStatus
}

// Healthy sets the probe server default service as serving.
func (h *ProbeServer) Healthy() {
	h.SetStatus(DefaultService, StatusServing)
	log.Debug().Bool("healthy", true).Msg("server is healthy")
}

// NotHealthy sets the probe server default service to not serving.
func (h *ProbeServer) NotHealthy() {
	h.SetStatus(DefaultService, StatusNotServing)
	log.Debug().Bool("healthy", false).Msg("server is not healthy")
}

// Sets the status of the specified service, notifying all watchers about the change in
// status. Callers should not set the status to Unknown or Service Unknown. If an empty
// string is sepcified as the service, the default service is used instead.
func (h *ProbeServer) SetStatus(service string, status HealthCheckResponse_ServingStatus) {
	if service == "" {
		service = DefaultService
	}

	h.Lock()
	defer h.Unlock()
	// Ensure the service and watchers map has been created
	h.checkstate(service)

	// Set the service state for future lookups
	h.services[service] = status

	// Notify the watchers that the status has changed
	for _, watcher := range h.watchers[service] {
		watcher <- status
	}
}

// Status returns the status of the specified service. If empty string is provided then
// the "default service" is used to lookup the status.
func (h *ProbeServer) Status(service string, stream bool) (status HealthCheckResponse_ServingStatus) {
	var ok bool
	if service == "" {
		service = DefaultService
	}

	h.RLock()
	defer h.RUnlock()
	if status, ok = h.services[service]; !ok {
		if stream {
			return HealthCheckResponse_SERVICE_UNKNOWN
		}
		return HealthCheckResponse_UNKNOWN
	}
	return status
}

// Add a watcher with a unique ID to listen for status changes for the specified service.
// A channel is returned that will broadcast status changes for that service.
func (h *ProbeServer) AddWatcher(id, service string) <-chan HealthCheckResponse_ServingStatus {
	watcher := make(chan HealthCheckResponse_ServingStatus, 1)
	if service == "" {
		service = DefaultService
	}

	h.Lock()
	defer h.Unlock()
	// Ensure the service and watchers map has been created
	h.checkstate(service)

	h.watchers[service][id] = watcher
	log.Trace().Str("service", service).Str("watcher", id).Msg("probe watcher added")
	return watcher
}

// Remove a watcher with the specified unique id and stop listening for status changes.
func (h *ProbeServer) DelWatcher(id string) {
	h.Lock()
	defer h.Unlock()
	for sname, service := range h.watchers {
		if watcher, ok := service[id]; ok {
			close(watcher)
			delete(service, id)
			log.Trace().Str("service", sname).Str("watcher", id).Msg("probe watcher deleted")
		}
	}
}

// Check implements the Health service interface and is a Unary response to a health
// check request for the specified service. If the specified service is empty the status
// of the default service is returned.
func (h *ProbeServer) Check(ctx context.Context, in *HealthCheckRequest) (out *HealthCheckResponse, err error) {
	out = &HealthCheckResponse{
		Status: h.Status(in.Service, false),
	}
	return out, nil
}

// Watch implements the Health service interface and provides server-side streaming
// updates of when the status of a specific service changes. If the specified service is
// empty then the status of the default service is returned.
func (h *ProbeServer) Watch(in *HealthCheckRequest, stream Health_WatchServer) (err error) {
	id := ulid.Make().String()
	watcher := h.AddWatcher(id, in.Service)
	ctx := stream.Context()

	// Send the first health check message
	if err = stream.Send(&HealthCheckResponse{Status: h.Status(in.Service, true)}); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return status.Error(codes.Aborted, err.Error())
	}

	// Wait for updates from the watcher and send those to the client
	for {
		select {
		case <-ctx.Done():
			return status.Error(codes.DeadlineExceeded, ctx.Err().Error())
		case serviceStatus := <-watcher:
			if err = stream.Send(&HealthCheckResponse{Status: serviceStatus}); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return status.Error(codes.Aborted, err.Error())
			}
		}
	}
}

// check state ensures that the maps have been initialized for the specified service.
func (h *ProbeServer) checkstate(service string) {
	if h.services == nil {
		h.services = make(map[string]HealthCheckResponse_ServingStatus)
	}

	if h.watchers == nil {
		h.watchers = make(map[string]map[string]chan<- HealthCheckResponse_ServingStatus)
	}

	if _, ok := h.watchers[service]; !ok {
		h.watchers[service] = make(map[string]chan<- HealthCheckResponse_ServingStatus)
	}
}
