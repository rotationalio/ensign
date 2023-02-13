package probez

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// Server is a quick way to get a probez.Handler service up and running, particularly
// for containers that do not necessarily serve HTTP requests.
type Server struct {
	Handler
	srv  *http.Server
	addr net.Addr
}

// New returns a probez.Server that is healthy but not ready.
func NewServer() *Server {
	srv := &Server{
		Handler: Handler{
			healthy: &atomic.Value{},
			ready:   &atomic.Value{},
		},
	}

	srv.Healthy()
	srv.NotReady()

	return srv
}

// Serve probe requests on the specified port, handling the /livez and /readyz
// endpoints according to the state of the server.
func (s *Server) Serve(addr string) (err error) {
	mux := http.NewServeMux()
	s.Mux(mux)

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
