package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
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
			Name:   "topics:list",
			Usage:  "list the topics, optionally filtered by a specific project",
			Before: connectDB,
			After:  closeDB,
			Action: topicsList,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "project",
					Aliases: []string{"p"},
					Usage:   "ulid of the project to filter on",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// Topic Management Commands
//===========================================================================

func topicsList(c *cli.Context) (err error) {
	var topics iterator.TopicIterator
	if project := c.String("project"); project != "" {
		var projectID ulid.ULID
		if projectID, err = ulid.Parse(project); err != nil {
			return cli.Exit(err, 1)
		}
		topics = meta.ListTopics(projectID)
	} else {
		topics = meta.ListAllTopics()
	}

	defer topics.Release()

	// Create a tab writer to output the topics list to
	// TODO: also allow writing to a CSV file
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	defer tabs.Flush()

	// Write the header of the tab table
	fmt.Fprintln(tabs, "Project ID\tTopic ID\tTopic\tEvents\tBytes")

	for topics.Next() {
		var topic *api.Topic
		if topic, err = topics.Topic(); err != nil {
			return cli.Exit(err, 1)
		}

		projectID, _ := topic.ParseProjectID()
		topicID, _ := topic.ParseTopicID()

		var events, dataSize uint64

		var info *api.TopicInfo
		if info, err = meta.TopicInfo(topicID); err == nil {
			events = info.Events
			dataSize = info.DataSizeBytes
		}

		fmt.Fprintf(tabs, "%s\t%s\t%s\t%d\t%d\n", projectID, topicID, topic.Name, events, dataSize)
	}

	if err = topics.Error(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

//===========================================================================
// Temporary Commands
//===========================================================================

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
