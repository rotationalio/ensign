package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	confire "github.com/rotationalio/confire/usage"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/urfave/cli/v2"
)

var (
	conf    config.Config
	timeout time.Duration = 15 * time.Second
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Creates a multi-command CLI application
	app := cli.NewApp()
	app.Name = "tenant"
	app.Version = pkg.Version()
	app.Usage = "run and manage a tenant server"
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "run the tenant server",
			Category: "server",
			Before:   configure,
			Action:   serve,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "config",
			Usage:    "print tenant configuration guide",
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
			Name:     "db:list",
			Usage:    "list all keys in a tenant namespace",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   listKeys,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "namespace",
					Aliases:  []string{"n"},
					Usage:    "namespace to list keys for",
					Required: true,
				},
			},
		},
		{
			Name:     "db:reindex",
			Usage:    "update all tenant-based indexes",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   reindex,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "show the effect of re-indexing without execution",
				},
			},
		},
		{
			Name:     "db:cleanup",
			Usage:    "remove all tenants and members that do not appear in the organization list",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   cleanup,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "orgs",
					Aliases:  []string{"f"},
					Usage:    "path to a CSV file containing a list of organization IDs to keep",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "show the effect of cleanup without execution",
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
	var srv *tenant.Server
	if srv, err = tenant.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
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
	if err := confire.Usagef("tenant", &conf, tabs, format); err != nil {
		return cli.Exit(err, 1)
	}
	tabs.Flush()
	return nil
}

//===========================================================================
// Client Commands
//===========================================================================

func listKeys(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Print keys in the namespace
	onListItem := func(item *trtl.KVPair) error {
		switch len(item.Key) {
		case 16:
			id, err := ulids.Parse(item.Key)
			if err != nil {
				return cli.Exit(err, 1)
			}
			fmt.Println(id.String())
		case 32:
			key := &db.Key{}
			if err = key.UnmarshalValue(item.Key); err != nil {
				return cli.Exit(err, 1)
			}

			var parent ulid.ULID
			if parent, err = key.ParentID(); err != nil {
				return cli.Exit(err, 1)
			}

			var object ulid.ULID
			if object, err = key.ObjectID(); err != nil {
				return cli.Exit(err, 1)
			}

			fmt.Println(parent.String(), object.String())
		default:
			return cli.Exit(fmt.Errorf("unexpected key length: %d", len(item.Key)), 1)
		}

		return nil
	}

	// Get all the keys in the namespace
	var next *pagination.Cursor
	for {
		if next, err = db.List(ctx, nil, c.String("namespace"), onListItem, next); err != nil {
			return cli.Exit(err, 1)
		}

		if next == nil {
			break
		}
	}

	return nil
}

func reindex(c *cli.Context) (err error) {
	dry := c.Bool("dry-run")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Fetch all current organization keys
	orgKeys := make(map[ulid.ULID]struct{})
	fetchKeys := func(item *trtl.KVPair) error {
		id, err := ulids.Parse(item.Key)
		if err != nil {
			return cli.Exit(err, 1)
		}
		orgKeys[id] = struct{}{}
		return nil
	}

	if _, err = db.List(ctx, nil, db.OrganizationNamespace, fetchKeys, nil); err != nil {
		return cli.Exit(err, 1)
	}

	// Fetch all tenants that do not have an organization key
	missingTenants := make(map[*db.Key]struct{})
	fetchTenants := func(item *trtl.KVPair) error {
		key := &db.Key{}
		if err = key.UnmarshalValue(item.Key); err != nil {
			return cli.Exit(err, 1)
		}

		var object ulid.ULID
		if object, err = key.ObjectID(); err != nil {
			return cli.Exit(err, 1)
		}

		if _, ok := orgKeys[object]; !ok {
			missingTenants[key] = struct{}{}
		}

		return nil
	}

	if _, err = db.List(ctx, nil, db.TenantNamespace, fetchTenants, nil); err != nil {
		return cli.Exit(err, 1)
	}

	// Fetch all projects that do not have an organization key
	missingProjects := make(map[db.Key]struct{})
	fetchProjects := func(item *trtl.KVPair) error {
		key := &db.Key{}
		if err = key.UnmarshalValue(item.Key); err != nil {
			return cli.Exit(err, 1)
		}

		var object ulid.ULID
		if object, err = key.ObjectID(); err != nil {
			return cli.Exit(err, 1)
		}

		if _, ok := orgKeys[object]; !ok {
			// Projects are stored by tenantID:projectID, however we need the orgID for the index
			project := &db.Project{}
			if err = project.UnmarshalValue(item.Value); err != nil {
				return cli.Exit(err, 1)
			}

			var orgKey db.Key
			if orgKey, err = db.CreateKey(project.OrgID, project.ID); err != nil {
				return cli.Exit(err, 1)
			}

			missingProjects[orgKey] = struct{}{}
		}

		return nil
	}

	if _, err = db.List(ctx, nil, db.ProjectNamespace, fetchProjects, nil); err != nil {
		return cli.Exit(err, 1)
	}

	// Create organization keys for all missing tenants and projects
	var migratedTenants, migratedProjects int
	var errs *multierror.Error
	for tenant := range missingTenants {
		var id ulid.ULID
		if id, err = tenant.ObjectID(); err != nil {
			return cli.Exit(err, 1)
		}

		var orgID ulid.ULID
		if orgID, err = tenant.ParentID(); err != nil {
			return cli.Exit(err, 1)
		}

		if dry {
			fmt.Println("Missing org key for tenant", id.String())
			continue
		}

		if err = db.PutOrgIndex(ctx, id, orgID); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		migratedTenants++
	}

	for project := range missingProjects {
		var id ulid.ULID
		if id, err = project.ObjectID(); err != nil {
			return cli.Exit(err, 1)
		}

		var orgID ulid.ULID
		if orgID, err = project.ParentID(); err != nil {
			return cli.Exit(err, 1)
		}

		if dry {
			fmt.Println("Missing org key for project", id.String())
			continue
		}

		if err = db.PutOrgIndex(ctx, id, orgID); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		migratedProjects++
	}

	if errs != nil {
		fmt.Println(errs.Error())
	}

	if dry {
		fmt.Printf("%d tenants and %d projects would be migrated\n", len(missingTenants), len(missingProjects))
	} else {
		fmt.Println("Migrated", migratedTenants, "tenants and", migratedProjects, "projects")
	}

	return nil
}

