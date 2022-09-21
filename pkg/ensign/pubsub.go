package ensign

import (
	"github.com/prometheus/client_golang/prometheus"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rs/zerolog/log"
)

const BufferSize = 10000

// PubSub is a simple dispatcher that has a publish queue that fans in events from
// different publisher streams and assigns event ids then fans the events out to send
// them to one or more subscriber streams with an outgoing buffer. Backpressure is
// applied to the publisher streams when the buffers get full.
type PubSub struct {
	inQ     chan inQ
	outQ    chan *api.Event
	counter uint32
	subs    []chan<- *api.Event
}

type inQ struct {
	res   chan rlid.RLID
	event *api.Event
}

func NewPubSub() (ps *PubSub) {
	ps = &PubSub{
		inQ:     make(chan inQ, BufferSize),
		outQ:    make(chan *api.Event, BufferSize),
		counter: 0,
		subs:    make([]chan<- *api.Event, 0),
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
		i.event.Id = id.String()
		ps.outQ <- i.event
		i.res <- id
		o11y.Events.With(prometheus.Labels{"node": "unk", "region": "unk"}).Add(1)
	}
}

// Handles outgoing events being sent to subscribers; loops on the outgoing queue and
// and sends the event to all connected subscribers. If there are no subscribers then
// the event is dropped.
func (ps *PubSub) sub() {
	for e := range ps.outQ {
		// TODO: add concurrency handling
		// HACK: current expectation is that we add all subscribers before publishing starts
		for _, sub := range ps.subs {
			sub <- e
		}
		log.Debug().Int("subs", len(ps.subs)).Str("id", e.Id).Bool("dropped", len(ps.subs) == 0).Msg("event handled")
	}
}

// Publish puts an event on the input queue and waits until the event is handled and an
// ID is assigned then returns the ID of the event to the caller.
func (ps *PubSub) Publish(event *api.Event) rlid.RLID {
	q := inQ{
		event: event,
		res:   make(chan rlid.RLID, 1),
	}
	ps.inQ <- q
	return <-q.res
}

// Subscribe creates returns a channel that the caller can use to fetch events off of
// the in memory queue from.
// TODO: add functionality to close and cleanup channels; it is currently unbounded.
func (ps *PubSub) Subscribe() <-chan *api.Event {
	c := make(chan *api.Event, BufferSize)
	ps.subs = append(ps.subs, c)
	return c
}
