package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
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
			Name:     "migrate:projects",
			Usage:    "migrate project owners in the database",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   migrateProjects,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "show the effect of migrating without execution",
				},
			},
		},
		{
			Name:     "purge:users",
			Usage:    "purge a CSV of users and their artifacts from the database",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   purgeUsers,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "show which records would be deleted without execution",
				},
				&cli.StringFlag{
					Name:     "csv",
					Aliases:  []string{"c"},
					Usage:    "path to a CSV file containing users to purge",
					Required: true,
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

func migrateProjects(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dry := c.Bool("dry-run")

	// Find the earliest owner in each organization
	orgOwners := make(map[ulid.ULID]*db.Member)
	onMembers := func(item *trtl.KVPair) error {
		// Skip members that are not owners
		member := &db.Member{}
		if err = member.UnmarshalValue(item.Value); err != nil {
			return cli.Exit(err, 1)
		}

		if member.Role != permissions.RoleOwner {
			return nil
		}

		// Add the user if they are an earlier owner
		if owner, ok := orgOwners[member.OrgID]; !ok {
			orgOwners[member.OrgID] = member
		} else if member.Created.Before(owner.Created) {
			orgOwners[member.OrgID] = member
		}

		return nil
	}

	if _, err = db.List(ctx, nil, nil, db.MembersNamespace, onMembers, nil); err != nil {
		return cli.Exit(err, 1)
	}

	// Iterate over all projects and find those with missing owners
	incompleteProjects := make(map[ulid.ULID]*db.Project)
	onProject := func(item *trtl.KVPair) error {
		// Skip projects that already have owners
		project := &db.Project{}
		if err = project.UnmarshalValue(item.Value); err != nil {
			return cli.Exit(err, 1)
		}

		if ulids.IsZero(project.OwnerID) {
			incompleteProjects[project.ID] = project
		}

		return nil
	}

	if _, err = db.List(ctx, nil, nil, db.ProjectNamespace, onProject, nil); err != nil {
		return cli.Exit(err, 1)
	}

	var migrated int
	var errs *multierror.Error
	for id, project := range incompleteProjects {
		// Get the owner of the organization
		var owner *db.Member
		var ok bool
		if owner, ok = orgOwners[project.OrgID]; !ok {
			errs = multierror.Append(errs, fmt.Errorf("no owner for organization %s", project.OrgID.String()))
			continue
		}

		if dry {
			fmt.Printf("Project %s is missing an owner, would fill with %s\n", id.String(), owner.ID.String())
			continue
		}

		// Do the project update
		project.OwnerID = owner.ID
		if err = db.UpdateProject(ctx, project); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		migrated++
	}

	if errs != nil {
		fmt.Println(errs.Error())
	}

	if dry {
		fmt.Printf("Would migrate %d projects\n", len(incompleteProjects))
	} else {
		fmt.Printf("Migrated %d/%d projects\n", migrated, len(incompleteProjects))
	}

	return nil
}

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

// Purge Members and Tenants from the database, this assumes a CSV file with the format:
// ID,OrgID,Name,Email
// It also assumes that the users have not been verified or invited to organizations.
func purgeUsers(c *cli.Context) (err error) {
	dry := c.Bool("dry-run")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Fetch all users from the CSV
	type user struct {
		ID      ulid.ULID
		OrgID   ulid.ULID
		Name    string
		Email   string
		Member  *db.Member
		Tenants []*db.Tenant
	}
	users := make([]*user, 0)
	userIDs := make(map[ulid.ULID]struct{})

	var f *os.File
	if f, err = os.Open(c.String("csv")); err != nil {
		return cli.Exit(err, 1)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// Skip the header row
	if _, err = reader.Read(); err != nil {
		return cli.Exit(err, 1)
	}

	var records [][]string
	if records, err = reader.ReadAll(); err != nil {
		return cli.Exit(err, 1)
	}

	for _, record := range records {
		user := &user{
			Name:  record[2],
			Email: record[3],
		}

		if user.ID, err = ulids.Parse(record[0]); err != nil {
			return cli.Exit(err, 1)
		}

		if user.OrgID, err = ulids.Parse(record[1]); err != nil {
			return cli.Exit(err, 1)
		}

		if _, ok := userIDs[user.ID]; ok {
			return cli.Exit(fmt.Errorf("duplicate user ID: %s", user.ID.String()), 1)
		}
		userIDs[user.ID] = struct{}{}

		// Lookup the member record in the database
		if user.Member, err = db.RetrieveMember(ctx, user.OrgID, user.ID); err != nil {
			if !errors.Is(err, db.ErrNotFound) {
				return cli.Exit(err, 1)
			}
		}

		// Lookup the tenants associated with the organization
		var cursor *pagination.Cursor
		if user.Tenants, cursor, err = db.ListTenants(ctx, user.OrgID, pagination.New("", "", 100)); err != nil {
			return cli.Exit(err, 1)
		}

		if cursor != nil {
			return cli.Exit(fmt.Errorf("pagination needed to list all the tenants for organization %s", user.OrgID.String()), 1)
		}

		users = append(users, user)
	}

	var (
		deletedMembers int
		deletedTenants int
		errs           *multierror.Error
	)
	for _, user := range users {
		if dry {
			if user.Member != nil {
				deletedMembers++
				fmt.Println("Would delete member", user.Member.ID.String(), user.Member.Name, user.Member.Email)
			}

			if user.Tenants != nil {
				for _, tenant := range user.Tenants {
					deletedTenants++
					fmt.Println("Would delete tenant", tenant.ID.String())
				}
			}
			fmt.Println()
		} else {
			if user.Member != nil {
				if err = db.DeleteMember(ctx, user.Member.OrgID, user.Member.ID); err != nil {
					errs = multierror.Append(errs, err)
				} else {
					deletedMembers++
				}
			}

			if user.Tenants != nil {
				for _, tenant := range user.Tenants {
					if err = db.DeleteTenant(ctx, tenant.OrgID, tenant.ID); err != nil {
						errs = multierror.Append(errs, err)
					} else {
						deletedTenants++
					}
				}
			}
		}
	}

	if dry {
		fmt.Println("Would delete", deletedMembers, "members and", deletedTenants, "tenants")
	} else {
		fmt.Println("Deleted", deletedMembers, "members and", deletedTenants, "tenants")
	}

	if errs != nil {
		fmt.Println(errs.Error())
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
