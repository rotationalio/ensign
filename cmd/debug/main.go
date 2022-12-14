package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	ensign "github.com/rotationalio/ensign/sdks/go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "ensign-debug"
	app.Version = pkg.Version()
	app.Before = setupLogger
	app.Usage = "client utilities to help debug an ensign server"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e"},
			Usage:   "endpoint of local ensign node to send requests to",
			Value:   "127.0.0.1:5356",
			EnvVars: []string{"ENSIGN_ENDPOINT"},
		},
		&cli.BoolFlag{
			Name:    "no-secure",
			Aliases: []string{"S"},
			Usage:   "do not connect with TLS credentials",
			EnvVars: []string{"ENSIGN_INSECURE"},
		},
		&cli.StringFlag{
			Name:    "verbosity",
			Aliases: []string{"L"},
			Usage:   "set the zerolog level",
			Value:   "info",
			EnvVars: []string{"ENSIGN_LOG_LEVEL"},
		},
		&cli.BoolFlag{
			Name:    "console",
			Aliases: []string{"C"},
			Usage:   "human readable console log instead of json",
			Value:   false,
			EnvVars: []string{"ENSIGN_CONSOLE_LOG"},
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
				&cli.Float64Flag{
					Name:    "rate",
					Aliases: []string{"r"},
					Usage:   "events to publish per second (-1 for as fast as possible)",
					Value:   30,
					EnvVars: []string{"ENSIGN_DEBUG_GENERATE_RATE"},
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "the size in bytes of the event data generated",
					Value:   128,
					EnvVars: []string{"ENSIGN_DEBUG_EVENT_SIZE"},
				},
			},
		},
		{
			Name:   "consume",
			Usage:  "subscribe to the stream and consume events",
			Before: connect,
			After:  disconnect,
			Action: consume,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "binulid",
			Usage:  "create a binary ULID to insert into SQLite",
			Action: binulid,
		},
		{
			Name:   "derkey",
			Usage:  "create a derived key to insert into SQLite",
			Action: derkey,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("could not execute cli app")
	}
}

var client *ensign.Client

func setupLogger(c *cli.Context) (err error) {
	switch strings.ToLower(c.String("verbosity")) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		return cli.Exit(fmt.Errorf("unknown log level %q", c.String("verbosity")), 1)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg
	zerolog.DurationFieldInteger = false
	zerolog.DurationFieldUnit = time.Millisecond

	if c.Bool("console") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		// Add the severity hook for GCP logging
		var gcpHook logger.SeverityHook
		log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
	}

	return nil
}

func connect(c *cli.Context) (err error) {
	opts := &ensign.Options{
		Endpoint: c.String("endpoint"),
		Insecure: c.Bool("no-secure"),
	}

	if client, err = ensign.New(opts); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func disconnect(c *cli.Context) (err error) {
	if err = client.Close(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func generate(c *cli.Context) (err error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var publisher ensign.Publisher
	if publisher, err = client.Publish(context.Background()); err != nil {
		return cli.Exit(err, 1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{}, 1)

	size := c.Int("size")
	if hz := c.Float64("rate"); hz > 0 {
		interval := time.Duration(float64(time.Second) / hz)
		log.Info().Float64("hz", hz).Dur("interval", interval).Int("size", size).Msg("starting rate limited publisher")
		ticker := time.NewTicker(interval)
		go func(done <-chan struct{}) {
			defer wg.Done()
			var msg uint64
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
				}

				msg++
				publisher.Publish(&api.Event{
					TopicId:  "generator",
					Mimetype: mimetype.ApplicationOctetStream,
					Type: &api.Type{
						Name:    "Random",
						Version: 1,
					},
					Data:    generateRandomBytes(size),
					Created: timestamppb.Now(),
				})

				if err = publisher.Err(); err != nil {
					log.Error().Err(err).Msg("could not publish event")
					return
				}

				log.Debug().Uint64("num", msg).Msg("event published")
			}
		}(done)
	} else {
		log.Info().Int("size", size).Msg("starting max rate publisher")
		go func(done <-chan struct{}) {
			defer wg.Done()
			var msg uint64
			for {
				select {
				case <-done:
					return
				default:
				}

				msg++
				publisher.Publish(&api.Event{
					TopicId:  "generator",
					Mimetype: mimetype.ApplicationOctetStream,
					Type: &api.Type{
						Name:    "Random",
						Version: 1,
					},
					Data:    generateRandomBytes(size),
					Created: timestamppb.Now(),
				})

				if err = publisher.Err(); err != nil {
					log.Error().Err(err).Msg("could not publish event")
					return
				}

				log.Debug().Uint64("num", msg).Msg("event published")
			}
		}(done)
	}

	<-quit
	log.Info().Msg("stopping")
	done <- struct{}{}
	wg.Wait()

	// Close the publish stream gracefully
	log.Info().Msg("closing the stream")
	if err = publisher.Close(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func consume(c *cli.Context) (err error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var subscriber ensign.Subscriber
	if subscriber, err = client.Subscribe(context.Background()); err != nil {
		return cli.Exit(err, 1)
	}

	var sub <-chan *api.Event
	if sub, err = subscriber.Subscribe(); err != nil {
		return cli.Exit(err, 1)
	}

	count := uint64(0)
primary:
	for {
		select {
		case event := <-sub:
			count++
			log.Debug().Str("id", event.Id).Int("size", len(event.Data)).Msg("event received")
			subscriber.Ack(event.Id)

			if count%1e3 == 0 {
				log.Info().Uint64("events", count).Msg("events received")
			}
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

func binulid(c *cli.Context) error {
	id := ulid.Make()
	data, err := id.MarshalBinary()
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(time.Now().UTC().Format(time.RFC3339Nano))
	fmt.Println(id.String())
	fmt.Println(hex.EncodeToString(data))
	return nil
}

func derkey(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit("specify password(s) to create derived key for", 1)
	}

	for i := 0; i < c.NArg(); i++ {
		pwdk, err := passwd.CreateDerivedKey(c.Args().Get(i))
		if err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Println(pwdk)
	}

	return nil
}
