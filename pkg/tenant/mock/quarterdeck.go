package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

const (
	StatusEP       = "/v1/status"
	RegisterEP     = "/v1/register"
	LoginEP        = "/v1/login"
	AuthenticateEP = "/v1/authenticate"
	RefreshEP      = "/v1/refresh"
	APIKeysEP      = "/v1/apikeys"
)

// Server embeds an httptest Server and provides additional methods for configuring
// mock responses and counting requests.
type Server struct {
	*httptest.Server
	requests map[string]int
	handlers map[string]http.HandlerFunc
}

// NewServer creates a new mock server for testing Quarterdeck interactions.
func NewServer() (*Server, error) {
	s := &Server{
		requests: make(map[string]int),
		handlers: make(map[string]http.HandlerFunc),
	}

	s.Server = httptest.NewServer(http.HandlerFunc(s.routeRequest))
	return s, nil
}

func (s *Server) Serve() {
	s.Server.Start()
}

func (s *Server) URL() string {
	return s.Server.URL
}

func (s *Server) Close() {
	s.Server.Close()
}

func (s *Server) routeRequest(w http.ResponseWriter, r *http.Request) {
	// Simple paths with no parameters
	path := r.URL.Path
	switch {
	case path == StatusEP:
		s.handlers[path](w, r)
	case path == RegisterEP:
		s.handlers[path](w, r)
	case path == LoginEP:
		s.handlers[path](w, r)
	case path == AuthenticateEP:
		s.handlers[path](w, r)
	case path == RefreshEP:
		s.handlers[path](w, r)
	case strings.Contains(path, APIKeysEP):
		s.handlers[path](w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

	s.requests[path]++
}

func useFixture(w http.ResponseWriter, r *http.Request, fixture interface{}) {
	var data []byte
	var err error
	if data, err = json.Marshal(fixture); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func fullPath(path, param string) string {
	return path + "/" + param
}

func (s *Server) OnStatus(f http.HandlerFunc) {
	s.handlers[StatusEP] = f
}

func (s *Server) OnRegister(f http.HandlerFunc) {
	s.handlers[RegisterEP] = f
}

func (s *Server) OnLogin(f http.HandlerFunc) {
	s.handlers[LoginEP] = f
}

func (s *Server) OnAuthenticate(f http.HandlerFunc) {
	s.handlers[AuthenticateEP] = f
}

func (s *Server) OnRefresh(f http.HandlerFunc) {
	s.handlers[RefreshEP] = f
}

func (s *Server) OnAPIKeys(param string, f http.HandlerFunc) {
	s.handlers[fullPath(APIKeysEP, param)] = f
}
