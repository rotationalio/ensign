---
title: "FAQ"
weight: 30
bookFlatSection: false
bookToc: true
bookHidden: false
bookCollapseSection: false
bookSearchExclude: false
---

## What is an event (or data) stream?
You can think of an event stream as a pipeline through which data flows, much like pipes that move water in your home. You can use an event stream to connect data sources that produce or emit data (events) to data sinks that consume or process the data. A data source or producer that emits data can be a database, a data warehouse, a data lake, machines, edge devices or sensors, microservices, or applications. A data sink or consumer can be a web or mobile application, a machine learning model, a custom dashboard, or other downstream microservices, sensors, machines, or devices like wearables.

With Ensign, you can set up any number of secure event streams &mdash; also called [topics]({{< ref "/glossary#topic" >}}) &mdash; for a project or use case. This allows you to customize how to move data in your application or organization, to where it is most useful. You don’t need specialized skills, tools, or infrastructure. Just an [API key]({{< ref "/glossary#api-key" >}}) and a few lines of code.

## What can I use an event stream for?
Event streams are applicable to almost any use case that benefits from the continuous flow of data.

Generally speaking, developers can use an event stream to build event-driven applications. Think about how satisfying it is as a customer when you see your data update in real time &mdash; whether it's being able to watch your package move closer and closer to being delivered, or your fitness tracker updating throughout your cardio session. These kind of rich, interactive, and personalized experiences depend on an underlying stream of data flowing in to the application. Check out our Ensign resources for developers [here]({{< ref "/examples/developers" >}}).

Data scientists can also find substantial value in event streams; in fact, some of the [earliest deployed machine learning apps](https://en.wikipedia.org/wiki/Naive_Bayes_spam_filtering) were event-driven! That’s right — back in the 90’s, email spam filters used Bayesian models to learn on the fly. Imagine how much more robust our models would be if they were not only trained on the freshest data, but could alert you to things like data drift as soon as they happened &mdash; you’d be able to react immediately as opposed to a batchwise process where you’d be lucky to catch the issue within a day! Check out our Ensign resources for data scientists [here]({{< ref "/examples/data_scientists" >}}).

Event streams can also play a huge role in process automation for data engineers and database administrators. Eventing can streamline ETL processes and enable database maintenance and scaling without the danger of data loss or application downtime. You can even use eventing to provision data more securely, updating event encryption to turn on and off access for downstream processes and users. Check out our Ensign resources for data engineers [here]({{< ref "/examples/data_engineers" >}}).

## What is change data capture (CDC)?
Databases are constantly changing, particularly transactional databases. However, databases must also support the underlying consistency models of the applications and organizations they support.

For example, when you look at your checking account balance, you expect to get a single amount, even though in reality that amount is continually going up and down (hopefully more up than down :-D). And let's say at the same time you were checking your account balance, your account manager also checks your balance &mdash; you would expect to both see the same value, even if you were checking from your phone while on vacation in Tokyo and your account manager is checking from her desk in Nebraska.

But it's tough for a database to be able to provide that level of consistency while also providing detailed insights about the up-and-down movement of your account balance. That's where change data capture comes in!

Ensign's event streams can provide bolt-on change data capture (CDC) by logging all changes to a database over time. In our bank example, an Ensign CDC log could be used for a lot of useful things &mdash; from training models to forecast negative balances and flag account holders before they incur fines, to real time alerting to protect customers from fraud.

## What does it mean to persist data?
Data persistence is a mechanism for saving data in a way that will protect it from failures like power outages, earthquakes, server crashes, etc.

Think about how stressful it would be if you were depositing money at an ATM, and the ATM screen shorted out and went black before your bank had a chance to update your account balance in their database! Or what if you were buying concert tickets online, and the ticket sales website crashed after your payment went through but before your Taylor Swift tickets were issued!

Considering how important data persistence is to fostering trust in consumer-facing contexts, you might be surprised to learn that most streaming tools don't provide data persistence! Some only save data for a few days before it is discarded, some must be specially configured to save data, and others do not have an ability to save data at all.

Ensign, on the other hand, persists all data by default &mdash; because it's better to be safe than sorry!

## In what way are my event streams and data secure?
We designed Ensign with security-first principles. All events are encrypted in motion and at rest.

To start, Ensign employs high performance symmetric key cryptography such as AES-GCM for encrypting events at rest. Your API keys contain a ClientID that uniquely identifies you and ClientSecret that proves you have permission to create and access event data. Your API keys are shown to you and only you, once and only once.

At the same time, Ensign never stores raw passwords to access the developer portal. We use the [Argon2](https://en.wikipedia.org/wiki/Argon2) key derivation algorithm to store passwords and API key secrets as hashes that add computation and memory requirements for validation.

Finally, for enterprise clients, we’re working on deploying Ensign in virtual private clouds (VPCs).

## What is the cost?
We are working on figuring that out. For now, there is no cost because we want developers and data scientists to experiment and build what they can imagine with eventing or streaming superpowers ;-D

Eventually, we will have to charge to be sustainable, but we do commit to always having a free tier for experimenters and builders.

## What can I build?
Glad you asked! Check out some ideas [here]({{< relref "examples">}}).