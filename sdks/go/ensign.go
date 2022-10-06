package ensign

import (
	"context"
	"crypto/tls"

	"github.com/kelseyhightower/envconfig"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const BufferSize = 128

// Client manages the credentials and connection to the ensign server.
type Client struct {
	opts *Options
	cc   *grpc.ClientConn
	api  api.EnsignClient
}

type Options struct {
	Endpoint     string `default:"flagship.rotational.app:443"`
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Insecure     bool   `default:"false"`
}

type Publisher struct {
	stream api.Ensign_PublishClient
	send   chan *api.Event
	recv   chan *api.Publication
	errc   chan error
}

type Subscriber struct {
	stream api.Ensign_SubscribeClient
	send   chan *api.Subscription
	recv   chan *api.Event
	errc   chan error
}

func New(opts *Options) (client *Client, err error) {
	if opts == nil {
		if err = envconfig.Process("ensign", &opts); err != nil {
			return nil, err
		}
	}

	dial := make([]grpc.DialOption, 0, 1)
	if opts.Insecure {
		dial = append(dial, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dial = append(dial, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	client = &Client{opts: opts}
	if client.cc, err = grpc.Dial(opts.Endpoint, dial...); err != nil {
		return nil, err
	}

	client.api = api.NewEnsignClient(client.cc)
	return client, nil
}

func (c *Client) Close() (err error) {
	defer func() {
		c.cc = nil
		c.api = nil
	}()

	if err = c.cc.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Publish(ctx context.Context) (pub *Publisher, err error) {
	pub = &Publisher{
		send: make(chan *api.Event, BufferSize),
		recv: make(chan *api.Publication, BufferSize),
		errc: make(chan error, 1),
	}
	if pub.stream, err = c.api.Publish(ctx); err != nil {
		return nil, err
	}

	// Start go routines
	go pub.sender()
	go pub.recver()

	return pub, nil
}

func (c *Client) Subscribe(ctx context.Context) (sub *Subscriber, err error) {
	sub = &Subscriber{
		send: make(chan *api.Subscription, BufferSize),
		recv: make(chan *api.Event, BufferSize),
		errc: make(chan error, 1),
	}
	if sub.stream, err = c.api.Subscribe(ctx); err != nil {
		return nil, err
	}

	// Start go routines
	go sub.sender()
	go sub.recver()

	return sub, nil
}

func (c *Publisher) Publish(e *api.Event) {
	c.send <- e
}

func (c *Publisher) Err() error {
	select {
	case err := <-c.errc:
		return err
	default:
	}
	return nil
}

func (c *Subscriber) Subscribe() *api.Event {
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

func (c *Subscriber) Err() error {
	select {
	case err := <-c.errc:
		return err
	default:
	}
	return nil
}

func (c *Publisher) sender() {
	for e := range c.send {
		if err := c.stream.Send(e); err != nil {
			c.errc <- err
			return
		}
	}
}

func (c *Publisher) recver() {
	for {
		ack, err := c.stream.Recv()
		if err != nil {
			c.errc <- err
			return
		}
		c.recv <- ack
	}
}

func (c *Subscriber) sender() {
	for e := range c.send {
		if err := c.stream.Send(e); err != nil {
			c.errc <- err
			return
		}
	}
}

func (c *Subscriber) recver() {
	for {
		e, err := c.stream.Recv()
		if err != nil {
			c.errc <- err
			return
		}
		c.recv <- e
	}
}
