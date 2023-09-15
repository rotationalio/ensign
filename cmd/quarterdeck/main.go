package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/csv"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	confire "github.com/rotationalio/confire/usage"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/quarterdeck/report"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/rows"
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
			Name:   "serve",
			Usage:  "run the quarterdeck server",
			Action: serve,
			Flags:  []cli.Flag{},
		},
		{
			Name:     "config",
			Usage:    "print quarterdeck configuration guide",
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
			Name:     "listorgs",
			Usage:    "list all organizations in the database",
			Category: "utility",
			Action:   listOrgs,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "csv",
					Aliases: []string{"c"},
					Usage:   "write the list as a CSV file to the specified path",
				},
			},
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
		{
			Name:     "tokenkey",
			Usage:    "generate an RSA token key pair and ksuid for JWT token signing",
			Category: "utility",
			Action:   generateTokenKey,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write keys out to (optional, will be saved as ulid.pem by default)",
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "number of bits for the generated keys",
					Value:   4096,
				},
			},
		},
		{
			Name:      "argon2",
			Usage:     "create a derived key to use as a fixture for testing",
			Category:  "debug",
			Action:    derkey,
			ArgsUsage: "password [password ...]",
			Flags:     []cli.Flag{},
		},
		{
			Name:     "keypair",
			Usage:    "create a fake apikey client ID and secret to use as a fixture for testing",
			Category: "debug",
			Action:   keypair,
			Flags:    []cli.Flag{},
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

//===========================================================================
// Utility Commands
//===========================================================================

func listOrgs(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// If csv path is specified, open the file for writing, otherwise use stdout
	var w rows.Writer
	if path := c.String("csv"); path != "" {
		var f *os.File
		if f, err = os.Create(path); err != nil {
			return cli.Exit(err, 1)
		}
		defer f.Close()

		w = csv.NewWriter(f)
		defer w.(*csv.Writer).Flush()
	} else {
		w = rows.NewTabRowWriter(tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0))
		defer w.(*rows.TabRowWriter).Flush()
	}

	// Write the header
	w.Write([]string{"ID", "Name", "Domain", "Projects", "Created", "Modified"})

	// Paginate over all organizations in the database, writing rows.
	cursor := pagination.New("", "", 100)
	for cursor != nil {
		var orgs []*models.Organization
		if orgs, cursor, err = models.ListAllOrgs(ctx, cursor); err != nil {
			return cli.Exit(err, 1)
		}

		for _, org := range orgs {
			w.Write([]string{org.ID.String(), org.Name, org.Domain, fmt.Sprintf("%d", org.ProjectCount()), org.Created, org.Modified})
		}
	}

	return nil
}

func usage(c *cli.Context) (err error) {
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	format := confire.DefaultTableFormat
	if c.Bool("list") {
		format = confire.DefaultListFormat
	}

	var conf config.Config
	if err := confire.Usagef("quarterdeck", &conf, tabs, format); err != nil {
		return cli.Exit(err, 1)
	}
	tabs.Flush()
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

func generateTokenKey(c *cli.Context) (err error) {
	// Create ULID and determine outpath
	keyid := ulid.Make()

	var out string
	if out = c.String("out"); out == "" {
		out = fmt.Sprintf("%s.pem", keyid)
	}

	// Generate RSA keys using crypto random
	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, c.Int("size")); err != nil {
		return cli.Exit(err, 1)
	}

	// Open file to PEM encode keys to
	var f *os.File
	if f, err = os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return cli.Exit(err, 1)
	}

	if err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("RSA key id: %s -- saved with PEM encoding to %s\n", keyid, out)
	return nil
}

//===========================================================================
// Debug Commands
//===========================================================================

func derkey(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit("specify password(s) to create argon2 derived key(s) from", 1)
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
