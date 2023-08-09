import json
import asyncio
import argparse
import base64

from pyensign.events import Event
from pyensign.ensign import Ensign

def decode_from_base64(encoded_data):
    """
    Decode a base64 encoded string to bytes.
    """
    return base64.b64decode(encoded_data)

def load_event_fixtures(path):
    events = []
    print()
    with open(path,'r') as f:
        data = json.load(f)
        for fixture in data:
            if fixture["mimetype"]==252:
                events.append(Event(json.dumps(decode_from_base64(fixture["data"])).encode("utf-8"),
                                mimetype=fixture["mimetype"]))
            else:
                events.append(Event(json.dumps(fixture["data"]).encode("utf-8"),
                                mimetype=fixture["mimetype"]))
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
    events = load_event_fixtures('fixtures_qa/one_type.json')
    for event in events:
        await ensign.publish("documents_one_type", event)

    # A topic that contains multiple mimetypes
    # More than 10 events - we want to test that only 10 events are returned when the FE queries the topic
    events = load_event_fixtures('fixtures_qa/two_type.json')
    for event in events:
        await ensign.publish("documents_two_type", event)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Staging Topics for QA")
    parser.add_argument('path', type=str, help='Enter the Ensig credentials')
    args = parser.parse_args()
    path = args.path
    asyncio.run(publish_fixtures(ensign_creds=path))