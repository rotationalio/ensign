package ensign

import (
	"sync"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
)

type subscriber struct {
	stream api.Ensign_SubscribeClient
	send   chan *api.Subscription
	recv   chan *api.Event
	stop   chan struct{}
	wg     sync.WaitGroup
	errc   chan error
}

var _ Subscriber = &subscriber{}

func (c *subscriber) Subscribe() *api.Event {
	// Block until event comes from ensign then send ack back immediately
	e := <-c.recv
	c.send <- &api.Subscription{
		Embed: &api.Subscription_Ack{
			Ack: &api.Ack{
				Id: e.Id,
			},
		},
	}
	return e
}

func (c *subscriber) Err() error {
	select {
	case err := <-c.errc:
		return err
	default:
	}
	return nil
}

func (c *subscriber) Close() error {
	// Cannot call CloseSend concurrently with send message.
	// Send stop signals to sender and recver routines
	c.stop <- struct{}{}
	close(c.send)

	c.wg.Wait()
	close(c.recv)
	return c.stream.CloseSend()
}

func (c *subscriber) sender() {
	defer c.wg.Done()
	for e := range c.send {
		if err := c.stream.Send(e); err != nil {
			c.errc <- err
			return
		}
	}
}

func (c *subscriber) recver() {
	defer c.wg.Done()
	for {
		select {
		case <-c.stop:
			return
		default:
		}

		e, err := c.stream.Recv()
		if err != nil {
			c.errc <- err
			return
		}
		c.recv <- e
	}
}
