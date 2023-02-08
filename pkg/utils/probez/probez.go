/*
Package probez provides a simple http server that implements Kubernetes readiness probes
(e.g. /livez and /readyz) that can be embedded in containers with long running processes
or added to gin services to ensure that they can control how Kubernetes views their
service state.
*/
package probez

import (
	"io"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

const (
	Healthz = "/healthz"
	Livez   = "/livez"
	Readyz  = "/readyz"
)

// The Probe Handler manages the health and readiness states of a container. When it is
// first created, it starts in a healthy, but not ready state. Users can mark the probe
// server as healthy, unhealthy, ready, or not ready, changing how it responds to HTTP
// Get requests at the /livez and /readyz endpoints respectively.
//
// Users should either instantiate a new probez.Server and serve it on a bind address or
// they should configure a probez.Handler to use an http server or a gin router for
// serving requests. When the the service is ready, they should mark the probez.Handler
// as ready, and when it is shutting down, the server should be marked as not ready and
// unhealthy.
type Handler struct {
	healthy *atomic.Value
	ready   *atomic.Value
}

var _ http.Handler = &Handler{}

// New returns a probez.Handler that is healthy but not ready.
func New() *Handler {
	srv := &Handler{
		healthy: &atomic.Value{},
		ready:   &atomic.Value{},
	}

	srv.Healthy()
	srv.NotReady()

	return srv
}

// Healthy sets the probe server to healthy so that it responds 204 No Content to
// liveness probes at the /livez endpoint.
func (h *Handler) Healthy() {
	h.healthy.Store(true)
}

// NotHealthy sets the probe server to unhealthy so that it responds 503 Unavailable to
// liveness probes at the /livez endpoint.
func (h *Handler) NotHealthy() {
	h.healthy.Store(false)
}

// IsHealthy returns if the Handler is healthy or not
func (h *Handler) IsHealthy() bool {
	return h.healthy.Load().(bool)
}

// Ready sets the probe server state to ready so that it responds 204 No Content to
// readiness probes at the /readyz endpoint. This operation is thread-safe.
func (h *Handler) Ready() {
	h.ready.Store(true)
}

// NotReady sets the probe server state to not ready so that it responds 503 Unavailable
// to readiness probes at the /readyz endpoint.
func (h *Handler) NotReady() {
	h.ready.Store(false)
}

// IsReady returns if the Handler is ready or not
func (h *Handler) IsReady() bool {
	return h.ready.Load().(bool)
}

// Use adds the server's routes to the specified Gin router.
func (h *Handler) Use(router *gin.Engine) {
	router.GET(Livez, gin.WrapF(h.Healthz))
	router.GET(Healthz, gin.WrapF(h.Healthz))
	router.GET(Readyz, gin.WrapF(h.Readyz))
}

// Mux adds the server's routes to the specified ServeMux.
func (h *Handler) Mux(mux *http.ServeMux) {
	mux.HandleFunc(Livez, h.Healthz)
	mux.HandleFunc(Healthz, h.Healthz)
	mux.HandleFunc(Readyz, h.Readyz)
}

// Healthz implements the kubernetes liveness check and responds to /healthz and /livez
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	if h.healthy == nil || !h.healthy.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
}

// Readyz implements the kubernetes readiness check and responds to /readyz requests.
func (h *Handler) Readyz(w http.ResponseWriter, _ *http.Request) {
	if h.ready == nil || !h.ready.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case Livez:
		h.Healthz(w, r)
	case Healthz:
		h.Healthz(w, r)
	case Readyz:
		h.Readyz(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
