package service_test

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/rotationalio/ensign/pkg/utils/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService is a mock for testing the service interfaces.
type TestService struct {
	sync.RWMutex
	calls        []string
	OnInitialize func() error
	OnSetup      func() error
	OnRoutes     func(router *gin.Engine) error
	OnStarted    func() error
	OnShutdown   func(ctx context.Context) error
}

func (s *TestService) Initialize() error {
	s.Incr("initialize")
	if s.OnInitialize != nil {
		return s.OnInitialize()
	}
	return nil
}

func (s *TestService) Setup() error {
	s.Incr("setup")
	if s.OnSetup != nil {
		return s.OnSetup()
	}
	return nil
}

func (s *TestService) Routes(router *gin.Engine) error {
	s.Incr("routes")
	if s.OnRoutes != nil {
		return s.OnRoutes(router)
	}
	return nil
}

func (s *TestService) Started() error {
	s.Incr("started")
	if s.OnStarted != nil {
		return s.OnStarted()
	}
	return nil
}

func (s *TestService) Shutdown(ctx context.Context) error {
	s.Incr("shutdown")
	if s.OnShutdown != nil {
		return s.OnShutdown(ctx)
	}
	return nil
}

func (s *TestService) Incr(method string) {
	s.Lock()
	defer s.Unlock()
	if s.calls == nil {
		s.calls = make([]string, 0, 5)
	}
	s.calls = append(s.calls, method)
}

func (s *TestService) Calls() []string {
	s.RLock()
	defer s.RUnlock()
	calls := make([]string, len(s.calls))
	copy(calls, s.calls)
	return calls
}

// SimpleService is a mock but only implements the required methods and none of the
// other interface methods defined in the service handler.
type SimpleService struct {
	sync.RWMutex
	calls      []string
	OnRoutes   func(router *gin.Engine) error
	OnShutdown func(ctx context.Context) error
}

func (s *SimpleService) Routes(router *gin.Engine) error {
	s.Incr("routes")
	if s.OnRoutes != nil {
		return s.OnRoutes(router)
	}
	return nil
}

func (s *SimpleService) Shutdown(ctx context.Context) error {
	s.Incr("shutdown")
	if s.OnShutdown != nil {
		return s.OnShutdown(ctx)
	}
	return nil
}

func (s *SimpleService) Incr(method string) {
	s.Lock()
	defer s.Unlock()
	if s.calls == nil {
		s.calls = make([]string, 0, 5)
	}
	s.calls = append(s.calls, method)
}

func (s *SimpleService) Calls() []string {
	s.RLock()
	defer s.RUnlock()
	calls := make([]string, len(s.calls))
	copy(calls, s.calls)
	return calls
}

func TestServer(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	server.Register(mock)

	require.Equal(t, gin.TestMode, gin.Mode())
	require.Len(t, mock.Calls(), 0, "no service calls should happen during registration")

	mock.OnInitialize = func() error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.False(t, server.IsReady())
		assert.True(t, server.StartTime().IsZero())
		return nil
	}

	mock.OnSetup = func() error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.False(t, server.IsReady())
		assert.True(t, server.StartTime().IsZero())
		return nil
	}

	mock.OnRoutes = func(router *gin.Engine) error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.False(t, server.IsReady())
		assert.True(t, server.StartTime().IsZero())

		router.GET("/foo", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"hello": "world"})
		})

		return nil
	}

	started := make(chan bool, 1)
	mock.OnStarted = func() error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.True(t, server.IsReady())
		assert.False(t, server.StartTime().IsZero())

		started <- true
		return nil
	}

	mock.OnShutdown = func(ctx context.Context) error {
		// Must use assert in go routines not require
		require.False(t, server.IsHealthy())
		require.False(t, server.IsReady())
		require.False(t, server.StartTime().IsZero())
		return nil
	}

	// Serve is a blocking function so call it in its own go routine
	errc := make(chan error)
	go func() {
		errc <- server.Serve()
	}()

	// Give the server time to startup in its own go routine
	select {
	case err := <-errc:
		require.NoError(t, err, "test could not continue")
	case <-started:
	}

	// Before the server is shutdown service should have had the following calls in order
	require.Equal(t, []string{"initialize", "setup", "routes", "started"}, mock.Calls())

	// Should be able to make a request to the route setup during Routes()
	rep, err := http.Get(server.URL() + "/foo")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)

	// Wait for the server to shutdown
	server.GracefulShutdown(context.Background())
	err = <-errc
	require.NoError(t, err, "server did not gracefully shutdown")

	// When the server is shutdown it should have also called the service shutdown.
	require.Equal(t, []string{"initialize", "setup", "routes", "started", "shutdown"}, mock.Calls())
}

