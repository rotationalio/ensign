---
title: "Python"
weight: 30
date: 2023-08-11T11:03:41-04:00
---

The Python SDK is the quickest way to get started with Ensign. In this example we'll create a simple Python project from scratch to publish and subscribe to Ensign!

### Prerequisites

- [create a free Ensign account and API key](https://rotational.app)
- [download and install Python](https://www.python.org/downloads/) according to your operating system

The Python SDK is currently compatible with Python 3.7, 3.8, 3.9, and 3.10.

### Project Setup

Create a project directory and install the official [Python SDK](https://pypi.org/project/pyensign/) using pip.

```bash
mkdir hello-ensign
cd hello-ensign
pip install pyensign
```

<a name="create-a-client"></a>
### Create a Client

Create a `main.py` file and create the Ensign client in code, which is similar to a database client like PostgreSQL or Mongo. The Python SDK is an [asyncio](https://docs.python.org/3/library/asyncio.html) API, which means we need to run the client methods as [coroutines](https://docs.python.org/3/library/asyncio-task.html#coroutines). Python's `async/await` syntax abstracts most of this for us.

```python
import asyncio
from datetime import datetime

from pyensign.events import Event
from pyensign.ensign import Ensign

async def main():
    client = Ensign()
    status = await client.status()
    print(status)

if __name__ == '__main__':
    asyncio.run(main())
```

The Python SDK requires a Client ID and Client Secret to communicate with Ensign. We recommend specifying them in the environment like so (replace with the values in your API key).

```bash
export ENSIGN_CLIENT_ID=DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa
export ENSIGN_CLIENT_SECRET=wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS
```

If you find yourself having to manage multiple API keys on the same machine, you can also specify a path to a JSON file with your credentials.

**my_project_key.json**
```json
{
    "ClientID": "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
    "ClientSecret": "wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS"
}
```

```python
client = Ensign(cred_path="my_project_key.json")
```

Run the code as a Python program.

```bash
python main.py
```

If you see a message like the following, then congratulations! You've successfully connected to Ensign!

```
status: 1
version: 0.12.5-beta.20 ([GIT HASH])
uptime: seconds: 130150
nanos: 862300696
```

### Make Some Data

Next, we need some data! Generally this is the place where you'd connect to your live data source (a database, weather data, etc). But to keep things simple, we'll just create a single event, which starts with a dictionary.

```python
data = {
    "sender": "Twyla",
    "timestamp": datetime.now(),
    "message": "Let's get this started!"
}
```

Next, we will convert our map into an event. Usually, an event contains some data (encoded to bytes), a mimetype which indicates how the data was encoded, and a schema type. The schema type consists of a name and a [semantic version](https://semver.org/) string.

```python
event = Event(
    json.dumps(data).encode("utf-8"),
    mimetype="application/json",
    schema_name="Generic",
    schema_version="1.0.0"
)
```

### Publish Your Event

Now we can publish your event by awaiting the `publish` method on the Ensign client we created above. You'll also need to pass in a topic name, which will be a string. If you aren't sure what topic to use, you can quickly [log into your Ensign dashboard](https://rotational.app) and look it up.

```python
await client.publish("quality-lemon-time", event)
```

You can publish many events at a time if you want!

```python
await client.publish("quality-lemon-time", event, event2, event3, event4)
```

### Create a Subscriber

So now you've published some events to a topic. We can consume those events with the `subscribe` method. `subscribe` works a bit differently than `publish`. Instead of immediately returning, it `yields` events to the caller. We can use the `async for` syntax to process the events as they come in on the stream.

```python
async for event in client.subscribe("quality-lemon-time"):
    msg = json.loads(event.data)
    print(msg["message"])
```

Try running the program again and see if you can get the message!

```bash
python main.go
```

```Let's get this started!```

## Next Steps

You're already well on your way to building your first event-driven microservice with Ensign!

If you're ready to see some more advanced examples with code, check out the [End-to-end Examples]({{< relref "examples">}}).

If you're looking for more on the basics of event-driven systems, check out [Eventing 101]({{< relref "eventing">}}).

Happy eventing!