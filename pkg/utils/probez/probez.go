/*
Package probez provides a simple http server that implements Kubernetes readiness probes
(e.g. /livez and /readyz) that can be embedded in containers with long running processes
or added to gin services to ensure that they can control how Kubernetes views their
service state.
*/
package probez

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// The Probe Server manages the health and readiness states of a container. When it is
// first created, it starts in a healthy, but not ready state. Users can mark the probe
// server as healthy, unhealthy, ready, or not ready, changing how it responds to HTTP
// Get requests at the /livez and /readyz endpoints respectively.
//
// Users should either instantiate a new probez.Server and serve it on a bind address or
// they should configure it to use a gin router for serving requests. When the the
// service is ready, they should mark the probez.Server as ready, and when it is
// shutting down, the server should be marked as not ready and unhealthy.
type Server struct {
	healthy *atomic.Value
	ready   *atomic.Value
	srv     *http.Server
	addr    net.Addr
}

// New returns a probez.Server that is healthy but not ready.
func New() *Server {
	srv := &Server{
		healthy: &atomic.Value{},
		ready:   &atomic.Value{},
	}

	srv.Healthy()
	srv.NotReady()

	return srv
}

// Healthy sets the probe server to healthy so that it responds 204 No Content to
// liveness probes at the /livez endpoint.
func (s *Server) Healthy() {
	s.healthy.Store(true)
}

// NotHealthy sets the probe server to unhealthy so that it responds 503 Unavailable to
// liveness probes at the /livez endpoint.
func (s *Server) NotHealthy() {
	s.healthy.Store(false)
}

// Ready sets the probe server state to ready so that it responds 204 No Content to
// readiness probes at the /readyz endpoint. This operation is thread-safe.
func (s *Server) Ready() {
	s.ready.Store(true)
}

// NotReady sets the probe server state to not ready so that it responds 503 Unavailable
// to readiness probes at the /readyz endpoint.
func (s *Server) NotReady() {
	s.ready.Store(false)
}

// Serve probe requests on the specified port, handling the /livez and /readyz
// endpoints according to the state of the server.
func (s *Server) Serve(addr string) (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/livez", s.Healthz)
	mux.HandleFunc("/healthz", s.Healthz)
	mux.HandleFunc("/readyz", s.Readyz)

	s.srv = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ErrorLog:          nil,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	var sock net.Listener
	if sock, err = net.Listen("tcp", addr); err != nil {
		return err
	}

	s.addr = sock.Addr()
	go func() {
		log.Debug().Str("addr", addr).Msg("starting the probez server")
		if err = s.srv.Serve(sock); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("probez server shutdown prematurely")
		}
	}()
	return nil
}

// Shutdown the http server and stop responding to probe requests.
func (s *Server) Shutdown(ctx context.Context) (err error) {
	if s.srv == nil {
		return nil
	}

	s.srv.SetKeepAlivesEnabled(false)
	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Debug().Msg("shutdown the probez server")
	return nil
}

// Handle adds the server's routes to the specified Gin router.
func (s *Server) Handle(router *gin.Engine) {
	router.GET("/healthz", gin.WrapF(s.Healthz))
	router.GET("/livez", gin.WrapF(s.Healthz))
	router.GET("/readyz", gin.WrapF(s.Readyz))
}

func (s *Server) Healthz(w http.ResponseWriter, _ *http.Request) {
	if s.healthy == nil || !s.healthy.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
}

func (s *Server) Readyz(w http.ResponseWriter, _ *http.Request) {
	if s.ready == nil || !s.ready.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
}

func (s *Server) URL() string {
	u := &url.URL{
		Scheme: "http",
		Host:   s.addr.String(),
	}

	if addr, ok := s.addr.(*net.TCPAddr); ok && addr.IP.IsUnspecified() {
		u.Host = fmt.Sprintf("127.0.0.1:%d", addr.Port)
	}
	return u.String()
}
