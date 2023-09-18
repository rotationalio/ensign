package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	confire "github.com/rotationalio/confire/usage"
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
	"github.com/vmihailenco/msgpack/v5"
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
			Name:     "migrate:profiles",
			Usage:    "migrate user profiles in the database",
			Category: "client",
			Before:   connectDB,
			After:    closeDB,
			Action:   migrateProfiles,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "orgs",
					Aliases:  []string{"o"},
					Usage:    "path to a CSV file containing the list of known organizations",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "show the effect of migrating without execution",
				},
				&cli.StringFlag{
					Name:    "report",
					Aliases: []string{"r"},
					Usage:   "path to a CSV file to write the migration report to",
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

type oldMember struct {
	OrgID        ulid.ULID       `msgpack:"org_id"`
	ID           ulid.ULID       `msgpack:"id"`
	Email        string          `msgpack:"email"`
	Name         string          `msgpack:"name"`
	Role         string          `msgpack:"role"`
	Status       db.MemberStatus `msgpack:"status"`
	Created      time.Time       `msgpack:"created"`
	Modified     time.Time       `msgpack:"modified"`
	DateAdded    time.Time       `msgpack:"date_added"`
	LastActivity time.Time       `msgpack:"last_activity"`
}

func migrateProfiles(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dry := c.Bool("dry-run")

	// Open the CSV file containing the list of organizations
	var orgsFile *os.File
	if orgsFile, err = os.Open(c.String("orgs")); err != nil {
		return cli.Exit(err, 1)
	}
	defer orgsFile.Close()

	// Open the file to write the report to
	var reportFile *os.File
	report := c.String("report")
	if report != "" {
		if reportFile, err = os.Create(c.String("report")); err != nil {
			return cli.Exit(err, 1)
		}
		defer reportFile.Close()
	}

	// Read the header
	var header []string
	reader := csv.NewReader(orgsFile)
	if header, err = reader.Read(); err != nil {
		return cli.Exit(err, 1)
	}

	// Ensure header fields are correct
	if header[0] != "ID" || header[1] != "Name" || header[2] != "Domain" {
		return cli.Exit(fmt.Errorf("unexpected header fields: %v", header), 1)
	}

	// Read the organizations from the CSV file
	type organization struct {
		Name   string
		Domain string
		owner  ulid.ULID
	}
	orgs := make(map[ulid.ULID]organization, 0)

	for {
		var record []string
		if record, err = reader.Read(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return cli.Exit(err, 1)
		}

		// Parse the organization ID
		var id ulid.ULID
		if id, err = ulids.Parse(record[0]); err != nil {
			return cli.Exit(err, 1)
		}

		orgs[id] = organization{
			Name:   record[1],
			Domain: record[2],
		}
	}

	// Iterate over all members in the database
	var (
		members, migrated int
		errs              *multierror.Error
	)
	onMember := func(item *trtl.KVPair) error {
		members++

		// Parse the old version of the member record
		old := &oldMember{}
		if err = msgpack.Unmarshal(item.Value, old); err != nil {
			errs = multierror.Append(errs, err)
			if reportFile != nil {
				reportFile.WriteString(fmt.Sprintf("error unmarshaling existing member record with key %s: %s\n", item.Key, err.Error()))
			}
			return nil
		}

		fmt.Printf("would migrate member %s:%s\n", old.OrgID.String(), old.ID.String())
		fmt.Printf("existing record: %v\n", old)
		if reportFile != nil {
			reportFile.WriteString(fmt.Sprintf("migrating member %s:%s\n", old.OrgID.String(), old.ID.String()))
			reportFile.WriteString(fmt.Sprintf("existing record: %v\n", old))
		}

		// Create the new version to store in the database
		member := &db.Member{
			OrgID:        old.OrgID,
			ID:           old.ID,
			Email:        old.Email,
			Name:         old.Name,
			Role:         old.Role,
			JoinedAt:     old.DateAdded,
			LastActivity: old.LastActivity,
			Created:      old.Created,
			Modified:     old.Modified,
		}

		// Ensure the organization exists
		var (
			org organization
			ok  bool
		)
		if org, ok = orgs[old.OrgID]; !ok {
			errs = multierror.Append(errs, fmt.Errorf("organization %s does not exist", old.OrgID.String()))
			if reportFile != nil {
				reportFile.WriteString(fmt.Sprintf("error migrating member %s:%s: %s\n", old.OrgID.String(), old.ID.String(), "organization does not exist"))
			}
			return nil
		}

		// Populate all missing fields
		if member.Organization == "" {
			member.Organization = org.Name
		}

		if member.Workspace == "" {
			member.Workspace = org.Domain
		}

		if member.ProfessionSegment == db.ProfessionSegmentUnspecified {
			member.ProfessionSegment = db.ProfessionSegmentPersonal
		}

		if len(member.DeveloperSegment) == 0 {
			member.DeveloperSegment = []db.DeveloperSegment{db.DeveloperSegmentSomethingElse}
		}

		// Because ULIDs are returned in time order and there is no way to remove team
		// members yet, the first listed member is the owner. Everybody else in the
		// organization is invited.
		if ulids.IsZero(org.owner) {
			org.owner = member.ID
			orgs[old.OrgID] = org
			member.Invited = false

			// Ensure JoinedAt is set for owners
			if member.JoinedAt.IsZero() {
				member.JoinedAt = member.Created
			}
		} else {
			// Sanity check that the invited member doesn't predate the owner
			if member.Created.Before(old.Created) {
				err = fmt.Errorf("db assumption invalid: member %s:%s was created before the owner", old.OrgID.String(), old.ID.String())
				errs = multierror.Append(errs, err)
				if reportFile != nil {
					reportFile.WriteString(fmt.Sprintf("error migrating member %s:%s: %s\n", old.OrgID.String(), old.ID.String(), err.Error()))
				}
				return nil
			}

			member.Invited = true
		}

		// Update the member record in the database
		if dry {
			fmt.Printf("Would migrate record: %v\n", member.ToAPI())
			fmt.Println()
		} else {
			if err = db.UpdateMember(ctx, member); err != nil {
				errs = multierror.Append(errs, err)
				if reportFile != nil {
					reportFile.WriteString(fmt.Sprintf("error migrating member %s:%s: %s\n", old.OrgID.String(), old.ID.String(), err.Error()))
				}
				return nil
			}

			fmt.Println(fmt.Sprintf("migrated record: %v", member.ToAPI()))
			fmt.Println()
			if reportFile != nil {
				reportFile.WriteString(fmt.Sprintf("migrated record: %v\n\n", member.ToAPI()))
			}
		}

		migrated++
		return nil
	}

	if _, err = db.List(ctx, nil, nil, db.MembersNamespace, onMember, nil); err != nil {
		return cli.Exit(err, 1)
	}

	if errs != nil {
		fmt.Println(errs.Error())
	}

	if dry {
		fmt.Printf("Would migrate %d members out of %d\n", migrated, members)
	} else {
		fmt.Printf("Migrated %d members out of %d\n", migrated, members)
	}

	return nil
}

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
