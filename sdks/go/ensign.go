package ensign

import (
	"context"
	"crypto/tls"
	"io"

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

// Options allows users to configure their connection to ensign.
type Options struct {
	Endpoint     string `default:"flagship.rotational.app:443"`
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Insecure     bool   `default:"false"`
}

// Publisher is a low level interface for sending events to a topic or a group of topics
// that have been defined in Ensign services.
type Publisher interface {
	io.Closer
	Errorer
	Publish(events ...*api.Event)
}

type Subscriber interface {
	io.Closer
	Errorer
	Subscribe() *api.Event
}

func New(opts *Options) (client *Client, err error) {
	if opts == nil {
		if err = envconfig.Process("ensign", &opts); err != nil {
			return nil, err
		}
	}

	client = &Client{opts: opts}
	if err = client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Connect(opts ...grpc.DialOption) (err error) {
	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0, 1)
		if c.opts.Insecure {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		}
	}

	if c.cc, err = grpc.Dial(c.opts.Endpoint, opts...); err != nil {
		return err
	}

	c.api = api.NewEnsignClient(c.cc)
	return nil
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

func (c *Client) Publish(ctx context.Context) (_ Publisher, err error) {
	pub := &publisher{
		send: make(chan *api.Event, BufferSize),
		recv: make(chan *api.Publication, BufferSize),
		stop: make(chan struct{}, 1),
		errc: make(chan error, 1),
	}
	if pub.stream, err = c.api.Publish(ctx); err != nil {
		return nil, err
	}

	// Start go routines
	pub.wg.Add(2)
	go pub.sender()
	go pub.recver()

	return pub, nil
}

func (c *Client) Subscribe(ctx context.Context) (_ Subscriber, err error) {
	sub := &subscriber{
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
