---
title: "Glossary"
weight: 30
date: 2023-05-17T17:03:41-04:00
---

When you're learning a new technology, there always a LOT of new lingo. We've tried to gather them all together here to help get you started:

#### **api key** <a name="api-key"></a>
"API" stands for "Application Programming Interface", which is a very broad term that refers (super high level) to the ways in which users or other applications can interact with an application.

Some applications (like Ensign) require permission to interact with, such as a password, token, or key.

You can get a free Ensign API key by visiting rotational.io/ensign. Your key will consist of two parts, a `ClientID` and a `ClientSecret`. The `ClientID` uniquely identifies you, and the `ClientSecret` proves that you have permission to create and access event data. You will need to pass both of these in to create an Ensign [client](#client) connection.

#### **asynchronous** <a name="asynchronous"></a>
An asynchronous microservice is one in which requests to a service and the subsequent responses are decoupled and can occur independently of each other.

This differs from the **synchronous** pattern, in which a client request (e.g. a query) is blocked from moving forward until a server response is received. Synchronous microservices can result in cascading failures and compounding latencies in applications.

Asynchronous microservices can make it a lot easier for teams to develop and deploy components independently.

Asynchronous microservices require an intermediary service usually known as a [broker](#broker) to hold messages emitted by a publisher that are awaiting retrieval from subscribers.

#### **broker** <a name="broker"></a>
An event broker is an intermediary service inside an asynchronous eventing system that stores events sent by publishers until they are received by all subscribers.

Brokers are also in charge of things like keeping events in the correct order, remembering which subscribers are listening to a topic stream, recording the last message each subscriber retrieved, etc.

In Ensign, brokers can save events permanently even after they have been retrieved (to support "time travel" &mdash; the ability to retroactively scan through an event stream to support analytics and machine learning).

#### **client** <a name="client"></a>
In order to write or read data from an underlying data system (like a database or event stream), you need a client to connect to the data system and interact with it as needed (such as reading and writing data). This connection often looks something like `conn = DBConnection(credentials)`, and after creating the `conn` variable, subsequent lines of code can leverage it to perform the kinds of data interactions you wish to make.

To establish a client in Ensign you need an [API key](#api-key).
If you add your `ClientID` and `ClientSecret` credentials to your bash profile, you can do the following and Ensign will read the credentials from your environment variables.

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

#### **event** <a name="events"></a>
In an event-driven or [microservice](#microservice) architecture, an event is the atomic element of data.

An event might look something like a dictionary, which is then wrapped in an object or struct that provides some schema information to help Ensign know how to serialize and deserialize your data.

```golang
order := make(map[string]string)
order["item"] = "large mushroom pizza"
order["customer_id"] = "984445"
order["customer_name"] = "Enson J. Otterton"
order["timestamp"] = time.Now().String()

evt := &ensign.Event{
    Mimetype: mimetype.ApplicationJSON,
    Type: &api.Type{
        Name:    "Generic",
		MajorVersion: 1,
		MinorVersion: 0,
		PatchVersion: 0,
    },
}

evt.Data, _ = json.Marshal(order)
```

#### **latency** <a name="latency"></a>
Latency can refer to both application-level communication lag (e.g. the time it takes for one part of the code to finish running before moving on to the next part) or to network communication lag (e.g. the time it takes for two remote servers on two different continents to send a single message back and forth).

Less latency is better, generally speaking.

In a microservices context, we can reduce application latency by using asynchronous communications and parallelizing functions so they run faster. Network latency can be reduced by creating more efficient communications between servers (e.g. using more scalable consensus algorithms).


#### **microservice** <a name="microservice"></a>
A microservice is a computer application composed of a collection of lightweight services, each of which is responsible for some discrete task.

Microservices can be coordinated to communicate via [events](#events).


#### **mime type** <a name="mimetype"></a>
A MIME (Multipurpose Internet Mail Extensions) type is a label that identifies a type of data, such as CSV, HTML, JSON, or protocol buffer.

MIME types allow an application to understand how to handle incoming and outgoing data.


#### **organization**
An Ensign organization is a collection of users who are working under the same Ensign [tenant](#tenant).


#### **publisher** <a name="publisher"></a>
In an event-driven microservice, a publisher is responsible for emitting [events](#events) to a [topic stream](#topic).

In Ensign, you can create a publisher once you have established a [client](#client). On publish, the client checks to see if it has an open publish stream created and if it doesn't, it opens a stream to the correct Ensign node.

```golang
client.Publish(yourTopic, yourEvent)
```

#### **real-time**
This is a tricky one because real-time can be used to mean different things. In some cases, "real-time" is used as a synonym for synchronous (i.e. the opposite of [asynchronous](#asynchronous)). However, the term is also used to mean "very fast" or ["low latency"](#latency).


#### **sdk**
SDK stands for "Software Development Kit". Software applications designed for a technical/developer audience frequently are considerate enough to provide user-facing SDKs in a few languages (e.g. Golang, Python, JavaScript). These SDKs give users a convenient way to interact with the application using a programming language with which they are familiar.

Ensign currently offers two SDKs: the [Golang SDK](https://github.com/rotationalio/ensign/blob/main/sdks/go/ensign.go) and a [Watermill API-compatible SDK](https://github.com/rotationalio/watermill-ensign/tree/main/pkg/ensign).

#### **stream**
An event stream is a flow composed of many, many individual pieces of data called [events](#events).

#### **subscriber**
In an event-driven context, a subscriber is a downstream component that is listening for [events](#events) published by a [publisher](#publisher) onto a [topic](#topic).

#### **tenant** <a name="tenant"></a>
A tenant is a user, group of users, team or company who share computing and/or storage resources.

#### **topic**
In event-driven microservices, a topic is a rough approximation of a traditional relational database table. In a relational DB, a table is a collection of related data fields arrayed as columns and rows. In an eventing context, a topic is a sequence of individual [events](#events) populated with the same fields (aka schema).
