package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/rotationalio/ensign/pkg/utils/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService is a mock for testing the service interfaces.
type TestService struct {
	Calls        []string
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
	if s.Calls == nil {
		s.Calls = make([]string, 0, 5)
	}
	s.Calls = append(s.Calls, method)
}

// SimpleService is a mock but only implements the required methods and none of the
// other interface methods defined in the service handler.
type SimpleService struct {
	Calls      []string
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
	if s.Calls == nil {
		s.Calls = make([]string, 0, 5)
	}
	s.Calls = append(s.Calls, method)
}

func TestServer(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &TestService{}

	server.Register(mock)

	require.Equal(t, gin.TestMode, gin.Mode())
	require.Len(t, mock.Calls, 0, "no service calls should happen during registration")

	mock.OnInitialize = func() error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.False(t, server.IsReady())
		assert.True(t, server.Started().IsZero())
		return nil
	}

	mock.OnSetup = func() error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.False(t, server.IsReady())
		assert.True(t, server.Started().IsZero())
		return nil
	}

	mock.OnRoutes = func(router *gin.Engine) error {
		// Must use assert in go routines not require
		assert.True(t, server.IsHealthy())
		assert.False(t, server.IsReady())
		assert.True(t, server.Started().IsZero())

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
		assert.False(t, server.Started().IsZero())

		started <- true
		return nil
	}

	mock.OnShutdown = func(ctx context.Context) error {
		// Must use assert in go routines not require
		require.False(t, server.IsHealthy())
		require.False(t, server.IsReady())
		require.False(t, server.Started().IsZero())
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
	require.Equal(t, []string{"initialize", "setup", "routes", "started"}, mock.Calls)

	// Should be able to make a request to the route setup during Routes()
	rep, err := http.Get(server.URL() + "/foo")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)

	// Wait for the server to shutdown
	server.Shutdown(context.Background())
	err = <-errc
	require.NoError(t, err, "server did not gracefully shutdown")

	// When the server is shutdown it should have also called the service shutdown.
	require.Equal(t, []string{"initialize", "setup", "routes", "started", "shutdown"}, mock.Calls)
}

func TestSimpleService(t *testing.T) {
	server := service.New(":0", service.WithMode(gin.TestMode))
	mock := &SimpleService{}
	server.Register(mock)

	// Serve is a blocking function so call it in its own go routine
	errc := make(chan error)
	go func() {
		errc <- server.Serve()
	}()

	// Give the server time to startup in its own go routine
	time.Sleep(15 * time.Millisecond)

	// Before the server is shutdown service should have had the following calls in order
	require.Equal(t, []string{"routes"}, mock.Calls)

	// Wait for the server to shutdown
	server.Shutdown(context.Background())
	err := <-errc
	require.NoError(t, err, "server did not gracefully shutdown")

	// When the server is shutdown it should have also called the service shutdown.
	require.Equal(t, []string{"routes", "shutdown"}, mock.Calls)
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
	require.Equal(t, []string{"initialize"}, mock.Calls)
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
	require.Equal(t, []string{"initialize", "setup"}, mock.Calls)

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
	require.Equal(t, []string{"initialize", "setup", "routes"}, mock.Calls)

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
	require.Equal(t, []string{"initialize", "setup", "routes", "started", "shutdown"}, mock.Calls)

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
	err := server.Shutdown(context.Background())

	require.ErrorIs(t, err, expectedErr)
	require.False(t, server.IsReady())
	require.Equal(t, []string{"initialize", "setup", "routes", "started", "shutdown"}, mock.Calls)

	_, err = http.Get(server.URL() + "/livez")
	require.Error(t, err, "server should be shutdown")
}
