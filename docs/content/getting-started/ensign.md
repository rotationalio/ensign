---
title: "Quickstart"
weight: 10
date: 2023-05-17T17:03:41-04:00
description: "Let's gooooooo!"
---

The first step is to get an Ensign API key by visiting [the sign-up page](https://rotational.app/register). Similar to getting a developer API key for [Youtube](https://developers.google.com/youtube/v3/getting-started), [Twitter](https://developer.twitter.com/en/docs/twitter-api/getting-started/getting-access-to-the-twitter-api) or [Data.gov](https://api.data.gov/docs/api-key/), you will need an API key to use Ensign and to follow along with the rest of this Quickstart guide.

<a name="ensign-keys"></a>
### Ensign API Keys

Your key consists of two parts, a `ClientID` and a `ClientSecret`. The `ClientID` uniquely identifies you, and the `ClientSecret` proves that you have permission to create and access event data.

| API Key Component Name | Length | Characters | Example |
|:------:|:------:|:------:|:------:|
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
go get -u github.com/rotationalio/go-ensign@main
```

<a name="create-a-client"></a>
### Create a Client

After you've made a new Go project for this example, create a `main.go` file and add the dependencies you'll need, which will include importing the Ensign API, SDK, and mimetypes.

Next, create an Ensign client, which is similar to establishing a connection to a database like PostgreSQL or Mongo. To create the client, use the `New` method and pass in an `ensign.Options` struct that specifies your Client ID and Client Secret (described in the section above on getting an API key).

```golang
package main

import (
    "fmt"
    "time"
    "context"

	api "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	ensign "github.com/rotationalio/go-ensign"
)

const myClientId = "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa"
const myClientSecret = "wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS"

client, err := ensign.New(
    ensign.WithCredentials(
        myClientId, myClientSecret
    )
)
if err != nil {
	panic(fmt.Errorf("could not create client: %s", err))
}
```

Congratulations, you now have an open connection to Ensign!

*Note: You're probably thinking that it's not a great idea to store credentials in your files! You are so right about that. Instead, if you add your `ClientID` and `ClientSecret` credentials to your bash profile, you can do the following instead and Ensign will read the credentials from your environment variables.*

```golang
package main

import (
    "fmt"

	ensign "github.com/rotationalio/go-ensign"
)

client, err := ensign.New()
if err != nil {
	panic(fmt.Errorf("could not create client: %s", err))
}
```


### Make Some Data

Next, we need some data! Generally this is the place where you'd connect to your live data source (a database, Twitter feed, weather data, etc). But to keep things simple, we'll just create a single event, which starts with a map.

```golang
data := make(map[string]string)
data["sender"] = "Twyla"
data["timestamp"] = time.Now().String()
data["message"] = "Let's get this started!"
```

Next, we will convert our map into an event, which will allow you to specify the mimetype of the message you intend to send (in this case, we'll say it's JSON), and the event type (which will be a generic event for this example).

```golang
e := &ensign.Event{
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
if e.Data, err = json.Marshal(data); err != nil {
    panic("could not marshal data to JSON: " + err.Error())
}
```

### Publish Your Event

Now we can publish your event by calling the `Publish` method on the Ensign client we created above. You'll also need to pass in a `TopicId`, which will be a string. If you aren't sure what `TopicId` to use, you can quickly [log into your Ensign dashboard](https://rotational.app) and look it up.

But the truth is that it's hard for humans to remember ULID and you have enough on your plate already. So, you can also use the name of the topic instead of the id. For this example, we'll pretend we want to publish to a topic named `"quality-lemon-time"`.

On publish, the Ensign `client` checks to see if it has an open publish stream created for that topic, and if it doesn't it opens a stream to the correct Ensign node.

```golang
client.Publish("quality-lemon-time", e)
```

You can publish many events at a time if you want!

```golang
client.Publish("quality-lemon-time", e, a, f, h, q, p)
```

If you `Publish` to a second topic, the Ensign `client` will create another new Publisher for you!

```golang
client.Publish("surprisingly-mashed-potatoes", e)
```


### Create a Subscriber

So now you have a `Publisher` going; now we need to consume those events using a `Subscriber`

```golang
sub, err := client.Subscribe("quality-lemon-time") // topic alias also works
if err != nil {
    panic(fmt.Errorf("could not create subscriber: %s", err))
}

for msg := range sub.C {
    var m superSecretMessage
    if err := json.Unmarshal(msg.Data, &m); err != nil {
        panic(fmt.Errorf("failed to unmarshal message: %s", err))
    }
    fmt.Println(m.Message)
}
```

## Next Steps

You're already well on your way to building your first event-driven microservice with Ensign!

If you're ready to see some more advanced examples with code, check out the [End-to-end Examples]({{< relref "examples">}}).

If you're looking for more on the basics of event-driven systems, check out [Eventing 101]({{< relref "eventing">}}).

Happy eventing!