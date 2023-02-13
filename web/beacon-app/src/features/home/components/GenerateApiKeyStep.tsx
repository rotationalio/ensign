import { AriaButton as Button, Card } from '@rotational/beacon-core';

export default function GenerateApiKeyStep() {
  return (
    <>
      <Card contentClassName="w-full min-h-[200px] border border-primary-900 rounded-md p-4">
        <Card.Header>
          <h1 className="font-bold">Step 3: Generate API Key</h1>
        </Card.Header>
        <Card.Body>
          <div className="mt-5 flex flex-col gap-8 md:flex-row">
            <p className="w-full md:w-4/5 lg:w-4/5">
              API keys enable you to securely connect your data sources to Ensign. Each key consists
              of two parts - a ClientID and a ClientSecret. You’ll need both to establish a client
              connection, create Ensign topics, publishers, and subscribers. Keep your API keys
              private -- if you misplace your keys, you can revoke them and generate new ones.
            </p>
            <div className="mr-8 grid w-full place-items-start gap-3 md:w-1/5 lg:w-1/5">
              <Button className="text-sm">Create API Key</Button>
            </div>
          </div>
        </Card.Body>
      </Card>
    </>
  );
}