func cleanup(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dry := c.Bool("dry-run")

	// Load the organizations from the CSV file
	var f *os.File
	if f, err = os.Open(c.String("orgs")); err != nil {
		return cli.Exit(err, 1)
	}

	// Ensure there is a header row
	var header []string
	reader := csv.NewReader(f)
	if header, err = reader.Read(); err != nil {
		return cli.Exit(err, 1)
	}

	// Find the ID column or assume the first column is the ID
	var idCol int
	for i, col := range header {
		if strings.ToLower(col) == "id" {
			idCol = i
			break
		}
	}

	orgs := make(map[ulid.ULID]struct{})
	for {
		var record []string
		if record, err = reader.Read(); err != nil {
			break
		}

		var id ulid.ULID
		if id, err = ulid.Parse(record[idCol]); err != nil {
			return cli.Exit(fmt.Errorf("could not parse org ID: %s", record[idCol]), 1)
		}

		orgs[id] = struct{}{}
	}

	if len(orgs) == 0 {
		return cli.Exit(fmt.Errorf("no organizations found in CSV file"), 1)
	}

	// Fetch tenants not in the organization list
	strandedTenants := make(map[ulid.ULID]*db.Tenant)
	for {
		var (
			tenants []*db.Tenant
			next    *pagination.Cursor
		)
		if tenants, next, err = db.ListTenants(ctx, ulid.ULID{}, next); err != nil {
			return cli.Exit(err, 1)
		}

		for _, tenant := range tenants {
			if _, ok := orgs[tenant.OrgID]; !ok {
				strandedTenants[tenant.ID] = tenant
			}
		}

		if next == nil {
			break
		}
	}

	// Fetch members not in the organization list
	strandedMembers := make(map[ulid.ULID]*db.Member)
	for {
		var (
			members []*db.Member
			next    *pagination.Cursor
		)
		if members, next, err = db.ListMembers(ctx, ulid.ULID{}, next); err != nil {
			return cli.Exit(err, 1)
		}

		for _, member := range members {
			if _, ok := orgs[member.OrgID]; !ok {
				strandedMembers[member.ID] = member
			}
		}

		if next == nil {
			break
		}
	}

	if dry {
		fmt.Println("The following tenants would be removed:")
		for _, tenant := range strandedTenants {
			fmt.Println(tenant.ID.String(), tenant.Name, tenant.Modified)
		}

		fmt.Println("The following members would be removed:")
		for _, member := range strandedMembers {
			fmt.Println(member.ID.String(), member.Email, member.Organization, member.Modified)
		}

		return nil
	} else {
		var errs error
		fmt.Println("Removing", len(strandedTenants), "tenants and", len(strandedMembers), "members")
		for _, tenant := range strandedTenants {
			fmt.Println("Removing tenant", tenant.ID.String(), tenant.Name, tenant.Modified)
			if err = db.DeleteTenant(ctx, tenant.OrgID, tenant.ID); err != nil {
				errs = errors.Join(errs, err)
			}
		}

		for _, member := range strandedMembers {
			fmt.Println("Removing member", member.ID.String(), member.Email, member.Organization, member.Modified)
			if err = db.DeleteMember(ctx, member.OrgID, member.ID); err != nil {
				errs = errors.Join(errs, err)
			}
		}

		if errs != nil {
			return cli.Exit(errs, 1)
		}
		fmt.Println("Successfully removed", len(strandedTenants), "tenants and", len(strandedMembers), "members")
	}

	return nil
}

//===========================================================================
// Helpers
//===========================================================================

func configure(c *cli.Context) (err error) {
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func connectDB(c *cli.Context) (err error) {
	// suppress output from the logger
	logger.Discard()

	// configure the environment to connect to the database
	if err = configure(c); err != nil {
		return err
	}
	conf.ConsoleLog = false

	// Connect to the trtl server
	if err = db.Connect(conf.Database); err != nil {
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
