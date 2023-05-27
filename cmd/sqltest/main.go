package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/mattn/go-sqlite3"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/rotationalio/ensign/pkg/utils/sqlite"
)

var (
	db  *sql.DB
	src *sqlite.Conn
)

const schema = `CREATE TABLE IF NOT EXISTS entries (
	id       INTEGER PRIMARY KEY,
	name     TEXT NOT NULL,
	blob     BLOB NOT NULL,
	created  TEXT NOT NULL,
	modified TEXT NOT NULL
);`

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "sqltest"
	app.Version = pkg.Version()
	app.Before = setupLogger
	app.Usage = "test a long running sqlite3 database"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "verbosity",
			Aliases: []string{"L"},
			Usage:   "set the zerolog level",
			Value:   "info",
			EnvVars: []string{"ENSIGN_LOG_LEVEL"},
		},
		&cli.BoolFlag{
			Name:    "console",
			Aliases: []string{"C"},
			Usage:   "human readable console log instead of json",
			Value:   true,
			EnvVars: []string{"ENSIGN_CONSOLE_LOG"},
		},
		&cli.StringFlag{
			Name:    "data",
			Aliases: []string{"d", "db"},
			Usage:   "path to write sqlite3 database to",
			Value:   "fixtures/longrun.db",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "run the long running sqlite3 test",
			Before: connectDB,
			After:  closeDB,
			Action: runTest,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "backups",
					Aliases: []string{"b"},
					Usage:   "directory to write backups to",
					Value:   "fixtures/backups",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("could not execute cli app")
	}
}

func setupLogger(c *cli.Context) (err error) {
	switch strings.ToLower(c.String("verbosity")) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		return cli.Exit(fmt.Errorf("unknown log level %q", c.String("verbosity")), 1)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg
	zerolog.DurationFieldInteger = false
	zerolog.DurationFieldUnit = time.Millisecond

	if c.Bool("console") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		// Add the severity hook for GCP logging
		var gcpHook logger.SeverityHook
		log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
	}

	return nil
}

func connectDB(c *cli.Context) (err error) {
	// Check if the file exists, if it doesn't exist it will be created and all
	// migrations will be applied to the database. Otherwise the code will attempt
	// to only apply migrations that have not yet been applied.
	empty := false
	if _, err := os.Stat(c.String("data")); os.IsNotExist(err) {
		empty = true
	}

	if db, err = sql.Open("ensign_sqlite3", c.String("data")); err != nil {
		return cli.Exit(err, 1)
	}
	db.Ping()

	var ok bool
	if src, ok = sqlite.GetLastConn(); !ok {
		return cli.Exit("could not get source sqlite3 connection", 1)
	}

	if empty {
		// Initialize the schema
		if _, err = db.Exec(schema); err != nil {
			return cli.Exit(err, 1)
		}
	}

	return nil
}

func closeDB(c *cli.Context) (err error) {
	if err = db.Close(); err != nil {
		return cli.Exit(err, 1)
	}

	db = nil
	src = nil
	return nil
}

func runTest(c *cli.Context) error {
	// Stop on CTRL+C
	nRoutines := 2
	done := make(chan struct{}, nRoutines)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Warn().Msg("shutting down")
		for i := 0; i < nRoutines; i++ {
			done <- struct{}{}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(nRoutines)

	// The first go routine randomly inserts and updates data every 750ms
	go func(done <-chan struct{}) {
		var inserts, updates uint64

		defer wg.Done()
		defer func() {
			log.Info().Uint64("inserts", inserts).Uint64("updates", updates).Msg("insert and updates stopped")
		}()

		// Ensure at least one row is in the database
		if err := insertRow(); err != nil {
			log.Error().Err(err).Msg("could not insert initial row into database")
			return
		}
		inserts++

		log.Info().Dur("interval", 750*time.Millisecond).Msg("starting insert and update go routine")
		ticker := time.NewTicker(750 * time.Millisecond)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
			}

			// Decide whether to insert or update
			if rand.Float64() < 0.38 {
				if err := insertRow(); err != nil {
					log.Error().Err(err).Msg("could not insert row into database")
					return
				}
				inserts++
				log.Debug().Uint64("updates", updates).Uint64("inserts", inserts).Msg("row inserted")
			} else {
				if err := updateRow(); err != nil {
					log.Error().Err(err).Msg("could not update row in database")
					return
				}
				updates++
				log.Debug().Uint64("updates", updates).Uint64("inserts", inserts).Msg("row updated")
			}
		}
	}(done)

	// The second go routine backs up the database every 5 minutes
	go func(done <-chan struct{}) {
		defer wg.Done()
		defer log.Info().Msg("backups stopped")

		backupDir := c.String("backups")
		log.Info().Str("backups", backupDir).Dur("interval", 5*time.Minute).Msg("starting backup go routine")
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
			}

			// Create a new backup file in the backup directory
			file := filepath.Join(backupDir, fmt.Sprintf("sqlite-%s.bak.db", time.Now().Format("200601021504")))

			// Perform the backup to the specified file
			if err := backup(file); err != nil {
				log.Error().Err(err).Int("nconns", sqlite.NumConns()).Str("file", file).Msg("could not create backup")
				return
			}
			log.Info().Int("nconns", sqlite.NumConns()).Str("file", file).Msg("backup created")
		}
	}(done)

	wg.Wait()
	return nil
}

func backup(dstPath string) (err error) {
	// Open a connection to the destination database
	var dstDB *sql.DB
	if dstDB, err = sql.Open("ensign_sqlite3", dstPath); err != nil {
		return err
	}
	defer dstDB.Close()
	dstDB.Ping()

	dstConn, ok := sqlite.GetLastConn()
	if !ok {
		return errors.New("could not get sqlite3 connection to dst database")
	}

	// Perform the backup
	var backup *sqlite3.SQLiteBackup
	if backup, err = dstConn.Backup("main", src, "main"); err != nil {
		return err
	}

	// For now, attempt complete backup rather than looping backup
	if _, err = backup.Step(-1); err != nil {
		return err
	}

	return backup.Finish()
}

func insertRow() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	name := sql.Named("name", fmt.Sprintf("%04x", rand.Int63()))
	created := sql.Named("created", time.Now().Format(time.RFC3339Nano))

	data := make([]byte, 512)
	if _, err = rand.Read(data); err != nil {
		return err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("INSERT INTO entries (name, blob, created, modified) VALUES (:name, :data, :created, :created)", name, sql.Named("data", data), created); err != nil {
		return err
	}

	return tx.Commit()
}

func updateRow() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	data := make([]byte, 512)
	if _, err = rand.Read(data); err != nil {
		return err
	}

	modified := sql.Named("modified", time.Now().Format(time.RFC3339Nano))

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("UPDATE entries SET blob=:data, modified=:modified WHERE id=(SELECT id FROM entries ORDER BY RANDOM() LIMIT 1) ", sql.Named("data", data), modified); err != nil {
		return err
	}

	return tx.Commit()
}
