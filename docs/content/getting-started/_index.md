---
title: "Quickstart"
weight: 5
date: 2023-05-17T17:03:41-04:00
description: "Getting Started with Ensign"
---

The first step is to get an Ensign API key by visiting [the sign-up page](https://rotational.app/register). Similar to getting a developer API key for [Youtube](https://developers.google.com/youtube/v3/getting-started) or [Data.gov](https://api.data.gov/docs/api-key/), you will need an API key to use Ensign and to follow along with the rest of this Quickstart guide.



**Step 1. Create a New Ensign Project**
- After creating an account, click the `Create` button under `Set Up a New Project`.
- Enter a `Project Name` (required) and provide an optional description, then click `Create Project` (You will be redirected to the Projects dashboard).
- Watch the video tutorial at [Ensign Project](https://www.youtube.com/watch?v=VskNgAVMORQ)

**Step 2. Create Ensign Topics**
- On the Project dashboard, click on your Project Name.
- Go to the `Design Your Data Flows: Set up Your Topics` section to create topics. 
- Click `New Topic` and enter a Topic Name (required), then click `Create Topic`. 
- Watch the video tutorial at [Ensign Topics](https://www.youtube.com/watch?v=1XuVPl_Ki4U)
> **Tips:**
> - Pick a name for each topic that encodes information about the data that's stored in that topic.
> - You can also include information about the data format (e.g. flights-json or hotels-html). Read more at [Naming Topics]({{< ref "/getting-started/topics" >}})
> - Topics can help serve as a data contract helping upstream publishers remember how to serialize the data and downstream subscribers how to parse it.

**Step 3. Create Ensign API Keys**
- Go to the `Permission Your Data Flows: Generate API Keys` section and click `New Key`.
- A `Generate API Key for project` pop-up will open, add a `Key Name,` and choose `Permissions` (Full Access or Custom Access), then click `Generate API Key`.
- From the `Your API Key` popup, copy both the `Client ID` and the `Client Secret` to open a client connection to Ensign. Read More at [Ensign API Keys](#ensign-api-keys) 
- Watch the video tutorial at [Creating Ensign API Keys](https://www.youtube.com/watch?v=KMejrUIouMw)
> **Important:**
> - Remember to download your keys and save this file.
> - Add the `Client ID` and `Client Secret` to your environment variables, or you can put the file into your Git repository via gitignore. Read More at [Authentication](#authentication)
> - Do not share your keys with anyone, and never commit them to a public GitHub repository.
> - If your keys get lost or compromised, don't worry, you can revoke them and create new ones.
> - API keys are defined at the project level, and everyone collaborating on the datasets will need their own API keys.
> - Once you add teammates to your project, they can create their own API keys using the same technique.
> - New API keys enable users or applications to publish new data to the project topics, or to subscribe to those topics.



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
