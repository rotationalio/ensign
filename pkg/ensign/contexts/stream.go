package contexts

import (
	"context"

	"google.golang.org/grpc"
)

// Stream allows users to override the context on a grpc.ServerStream handler so that
// it returns a new context rather than the old context. It is advised to use the
// original stream's context as the new context's parent but this method does not
// enforce it and instead simply returns the context specified.
func Stream(s grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return &stream{s, ctx}
}

type stream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *stream) Context() context.Context {
	return s.ctx
}
