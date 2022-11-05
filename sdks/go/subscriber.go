package ensign

import (
	"sync"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
)

type subscriber struct {
	sync.RWMutex
	stream api.Ensign_SubscribeClient
	send   chan *api.Subscription
	recv   []chan<- *api.Event
	stop   chan struct{}
	wg     sync.WaitGroup
	errc   chan error
}

var _ Subscriber = &subscriber{}

func (c *subscriber) Subscribe() (<-chan *api.Event, error) {
	sub := make(chan *api.Event, BufferSize)
	c.Lock()
	defer c.Unlock()
	c.recv = append(c.recv, sub)
	return sub, nil
}

func (c *subscriber) Ack(id string) error {
	c.send <- &api.Subscription{
		Embed: &api.Subscription_Ack{
			Ack: &api.Ack{
				Id: id,
			},
		},
	}
	return nil
}

func (c *subscriber) Nack(id string, err error) error {
	nack := &api.Nack{
		Id: id,
	}
	if err != nil {
		nack.Error = err.Error()
	}

	c.send <- &api.Subscription{
		Embed: &api.Subscription_Nack{
			Nack: nack,
		},
	}
	return nil
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
	for _, sub := range c.recv {
		close(sub)
	}
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

		c.RLock()
		for _, sub := range c.recv {
			sub <- e
		}
		c.RUnlock()
	}
}
