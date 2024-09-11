package bufconn

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufsize  = 1024 * 1024
	Endpoint = "passthrough://bufnet"
)

// Listener handles gRPC connections using an in-memory buffer that is useful for
// testing to prevent actual TCP network requests. Using a bufconn connection provides
// the most realistic gRPC server for tests that include serialization and
// deserialization of protocol buffers and an actual wire transfer between client and
// server. We prefer to use the bufconn over simply making method calls to the handlers.
type Listener struct {
	sock   *bufconn.Listener
	target string
}

// New creates a bufconn listener ready to attach servers and clients to. To provide a
// different target name (e.g. for mTLS buffers) use the WithTarget() dial option. You
// can also specify a different buffer size using WithBufferSize() or pass in an already
// instantiated bufconn.Listener using WithBuffer().
func New(opts ...DialOption) *Listener {
	sock := &Listener{}
	for _, opt := range opts {
		opt(sock)
	}

	if sock.target == "" {
		sock.target = Endpoint
	}

	if sock.sock == nil {
		sock.sock = bufconn.Listen(bufsize)
	}

	return sock
}

// Sock returns the server side of the bufconn connection.
func (l *Listener) Sock() net.Listener {
	return l.sock
}

// Close the bufconn listener and prevent either clients or servers from communicating.
func (l *Listener) Close() error {
	return l.sock.Close()
}

// Connect returns the client side of the bufconn connection.
func (l *Listener) Connect(opts ...grpc.DialOption) (cc *grpc.ClientConn, err error) {
	opts = append([]grpc.DialOption{grpc.WithContextDialer(l.Dialer)}, opts...)
	if cc, err = grpc.NewClient(l.target, opts...); err != nil {
		return nil, err
	}
	return cc, nil
}

// Dialer implements the ContextDialer interface for use with grpc.DialOptions
func (l *Listener) Dialer(context.Context, string) (net.Conn, error) {
	return l.sock.Dial()
}

// DialOption -- optional arguments for constructing a Listener.
type DialOption func(*Listener)

// WithTarget allows the user to change the "endpoint" of the connection from "bufnet"
// to some other endpoint. This is useful for tests that include mTLS to ensure that
// TLS certificate handling is correct.
func WithTarget(target string) DialOption {
	return func(lis *Listener) {
		lis.target = target
	}
}

// The default buffer size is 1MiB -- if ou need a larger buffer to send larger messages
// then specify this dial option with a larger size.
func WithBufferSize(size int) DialOption {
	return func(lis *Listener) {
		lis.sock = bufconn.Listen(size)
	}
}

// Allows you to pass an already instantiated grpc bufconn.Listener into the Listener.
func WithBuffer(sock *bufconn.Listener) DialOption {
	return func(lis *Listener) {
		lis.sock = sock
	}
}
