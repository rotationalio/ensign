---
title: "Topic Management"
weight: 80
date: 2023-11-01T15:27:04-04:00
draft: true
---

Topics are at the heart of Ensign.  You can organize and store all of your data in topics within your projects much like how you can do so currently within tables in a traditional relational database.  

The Ensign SDKs provide several features that you can use to manage your topics.  The examples use `pyensign`, the Python SDK.

A `Topic` object has the following attributes:
- **id**: Unique ulid value associated with the topic
- **name**: Name of the topic
- **project_id**: Id of the project (database) that the topic belongs to
- **event_offset_id**: Id of the offset associated with the last published event in the topic
- **events**: Number of events in the topic
- **duplicates**: Number of duplicate events in the topic
- **data_size_bytes**: Data size in bytes
- **types**: Types of events stored in the topic.  It is possible to have multiple types of events with different schemas and mimetypes.

You can retrieve details associated with all topics within your project by running the following (_insert the credentials associated with your project_):

```python
ensign = Ensign(client_id=YOUR_CLIENT_ID, client_secret=YOUR_CLIENT_SECRET)
await ensign.info()
```

#### Automated Deduplication

Starting with version `v0.12.0`, Ensign has the functionality to perform automated deduplication on topics!  If duplicate data is going to cause issues in your application, you no longer have to write custom code to handle it.  

There are a lot of ways to determine if two events are duplicates of each other. Ensign uses user-defined policies to figure out duplication. Policies you can set are:

- _Strict_: Two events with identical metadata, data, mimetype, and type (though provenance, region, publisher, encryption, etc. may differ).
- _Datagram_: Two events with identical data regardless of metadata, mimetype, or type information
- _Key Grouped_: Two events with identical data and the same value for a user specified key or keys in the metadata
- _Unique Key Constraint_: Two events with the same value for specified key or keys in the metadata (unique key index)
- _Unique Field Constraint_*: Two events with the same value for a field or fields in the data (unique field index)
- _None_: We store all events no matter what

*Note that the unique field constraint requires us to be able to process your data -- which we won't be able to do until we have a schema registry. So although this is _technically_ a deduplication option, in practice it is not usable and will return not implemented errors. 

One quick caveat: deduplication happens as a background process right now, not in realtime. However in our next release, we will add real-time deduplication and data quality checks!

The default strategy for topics is _None_.  You can also specify the offset position from which you want to apply the deduplication strategy.  The options for offest position are as follows:

- _Earliest_: The default deduplication offset policy: the first record in the log is identified as the canonical event and all duplicate events point to this record. 

- _Latest_: An offset policy where earlier duplicates are marked as duplicates and the
latest event is identified as the canonical event. This policy slows down event processing, but allows queries to see duplicate values sooner.

The following is an example of how you can change the deduplication policy of your topic.  This example changes the strategy from _None_ to _Datagram_ and does not change the offset policy.  You will need the topic id to set the deduplication strategy.

```python
ensign = Ensign(client_id=YOUR_CLIENT_ID, client_secret=YOUR_CLIENT_SECRET)
await ensign.set_topic_deduplication_policy(topic_id, strategy=3)
```

It will take some time for the duplicates to be removed from the topic.  The topic will be in the `PENDING` status while the operation is taking place and will change to the `READY` status once the operation is complete.  In order to view the topic state, you can run the following command:

```python
ensign = Ensign(client_id=YOUR_CLIENT_ID, client_secret=YOUR_CLIENT_SECRET)
topics = await ensign.get_topics()
print(topics)
```

This prints a list of topics associated with the project and for each topic, you can see the following information: status, deduplication policy, created timestamp, and modified timestamp.