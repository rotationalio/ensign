package broker

import (
	"sync"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const BufferSize = 16384

func New(events store.EventStore) *Broker {
	return &Broker{
		wg:     &sync.WaitGroup{},
		pubs:   make(map[rlid.RLID]chan<- PublishResult),
		subs:   make(map[rlid.RLID]subscription),
		rlids:  &rlid.LockedSequence{},
		events: events,
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
	inQ    chan<- incoming                    // input queue - incoming events from publishers are written here.
	wg     *sync.WaitGroup                    // wait for go routines to finish on shutdown.
	pubmu  sync.RWMutex                       // guards the pubs map and the broker state
	pubs   map[rlid.RLID]chan<- PublishResult // registered publishers with an event callback channel.
	submu  sync.RWMutex                       // guards the subs map and the broker state
	subs   map[rlid.RLID]subscription         // registered subscribers with an outgoing event queue.
	rlids  *rlid.LockedSequence               // used to generate publisher and subscriber IDs
	events store.EventStore                   // used to store events to disk
}

// Run the broker; any fatal errors will be sent on the specified channel.
func (b *Broker) Run(errc chan<- error) {
	b.Lock()
	defer b.Unlock()

	// If the broker is already running, ignore
	if b.isRunning() {
		return
	}

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
		// Create the publish result with the localID for handling
		result := PublishResult{LocalID: incoming.event.LocalId}

		// TODO: sequence RLIDs over topic offset instead of globally.
		incoming.event.Id = seq.Next().Bytes()

		// Write event to disk
		// NOTE: the insert will nil out the localID
		if err := b.events.Insert(incoming.event); err != nil {
			sentry.Error(nil).Err(err).Msg("could not insert event into database")
			result.Code = api.Nack_INTERNAL
			b.result(incoming, result)
			continue
		}

		// TODO: consensus

		// TODO: update topic metadata

		// Send event on the outgoing queue
		incoming.event.Committed = timestamppb.Now()
		outQ <- incoming.event

		// Send ack back to the publisher
		b.result(incoming, result)

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

		b.submu.RLock()
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
		b.submu.RUnlock()
		log.Trace().Int("subs", sends).Bytes("id", event.Id).Int("dropped", nsubs-sends).Msg("event handled")
	}
}

// Gracefully shutdown the broker. If a consensus or write operation is underway, then
// shutdown blocks until it is concluded. The broker then stops handling incoming events
// from publishers and closes all registered publishers and subscribers. This has the
// effect of closing any open event stream handlers.
func (b *Broker) Shutdown() error {
	// Acquire a lock to close the the inQ channel and signal that we're no longer running.
	b.Lock()

	// If the broker is not running, ignore
	if !b.isRunning() {
		b.Unlock()
		return nil
	}

	// Stop the internal go routines from handling any events
	// Make sure to mark inQ as nil so that no further events are published nor will any
	// publishers or subscribers be added to the broker.
	close(b.inQ)
	b.inQ = nil

	// Unlock and wait for the incoming and outgoing go routines to stop processing.
	// NOTE: Edge case: if Run() is called again weirdness can occur.
	b.Unlock()
	b.wg.Wait()

	// Relock to finalize the shutdown
	b.Lock()
	defer b.Unlock()

	// Close all publishers to stop receiving events
	for pubID, ch := range b.pubs {
		close(ch)
		delete(b.pubs, pubID)
	}

	// Close all subscribers/consumer groups
	for subID, subscription := range b.subs {
		close(subscription.out)
		delete(b.subs, subID)
	}
	return nil
}

// Register a publisher to receive an ack/nack channel for events that are
// published using the publisher ID specified. If the broker is not running, an error
// is returned so that the publisher can shutdown the stream.
func (b *Broker) Register() (rlid.RLID, <-chan PublishResult, error) {
	cb := make(chan PublishResult, BufferSize)
	publisherID := b.rlids.Next()

	b.pubmu.Lock()
	defer b.pubmu.Unlock()
	if !b.isRunning() {
		return rlid.RLID{}, nil, ErrBrokerNotRunning
	}

	b.pubs[publisherID] = cb
	return publisherID, cb, nil
}

// Publish an event from the specified publisher. When the event is committed, an
// acknowledgement or error is sent on the channel specified when registering.
func (b *Broker) Publish(publisherID rlid.RLID, event *api.EventWrapper) error {
	// The readlock synchronizes access to isRunning and to inQ to make sure we're not
	// sending on a closed channel.
	b.pubmu.RLock()
	defer b.pubmu.RUnlock()

	// If not running, error (prevent panics from send on closed channel)
	if !b.isRunning() {
		return ErrBrokerNotRunning
	}

	b.inQ <- incoming{publisherID, event}
	return nil
}

// Subscribe to events filtered by topic ids. All recent events will be sent on the
// event wrapper channel once they are committed. If the broker is not running an error
// is returned so that the consumer group can shutdown the stream.
func (b *Broker) Subscribe(topics ...ulid.ULID) (rlid.RLID, <-chan *api.EventWrapper, error) {
	subscriberID := b.rlids.Next()
	events := make(chan *api.EventWrapper, BufferSize)
	sub := subscription{
		topics: make(map[ulid.ULID]struct{}, len(topics)),
		out:    events,
	}

	for _, topic := range topics {
		sub.topics[topic] = struct{}{}
	}

	b.submu.Lock()
	defer b.submu.Unlock()

	// If the broker is not running, ignore
	if !b.isRunning() {
		return rlid.RLID{}, nil, ErrBrokerNotRunning
	}

	b.subs[subscriberID] = sub
	return subscriberID, events, nil
}

// Close either a publisher or subscriber so no events will be sent from the broker.
func (b *Broker) Close(id rlid.RLID) error {
	b.pubmu.Lock()
	if cb, ok := b.pubs[id]; ok {
		close(cb)
		delete(b.pubs, id)

		b.pubmu.Unlock()
		return nil
	}
	b.pubmu.Unlock()

	b.submu.Lock()
	if sub, ok := b.subs[id]; ok {
		close(sub.out)
		delete(b.subs, id)
		b.submu.Unlock()
		return nil
	}
	b.submu.Unlock()

	return ErrUnknownID
}

func (b *Broker) NumPublishers() int {
	b.pubmu.RLock()
	defer b.pubmu.RUnlock()
	return len(b.pubs)
}

func (b *Broker) NumSubscribers() int {
	b.submu.RLock()
	defer b.submu.RUnlock()
	return len(b.subs)
}

// Returns true if the broker has been started, false otherwise. Not thread-safe.
func (b *Broker) isRunning() bool {
	return b.inQ != nil
}

// Send a result ack or nack for an event back to the publisher.
func (b *Broker) result(in incoming, result PublishResult) {
	b.pubmu.RLock()
	if cb, ok := b.pubs[in.pubID]; ok {
		// Non-blocking send so non-responding publishers don't hurt performance.
		select {
		case cb <- result:
		default:
		}
	}
	b.pubmu.RUnlock()
}

// Locks both the pubmu and submu mutexes.
func (b *Broker) Lock() {
	b.pubmu.Lock()
	b.submu.Lock()
}

// Unlocks both the pubmu and submu mutexex.
func (b *Broker) Unlock() {
	b.submu.Unlock()
	b.pubmu.Unlock()
}

// RLocks the pubmu and submu mutexes.
func (b *Broker) RLock() {
	b.pubmu.RLock()
	b.submu.RLock()
}

// RUnlocks both the pubmu and submu mutexes.
func (b *Broker) RUnlock() {
	b.submu.RUnlock()
	b.pubmu.RUnlock()
}
