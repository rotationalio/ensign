---
title: "End-to-End Examples"
weight: 30
date: 2023-05-17T17:03:41-04:00
---

This section of the documentation provides end-to-end examples using Ensign to help get you started!

- [Ensign for Data Engineers]({{< ref "/examples/data_engineers" >}}): This end-to-end example demonstrates how to create a live pub-sub data workflow with Ensign. Create a publisher to call a [Weather API](https://www.weatherapi.com) and emit the data to a topic stream, then read the stream in real time with a subscriber.

- [Ensign for Data Scientists]({{< ref "/examples/data_scientists" >}}): What does event-driven data science look like? In this example, see how to create an Ensign subscriber to [Baleen](https://github.com/rotationalio/baleen), a live RSS ingestion engine, and use the incoming data to perform streaming HTML parsing, entity extraction, and sentiment analysis.