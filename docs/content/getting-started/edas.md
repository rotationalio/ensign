---
title: "Designing Event-Driven Architectures"
weight: 20
date: 2023-05-23T10:51:25-04:00
description: "Designing your first EDA"
---


An event-driven architecture (EDA) is a plan for how data will flow through your application. It can be helpful to decompose these architectures into different handlers that are responsible for performing operations on the data (e.g. ingestion, statistical inference, prediction) and routing it between the layers of your application via the topics.

{{< image src="img/sample-eda.png" alt="Producers and consumers route data between the layers of a sample application from ingestion to analytics to the Web UI." zoomable="true" >}}

Handlers usually fall into two categories, "Producers" and "Consumers". Producers are responsible for writing data to topics, while consumers read data from those topics and perform some transformation on it (e.g. feature extraction, normalization, standardization, de-noising, model training). Some layers of your application may include both a producer and a consumer, or even multiple consumers and producers!

In the figure above, we have an architecture for a lightweight Python web-based application that uses raw data from a streaming weather API, trains an online model, predicts the weather for tomorrow, and displays it alongside a timeseries plot of the last two weeks of weather reports.