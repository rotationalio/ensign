---
title: "Quickstart"
weight: 5
date: 2023-05-17T17:03:41-04:00
description: "Getting Started with Ensign"
---

The first step is to get an Ensign API key by visiting [the sign-up page](https://rotational.app/register). Similar to getting a developer API key for [Youtube](https://developers.google.com/youtube/v3/getting-started) or [Data.gov](https://api.data.gov/docs/api-key/), you will need an API key to use Ensign and to follow along with the rest of this Quickstart guide.



**Step 1. Create a New Ensign Project**
- After creating an account, click the `Create` button under `Set Up a New Project`.
- Enter a `Project Name` (required) and provide an optional description, then click `Create Project`. You will be redirected to the Projects dashboard.
- For more details, watch the [Project Creation Tutorial](https://www.youtube.com/watch?v=VskNgAVMORQ)

**Step 2. Create Ensign Topics**
- On the Project dashboard, click on your Project Name.
- Click `New Topic` and enter a Topic Name (required), then click `Create Topic`. 
- For more details, watch the [Topic Creation Tutorial](https://www.youtube.com/watch?v=1XuVPl_Ki4U)
> **Tips:**
> - Pick a name for each topic that encodes information about the data you intend to store in that topic.
> - You can also include information about the data format (e.g. flights-json or hotels-html). Read more [here]({{< ref "/getting-started/topics" >}})
> - Topics can help serve as a data contract helping upstream publishers remember how to serialize the data and downstream subscribers how to parse it.

**Step 3. Create Ensign API Keys**
- Click the `New Key` button.
- Create a name for the key and select its permissions (Full Access or Custom Access), then click `Generate API Key`.
- From the `Your API Key` popup, copy both the `Client ID` and the `Client Secret` to open a client connection to Ensign. Read More at [Ensign API Keys](#ensign-api-keys) 
- Watch the video tutorial at [Creating Ensign API Keys](https://www.youtube.com/watch?v=KMejrUIouMw)
> **Important:**
> - Remember to download your keys and save this file.
> - Do not share your keys with anyone, and never commit them to a public GitHub repository. Read more about best practices for API keys [here](#authentication).
> - If your keys get lost or compromised, don't worry, you can revoke them and create new ones.
> - API keys grant access to all topics in the project.
> - Team members with the `Member`, `Owner`, or `Admin` role can create API keys, an `Observer` can not create API keys.



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
