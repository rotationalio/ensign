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

Generally speaking, developers can use an event stream to build event-driven applications, rich, interactive front-end experiences, personalized experiences, and more.

Data scientists can try their hand at online learning and prediction..

Data engineers can streamline the ETL process and securely provision data to the process where it’s valued &mdash; and needed &mdash; most.

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