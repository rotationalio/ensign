---
title: "End-to-End Examples"
weight: 20
bookFlatSection: false
bookToc: true
bookHidden: false
bookCollapseSection: false
bookSearchExclude: false
---

# End-to-End Examples

This section of the documentation provides end-to-end examples using Ensign to help get you started!

- [Ensign for Application Developers]({{< ref "/examples/developers" >}}): In this end-to-end example, see how to curate a custom Twitter feed using Ensign. Create a publisher to start emitting tweets to a topic stream in just a few minutes! See how to create one or more asynchronous subscribers that can read off the stream and process the data as needed.

- [Ensign for Data Engineers]({{< ref "/examples/data_engineers" >}}): This end-to-end example demonstrates how to retrieve and save weather data using Ensign and [Watermill](https://watermill.io). Create a publisher to call the [Weather API](https://www.weatherapi.com) and emit the data to a topic stream and use Watermill's [router](https://watermill.io/docs/messages-router/) and [SQL Pub/Sub](https://watermill.io/pubsubs/sql/) to save the data into a PostgreSQL database.

- [Ensign for Data Scientists]({{< ref "/examples/data_scientists" >}}): What does event-driven data science look like? In this example, see how to create an Ensign subscriber to Baleen, a live RSS ingestion engine, and use the incoming data to perform streaming HTML parsing, entity extraction, and sentiment analysis.