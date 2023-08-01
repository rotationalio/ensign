package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
)

const (
	StatusEP        = "/v1/status"
	RegisterEP      = "/v1/register"
	LoginEP         = "/v1/login"
	AuthenticateEP  = "/v1/authenticate"
	RefreshEP       = "/v1/refresh"
	SwitchEP        = "/v1/switch"
	VerifyEP        = "/v1/verify"
	APIKeysEP       = "/v1/apikeys"
	ProjectsEP      = "/v1/projects"
	OrganizationsEP = "/v1/organizations"
	UsersEP         = "/v1/users"
	InvitesEP       = "/v1/invites"
)

// Server embeds an httptest Server and provides additional methods for configuring
// mock responses and counting requests. By default handlers will panic, it's the
// responsibility of the test writer to configure the behavior of each handler that
// will be invoked by using the appropriate On* method and passing in the desired
// HandlerOption(s). If no HandlerOption is specified, the default behavior is to
// return a 200 OK response with an empty body.
type Server struct {
	sync.RWMutex
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

func (s *Server) handlerFunc(requestKey string) http.HandlerFunc {
	s.RLock()
	defer s.RUnlock()

	if handler, ok := s.handlers[requestKey]; ok {
		return handler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		reply := api.Reply{
			Error: fmt.Sprintf("mock handler not registered for request %q", requestKey),
		}
		json.NewEncoder(w).Encode(reply)
	}
}

func (s *Server) incrementRequest(path string) {
	s.Lock()
	defer s.Unlock()
	s.requests[path]++
}

func (s *Server) routeRequest(w http.ResponseWriter, r *http.Request) {
	requestKey := methodPath(r.Method, r.URL.Path)
	s.handlerFunc(requestKey)(w, r)
	s.incrementRequest(requestKey)
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
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.Reply{Error: "missing authorization header"})
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

func methodPath(method, path string) string {
	return method + " " + path
}

func fullPath(path, param string) string {
	if param != "" {
		param = "/" + param
	}
	return path + param
}

func (s *Server) setHandler(method, path string, opts ...HandlerOption) {
	s.Lock()
	defer s.Unlock()
	s.handlers[methodPath(method, path)] = handler(opts...)
}

// Endpoint handlers
func (s *Server) OnStatus(opts ...HandlerOption) {
	s.setHandler(http.MethodGet, StatusEP, opts...)
}

func (s *Server) OnRegister(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, RegisterEP, opts...)
}

func (s *Server) OnLogin(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, LoginEP, opts...)
}

func (s *Server) OnAuthenticate(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, AuthenticateEP, opts...)
}

func (s *Server) OnRefresh(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, RefreshEP, opts...)
}

func (s *Server) OnSwitch(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, SwitchEP, opts...)
}

func (s *Server) OnVerify(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, VerifyEP, opts...)
}

func (s *Server) OnAPIKeysList(opts ...HandlerOption) {
	s.setHandler(http.MethodGet, APIKeysEP, opts...)
}

func (s *Server) OnAPIKeysCreate(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, APIKeysEP, opts...)
}

func (s *Server) OnAPIKeysDetail(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodGet, fullPath(APIKeysEP, id), opts...)
}

func (s *Server) OnAPIKeysDelete(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodDelete, fullPath(APIKeysEP, id), opts...)
}

func (s *Server) OnAPIKeysUpdate(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodPut, fullPath(APIKeysEP, id), opts...)
}

func (s *Server) OnAPIKeysPermissions(opts ...HandlerOption) {
	s.setHandler(http.MethodGet, fullPath(APIKeysEP, "permissions"), opts...)
}

func (s *Server) OnProjectsList(opts ...HandlerOption) {
	s.setHandler(http.MethodGet, ProjectsEP, opts...)
}

func (s *Server) OnProjectsCreate(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, ProjectsEP, opts...)
}

func (s *Server) OnProjectsAccess(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, fullPath(ProjectsEP, "access"), opts...)
}

func (s *Server) OnProjectsDetail(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodGet, fullPath(ProjectsEP, id), opts...)
}

func (s *Server) OnOrganizations(param string, opts ...HandlerOption) {
	s.setHandler(http.MethodGet, fullPath(OrganizationsEP, param), opts...)
}

func (s *Server) OnUsersDetail(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodGet, fullPath(UsersEP, id), opts...)
}

func (s *Server) OnUsersUpdate(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodPut, fullPath(UsersEP, id), opts...)
}

func (s *Server) OnUsersList(opts ...HandlerOption) {
	s.setHandler(http.MethodGet, UsersEP, opts...)
}

