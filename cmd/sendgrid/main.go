package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/urfave/cli/v2"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Creates a multi-command CLI application
	app := cli.NewApp()
	app.Name = "sendgrid"
	app.Version = pkg.Version()
	app.Usage = "interact with the sendgrid API"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:     "lists",
			Usage:    "fetch marketing lists from sendgrid",
			Category: "api",
			Action:   fetchLists,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "apikey",
					Aliases: []string{"k"},
					Usage:   "sendgrid API key",
				},
			},
		},
		{
			Name:     "defs",
			Usage:    "fetch field definitions from sendgrid",
			Category: "api",
			Action:   fetchDefs,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "apikey",
					Aliases: []string{"k"},
					Usage:   "sendgrid API key",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// SendGrid API Commands
//===========================================================================

func fetchLists(c *cli.Context) (err error) {
	key := apiKey(c)
	if key == "" {
		return cli.Exit("specify sendgrid api key in command line or environment", 1)
	}

	// Fetch a page of sendgrid marketing lists
	// Note: Default page size is 100 so we should get everything in one request
	var rep string
	if rep, err = sendgrid.MarketingLists(key, ""); err != nil {
		return cli.Exit(err, 1)
	}

	// Print the response
	printJSON(rep)
	return nil
}

func fetchDefs(c *cli.Context) (err error) {
	key := apiKey(c)
	if key == "" {
		return cli.Exit("specify sendgrid api key in command line or environment", 1)
	}

	// Fetch the field definitions
	var rep string
	if rep, err = sendgrid.FieldDefinitions(key); err != nil {
		return cli.Exit(err, 1)
	}

	// Print the response
	printJSON(rep)
	return nil
}

//===========================================================================
// Helpers
//===========================================================================

// Get the API key from the command line or environment
func apiKey(c *cli.Context) string {
	key := c.String("apikey")
	if key == "" {
		key = os.Getenv("SENDGRID_API_KEY")
	}

	return key
}

// Print a JSON string response to stdout and add indents for readability
func printJSON(rep string) {
	out := &bytes.Buffer{}
	if err := json.Indent(out, []byte(rep), "", "  "); err != nil {
		cli.Exit(err, 1)
	}

	fmt.Println(out.String())
}
