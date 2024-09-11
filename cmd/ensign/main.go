package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/joho/godotenv"
	confire "github.com/rotationalio/confire/usage"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/ensign"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "ensign"
	app.Version = pkg.Version()
	app.Usage = "run and manage an ensign server node"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "run the ensign server",
			Category: "server",
			Action:   serve,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "config",
			Usage:    "print ensign configuration guide",
			Category: "utility",
			Action:   usage,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "print in list mode instead of table mode",
				},
			},
		},
		{
			Name:     "status",
			Usage:    "check the status of a running ensign server",
			Category: "ops",
			Action:   status,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Usage:   "endpoint to make status request to",
					Value:   "localhost:5356",
				},
				&cli.BoolFlag{
					Name:    "no-secure",
					Aliases: []string{"S"},
					Usage:   "do not connect with TLS credentials",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// Server Commands
//===========================================================================

func serve(c *cli.Context) (err error) {
	// Load the configuration from a file or from the environment.
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	// Create and serve the Ensign server
	var srv *ensign.Server
	if srv, err = ensign.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

//===========================================================================
// Ops Commands
//===========================================================================

func usage(c *cli.Context) (err error) {
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	format := confire.DefaultTableFormat
	if c.Bool("list") {
		format = confire.DefaultListFormat
	}

	var conf config.Config
	if err := confire.Usagef("ensign", &conf, tabs, format); err != nil {
		return cli.Exit(err, 1)
	}
	tabs.Flush()
	return nil
}

func status(c *cli.Context) (err error) {
	opts := make([]grpc.DialOption, 0, 1)
	endpoint := c.String("endpoint")

	if c.Bool("no-secure") {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.NewClient(endpoint, opts...); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rep *api.ServiceState
	client := api.NewEnsignClient(cc)
	if rep, err = client.Status(ctx, &api.HealthCheck{}); err != nil {
		return cli.Exit(err, 1)
	}

	pbjson := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = pbjson.Marshal(rep); err != nil {
		return cli.Exit(err, 1)
	}
	fmt.Println(string(data))
	return nil
}
