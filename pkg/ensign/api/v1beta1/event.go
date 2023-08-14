package api

import (
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/protobuf/proto"
)

func (w *EventWrapper) Wrap(e *Event) (err error) {
	if w.Event, err = proto.Marshal(e); err != nil {
		return err
	}
	return nil
}

func (w *EventWrapper) Unwrap() (e *Event, err error) {
	if len(w.Event) == 0 {
		return nil, errors.New("event wrapper contains no event")
	}

	e = &Event{}
	if err = proto.Unmarshal(w.Event, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (w *EventWrapper) ParseTopicID() (ulid.ULID, error) {
	return ulids.Parse(w.TopicId)
}

// Returns the user-specified client ID if set, otherwise returns the publisher ID.
func (p *Publisher) ResolveClientID() string {
	if p.ClientId != "" {
		return p.ClientId
	}
	return p.PublisherId
}
