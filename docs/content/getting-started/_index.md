---
title: "Getting Started"
weight: 10
bookFlatSection: false
bookToc: true
bookHidden: false
bookCollapseSection: false
bookSearchExclude: false
---

# Welcome!

Ready to get started with eventing? Let's go!


## What is Ensign?

Ensign is a new eventing tool that make it fast, convenient, and fun to create  event-driven microservices without needing a big team of devOps or platform engineers. All you need is a free API key to get started.

## Getting Started

The first step is to get an Ensign API key by visiting rotational.io/ensign. Similar to getting a developer API key for [Youtube](https://developers.google.com/youtube/v3/getting-started), [Twitter](https://developer.twitter.com/en/docs/twitter-api/getting-started/getting-access-to-the-twitter-api) or [Data.gov](https://api.data.gov/docs/api-key/), you will need an API key to use Ensign and to follow along with the rest of this Quickstart guide.

Your key consists of two parts, a `ClientID` and a `ClientSecret`. The `ClientID` uniquely identifies you, and the `ClientSecret` proves that you have permission to create and access event data.

| API Key Component Name | Length | Characters             | Example                                                          |
|-------------------|--------|------------------------|------------------------------------------------------------------|
| ClientID          | 32     | alphabetic (no digits) | DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa                                 |
| ClientSecret      | 64     | alphanumeric           | wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS |


Together, the `ClientID` and `ClientSecret` uniquely identify you. They enable you to create Ensign topics, publishers, and subscribers, which will be the building blocks of your microservice! Keep in mind that the `ClientID` and `ClientSecret` should be kept private and not shared.

### Prerequisites

Ensign's SDK currently supports Golang (Python and Javascript coming soon!).
If you haven't already:
- [download and install Golang](https://go.dev/doc/install) according to your operating system
- [set up your GOPATH and workspace](https://go.dev/doc/gopath_code)

### Install Ensign

In your command line, type the following to install the ensign API, SDK, and library code for Go:

```bash
go install github.com/rotationalio/ensign@latest
```

### Create a Client

After you've made a new Go project for this example, create a `main.go` file and add the dependencies you'll need, which will include importing the Ensign API, SDK, and mimetypes.

Next, create an Ensign client, which is similar to establishing a connection to a database like PostgreSQL or Mongo. To create the client, use the `New` method and pass in an `ensign.Options` struct that specifies your Client ID and Client Secret (described in the section above on getting an API key).

```golang
package main

import (
    "fmt"
    "time"
    "context"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/mimetype/v1beta1"
	ensign "github.com/rotationalio/ensign/sdks/go"
)


client, err := ensign.New(&ensign.Options{
	ClientID: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
	ClientSecret: "wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS",
})
if err != nil {
	fmt.Errorf("could not create client: %s", err)
}
```

Congratulations, you now have an open connection to Ensign!

### Create a Publisher

The next step is to start publishing data onto your event stream. Start by creating a publisher using the `Publish` method:

```golang
pub, err := client.Publish(context.Background())
if err != nil {
    fmt.Errorf("could not create publisher: %s", err)
}
```


Next, we need some data! Generally this is the place where you'd connect to your live data source (a database, Twitter feed, weather data, etc). But to keep things simple, we'll just create a single event, which starts with a map.

```golang
data := make(map[string]string)
data["sender"] = "Twyla"
data["timestamp"] = time.Now().String()
data["message"] = "Let's get this started!"
```

Next, we will convert our map into an event, which will allow you to specify the mimetype of the message you intend to send (in this case, we'll say it's JSON), and the event type (which will be a generic event for this example). You'll also need to pass in a `TopicId`, which will be a string. If you aren't sure what `TopicId` to use, you can quickly [log into your Ensign dashboard](https://rotational.io/ensign/) and look it up. For this example, we'll pretend it's `"quality-lemon-time"`:

```golang
e := &api.Event{
    TopicId:  "quality-lemon-time",
    Mimetype: mimetype.ApplicationJSON,
    Type: &api.Type{
        Name:    "Generic",
        Version: 1,
    },
}
```

Next, we'll marshall our dictionary into the `Data` attribute of our sample event, and publish it by calling the `Publish` method on the publisher we created above:

```golang
e.Data, _ = json.Marshal(data)
pub.Publish(e)
```

### Create a Subscriber

Creating a subscriber is a bit more straightforward:

```golang
sub, err := client.Subscribe(context.Background())
if err != nil {
    fmt.Errorf("could not create subscriber: %s", err)
}

msg := sub.Subscribe()
fmt.Sprintln(msg.Data)
```

## Next Steps

You're already well on your way to building your first event-driven microservice with Ensign!

If you're ready to see some more advanced examples with code, check out the [End-to-end Examples]({{< relref "examples">}}).

If you're looking for more on the basics of event-driven systems, check out [Eventing 101]({{< relref "eventing">}}).

Happy eventing!