package buffer

import (
	"context"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
)

// Buffers are chunks of memory that are used to store events while they are being
// processed. For example a publisher buffer may hold events while they are being
// committed, written to disk, and copied to a dispatcher buffer. A dispatcher buffer
// is used to retrieve events from disk or from the publisher buffer and cache them
// while the event is being streamed to the subscriber. Buffers are intended to help
// Ensign manage the memory utilization of a single node.
//
// There are two basic operations for the buffer: writing and reading events. Buffers
// act as FIFO queues: a write pushes the event to the back of the queue and a Read
// pops the event off of the front of the queue. Buffers are intended for use by only
// a single go routine (otherwise you should just use channels) and are not thread safe.
//
// TODO: Should the interface use generics or should it be type specific for events?
// TODO: Should buffers hold values or pointers, how do we balance memory and performance?
type Buffer interface {
	// Write pushes an event to the back of the queue. Write may block until the buffer
	// has room for the event; use the context to specify timeouts if needed. If the
	// buffer is full or for some reason cannot push the event onto the queue it should
	// return an error, including a timeout error if necessary.
	Write(context.Context, *api.EventWrapper) error

	// Read pops an event off the top of the queue and returns it. Note that the event
	// returned is released from the buffer after it is read and cannot be read again.
	Read(context.Context) (*api.EventWrapper, error)
}
