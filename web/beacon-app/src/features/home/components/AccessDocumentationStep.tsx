import { AriaButton as Button, Card } from '@rotational/beacon-core';

export default function AccessDocumentationStep() {
  return (
    <>
      <Card contentClassName="w-full min-h-[200px] border border-primary-900 rounded-md p-4">
        <Card.Header>
          <h1 className="font-bold">Step 3: View Documentation</h1>
        </Card.Header>
        <Card.Body>
          <div className="mt-5 flex flex-col gap-8 md:flex-row">
            <p className="w-full md:w-4/5 lg:w-4/5">
              Love seeing examples with real code? Prefer watching tutorial videos? Still learning
              the basics? We’ve got you covered!
            </p>
            <div className="mr-8 grid w-full place-items-center gap-3 md:w-1/5 lg:w-1/5">
              <a
                href="https://ensign.rotational.dev/getting-started/"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button className="text-sm">View Docs</Button>
              </a>
            </div>
          </div>
        </Card.Body>
      </Card>
    </>
  );
}
