package db

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rotationalio/ensign/pkg"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	lastAppliedSQL     = `SELECT id FROM migrations ORDER BY id DESC LIMIT 1;`
	insertMigrationSQL = `INSERT INTO migrations (id, name, version, created) VALUES ($1, $2, $3, datetime('now'));`
)

// Initialize schema applies any unapplied migrations to the database and should be run
// when the database is first connected to. If empty is true then the migration table is
// created and all migrations are applied. If it is not true then the current migration
// of the database is fetched and all unapplied migrations are applied.
func InitializeSchema(empty bool) (err error) {
	lastApplied := -1
	if !empty {
		// Fetch the latest migration applied to the database
		if err = conn.QueryRow(lastAppliedSQL).Scan(&lastApplied); err != nil {
			return fmt.Errorf("could not fetch last applied migration: %s", err)
		}
	}

	var migrations []*Migration
	if migrations, err = Migrations(); err != nil {
		return err
	}

	for _, migration := range migrations {
		if migration.ID > lastApplied {
			var query string
			if query, err = migration.SQL(); err != nil {
				return err
			}

			if _, err = conn.Exec(query); err != nil {
				return fmt.Errorf("could not apply schema %d: %s", migration.ID, err)
			}

			if _, err = conn.Exec(insertMigrationSQL, migration.ID, migration.Name, pkg.Version()); err != nil {
				return fmt.Errorf("could not insert migration record for %d: %s", migration.ID, err)
			}
		}
	}

	return nil
}

// Migrations contains the SQL commands from the migrations directory and is used to
// ensure that the database has the most current and up to date schema.
//
//go:embed migrations/*.sql
var migrations embed.FS

// Process migration file names
var (
	caser  = cases.Title(language.English)
	pathre = regexp.MustCompile(`^(\d+)_(\w+)\.sql$`)
)

// Migration is used to represent both a SQL migration from the embedded file system and
// a migration record in the database. These records are compared to ensure the database
// is as up to date as possible before the application starts.
type Migration struct {
	ID      int       // The unique sequence ID of the migration
	Name    string    // The human readable name of the migration
	Version string    // The package version when the migration was applied
	Created time.Time // The timestamp when the migration was applied
	Path    string    // The path of the migration in the filesystem
}

// Migrations returns the migration files from the embedded file system.
func Migrations() (data []*Migration, err error) {
	var entries []fs.DirEntry
	if entries, err = migrations.ReadDir("migrations"); err != nil {
		return nil, err
	}

	data = make([]*Migration, 0, len(entries))
	for _, entry := range entries {
		groups := pathre.FindStringSubmatch(entry.Name())
		if len(groups) != 3 {
			return nil, fmt.Errorf("could not parse %q into migration", entry.Name())
		}

		var id int
		if id, err = strconv.Atoi(groups[1]); err != nil {
			return nil, fmt.Errorf("could not parse %q into migration", entry.Name())
		}

		data = append(data, &Migration{
			ID:   id,
			Name: caser.String(strings.Join(strings.Split(groups[2], "_"), " ")),
			Path: entry.Name(),
		})
	}
	return data, nil
}

// SQL loads the schema sql query from the embedded file on disk.
func (m *Migration) SQL() (_ string, err error) {
	if m.Path == "" {
		return "", fmt.Errorf("cannot fetch sql for migration %d", m.ID)
	}

	var data []byte
	if data, err = migrations.ReadFile(filepath.Join("migrations", m.Path)); err != nil {
		return "", fmt.Errorf("could not read sql for migration %d: %s", m.ID, err)
	}

	return string(data), nil
}