func (s *Server) OnUsersRemove(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodDelete, fullPath(UsersEP, id), opts...)
}

func (s *Server) OnUsersRemoveConfirm(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodPost, fullPath(UsersEP, id+"/confirm"), opts...)
}

func (s *Server) OnUsersRoleUpdate(id string, opts ...HandlerOption) {
	s.setHandler(http.MethodPost, fullPath(UsersEP, id), opts...)
}

func (s *Server) OnInvitesPreview(token string, opts ...HandlerOption) {
	s.setHandler(http.MethodGet, fullPath(InvitesEP, token), opts...)
}

func (s *Server) OnInvitesCreate(opts ...HandlerOption) {
	s.setHandler(http.MethodPost, InvitesEP, opts...)
}

func (s *Server) count(requestKey string) int {
	s.RLock()
	defer s.RUnlock()
	return s.requests[requestKey]
}

// Request counters
func (s *Server) StatusCount() int {
	return s.count(methodPath(http.MethodGet, StatusEP))
}

func (s *Server) RegisterCount() int {
	return s.count(methodPath(http.MethodPost, RegisterEP))
}

func (s *Server) LoginCount() int {
	return s.count(methodPath(http.MethodPost, LoginEP))
}

func (s *Server) AuthenticateCount() int {
	return s.count(methodPath(http.MethodPost, AuthenticateEP))
}

func (s *Server) RefreshCount() int {
	return s.count(methodPath(http.MethodPost, RefreshEP))
}

func (s *Server) SwitchCount() int {
	return s.count(methodPath(http.MethodPost, SwitchEP))
}

func (s *Server) VerifyCount() int {
	return s.count(methodPath(http.MethodPost, VerifyEP))
}

func (s *Server) APIKeysListCount() int {
	return s.count(methodPath(http.MethodGet, APIKeysEP))
}

func (s *Server) APIKeysDetailCount(id string) int {
	return s.count(methodPath(http.MethodGet, fullPath(APIKeysEP, id)))
}

func (s *Server) APIKeysCreateCount() int {
	return s.count(methodPath(http.MethodPost, APIKeysEP))
}

func (s *Server) APIKeysDeleteCount(id string) int {
	return s.count(methodPath(http.MethodDelete, fullPath(APIKeysEP, id)))
}

func (s *Server) APIKeysUpdateCount(id string) int {
	return s.count(methodPath(http.MethodPut, fullPath(APIKeysEP, id)))
}

func (s *Server) APIKeysPermissionsCount() int {
	return s.count(methodPath(http.MethodPost, fullPath(APIKeysEP, "permissions")))
}

func (s *Server) ProjectsListCount() int {
	return s.count(methodPath(http.MethodGet, ProjectsEP))
}

func (s *Server) ProjectsCreateCount() int {
	return s.count(methodPath(http.MethodPost, ProjectsEP))
}

func (s *Server) ProjectsAccessCount() int {
	return s.count(methodPath(http.MethodPost, fullPath(ProjectsEP, "access")))
}

func (s *Server) ProjectsDetailCount(id string) int {
	return s.count(methodPath(http.MethodGet, fullPath(ProjectsEP, id)))
}

func (s *Server) OrganizationsCount(param string) int {
	return s.count(methodPath(http.MethodGet, fullPath(OrganizationsEP, param)))
}

func (s *Server) UsersDetailCount(id string) int {
	return s.count(methodPath(http.MethodGet, fullPath(UsersEP, id)))
}

func (s *Server) UsersUpdateCount(id string) int {
	return s.count(methodPath(http.MethodPut, fullPath(UsersEP, id)))
}

func (s *Server) UsersListCount() int {
	return s.count(methodPath(http.MethodGet, UsersEP))
}

func (s *Server) UsersRemoveCount(id string) int {
	return s.count(methodPath(http.MethodDelete, fullPath(UsersEP, id)))
}

func (s *Server) UsersRemoveConfirmCount(id string) int {
	return s.count(methodPath(http.MethodPost, fullPath(UsersEP, id+"/confirm")))
}

func (s *Server) UsersRoleUpdateCount(id string) int {
	return s.count(methodPath(http.MethodPost, fullPath(UsersEP, id)))
}

func (s *Server) InvitesPreviewCount(token string) int {
	return s.count(methodPath(http.MethodGet, fullPath(InvitesEP, token)))
}

func (s *Server) InvitesCreateCount() int {
	return s.count(methodPath(http.MethodPost, InvitesEP))
}
