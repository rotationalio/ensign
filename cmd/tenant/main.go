package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
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
			Name:     "db:update-members",
			Usage:    "update team members in the tenant database from a CSV file",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   updateMembers,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "file",
					Aliases:  []string{"f"},
					Usage:    "path to CSV file with member ID to email mappings",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "show the effect of updating members without execution",
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
		if next, err = db.List(ctx, nil, nil, c.String("namespace"), onListItem, next); err != nil {
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

	if _, err = db.List(ctx, nil, nil, db.OrganizationNamespace, fetchKeys, nil); err != nil {
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

	if _, err = db.List(ctx, nil, nil, db.TenantNamespace, fetchTenants, nil); err != nil {
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

	if _, err = db.List(ctx, nil, nil, db.ProjectNamespace, fetchProjects, nil); err != nil {
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

func updateMembers(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dry := c.Bool("dry-run")
	userEmails := make(map[ulid.ULID]string)

	// Load user emails from the CSV
	var f *os.File
	if f, err = os.Open(c.String("file")); err != nil {
		return cli.Exit(err, 1)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = '|'

	for {
		var record []string
		if record, err = r.Read(); err != nil {
			if err == io.EOF {
				break
			}

			return cli.Exit(err, 1)
		}

		if len(record) != 2 {
			return cli.Exit("invalid record", 1)
		}

		var data []byte
		if data, err = hex.DecodeString(record[0]); err != nil {
			return cli.Exit(err, 1)
		}

		var id ulid.ULID
		if id, err = ulids.Parse(data); err != nil {
			return cli.Exit(err, 1)
		}

		userEmails[id] = record[1]
	}

	// Fetch all members with missing data
	missingEmails := make(map[ulid.ULID]*db.Member)
	fetchMembers := func(item *trtl.KVPair) error {
		member := &db.Member{}
		if err = member.UnmarshalValue(item.Value); err != nil {
			return cli.Exit(err, 1)
		}

		if member.Email == "" {
			missingEmails[member.ID] = member
		}

		return nil
	}

	if _, err = db.List(ctx, nil, nil, db.MembersNamespace, fetchMembers, nil); err != nil {
		return cli.Exit(err, 1)
	}

	var errs *multierror.Error
	var updated int
	for id, member := range missingEmails {
		old := *member

		// Update the email from the CSV
		var ok bool
		if member.Email, ok = userEmails[id]; !ok {
			errs = multierror.Append(errs, fmt.Errorf("could not find email for member %s", id.String()))
			continue
		}

		// Populate missing fields for organization owners
		if member.Status == db.MemberStatusPending && member.Role == perms.RoleOwner {
			// Make sure that this is the only owner, otherwise it's ambiguous which one was the original
			var owners int
			countOwners := func(item *trtl.KVPair) error {
				member := &db.Member{}
				if err = member.UnmarshalValue(item.Value); err != nil {
					return cli.Exit(err, 1)
				}

				if member.Role == perms.RoleOwner {
					owners++
				}

				return nil
			}

			if _, err = db.List(ctx, member.OrgID[:], nil, db.MembersNamespace, countOwners, nil); err != nil {
				return cli.Exit(err, 1)
			}

			if owners != 1 {
				errs = multierror.Append(errs, fmt.Errorf("member %s is not the only owner in org %s, can't proceed with status change", id.String(), member.OrgID.String()))
				continue
			}

			member.Status = db.MemberStatusConfirmed

			// Update missing timestamps for organization owners
			if member.DateAdded.IsZero() {
				member.DateAdded = member.Created
			}

			if member.LastActivity.IsZero() {
				member.LastActivity = member.Created
			}
		}

		fmt.Printf("Member %s will be updated\n", id.String())
		fmt.Printf("  Old: %+v\n", old)
		fmt.Printf("  New: %+v\n", *member)
		fmt.Println()

		if dry {
			continue
		}

		if err = db.UpdateMember(ctx, member); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		updated++
	}

	if errs != nil {
		fmt.Println(errs.Error())
	}

	if dry {
		fmt.Printf("%d members would be updated\n", len(missingEmails))
	} else {
		fmt.Println("Updated", updated, "members")
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

	// Connect tot he trtl server
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

/*
func connectQD(c *cli.Context) (err error) {

}*/
