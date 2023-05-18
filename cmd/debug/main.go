package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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
			Name:   "binulid",
			Usage:  "create a binary ULID to insert into SQLite",
			Action: binulid,
		},
		{
			Name:   "derkey",
			Usage:  "create a derived key to insert into SQLite",
			Action: derkey,
		},
		{
			Name:   "keypair",
			Usage:  "create an api key client id and secret to insert into SQLite",
			Action: keypair,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("could not execute cli app")
	}
}

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

func binulid(c *cli.Context) error {
	ts := time.Now().UTC()
	id := ulids.FromTime(ts)
	data, err := id.MarshalBinary()
	if err != nil {
		return cli.Exit(err, 1)
	}

	out := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', tabwriter.AlignRight|tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(out, "ULID\t%s\t\n", id.String())
	fmt.Fprintf(out, "Time\t%s\t\n", ts.Format(time.RFC3339Nano))
	fmt.Fprintf(out, "Hex Bytes\t%s\t\n", hex.EncodeToString(data))
	fmt.Fprintf(out, "b64 Bytes\t%s\t\n", base64.RawStdEncoding.EncodeToString(data))
	out.Flush()
	return nil
}

func derkey(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit("specify password(s) to create derived key(s) from", 1)
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

func keypair(c *cli.Context) error {
	clientID := keygen.KeyID()
	secret := keygen.Secret()
	fmt.Printf("%s.%s\n", clientID, secret)
	return nil
}
