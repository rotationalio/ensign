import json
import asyncio

from pyensign.events import Event
from pyensign.ensign import Ensign

def load_event_fixtures(path):
    events = []
    print()
    with open(path,'r') as f:
        data = json.load(f)
        for fixture in data:
            events.append(Event(json.dumps(fixture["data"]).encode("utf-8"), mimetype=fixture["mimetype"]))
    return events

async def publish_fixtures(ensign_creds):
    """
    Read the JSON files from the fixtures directory and publish events to the test topics.
    """
    ensign = Ensign(
        cred_path=ensign_creds,
        endpoint="staging.ensign.world:443", 
        auth_url="https://auth.ensign.world"
        )

    # A topic that contains only one mimetype
    # Less than 10 events
    events = load_event_fixtures('fixtures/one_type.json')
    for event in events:
        await ensign.publish("documents_one_type", event)

    # A topic that contains multiple mimetypes
    # More than 10 events - we want to test that only 10 events are returned when the FE queries the topic
    events = load_event_fixtures('fixtures/two_type.json')
    for event in events:
        await ensign.publish("documents_two_type", event)

if __name__ == "__main__":
    asyncio.run(publish_fixtures(ensign_creds='secret/qa_creds.json'))