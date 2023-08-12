package broker

import (
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PublishResult is sent back to the publisher stream that created an event to let it
// send an ack/nack back to to the client. If the event was correctly committed and
// emitted the result will contain an Ack message. If the event was unable to be
// processed by the broker (e.g. not committed, not written, etc.) then the result will
// contain a Nack message.
type PublishResult struct {
	LocalID   []byte        // The localID on the event sent by the publisher for client-side correlation
	Committed time.Time     // The timestamp the event was committed (if it was committed)
	Code      api.Nack_Code // The reason why the result errored; if not unknown the result is treated as an error
	Error     string        // An error message, should be set if the nack code is set
}

// Returns true if the reply is an ack, false if it is a nack
func (p PublishResult) IsAck() bool {
	return p.Code <= 0
}

// Returns true if the reply is a nack, fals if it is an ack
func (p PublishResult) IsNack() bool {
	return p.Code > 0
}

// Reply composes a publisher reply to return to the client via the publish stream. If
// the Code is > 0 (e.g. is not NACK_UNKNOWN) then a Nack is returned; otherwise an Ack
// is returned. This method performs no data validation other than to set the error
// message to a standard message for the code if it is a nack.
func (p PublishResult) Reply() *api.PublisherReply {
	// Return a Nack if there is an error code
	if p.IsNack() {
		return &api.PublisherReply{
			Embed: &api.PublisherReply_Nack{
				Nack: p.Nack(),
			},
		}
	}

	return &api.PublisherReply{
		Embed: &api.PublisherReply_Ack{
			Ack: p.Ack(),
		},
	}
}

func (p PublishResult) Ack() *api.Ack {
	return &api.Ack{
		Id:        p.LocalID,
		Committed: timestamppb.New(p.Committed),
	}
}

func (p PublishResult) Nack() *api.Nack {
	err := p.Error
	if err == "" {
		// Use "standard" nack error messages
		err = api.DefaultNackMessage(p.Code)
	}

	return &api.Nack{
		Id:    p.LocalID,
		Code:  p.Code,
		Error: err,
	}
}

// An incoming event is one that needs to be processed by the event handler and contains
// the publisher ID so that the result is sent back to the correct publisher.
type incoming struct {
	pubID rlid.RLID
	event *api.EventWrapper
}

// A subscription includes the topic filter for events and the channel to send those
// events on so that they can get back to the subscriber.
type subscription struct {
	topics map[ulid.ULID]struct{}
	out    chan<- *api.EventWrapper
}
