---
title: "Golang"
weight: 50
date: 2023-08-11T09:03:41-04:00
---

The Go SDK is a great choice for developers who wish to integrate Ensign into their existing projects. Or if you're starting from scratch and want to take advantage of Golang's static typing and high-performance multiprocessing, that's cool too.

In this example we'll create a simple Go project from scratch to publish and subscribe to Ensign!

### Prerequisites

- [create a free Ensign account and API key](https://rotational.app)
- [download and install Golang](https://go.dev/doc/install) according to your operating system

### Project Setup

As with any new Go project, start by creating a directory and module to start adding dependencies to.

```bash
mkdir hello-ensign
cd hello-ensign
go mod init example.com/hello/ensign
```

The next step is to install the official Go SDK for Ensign.

```bash
go get github.com/rotationalio/go-ensign
```

<a name="create-a-client"></a>
### Create a Client

Create a `main.go` file and create the Ensign client in code, which is similar to a database client like PostgreSQL or Mongo.

```golang
package main

import (
	"context"
	"fmt"

	ensign "github.com/rotationalio/go-ensign"
)

func main() {
    // Create an Ensign client
	client, err := ensign.New()
	if err != nil {
		panic(fmt.Errorf("could not create client: %s", err))
	}
	defer client.Close()

    // Fetch status from Ensign
	ctx := context.Background()
	state, err := client.Status(ctx)
	if err != nil {
		panic(fmt.Errorf("could not get status from Ensign: %s", err))
	}
	fmt.Println(state.Status, state.Version)
}
```

The Go SDK requires a Client ID and Client Secret to communicate with Ensign. We recommend specifying them in the environment like so (replace with the values in your API key).

```bash
export ENSIGN_CLIENT_ID=DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa
export ENSIGN_CLIENT_SECRET=wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS
```

If you find yourself having to manage multiple API keys on the same machine, you can also specify a path to a JSON file with your credentials.

**my_project_key.json**
```json
{
    "ClientID": "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
    "ClientSecret": "wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS"
}
```

```golang
client, err := ensign.New(ensign.WithLoadCredentials("my_project_key.json"))
```

Run the code as a Go program.

```bash
go run main.go
```

If you see a message like the following, then congratulations! You've successfully connected to Ensign!

```HEALTHY 0.12.8-beta.23 ([GIT HASH])```

### Make Some Data

Next, we need some data! Generally this is the place where you'd connect to your live data source (a database, weather data, etc). But to keep things simple, we'll just create a single event, which starts with a map.

```golang
data := make(map[string]string)
data["sender"] = "Twyla"
data["timestamp"] = time.Now().String()
data["message"] = "Let's get this started!"
```

Next, we will convert our map into an event, which will allow you to specify the mimetype of the message you intend to send (in this case, we'll say it's JSON), and the event type (which will be a generic event for this example).

```golang
event := &ensign.Event{
    Mimetype: mimetype.ApplicationJSON,
    Type: &api.Type{
		Name:         "Generic",
		MajorVersion: 1,
		MinorVersion: 0,
		PatchVersion: 0,
    },
}
```

Next, we'll marshal our dictionary into the `Data` attribute of our sample event

```golang
if event.Data, err = json.Marshal(data); err != nil {
    panic("could not marshal data to JSON: " + err.Error())
}
```

### Publish Your Event

Now we can publish your event by calling the `Publish` method on the Ensign client we created above. You'll also need to pass in a topic name, which will be a string. If you aren't sure what topic to use, you can quickly [log into your Ensign dashboard](https://rotational.app) and look it up.

```golang
client.Publish("quality-lemon-time", event)
```

You can publish many events at a time if you want!

```golang
client.Publish("quality-lemon-time", event, event2, event3, event4)
```

### Create a Subscriber

So now you've published some events to a topic. We can consume those events using the `Subscribe` method. `Subscribe` works a bit differently than `Publish`; it returns a `Subscription` with a Go channel that you can read events from.

```golang
sub, err := client.Subscribe("quality-lemon-time")
if err != nil {
    panic(fmt.Errorf("could not create subscriber: %s", err))
}

for event := range sub.C {
    var m map[string]string
    if err := json.Unmarshal(event.Data, &m); err != nil {
        panic(fmt.Errorf("failed to unmarshal message: %s", err))
    }
    fmt.Println(m["message"])
}
```

Try running the program again and see if you can get the message!

```bash
go run main.go
```

```Let's get this started!```

## Next Steps

You're already well on your way to building your first event-driven microservice with Ensign!

If you're ready to see some more advanced examples with code, check out the [End-to-end Examples]({{< relref "examples">}}).

If you're looking for more on the basics of event-driven systems, check out [Eventing 101]({{< relref "eventing">}}).

Happy eventing!