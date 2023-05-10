package buffer

import (
	"context"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
)

// A channel buffer is just what it says on the tin, it's a channel that you can send
// and receive events on. The reason for this type alias is to implement the Buffer
// interface so that outside callers can use it as an event buffer.
//
// To create this buffer use make as you would a channel, e.g. make(buffer.Channel, 1)
// will create a channel with a buffer size of 1.
type Channel chan *api.EventWrapper

// Compile time check that Channel implements the Buffer interface.
var _ Buffer = make(Channel, 1)

func (c Channel) Read(context.Context) (*api.EventWrapper, error) {
	return <-c, nil
}

func (c Channel) Write(_ context.Context, event *api.EventWrapper) error {
	c <- event
	return nil
}
