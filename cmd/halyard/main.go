package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/urfave/cli/v2"
)

var (
	conf config.Config
	data store.EventStore
	meta store.MetaStore
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "halyard"
	app.Version = pkg.Version()
	app.Usage = "utilities and administrative commands for ensign nodes"
	app.Flags = []cli.Flag{}
	app.Before = configure
	app.Commands = []*cli.Command{
		{
			Name:   "migrate:duplicatepolicy",
			Usage:  "set the topic duplication policies to none",
			Before: connectDB,
			After:  closeDB,
			Action: migrateDuplicatePolicies,
			Flags:  []cli.Flag{},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// Temporary Commands
//===========================================================================

func migrateDuplicatePolicies(c *cli.Context) (err error) {
	topics := meta.ListAllTopics()
	defer topics.Release()

	rowsAffected := 0
	for topics.Next() {
		var topic *api.Topic
		if topic, err = topics.Topic(); err != nil {
			return cli.Exit(fmt.Errorf("could not unmarshal topic %s: %w", topics.Key(), err), 1)
		}

		if topic.Deduplication != nil && topic.Deduplication.Strategy != api.Deduplication_UNKNOWN {
			continue
		}

		topic.Deduplication = &api.Deduplication{
			Strategy: api.Deduplication_NONE,
			Offset:   api.Deduplication_OFFSET_EARLIEST,
		}

		if err = meta.UpdateTopic(topic); err != nil {
			return cli.Exit(fmt.Errorf("could not update topic %s: %w", topics.Key(), err), 1)
		}

		rowsAffected++
	}

	if err = topics.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("duplicate policy migration complete, %d rows affected\n", rowsAffected)
	return nil
}

//===========================================================================
// Helper Commands
//===========================================================================

func configure(c *cli.Context) (err error) {
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if !conf.Maintenance {
		fmt.Println("\033[33m\033[1mWARNING: Ensign is not in maintenance mode.\033[0m")
	}

	return nil
}

func connectDB(c *cli.Context) (err error) {
	if data, meta, err = store.Open(conf.Storage); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func closeDB(c *cli.Context) (err error) {
	if err = errors.Join(data.Close(), meta.Close()); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}
