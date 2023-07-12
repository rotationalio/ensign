package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/report"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/urfave/cli/v2"
)

var conf config.Config

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
			Name:     "serve",
			Usage:    "run the quarterdeck server",
			Category: "server",
			Action:   serve,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "report",
			Usage:    "runs the daily PLG report",
			Category: "utility",
			Action:   runDailyUsersReport,
			Before:   connectDB,
			After:    closeDB,
			Flags:    []cli.Flag{},
		},
		{
			Name:      "revoke",
			Usage:     "manually revoke an API key the hard way",
			Category:  "utility",
			Action:    revoke,
			Before:    connectDB,
			After:     closeDB,
			ArgsUsage: "clientID [clientID ...]",
			Flags:     []cli.Flag{},
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

func runDailyUsersReport(c *cli.Context) (err error) {
	emailer := &Emailer{conf: conf}
	if emailer.sendgrid, err = emails.New(conf.SendGrid); err != nil {
		return cli.Exit(err, 1)
	}

	var daily *report.DailyUsers
	if daily, err = report.NewDailyUsers(emailer); err != nil {
		return cli.Exit(err, 1)
	}

	if err = daily.Report(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func revoke(c *cli.Context) (err error) {
	if c.NArg() < 1 {
		return cli.Exit("specify at least one clientID of an API key to revoke", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	args := c.Args()
	for i := 0; i < c.NArg(); i++ {
		// Lookup the clientID in the database.
		var apikey *models.APIKey
		if apikey, err = models.GetAPIKey(ctx, args.Get(i)); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				fmt.Printf("could not find key with client ID %q\n", args.Get(i))
			} else {
				return cli.Exit(fmt.Errorf("could not retrieve key with client ID %q: %w", args.Get(i), err), 1)
			}
		}

		// Delete (revoke) the API key from the database.
		if err = models.DeleteAPIKey(ctx, apikey.ID, apikey.OrgID); err != nil {
			return cli.Exit(fmt.Errorf("could not revoke key with client ID %q: %w", apikey.KeyID, err), 1)
		}
		fmt.Printf("revoked key %q (%s) with clientID %s\n", apikey.Name, apikey.ID, apikey.KeyID)
	}

	return nil
}

//===========================================================================
// Before and After Actions
//===========================================================================

func connectDB(c *cli.Context) (err error) {
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	// Ensure we're connected to the database
	if err = db.Connect(conf.Database.URL, true); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func closeDB(c *cli.Context) (err error) {
	if err = db.Close(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

//===========================================================================
// Emailer Utility
//===========================================================================

// Emailer is a helper so that CLI programs can send emails and make reports.
type Emailer struct {
	conf     config.Config
	sendgrid *emails.EmailManager
}

// Send the daily users report to the Rotational admins.
// This method overwrites the email data on the report with the configured sender and
// recipient so it should not be specified by the user (e.g. the user should only supply
// the report data for the email template).
func (s *Emailer) SendDailyUsers(data *emails.DailyUsersData) (err error) {
	data.EmailData = emails.EmailData{
		Sender:    s.conf.SendGrid.MustFromContact(),
		Recipient: s.conf.SendGrid.MustAdminContact(),
	}

	data.Domain = s.conf.Reporting.Domain
	data.EnsignDashboardLink = s.conf.Reporting.DashboardURL

	var msg *mail.SGMailV3
	if msg, err = emails.DailyUsersEmail(*data); err != nil {
		return err
	}

	// Attach the report as json
	var attachment []byte
	if attachment, err = json.MarshalIndent(data, "", " "); err != nil {
		return err
	}

	if err = emails.AttachJSON(msg, attachment, fmt.Sprintf("daily_users_%s.json", data.Date.Format("20060102"))); err != nil {
		return err
	}

	// Attach the new accounts CSV if there are any available
	if len(data.NewAccounts) > 0 {
		var accounts []byte
		if accounts, err = data.NewAccountsCSV(); err != nil {
			return err
		}

		if err = emails.AttachCSV(msg, accounts, fmt.Sprintf("new_accounts_%s.json", data.Date.Format("20060102"))); err != nil {
			return err
		}
	}

	// Send the email
	return s.sendgrid.Send(msg)
}
