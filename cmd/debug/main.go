package main

import (
	"context"
	"crypto/rand"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/mimetype/v1beta1"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "ensign-debug"
	app.Version = pkg.Version()
	app.Usage = "client utilities to help debug an ensign server"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e"},
			Usage:   "endpoint of local ensign node to send requests to",
			Value:   "127.0.0.1:7777",
			EnvVars: []string{"ENSIGN_ENDPOINT"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "generate",
			Usage:  "generate a constant or fixed length stream of random event data",
			Before: connect,
			After:  disconnect,
			Action: generate,
			Flags: []cli.Flag{
				&cli.DurationFlag{
					Name:    "interval",
					Aliases: []string{"i"},
					Usage:   "the amount of time between events being published",
					Value:   250 * time.Millisecond,
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "the size in bytes of the event data generated",
					Value:   1024,
				},
			},
		},
		{
			Name:   "consume",
			Usage:  "subscribe to the straem an dconsume events",
			Before: connect,
			After:  disconnect,
			Action: consume,
			Flags:  []cli.Flag{},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("could not execute cli app")
	}
}

var (
	cc     *grpc.ClientConn
	client api.EnsignClient
)

func connect(c *cli.Context) (err error) {
	if cc, err = grpc.Dial(c.String("endpoint"), grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return cli.Exit(err, 1)
	}
	client = api.NewEnsignClient(cc)
	return nil
}

func disconnect(c *cli.Context) (err error) {
	defer func() {
		cc = nil
		client = nil
	}()

	if err = cc.Close(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func generate(c *cli.Context) (err error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var stream api.Ensign_PublishClient
	if stream, err = client.Publish(context.Background()); err != nil {
		return cli.Exit(err, 1)
	}

	// Create send and receive channels and go routines to manage the stream
	send := make(chan *api.Event, 2)
	recv := make(chan *api.Publication, 2)
	errc := make(chan error, 1)

	go func(send <-chan *api.Event) {
		for e := range send {
			if err := stream.Send(e); err != nil {
				errc <- err
				return
			}
		}
	}(send)

	go func(recv chan<- *api.Publication) {
		for {
			ack, err := stream.Recv()
			if err != nil {
				errc <- err
				return
			}
			recv <- ack
		}
	}(recv)

	size := c.Int("size")
	ticker := time.NewTicker(c.Duration("interval"))

primary:
	for {
		select {
		case ack := <-recv:
			log.Info().Str("id", ack.GetAck().GetId()).Msg("ack")
		case <-ticker.C:
			send <- &api.Event{
				TopicId:  "generator",
				Mimetype: mimetype.ApplicationOctetStream,
				Type: &api.Type{
					Name:    "Random",
					Version: 1,
				},
				Data:    generateRandomBytes(size),
				Created: timestamppb.Now(),
			}
		case err = <-errc:
			return cli.Exit(err, 1)
		case <-quit:
			log.Info().Msg("closing the stream")
			close(send)
			break primary
		}
	}

	// Close the publish stream gracefully
	if err = stream.CloseSend(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func consume(c *cli.Context) (err error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var stream api.Ensign_SubscribeClient
	if stream, err = client.Subscribe(context.Background()); err != nil {
		return cli.Exit(err, 1)
	}

	// Create a recv channel to manage the incoming stream
	errc := make(chan error, 1)
	recv := make(chan *api.Event, 2)
	go func(recv chan<- *api.Event, errc chan<- error) {
		for {
			event, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errc <- err
				}
				return
			}
			recv <- event
		}
	}(recv, errc)

primary:
	for {
		select {
		case event := <-recv:
			log.Info().Str("id", event.Id).Int("size", len(event.Data)).Msg("event received")
			stream.Send(&api.Subscription{Embed: &api.Subscription_Ack{Ack: &api.Ack{Id: event.Id}}})
		case err = <-errc:
			return cli.Exit(err, 1)
		case <-quit:
			break primary
		}
	}

	return nil
}

func generateRandomBytes(n int) (b []byte) {
	b = make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
