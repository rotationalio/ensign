package broker

import (
	"fmt"
	"sync"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
)

const BufferSize = 16384

func New() *Broker {
	return &Broker{
		wg:    &sync.WaitGroup{},
		pubs:  make(map[rlid.RLID]chan<- bool),
		subs:  make(map[rlid.RLID]chan<- *api.EventWrapper),
		rlids: rlid.Sequence(0),
	}
}

// Every Ensign node is composed of a single Broker routine that collects events from
// publisher handlers, commits the events through consensus, writes the events to disk,
// and ensures that registered consumer groups receive any published events they are
// subscribed to. Essentially, the Broker fans in events from multiple, concurrent
// publisher streams, ensures the events are replicated and written, then fans out the
// events to one or more subscriber streams. The Broker uses an internal buffer that
// applies backpressure to the publisher streams when the buffer is full.
type Broker struct {
	sync.RWMutex
	inQ   chan<- inQ                             // input queue - incoming events from publishers are written here.
	wg    *sync.WaitGroup                        // wait for go routines to finish on shutdown
	pubs  map[rlid.RLID]chan<- bool              // registered publishers with an event callback channel.
	subs  map[rlid.RLID]chan<- *api.EventWrapper // registered subscribers with an outgoing event queue.
	rlids rlid.Sequence                          // used to generate publisher and subscriber IDs
}

// Run the broker; any fatal errors will be sent on the specified channel.
func (b *Broker) Run(errc chan<- error) {
	b.Lock()
	defer b.Unlock()

	queue := make(chan inQ, BufferSize)
	b.inQ = queue

	b.wg.Add(1)
	go func(inQ <-chan inQ) {
		defer b.wg.Done()

		var counter uint32
		for event := range inQ {
			counter++
			id := rlid.Make(counter)
			event.event.Id = id.Bytes()

			// TODO: concurrency issue here
			b.pubs[event.pubID] <- true

			o11y.Events.WithLabelValues("unk", "unk").Inc()
		}
	}(queue)
}

// Gracefully shutdown the broker. If a consensus or write operation is underway, then
// shutdown blocks until it is concluded. The broker then stops handling incoming events
// from publishers and closes all registered publishers and subscribers. This has the
// effect of closing any open event stream handlers.
func (b *Broker) Shutdown() error {
	b.Lock()
	defer b.Unlock()

	// If the broker is not running, ignore
	if !b.isRunning() {
		return nil
	}

	// Close all publishers to stop receiving events
	for pubID, ch := range b.pubs {
		close(ch)
		delete(b.pubs, pubID)
	}

	// Stop the internal go routines from handling any events
	close(b.inQ)
	b.wg.Wait()
	b.inQ = nil

	// Close all subscribers/consumer groups
	for subID, ch := range b.subs {
		close(ch)
		delete(b.subs, subID)
	}
	return nil
}

// Register a publisher to receive an ack/nack channel for events that are
// published using the publisher ID specified.
func (b *Broker) Register() (rlid.RLID, <-chan bool) {
	cb := make(chan bool, 1)

	b.Lock()
	publisherID := b.rlids.Next()
	b.pubs[publisherID] = cb
	b.Unlock()

	return publisherID, cb
}

// Publish an event from the specified publisher. When the event is committed, an
// acknowledgement or error is sent on the channel specified when registering.
func (b *Broker) Publish(publisherID rlid.RLID, event *api.EventWrapper) error {
	// TODO: if not running, error
	b.inQ <- inQ{publisherID, event}
	return nil
}

// Subscribe to events filtered by topic ids. All recent events will be sent on the
// event wrapper channel once they are committed.
func (b *Broker) Subscribe(topics ...ulid.ULID) (rlid.RLID, <-chan *api.EventWrapper) {
	events := make(chan *api.EventWrapper, 1)

	b.Lock()
	subscriberID := b.rlids.Next()
	b.subs[subscriberID] = events
	b.Unlock()

	return subscriberID, events
}

// Close either a publisher or subscriber so no events will be sent from the broker.
func (b *Broker) Close(id rlid.RLID) error {
	b.Lock()
	defer b.Unlock()

	if cb, ok := b.pubs[id]; ok {
		close(cb)
		delete(b.pubs, id)
		return nil
	}

	if evts, ok := b.subs[id]; ok {
		close(evts)
		delete(b.subs, id)
		return nil
	}

	return fmt.Errorf("no broker with id %q", id)
}

func (b *Broker) NumPublishers() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.pubs)
}

func (b *Broker) NumSubscribers() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.subs)
}

// Returns true if the broker has been started, false otherwise. Not thread-safe.
func (b *Broker) isRunning() bool {
	return b.inQ != nil
}

type inQ struct {
	pubID rlid.RLID
	event *api.EventWrapper
}
