---
title: "Ensign for Application Developers"
weight: 30
date: 2023-05-17T17:03:41-04:00
---

Hi there! This tutorial is targeted towards Golang application developers. If you are interested in or currently writing event-driven applications in Go you are in the right place! In this code-driven tutorial we will use the Ensign Golang SDK to publish curated tweets to an event stream and retrieve them in real time.

If you came here for the code the full example is available [here](https://github.com/rotationalio/ensign-examples/tree/main/go/tweets).

## Prerequisites

To follow along with this tutorial you'll need to:

- [Generate an API key to access Ensign]({{< ref "/getting-started/ensign#getting-started" >}})
- [Set up a developer account with Twitter (it's free)](https://developer.twitter.com/en/docs/twitter-api/getting-started/getting-access-to-the-twitter-api)
- [Add a phone number to your Twitter developer account](https://help.twitter.com/en/managing-your-account/how-to-add-a-phone-number-to-your-account)
- [Set up your GOPATH and workspace](https://go.dev/doc/gopath_code)

## Project Setup

The first thing we need to do is setup an environment to run our code. Let's create a blank module with a suitable name for our project:

```bash
$ mkdir tweets
$ go mod init github.com/rotationalio/ensign-examples/go/tweets
```

Next we'll need to install the Go SDK client and its dependencies from the GitHub [repo](https://github.com/rotationalio/ensign). In this tutorial we also use the [go-twitter](https://github.com/g8rswimmer/go-twitter) client to interact with the twitter API (although you can also create the requests yourself)!

```bash
$ go get -u github.com/rotationalio/ensign/sdks/go@latest
$ go get -u github.com/g8rswimmer/go-twitter/v2@latest
```

Our project needs a [publisher]({{< ref "/eventing/glossary#publisher" >}}) to write events to Ensign and a [subscriber]({{< ref "/eventing/glossary#subscriber" >}}) to read those events (asynchronously). In a real application these would most likely be independent microservices that run in different execution contexts (e.g. containers in a k8s cluster or even across different regions). Let's create separate packages for the two command line applications as well as a shared package for our event schemas.

```bash
$ mkdir publish
$ mkdir subscribe
$ mkdir schemas
```

## Sourcing Tweets

In event-driven systems, events are the main unit of data. In production applications events might be sourced from user actions, IoT devices, webhooks, or act as control signals between microservices.

For this example our data source is curated tweets from twitter. Create a file called `main.go` in the `publish` directory and add the following code to it.

```golang
package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"context"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func main() {
    var (
		err   error
		token string
	)

	if token = os.Getenv("TWITTER_API_BEARER_TOKEN"); token == "" {
		panic("TWITTER_API_BEARER_TOKEN environment variable is required")
	}

	query := flag.String("query", "distributed systems", "Twitter search query")
	flag.Parse()

	tweets := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *tweets.TweetRecentSearchResponse
	if rep, err = client.TweetRecentSearch(ctx, *query, twitter.TweetRecentSearchOpts{}); err != nil {
		panic(err)
	}

	for _, errs := range rep.Raw.Errors {
		fmt.Printf("Error: %s\n", errs.Detail)
	}

	for _, tweet := range rep.Raw.Tweets {
		fmt.Printf("%s: %s\n", tweet.AuthorID, tweet.Text)
	}
}
```

This is a simple command line application that will retrieve a single page of search results from twitter and print them out. Feel free to build the program and run it with any search query to make sure it works!

```bash
$ export TWITTER_API_BEARER_TOKEN=# Your Twitter API bearer token goes here
$ cd publish
$ go build -o publish main.go
$ ./publish --query "distributed systems"
```

## Creating a Publisher

Now that we have a data source, the next step is to create an Ensign client using the Client ID and Client Secret pair you received when [generating your API key]({{< ref "/getting-started/ensign#getting-started" >}}).

```golang
import (
    ...
	twitter "github.com/g8rswimmer/go-twitter/v2"
	ensign "github.com/rotationalio/go-ensign"
)
```

```golang

const DistSysTweets = "distsys-tweets"


func main() {
	var (
		err   error
		token string
	)

	if token = os.Getenv("TWITTER_API_BEARER_TOKEN"); token == "" {
		panic("TWITTER_API_BEARER_TOKEN environment variable is required")
	}

	query := flag.String("query", "distributed systems", "Twitter search query")
	flag.Parse()

	// ENSIGN_CLIENT_ID and ENSIGN_CLIENT_SECRET environment variables must be set
	var client *ensign.Client
	if client, err = ensign.New(&ensign.Options{
		ClientID:     os.Getenv("ENSIGN_CLIENT_ID"),
		ClientSecret: os.Getenv("ENSIGN_CLIENT_SECRET"),
	}); err != nil {
		panic("failed to create Ensign client: " + err.Error())
	}

	// Check to see if topic exists and create it if not
	exists, err := client.TopicExists(context.Background(), DistSysTweets)
	if err != nil {
		panic(fmt.Errorf("unable to check topic existence: %s", err))
	}

	var topicID string
	if !exists {
		if topicID, err = client.CreateTopic(context.Background(), DistSysTweets); err != nil {
			panic(fmt.Errorf("unable to create topic: %s", err))
		}
	} else {
		topics, err := client.ListTopics(context.Background())
		if err != nil {
			panic(fmt.Errorf("unable to retrieve project topics: %s", err))
		}

		for _, topic := range topics {
			if topic.Name == DistSysTweets {
				var topicULID ulid.ULID
				if err = topicULID.UnmarshalBinary(topic.Id); err != nil {
					panic(fmt.Errorf("unable to retrieve requested topic: %s", err))
				}
				topicID = topicULID.String()
			}
		}
	}
...
```

In the Go SDK, creating a `Publisher` interface from the client is straightforward.

```golang
	var pub ensign.Publisher
	if pub, err = client.Publish(context.Background()); err != nil {
		panic("failed to create publisher from Ensign client: " + err.Error())
	}
```

## Publishing Events

In Ensign, events include a lot more than the data itself. As we can see from the [protocol buffer](https://github.com/rotationalio/ensign/blob/main/pkg/api/v1beta1/event.pb.go), events are self-descriptive and are quite flexible.

```golang
type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TopicId       string                 `protobuf:"bytes,2,opt,name=topic_id,json=topicId,proto3" json:"topic_id,omitempty"`
	Mimetype      v1beta1.MIME           `protobuf:"varint,3,opt,name=mimetype,proto3,enum=mimetype.v1beta1.MIME" json:"mimetype,omitempty"`
	Type          *Type                  `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Key           []byte                 `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	Data          []byte                 `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	Encryption    *Encryption            `protobuf:"bytes,7,opt,name=encryption,proto3" json:"encryption,omitempty"`
	Compression   *Compression           `protobuf:"bytes,8,opt,name=compression,proto3" json:"compression,omitempty"`
	Geography     *Region                `protobuf:"bytes,9,opt,name=geography,proto3" json:"geography,omitempty"`
	Publisher     *Publisher             `protobuf:"bytes,10,opt,name=publisher,proto3" json:"publisher,omitempty"`
	UserDefinedId string                 `protobuf:"bytes,11,opt,name=user_defined_id,json=userDefinedId,proto3" json:"user_defined_id,omitempty"`
	Created       *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=created,proto3" json:"created,omitempty"`
	Committed     *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=committed,proto3" json:"committed,omitempty"`
}
```

For this tutorial we are mostly concerned with the following fields.

- `TopicId`: Events are organized into [topics]({{< ref "/eventing/glossary#topic" >}}) and events in a topic usually follow a similar schema
- `Mimetype`: In Ensign all event data is generic "blob" data to allow for heterogenous event streams. The mimetype allows subcribers to deserialize data back into an understandable format.
- `Type`: Events in Ensign are tagged with schema type and versioning info to allow publishers and subscribers to lookup schemas in a shared registry. This is important because certain serialization methods (e.g. protobuf, parquet) require explicit schemas for deserialization and schema-less methods (e.g. JSON) can be enhanced with versioning.

In this example we can get away with structured JSON. In production workflows we would most likely want to store the definition in a schema registry but for now let's add it to `tweets.go` in the `schemas` directory so both our producer and subscriber can access it.

```golang
package schemas

type Tweet struct {
	Author    string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}
```

Now that we know how to serialize JSON, in the tweet loop instead of printing to the console let's go ahead and publish some events.

```golang
	for _, tweet := range rep.Raw.Tweets {
		e := &api.Event{
			TopicId:  "tweets",
			Mimetype: mimetype.ApplicationJSON,
			Type: &api.Type{
				Name:    "tweet",
				Version: 1,
			},
		}

        tweetObj := &schemas.Tweet{
            Author:    tweet.AuthorID,
            Text:      tweet.Text,
            CreatedAt: tweet.CreatedAt,
        }
		if e.Data, err = json.Marshal(tweetObj); err != nil {
			panic("could not marshal tweet to JSON: " + err.Error())
		}

        // Publish the event to the Ensign topic
		pub.Publish(topicID, e)

        // Check for errors
        if err = pub.Err(); err != nil {
			panic("failed to publish event(s): " + err.Error())
		}
	}
```

If your IDE did not resolve the imports for you, you will need to specify them manually:

```golang
import (
    ...
	api "github.com/rotationalio/go-ensign/api/v1beta1"
    mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
    ...
)
```

*Note that `pub.Publish(e)` does not return an immediate error, it's an asynchronous operation so if we want to check for errors we have to do so after the fact. This means that we can't be sure which event actually triggered the error.*

Finally, to make our publisher feel like a real service, we can add an outer loop with a ticker so that the program periodically pulls the most recent tweets our search query of choice. Another useful improvement might be to utilize the `SinceID` on the twitter search options so that we aren't producing duplicate tweets!

```golang
	ticker := time.NewTicker(10 * time.Second)
	sinceID := ""
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("searching for tweets...")
			opts := twitter.TweetRecentSearchOpts{
				SortOrder: twitter.TweetSearchSortOrderRecency,
				SinceID:   sinceID,
			}

			var rep *twitter.TweetRecentSearchResponse
			if rep, err = tweets.TweetRecentSearch(ctx, *query, opts); err != nil {
				panic(err)
			}

			for _, errs := range rep.Raw.Errors {
				fmt.Printf("Error: %s\n", errs.Detail)
			}

			for _, tweet := range rep.Raw.Tweets {
				e := &api.Event{
					TopicId:  "tweets",
					Mimetype: mimetype.ApplicationJSON,
					Type: &api.Type{
						Name:    "Generic",
						Version: 1,
					},
				}

				if e.Data, err = json.Marshal(tweet); err != nil {
					panic("could not marshal tweet to JSON: " + err.Error())
				}

		        // Publish the event to the Ensign topic
				pub.Publish(topicID, e)

				if err = pub.Err(); err != nil {
					panic("failed to publish event(s): " + err.Error())
				}

				fmt.Printf("published tweet with ID: %s\n", tweet.ID)
			}

			if len(rep.Raw.Tweets) > 0 {
				sinceID = rep.Raw.Tweets[0].ID
			}
		}
	}
```

At this point our publisher will be able to request some new tweets from Twitter every 10 seconds and publish them as events to the `tweets` topic. Go ahead and try it out!

```bash
$ export ENSIGN_CLIENT_ID=# Your Ensign Client ID goes here
$ export ENSIGN_CLIENT_SECRET=# Your Ensign Client Secret goes here
$ go build -o publish main.go
$ ./publish --query "otters"
```
*Note: Here the Ensign Client ID and Client Secret are retrieved from environment variables but it's also possible to specify them in [code]({{< ref "/getting-started/ensign#create-a-client" >}})*

## Creating a subscriber

Similarly to the `Publisher`, a `Subscriber` interface can be created from an Ensign client. Once created, the `Subscriber` allows us to read events directly from a Go channel. Create a `main.go` in the `subscribe` directory and add the following code to it.

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rotationalio/ensign-examples/go/tweets/schemas"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	ensign "github.com/rotationalio/go-ensign"
)

func main() {
	var (
		err    error
		client *ensign.Client
	)

	// ENSIGN_CLIENT_ID and ENSIGN_CLIENT_SECRET environment variables must be set
	if client, err = ensign.New(&ensign.Options{
		ClientID:     os.Getenv("ENSIGN_CLIENT_ID"),
		ClientSecret: os.Getenv("ENSIGN_CLIENT_SECRET"),
	}); err != nil {
		panic("failed to create Ensign client: " + err.Error())
	}

	// Create a subscriber from the client
	var sub ensign.Subscriber
	if sub, err = client.Subscribe(context.Background(), topicID); err != nil {
		panic("failed to create subscriber from client: " + err.Error())
	}
	defer sub.Close()

	// Create the event stream as a channel
	var events <-chan *api.Event
	if events, err = sub.Subscribe(); err != nil {
		panic("failed to create subscribe stream: " + err.Error())
	}

	// Events are processed as they show up on the channel
	for event := range events {
		tweet := &schemas.Tweet{}
		if err = json.Unmarshal(event.Data, tweet); err != nil {
			panic("failed to unmarshal event: " + err.Error())
		}

		fmt.Printf("received tweet %s\n", tweet.ID)
		fmt.Println(tweet.Text)
		fmt.Println()
	}
}
```

At this point you should be able to build and the run the subscriber in a second command window to retrieve tweet events in real time!

```bash
$ export ENSIGN_CLIENT_ID=# Your Ensign Client ID
$ export ENSIGN_CLIENT_SECRET=# Your Ensign Client Secret
$ cd subscribe
$ go build -o subscribe main.go
$ ./subscribe
```

## What Next?

Hopefully this gets you on the right track and inspires some ideas for event-driven applications. If this example were to become a real application, here are some things we might consider.

### Event Schemas

Remember that an Ensign event encodes a lot of metadata. When dealing with more strutured or versioned serialization formats such as protobuf, we definitely want to consider adding some logic to the subscriber to lookup the event schema in the schema registry or a local cache with the `event.Type` field.

### Additional Topic Streams

With Ensign it's easy to scale up by adding new topics. We might want to have different topics for error observability (e.g. if the Twitter API changes or schemas unexpectedly change), metrics capturing, or different types of Twitter queries.

### Downstream Processing

Once we have an event stream, what do we do with it? A traditional approach is to capture data into a database for persistence and to make it easy to materialize data views for application users. This is certainly possible with Ensign. However, Ensign also offers persistence of event streams, which makes it possible to perform historical queries on the streams themselves.