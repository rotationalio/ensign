package service

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

// Option allows users to set optional arguments on the server when creating it.
type Option func(*options)

type options struct {
	mode   string
	server *http.Server
	router *gin.Engine
}

func WithMode(mode string) Option {
	return func(opts *options) {
		opts.mode = mode
	}
}

func WithServer(srv *http.Server) Option {
	return func(opts *options) {
		opts.server = srv
	}
}

func WithTestServer(srv *httptest.Server) Option {
	return func(opts *options) {
		opts.server = &http.Server{
			Addr:              srv.Config.Addr,
			Handler:           srv.Config.Handler,
			TLSConfig:         srv.TLS,
			ErrorLog:          nil,
			ReadHeaderTimeout: 20 * time.Second,
			WriteTimeout:      20 * time.Second,
			IdleTimeout:       30 * time.Second,
		}
	}
}

func WithRouter(router *gin.Engine) Option {
	return func(opts *options) {
		opts.router = router
	}
}

func newOptions(opts ...Option) *options {
	conf := &options{}
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}
