package contexts_test

import (
	"context"
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestStream(t *testing.T) {
	mock := &MockStream{}
	stream := contexts.Stream(mock, context.WithValue(mock.Context(), contexts.KeyUnknown, "bar"))

	ctx := stream.Context()
	require.Equal(t, "bar", ctx.Value(contexts.KeyUnknown).(string))

	mock.cancel()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}

type MockStream struct {
	grpc.ServerStream
	ctx    context.Context
	cancel context.CancelFunc
}

func (m *MockStream) Context() context.Context {
	if m.ctx == nil {
		m.ctx, m.cancel = context.WithCancel(context.Background())
	}
	return m.ctx
}
