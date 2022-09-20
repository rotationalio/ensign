package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
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
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "bindaddr",
					Aliases: []string{"a"},
					Usage:   "address to bind the ensign server to",
				},
				&cli.StringFlag{
					Name:    "conf-file",
					Aliases: []string{"c"},
					Usage:   "path to a configuration file on disk",
				},
			},
		},
		{
			Name:     "config",
			Usage:    "inspect and validate the ensign configuration",
			Category: "ops",
			Action:   inspectConfig,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "conf-file",
					Aliases: []string{"c"},
					Usage:   "path to a configuration file on disk",
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
					Value:   "localhost:7777",
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
	if path := c.String("conf-file"); path != "" {
		if conf, err = config.Load(path); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		if conf, err = config.New(); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Override configuration based on CLI flags
	if addr := c.String("bindaddr"); addr != "" {
		conf.BindAddr = addr
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

func inspectConfig(c *cli.Context) (err error) {
	// Load the configuration from a file or from the environment.
	var conf config.Config
	if path := c.String("conf-file"); path != "" {
		if conf, err = config.Load(path); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		if conf, err = config.New(); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// If the configuration is valid print it as YAML
	if err = yaml.NewEncoder(os.Stdout).Encode(&conf); err != nil {
		return err
	}
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cc *grpc.ClientConn
	if cc, err = grpc.DialContext(ctx, endpoint, opts...); err != nil {
		return cli.Exit(err, 1)
	}

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
