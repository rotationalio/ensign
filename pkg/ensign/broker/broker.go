package broker

import (
	"sync"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const BufferSize = 16384

func New() *Broker {
	return &Broker{
		wg:    &sync.WaitGroup{},
		pubs:  make(map[rlid.RLID]chan<- PublishResult),
		subs:  make(map[rlid.RLID]subscription),
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
	inQ   chan<- incoming                    // input queue - incoming events from publishers are written here.
	wg    *sync.WaitGroup                    // wait for go routines to finish on shutdown.
	pubs  map[rlid.RLID]chan<- PublishResult // registered publishers with an event callback channel.
	subs  map[rlid.RLID]subscription         // registered subscribers with an outgoing event queue.
	rlids rlid.Sequence                      // used to generate publisher and subscriber IDs
}

// Run the broker; any fatal errors will be sent on the specified channel.
func (b *Broker) Run(errc chan<- error) {
	b.Lock()
	defer b.Unlock()

	// TODO: fetch list of topics

	inQ := make(chan incoming, BufferSize)
	outQ := make(chan *api.EventWrapper, BufferSize)
	b.inQ = inQ

	b.wg.Add(2)
	go b.handleIncoming(inQ, outQ)
	go b.handleOutgoing(outQ)
}

func (b *Broker) handleIncoming(inQ <-chan incoming, outQ chan<- *api.EventWrapper) {
	defer b.wg.Done()
	defer close(outQ)

	seq := rlid.Sequence(0)
	for incoming := range inQ {
		// TODO: sequence RLIDs over topic offset instead of globally.
		incoming.event.Id = seq.Next().Bytes()
		incoming.event.Committed = timestamppb.Now()

		// TODO: write event to disk
		// TODO: consensus
		// TODO: update topic metadata

		// Send event on the outgoing queue
		outQ <- incoming.event

		// Send ack back to the publisher
		b.RLock()
		if cb, ok := b.pubs[incoming.pubID]; ok {
			res := PublishResult{
				LocalID:   incoming.event.LocalId,
				Committed: incoming.event.Committed.AsTime(),
			}

			// Non-blocking send so non-responding publishers don't hurt performance.
			select {
			case cb <- res:
			default:
			}
		}
		b.RUnlock()

		// Update metrics with number events
		// TODO: update label values with topic name, publisher ID, node, and region
		if o11y.Events != nil {
			o11y.Events.WithLabelValues("unk", "unk").Inc()
		}
	}
}

func (b *Broker) handleOutgoing(outQ <-chan *api.EventWrapper) {
	defer b.wg.Done()
	for event := range outQ {
		sends := 0
		nsubs := 0

		// Compute the topicID for the event
		// TODO: how to handle topicID parsing errors?
		topicID, _ := event.ParseTopicID()

		b.RLock()
		for _, sub := range b.subs {
			// Match the topic filter
			if _, ok := sub.topics[topicID]; !ok {
				continue
			}

			// Non-blocking send to prevent slow subscribers from interupting performance
			nsubs++
			select {
			case sub.out <- event:
				sends++
			default:
			}
		}
		b.RUnlock()
		log.Trace().Int("subs", sends).Bytes("id", event.Id).Int("dropped", nsubs-sends).Msg("event handled")
	}
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
	for subID, subscription := range b.subs {
		close(subscription.out)
		delete(b.subs, subID)
	}
	return nil
}

// Register a publisher to receive an ack/nack channel for events that are
// published using the publisher ID specified.
func (b *Broker) Register() (rlid.RLID, <-chan PublishResult) {
	cb := make(chan PublishResult, 1)

	b.Lock()
	publisherID := b.rlids.Next()
	b.pubs[publisherID] = cb
	b.Unlock()

	return publisherID, cb
}

// Publish an event from the specified publisher. When the event is committed, an
// acknowledgement or error is sent on the channel specified when registering.
func (b *Broker) Publish(publisherID rlid.RLID, event *api.EventWrapper) error {
	// The readlock synchronizes access to isRunning and to inQ to make sure we're not
	// sending on a closed channel.
	b.RLock()
	defer b.RUnlock()

	// If not running, error (prevent panics from send on closed channel)
	if !b.isRunning() {
		return ErrBrokerNotRunning
	}

	b.inQ <- incoming{publisherID, event}
	return nil
}

// Subscribe to events filtered by topic ids. All recent events will be sent on the
// event wrapper channel once they are committed.
func (b *Broker) Subscribe(topics ...ulid.ULID) (rlid.RLID, <-chan *api.EventWrapper) {
	events := make(chan *api.EventWrapper, 1)
	sub := subscription{
		topics: make(map[ulid.ULID]struct{}, len(topics)),
		out:    events,
	}

	for _, topic := range topics {
		sub.topics[topic] = struct{}{}
	}

	b.Lock()
	subscriberID := b.rlids.Next()
	b.subs[subscriberID] = sub
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

	if sub, ok := b.subs[id]; ok {
		close(sub.out)
		delete(b.subs, id)
		return nil
	}

	return ErrUnknownID
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
