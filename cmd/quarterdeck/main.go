package main

import (
	"os"
	"log"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/urfave/cli/v2"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "quarterdeck"
	app.Version = pkg.Version()
	app.Usage = "run and manage a quarterdeck server"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name: "serve",
			Usage: "run the quarterdeck server",
			Category: "server",
			Action: serve,
			Flags: []cli.Flag{},
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
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	var srv *quarterdeck.Server
	if srv, err = quarterdeck.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}