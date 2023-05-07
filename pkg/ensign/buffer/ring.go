package buffer

import (
	"context"
	"sync"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
)

type Ring struct {
	head  int                 // the index of the previously read element in the queue
	tail  int                 // the index of the previously written element in the queue
	size  int                 // the size of the queue to prevent calls to len()
	queue []*api.EventWrapper // the slice of events being buffered
}

func NewRing(size int) *Ring {
	return &Ring{
		head:  -1,
		tail:  -1,
		size:  size,
		queue: make([]*api.EventWrapper, size),
	}
}

// Compile time check that Ring implements the Buffer interface.
var _ Buffer = &Ring{}

func (b *Ring) Read(ctx context.Context) (e *api.EventWrapper, _ error) {
	if b.head == -1 {
		return nil, ErrBufferEmpty
	}

	// Condition for only one element
	if b.head == b.tail {
		e = b.queue[b.head]
		b.head = -1
		b.tail = -1
		return e, nil
	}

	e = b.queue[b.head]
	b.head = (b.head + 1) % b.size
	return e, nil
}

func (b *Ring) Write(ctx context.Context, event *api.EventWrapper) error {
	// Check if queue is empty
	if b.head == -1 {
		b.head = 0
		b.tail = 0
		b.queue[b.tail] = event
		return nil
	}

	tail := (b.tail + 1) % b.size
	if tail == b.head {
		return ErrBufferFull
	}

	b.tail = tail
	b.queue[b.tail] = event
	return nil
}

func NewLockingRing(size int) *LockingRing {
	return &LockingRing{
		Ring: Ring{
			head:  -1,
			tail:  -1,
			size:  size,
			queue: make([]*api.EventWrapper, size),
		},
	}
}

// A thread-safe ring buffer with a mutex guarding reads and writes.
type LockingRing struct {
	Ring
	sync.Mutex
}

func (b *LockingRing) Read(ctx context.Context) (e *api.EventWrapper, err error) {
	b.Lock()
	e, err = b.Ring.Read(ctx)
	b.Unlock()
	return e, err
}

func (b *LockingRing) Write(ctx context.Context, event *api.EventWrapper) (err error) {
	b.Lock()
	err = b.Ring.Write(ctx, event)
	b.Unlock()
	return err
}
