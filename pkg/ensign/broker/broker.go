package broker

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
)

const BufferSize = 16384

// Every Ensign node is composed of a single broker routine that collects events from
// publisher handlers, commits the events through consensus, writes the events to disk,
// and ensures that registered consumer groups receive any published events they are
// subscribed to.
type Broker interface {
	// Run the broker; any fatal errors will be sent on the specified channel.
	Run(errc chan<- error)

	// Register a publisher to receive an ack/nack channel for events that are
	// published using the publisher ID specified.
	Register(publisherID rlid.RLID) <-chan bool

	// Publish an event from the specified publisher. When the event is committed, an
	// acknowledgement or error is sent on the channel specified when registering.
	Publish(publisherID rlid.RLID, event *api.EventWrapper) error

	// Subscribe to events filtered by topic ids. All recent events will be sent on the
	// event wrapper channel once they are committed.
	Subscribe(subscriberID rlid.RLID, topics ...ulid.ULID) <-chan *api.EventWrapper

	// Close either a publisher or subscriber so no events will be sent from the broker.
	Close(id rlid.RLID) error
}

func New() Broker {
	return &broker{
		pubs: make(map[rlid.RLID]chan<- bool),
		subs: make(map[rlid.RLID]chan<- *api.EventWrapper),
	}
}

type broker struct {
	inQ  chan<- inQ
	pubs map[rlid.RLID]chan<- bool
	subs map[rlid.RLID]chan<- *api.EventWrapper
}

// Run the broker; any fatal errors will be sent on the specified channel.
func (b *broker) Run(errc chan<- error) {
	queue := make(chan inQ, BufferSize)
	b.inQ = queue

	go func(inQ <-chan inQ) {
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

// Register a publisher to receive an ack/nack channel for events that are
// published using the publisher ID specified.
func (b *broker) Register(publisherID rlid.RLID) <-chan bool {
	cb := make(chan bool, 1)
	b.pubs[publisherID] = cb
	return cb
}

// Publish an event from the specified publisher. When the event is committed, an
// acknowledgement or error is sent on the channel specified when registering.
func (b *broker) Publish(publisherID rlid.RLID, event *api.EventWrapper) error {
	// TODO: if not running, error
	b.inQ <- inQ{publisherID, event}
	return nil
}

// Subscribe to events filtered by topic ids. All recent events will be sent on the
// event wrapper channel once they are committed.
func (b *broker) Subscribe(subscriberID rlid.RLID, topics ...ulid.ULID) <-chan *api.EventWrapper {
	events := make(chan *api.EventWrapper, 1)
	b.subs[subscriberID] = events
	return events
}

// Close either a publisher or subscriber so no events will be sent from the broker.
func (b *broker) Close(id rlid.RLID) error {
	if cb, ok := b.pubs[id]; ok {
		close(cb)
		return nil
	}

	if evts, ok := b.subs[id]; ok {
		close(evts)
		return nil
	}

	return fmt.Errorf("no broker with id %q", id)
}

type inQ struct {
	pubID rlid.RLID
	event *api.EventWrapper
}
