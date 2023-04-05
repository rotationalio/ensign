package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
)

const (
	StatusEP        = "/v1/status"
	RegisterEP      = "/v1/register"
	LoginEP         = "/v1/login"
	AuthenticateEP  = "/v1/authenticate"
	RefreshEP       = "/v1/refresh"
	VerifyEP        = "/v1/verify"
	APIKeysEP       = "/v1/apikeys"
	ProjectsEP      = "/v1/projects"
	OrganizationsEP = "/v1/organizations"
	UsersEP         = "/v1/users"
)

// Server embeds an httptest Server and provides additional methods for configuring
// mock responses and counting requests. By default handlers will panic, it's the
// responsibility of the test writer to configure the behavior of each handler that
// will be invoked by using the appropriate On* method and passing in the desired
// HandlerOption(s). If no HandlerOption is specified, the default behavior is to
// return a 200 OK response with an empty body.
type Server struct {
	*httptest.Server
	auth     *authtest.Server
	requests map[string]int
	handlers map[string]http.HandlerFunc
}

// NewServer creates and starts a new mock server for testing Quarterdeck interactions.
func NewServer() (s *Server, err error) {
	s = &Server{
		requests: make(map[string]int),
		handlers: make(map[string]http.HandlerFunc),
	}
	s.Server = httptest.NewServer(http.HandlerFunc(s.routeRequest))

	// Start an authtest server to test authentication responses
	if s.auth, err = authtest.NewServer(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) URL() string {
	return s.Server.URL
}

func (s *Server) Reset() {
	s.requests = make(map[string]int)
	s.handlers = make(map[string]http.HandlerFunc)
}

func (s *Server) Close() {
	s.auth.Close()
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
	case path == VerifyEP:
		s.handlers[path](w, r)
	case strings.Contains(path, APIKeysEP):
		s.handlers[path](w, r)
	case strings.Contains(path, ProjectsEP):
		s.handlers[ProjectsEP](w, r)
	case strings.Contains(path, OrganizationsEP):
		s.handlers[path](w, r)
	case strings.Contains(path, UsersEP):
		s.handlers[path](w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.requests[path]++
}

// HandlerOption allows users of the mock to configure specific endpoint handler
// behavior or override it entirely.
type HandlerOption func(*handlerOptions)

type handlerOptions struct {
	handler http.HandlerFunc
	status  int
	err     string
	fixture interface{}
	auth    bool
}

// Helper to apply the supplied options, panics if there is an error
func handler(opts ...HandlerOption) http.HandlerFunc {
	// Default handler returns a 200 OK with no body
	conf := handlerOptions{
		status:  http.StatusOK,
		fixture: nil,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	// Specified handler overrides all other options
	if conf.handler != nil {
		return conf.handler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case conf.auth && r.Header.Get("Authorization") == "":
			// TODO: Validate the auth token
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing authorization header"))
		case conf.fixture != nil:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(conf.status)
			json.NewEncoder(w).Encode(conf.fixture)
		case conf.err != "":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(conf.status)
			reply := api.Reply{
				Error: conf.err,
			}
			json.NewEncoder(w).Encode(reply)
		default:
			w.WriteHeader(conf.status)
		}
	}
}

// Configure the status code to be returned by the handler
func UseStatus(status int) HandlerOption {
	return func(o *handlerOptions) {
		o.status = status
	}
}

// Configure a basic error reply to be returned by the handler
func UseError(status int, err string) HandlerOption {
	return func(o *handlerOptions) {
		o.status = status
		o.err = err
	}
}

// Configure a JSON fixture to be returned by the handler
func UseJSONFixture(fixture interface{}) HandlerOption {
	return func(o *handlerOptions) {
		o.fixture = fixture
	}
}

// Use the given handler, overriding all other options
func UseHandler(f http.HandlerFunc) HandlerOption {
	return func(o *handlerOptions) {
		o.handler = f
	}
}

// Return a 401 response if the request is not authenticated
func RequireAuth() HandlerOption {
	return func(o *handlerOptions) {
		o.auth = true
	}
}

func fullPath(path, param string) string {
	if param != "" {
		param = "/" + param
	}
	return path + param
}

// Endpoint handlers
func (s *Server) OnStatus(opts ...HandlerOption) {
	s.handlers[StatusEP] = handler(opts...)
}

func (s *Server) OnRegister(opts ...HandlerOption) {
	s.handlers[RegisterEP] = handler(opts...)
}

func (s *Server) OnLogin(opts ...HandlerOption) {
	s.handlers[LoginEP] = handler(opts...)
}

func (s *Server) OnAuthenticate(opts ...HandlerOption) {
	s.handlers[AuthenticateEP] = handler(opts...)
}

func (s *Server) OnRefresh(opts ...HandlerOption) {
	s.handlers[RefreshEP] = handler(opts...)
}

func (s *Server) OnVerify(opts ...HandlerOption) {
	s.handlers[VerifyEP] = handler(opts...)
}

func (s *Server) OnAPIKeys(param string, opts ...HandlerOption) {
	s.handlers[fullPath(APIKeysEP, param)] = handler(opts...)
}

func (s *Server) OnProjects(opts ...HandlerOption) {
	s.handlers[ProjectsEP] = handler(opts...)
}

func (s *Server) OnOrganizations(param string, opts ...HandlerOption) {
	s.handlers[fullPath(OrganizationsEP, param)] = handler(opts...)
}

func (s *Server) OnUsers(param string, opts ...HandlerOption) {
	s.handlers[fullPath(UsersEP, param)] = handler(opts...)
}

// Request counters
func (s *Server) StatusCount() int {
	return s.requests[StatusEP]
}

func (s *Server) RegisterCount() int {
	return s.requests[RegisterEP]
}

func (s *Server) LoginCount() int {
	return s.requests[LoginEP]
}

func (s *Server) AuthenticateCount() int {
	return s.requests[AuthenticateEP]
}

func (s *Server) RefreshCount() int {
	return s.requests[RefreshEP]
}

func (s *Server) VerifyCount() int {
	return s.requests[VerifyEP]
}

func (s *Server) APIKeysCount(param string) int {
	return s.requests[fullPath(APIKeysEP, param)]
}

func (s *Server) ProjectsCount() int {
	return s.requests[ProjectsEP]
}

func (s *Server) OrganizationsCount(param string) int {
	return s.requests[fullPath(OrganizationsEP, param)]
}

func (s *Server) UsersCount(param string) int {
	return s.requests[fullPath(UsersEP, param)]
}
