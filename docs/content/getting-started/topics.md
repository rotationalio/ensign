---
title: "Naming Topics"
weight: 30
date: 2023-05-23T11:17:31-04:00
description: "Best practices for naming your topics"
---

What's the best way to name your topics? This is an **excellent** question!

Remember, a topic is just like a table in a traditional relational database, so it can be helpful to think about that as you name them.

Each team may have slightly different naming conventions, and the most important thing when it comes to naming is that your teammates understand the names you use!

That said, our favorite technique is to give each data source and type its own topic, for instance:

**`user_logins_plaintext`**:
We might expect this topic to contain data about user logins that could be stored as plaintext, meaning it doesn't contain any publicly identifiable information (PII).

**`product_reviews_xml`**:
Here the topic likely contains multi-field product reviews that might include text content, numeric ratings (e.g. stars), etc., stored as XML.

**`weather_reports_json`**:
With this topic, you could expect the data to be weather reports formatted as JSON data.

**`model_results_pickle`**:
This topic might contain machine learning models that have been trained and serialized in the [Python pickle](https://docs.python.org/3/library/pickle.html) format.

Adding the type at the end of the topic name might not always be necessary, but it can be a very helpful way for the Publishers to communicate to the Subscribers what the [MIME type]({{< ref "/eventing/glossary#mimetype" >}}) of the data will be.

Think about how it will feel to query the topic later (which you can do with [EnSQL]({{< ref "/ensql" >}})!), e.g.

```sql
SELECT DISTINCT reviewer_id FROM product_reviews_xml WHERE timestamp < 2022-05-11T13:00:00-04:00
```
