---
title: "FAQ"
weight: 120
date: 2023-08-21T17:09:36-04:00
---

Got a question you don't see here? Let us know so we can add it! support@rotational.io

<!--more-->

### How do I sign up for an Ensign account?
We are so glad you are interested in checking out Ensign!  You can register [here](https://rotational.app/register) and check out the [quickstart guide]({{< relref "getting-started">}}) to get started!

### What are some good use cases for Ensign?
Ensign can be used for many different types of real-time applications.  We put together this [guide]({{< ref "/eventing/use_cases" >}}) to help you come up with ideas for your next project.

### What can I build - how do I get started?
Glad you asked! Check out some ideas [here]({{< relref "examples">}}).

### In what agencies or companies is this kind of work relevant? Where can I work if I want to do projects like this?
Ensign serves a broad set of use cases, some of which we have outlined in this [guide]({{< ref "/eventing/use_cases" >}}).  As you can see, ensign can be used by any agency or company that is looking to leverage its data.

### Can you share any data flow diagrams as a starting point?
Sure thing!  Check out this [link]({{< ref "/getting-started/edas" >}}) that contains a data flow diagram and a primer on how to design event-driven architectures for building streaming applications.

### What is EnSQL and how do I use it?
EnSQL is a lightweight structured query language that helps you query your topics just like you would any table.  Moreover, you have the added benefit of being able to retrieve changes to data over time rather than being limited to the current version of the data.  This [reference]({{< relref "ensql">}}) provides an in-depth tutorial on how to use EnSQL.

### Does PyEnsign integrate with leading streaming platforms (e.g. Apache Kafka, AWS Kinesis etc.)?
Absolutely!  Since Ensign is also a streaming engine, you can use the Pub/Sub framework to interact with Kafka and Kinesis.  For example, you can create an Ensign subscriber that subscribes to a Kafka topic or ingests from a Kinesis stream and does some data transformations.  On the flip side, you can have an Ensign publisher that publishes data to a Kafka topic or Kinesis stream.  Ensign is cloud agnostic and therefore it can be used in any environment with no additional infrastructure requirements.  All you have to do is [create an Ensign account](https://rotational.app/register) and `pip install pyensign` and write your application code.

### Why use Ensign for model topics instead of using something like Artifactory since you are versioning binaries?
If you use Ensign, you eliminate the need for a different tool, and you don't need to spend additional time learning and configuring another component in your environment.  The additional benefits are reduced latency because the model and the data are within the ensign ecosystem and easier troubleshooting when things go wrong due to a less complex infrastructure.  

### What is the cost?
We are working on figuring that out. For now, there is no cost because we want developers and data scientists to experiment and build what they can imagine with eventing or streaming superpowers ;-D

Eventually, we will have to charge to be sustainable, but we do commit to always having a free tier for experimenters and builders.

### What is an event (or data) stream?
You can think of an event stream as a pipeline through which data flows, much like pipes that move water in your home. You can use an event stream to connect data sources that produce or emit data (events) to data sinks that consume or process the data. A data source or producer that emits data can be a database, a data warehouse, a data lake, machines, edge devices or sensors, microservices, or applications. A data sink or consumer can be a web or mobile application, a machine learning model, a custom dashboard, or other downstream microservices, sensors, machines, or devices like wearables.

With Ensign, you can set up any number of secure event streams &mdash; also called [topics]({{< ref "/glossary#topic" >}}) &mdash; for a project or use case. This allows you to customize how to move data in your application or organization, to where it is most useful. You don’t need specialized skills, tools, or infrastructure. Just an [API key]({{< ref "/glossary#api-key" >}}) and a few lines of code.

### What are Ensign "projects"?
How does your organization group its datasets? In traditional storage, related data is stored in a `database`. An Ensign `project` defines a collection of datasets related by use case.

### What are Ensign "topics"?
How does your organization normalize data? Traditional relational database management systems (RDBMS) break data down into `tables` that describe objects. Ensign `topics` are the same, but also capture all changes to every object (in chronological order!).

### What should I name my topics? Are there good naming conventions?
First and foremost, consider your teammates when you name your topics.  Come up with names that facilitate communication, not cause confusion.  Check out this handy [resource]({{< ref "/getting-started/topics" >}}) we put together to help you.

### What are Ensign "workspaces"?
Who has permission to connect to the database? In Ensign, you create a `workspace` to grant roles and adjust permissions for each teammate. You can have a workspace of one, but it's not as fun 😉

### What are Ensign "tenants"?
Tenant is primarily an accounting term; an organization can have multiple tenants. Tenants are particularly useful for setting up different environments, such as Dev/Prod/Staging.

### What can I use an event stream for?
Event streams are applicable to almost any use case that benefits from the continuous flow of data.

Generally speaking, developers can use an event stream to build event-driven applications. Think about how satisfying it is as a customer when you see your data update in real time &mdash; whether it's being able to watch your package move closer and closer to being delivered, or your fitness tracker updating throughout your cardio session. These kind of rich, interactive, and personalized experiences depend on an underlying stream of data flowing in to the application.

Data scientists can also find substantial value in event streams; in fact, some of the [earliest deployed machine learning apps](https://en.wikipedia.org/wiki/Naive_Bayes_spam_filtering) were event-driven! That’s right — back in the 90’s, email spam filters used Bayesian models to learn on the fly. Imagine how much more robust our models would be if they were not only trained on the freshest data, but could alert you to things like data drift as soon as they happened &mdash; you’d be able to react immediately as opposed to a batchwise process where you’d be lucky to catch the issue within a day! Check out our Ensign resources for data scientists [here]({{< ref "/examples/data_scientists" >}}).

Event streams can also play a huge role in process automation for data engineers and database administrators. Eventing can streamline ETL processes and enable database maintenance and scaling without the danger of data loss or application downtime. You can even use eventing to provision data more securely, updating event encryption to turn on and off access for downstream processes and users. Check out our Ensign resources for data engineers [here]({{< ref "/examples/data_engineers" >}}).

### What is change data capture (CDC)?
Databases are constantly changing, particularly transactional databases. However, databases must also support the underlying consistency models of the applications and organizations they support.

For example, when you look at your checking account balance, you expect to get a single amount, even though in reality that amount is continually going up and down (hopefully more up than down :-D). And let's say at the same time you were checking your account balance, your account manager also checks your balance &mdash; you would expect to both see the same value, even if you were checking from your phone while on vacation in Tokyo and your account manager is checking from her desk in Nebraska.

But it's tough for a database to be able to provide that level of consistency while also providing detailed insights about the up-and-down movement of your account balance. That's where change data capture comes in!

Ensign's event streams can provide bolt-on change data capture (CDC) by logging all changes to a database over time. In our bank example, an Ensign CDC log could be used for a lot of useful things &mdash; from training models to forecast negative balances and flag account holders before they incur fines, to real time alerting to protect customers from fraud.

### What does it mean to persist data?
Data persistence is a mechanism for saving data in a way that will protect it from failures like power outages, earthquakes, server crashes, etc.

Think about how stressful it would be if you were depositing money at an ATM, and the ATM screen shorted out and went black before your bank had a chance to update your account balance in their database! Or what if you were buying concert tickets online, and the ticket sales website crashed after your payment went through but before your Taylor Swift tickets were issued!

Considering how important data persistence is to fostering trust in consumer-facing contexts, you might be surprised to learn that most streaming tools don't provide data persistence! Some only save data for a few days before it is discarded, some must be specially configured to save data, and others do not have an ability to save data at all.

Ensign, on the other hand, persists all data by default &mdash; because it's better to be safe than sorry!

### In what way are my event streams and data secure?
We designed Ensign with security-first principles, including Transparent Data Encryption (TDE), symmetric key cryptography, key rotation, JWT-based authentication and authorization, and Argon2 hashing. All events are encrypted by default, and customers are empowered to control how cryptography occurs for their own data within the Ensign system.

When you set up your Ensign tenant, you're issued server-side global keys. Your API keys are shown to you and only you, once and only once. These keys are stored in a globally replicated key management system and you can revoke them using the management UI. Once you revoke keys, data encrypted with those keys is no longer accessible.

Ensign never stores raw passwords to access the developer portal. We use the [Argon2](https://en.wikipedia.org/wiki/Argon2) key derivation algorithm to store passwords and API key secrets as hashes that add computation and memory requirements for validation.

For enterprise clients, we’re working on deploying Ensign in virtual private clouds (VPCs).

### Where does the data go/live?
Ensign is a hosted solution and so your data is stored securely on our servers and will not get deleted unless you instruct us to do so.  For enterprise clients, we will deploy the storage solution within virtual private clouds (VPCs).

### How do you deploy models?
In a streaming application, you can deploy a model by publishing it to a model topic from which a subscriber can retrieve the model and use it for generating predictions.  

For example, we can create a publisher that will periodically send requests to an API to receive data to train a classification model.  Once we are satisfied with the model's performance and choose to deploy, we can pickle the model and the performance metrics and then publish the event to a `model` topic.  We can then create a subscriber on the other end that is listening for new messages on this topic as well as a `score_data` topic that sends new data that will be used to generate predictions.  Once the subscriber receives the message, it will extract the model from the event.  Now, it is ready to use this model to generate productions on data it receives from the `score_data` topic.  

The added benefit of this architecture is that the model topic serves as a model registry that can be used to version control models and because we are storing metrics, it is easy to evaluate different model versions.

Check out this [repository](https://github.com/rotationalio/online-model-examples/tree/main/other_examples) for an implementation of this design pattern.

### What data transmission considerations are there as you move to production with Ensign?
Ensign is a streaming engine, and therefore, events cannot be larger than 5 megabytes.  If you have data that is larger, you can compress it and/or break it down before publishing and provide a way for the subscriber on the other end to put it together for processing.  If you are attempting to publish a model that is simply too large to be streamed, then one option is to store it on AWS S3 or another storage solution and include the link to the location in the event that is published to the model topic.  A similar strategy can be employed for other types of data.

Also, ensign is designed to persist data by default, and so the data won't be deleted unless you tell us to do so.  

### What are protocol buffers and why does Ensign use them?
Eventing is all about machines talking to each other and protocol buffers help machines quickly send and receive messages to and from each another.  Potocol buffers are described in detail at the bottom of this [page]({{< relref "sdk">}}).

### Why am I not seeing any events coming to my subscriber?
In order for an event to be seen by a subscriber, two things need to happen.  First, the ensign server needs to have successfully received the event from the ensign client.  The server communicates the successful receipt by sending an `ack` message.  Second, the publisher needs to receive this `ack` message before publishing the event to the topic.  If there is a failure in any of the steps, the event won't be seen by the subscriber.  For this reason, it is recommended that users write code to handle acks and nacks (successful and unsuccessful receipt of messages sent by the ensign server to the ensign client).  

Also, ensure that the subscriber is stood up first to receive messages because it is possible for messages to have already been sent to a topic before the subscriber started listening to them.  

Furthermore, make sure that a sleep interval is set on the publisher because it is possible that there can be delays in receiving data from either the ensign server or the data source (e.g. an API).  If the sleep interval is not set, it is possible that the publisher terminates before an acknowledgement from the ensign server or before receiving data from the data source. In such cases, the event does not get added to the topic and so the subscriber will not see it.
