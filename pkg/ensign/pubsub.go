package ensign

import (
	"github.com/google/uuid"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/buffer"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rs/zerolog/log"
)

const BufferSize = 16384

// PubSub is a simple dispatcher that has a publish queue that fans in events from
// different publisher streams and assigns event ids then fans the events out to send
// them to one or more subscriber streams with an outgoing buffer. Backpressure is
// applied to the publisher streams when the buffers get full.
type PubSub struct {
	inQ     chan inQ
	outQ    buffer.Channel
	counter uint32
	subs    map[uuid.UUID]buffer.Channel
	subQ    chan subQ
	finQ    chan uuid.UUID
}

type inQ struct {
	res   chan rlid.RLID
	event *api.EventWrapper
}

type subQ struct {
	id  uuid.UUID
	buf buffer.Channel
}

func NewPubSub() (ps *PubSub) {
	ps = &PubSub{
		inQ:     make(chan inQ, BufferSize),
		outQ:    make(chan *api.EventWrapper, BufferSize),
		counter: 0,
		subs:    make(map[uuid.UUID]buffer.Channel),
		subQ:    make(chan subQ),
		finQ:    make(chan uuid.UUID),
	}
	go ps.pub()
	go ps.sub()
	return ps
}

// Handles incoming events being published; loops on the incoming queue, assigns a
// monotonically increasing event ID then puts the even on the outgoing queue.
func (ps *PubSub) pub() {
	for i := range ps.inQ {
		ps.counter++
		id := rlid.Make(ps.counter)
		i.event.Id = id.Bytes()
		ps.outQ <- i.event
		i.res <- id
		o11y.Events.WithLabelValues("unk", "unk").Inc()
	}
}

// Handles outgoing events being sent to subscribers; loops on the outgoing queue and
// and sends the event to all connected subscribers. If there are no subscribers then
// the event is dropped.
func (ps *PubSub) sub() {
	for {
		select {
		// Handle events
		case e := <-ps.outQ:
			sends := 0
			for _, sub := range ps.subs {
				// Non-blocking send to prevent slow subscribers from interupting performance
				select {
				case sub <- e:
					sends++
				default:
				}
			}
			log.Debug().Int("subs", sends).Bytes("id", e.Id).Bool("dropped", sends == 0).Msg("event handled")
		// Handle new subscribers
		case sub := <-ps.subQ:
			ps.subs[sub.id] = sub.buf
			// Handle closing subscribers
		case id := <-ps.finQ:
			if buf, ok := ps.subs[id]; ok {
				close(buf)
				delete(ps.subs, id)
			}
		}
	}
}

// Publish puts an event on the input queue and waits until the event is handled and an
// ID is assigned then returns the ID of the event to the caller.
func (ps *PubSub) Publish(event *api.EventWrapper) rlid.RLID {
	q := inQ{
		event: event,
		res:   make(chan rlid.RLID, 1),
	}
	ps.inQ <- q
	return <-q.res
}

// Subscribe creates returns a channel that the caller can use to fetch events off of
// the in memory queue from. It also returns an ID so that the caller can close and
// cleanup the channel when it is done listenting for events.
func (ps *PubSub) Subscribe() (id uuid.UUID, c buffer.Channel) {
	id = uuid.New()
	c = make(buffer.Channel, BufferSize)

	ps.subQ <- subQ{id, c}
	return id, c
}

// Finish closes a subscribe channel and removes it so that the PubSub no longer sends.
func (ps *PubSub) Finish(id uuid.UUID) {
	ps.finQ <- id
}
