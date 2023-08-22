---
title: "Event Brokering"
weight: 3
date: 2023-05-17T17:03:41-04:00
---

At the core of Ensign is the **Broker**.

Event brokers are what differentiate eventing systems (e.g. Kafka, Pulsar, Redpanda, Google PubSub, Ensign) from synchronous messaging systems (e.g. RabbitMQ, Ably, Amazon SQS).

## What Does a Broker Do?

Brokers are responsible for a lot.

{{< image src="img/broker.png" alt="Brokers are responsible for persisting data, ordering events, and remembering subscriber offsets." zoomable="true" >}}

Here are some of the main responsibilities of an Ensign broker:

1. **Persisting Data**

The Ensign Broker persists data written to any topic by the Publisher so that multiple Subscribers can read data from that same topic. Persisting data also means that Ensign has geodistributed backups of data, so it's safe from things like earthquakes, floods, and other catastrophes.

In a synchronous messaging system like RabbitMQ, messages are discarded after they are acknowledged, so there is no expectation of persistence whatsoever.

While tools like Kafka, Pulsar, and Google PubSub can be configured by editing YAML files to persist data indefinitely, enabling persistence has two consequences. First, it significantly reduces throughput and increases latency, because it requires the broker to do more work than usual. Enabling persistence by default also increases costs in these other systems, because storing data costs money.

2. **Keeping Events in the Right Order**

The Ensign Broker keeps all events in order for every topic, even if there are publishers writing to that topic from two different sides of the Earth. Knowing that your events will always be totally ordered provides a powerful semantic that teams can use to design applications that will always respect the same ordering of events.

Ensign events are ordered by the Broker using an RLID. An RLID is a totally ordered, 80-byte data structure that encodes both time and a monotonically increasing sequence number using [Crockford's base32](https://www.crockford.com/base32.html). It is inspired by ULID and Snowflake IDs.

To date, Ensign is the only eventing system in the world that guarantees totally ordered concurrent events within topics.

3. **Remembering the Offsets**

The Ensign Broker remembers subscriber offsets. This means that an application reading data from Ensign does not have to maintain any state related to Ensign &mdash; this is why Ensign is particularly convenient for teams that deploy stateless applications.

Consider two subscribers connecting to their Ensign Broker, Subscriber A that is connecting for the first time, and Subscriber B that has been dormant/inactive for some period of time. Presumably, Subscriber A and Subscriber B will want to start reading the data from different points in the topic stream; Subscriber A might want to read in all the events from the very beginning, while Subscriber B might prefer to start back from where it left off so that it can recover its function using the minimum needed computation or memory. Ensign's Broker maintains a mapping between topics and their subscribers.

## What Makes the Ensign Broker Unique?

There are several things that are unique in the implementation of the Ensign Broker that diverge from similar systems. Here are a few differentiators:

- **In Ensign, Data is persisted by default**: In most similar systems, persistence is disabled by default, must be activated using configuration/YAML, and will instantly and significantly increase the costs.
- **The Ensign Broker stores the Subscriber offsets**: This means Ensign Subscribers can be fully stateless.
- **Consensus is flexible**: Similar systems use Raft or Zookeeper for consensus. Ensign is built so that decision-making in the Broker can automatically shift into a hierarchical consensus mode. [Hierarchical consensus](https://dl.acm.org/doi/10.1145/3087801.3087853) enables small local consensus groups to function independently when possible, enabling far more geographic scaling than would be possible in any other system.

In combination, these differences mean that Ensign fulfills the criteria of both a database and an eventing platform, and is far safer for data (aka more fault tolerant) than any comparable eventing system.