func TestSimpleService(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &SimpleService{}
	server.Register(mock)

	started := make(chan bool)
	mock.OnRoutes = func(*gin.Engine) error {
		started <- true
		return nil
	}

	// Serve is a blocking function so call it in its own go routine
	errc := make(chan error)
	go func() {
		errc <- server.Serve()
	}()

	// Give the server time to startup in its own go routine
	<-started

	// Before the server is shutdown service should have had the following calls in order
	require.Equal(t, []string{"routes"}, mock.Calls())

	// Wait for the server to shutdown
	server.GracefulShutdown(context.Background())
	err := <-errc
	require.NoError(t, err, "server did not gracefully shutdown")

	// When the server is shutdown it should have also called the service shutdown.
	require.Equal(t, []string{"routes", "shutdown"}, mock.Calls())
}

func TestServerNoRegistration(t *testing.T) {
	// Should not be able to serve without registering
	server := service.New(":0", service.WithMode(gin.TestMode))
	err := server.Serve()
	require.ErrorIs(t, err, service.ErrNoServiceRegistered)
}

func TestInitializeError(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	expectedErr := errors.New("something wicked this way comes")
	server.Register(mock)
	mock.OnInitialize = func() error {
		return expectedErr
	}

	err := server.Serve()
	require.ErrorIs(t, err, expectedErr)
	require.False(t, server.IsReady())
	require.Equal(t, []string{"initialize"}, mock.Calls())
}

func TestSetupError(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	expectedErr := errors.New("something wicked this way comes")
	server.Register(mock)
	mock.OnSetup = func() error {
		return expectedErr
	}

	err := server.Serve()
	require.ErrorIs(t, err, expectedErr)
	require.False(t, server.IsReady())
	require.Equal(t, []string{"initialize", "setup"}, mock.Calls())

	_, err = http.Get(server.URL() + "/livez")
	require.Error(t, err, "server should be shutdown")
}

func TestRoutesError(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	expectedErr := errors.New("something wicked this way comes")
	server.Register(mock)
	mock.OnRoutes = func(*gin.Engine) error {
		return expectedErr
	}

	err := server.Serve()
	require.ErrorIs(t, err, expectedErr)
	require.False(t, server.IsReady())
	require.Equal(t, []string{"initialize", "setup", "routes"}, mock.Calls())

	_, err = http.Get(server.URL() + "/livez")
	require.Error(t, err, "server should be shutdown")
}

func TestStartedError(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	expectedErr := errors.New("something wicked this way comes")
	server.Register(mock)
	mock.OnStarted = func() error {
		return expectedErr
	}

	mock.OnShutdown = func(context.Context) error {
		return errors.New("secondary error")
	}

	err := server.Serve()
	require.ErrorIs(t, err, expectedErr)
	require.False(t, server.IsReady())
	require.Equal(t, []string{"initialize", "setup", "routes", "started", "shutdown"}, mock.Calls())

	merr, ok := err.(*multierror.Error)
	require.True(t, ok)
	require.Len(t, merr.Errors, 2)

	_, err = http.Get(server.URL() + "/livez")
	require.Error(t, err, "server should be shutdown")
}

func TestShutdownError(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	started := make(chan bool)
	expectedErr := errors.New("something wicked this way comes")
	server.Register(mock)
	mock.OnStarted = func() error {
		started <- true
		return nil
	}

	mock.OnShutdown = func(context.Context) error {
		return expectedErr
	}

	go func() {
		server.Serve()
	}()

	<-started
	err := server.GracefulShutdown(context.Background())

	require.ErrorIs(t, err, expectedErr)
	require.False(t, server.IsReady())
	require.Equal(t, []string{"initialize", "setup", "routes", "started", "shutdown"}, mock.Calls())

	_, err = http.Get(server.URL() + "/livez")
	require.Error(t, err, "server should be shutdown")
}
