---
title: "Quickstart"
weight: 10
date: 2023-05-17T17:03:41-04:00
description: "Let's gooooooo!"
---

The first step is to get an Ensign API key by visiting [the sign-up page](https://rotational.app/register). Similar to getting a developer API key for [Youtube](https://developers.google.com/youtube/v3/getting-started) or [Data.gov](https://api.data.gov/docs/api-key/), you will need an API key to use Ensign and to follow along with the rest of this Quickstart guide.

<a name="ensign-keys"></a>
### Ensign API Keys

An API key consists of two parts, a `ClientID` and a `ClientSecret`. The `ClientID` uniquely identifies a project, and the `ClientSecret` allows you to create and access event data within that project.

| API Key Component Name | Length | Characters | Example |
|:------:|:------:|:------:|:------:|
| ClientID          | 32     | alphabetic (no digits) | DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa                                 |
| ClientSecret      | 64     | alphanumeric           | wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS |


Together, the `ClientID` and `ClientSecret` provide access to a project. They enable you to create Ensign topics, publishers, and subscribers, which will be the building blocks of your microservice! Anybody with the `ClientID` and `ClientSecret` has access to your project, so these values should be kept private and not shared.

### SDKs

The SDKs are where you use your API key and allow you to integrate Ensign into your projects. There are currently two SDKs supported.

- [PyEnsign](https://github.com/rotationalio/pyensign) is the official SDK for Python
- [go-ensign](https://github.com/rotationalio/go-ensign) is the official SDK for Golang

### Authentication

Using the SDKs requires authenticating with your API key. By default the SDKs will read the Client ID and Client Secret from your environment.

```bash
export ENSIGN_CLIENT_ID=DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa
export ENSIGN_CLIENT_SECRET=wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS
```

### Publish and Subscribe

The most common use of the SDKs is to publish events to topics, and subscribe to events from somewhere else. Topics can be created from your [dashboard](https://rotational.app) or directly from the SDKs. See the SDK guides for a quick example of how to get started.

- [Python SDK]({{< ref "/sdk/python" >}})
- [Go SDK]({{< ref "/sdk/golang" >}})
