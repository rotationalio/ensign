import { Trans } from '@lingui/macro';
import { Button, Card, Heading } from '@rotational/beacon-core';

export default function AccessDocumentationStep() {
  return (
    <>
      <Card contentClassName="w-full min-h-[200px] border-2 rounded-md p-4">
        <Card.Header>
          <Heading as="h3" className="px-2 font-bold">
            <Trans>Step 3: View Documentation</Trans>
          </Heading>
        </Card.Header>
        <Card.Body>
          <div className="mt-5 flex flex-col gap-8 px-3 md:flex-row">
            <p className="w-full sm:w-4/5">
              <Trans>
                Love seeing examples with real code? Prefer watching tutorial videos? Still learning
                the basics? We’ve got you covered!
              </Trans>
            </p>
            <div className="sm:w-1/5">
              <a
                href="https://ensign.rotational.dev/getting-started/"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button className="h-[44px] w-[165px] text-sm">
                  <Trans>View Docs</Trans>
                </Button>
              </a>
            </div>
          </div>
        </Card.Body>
      </Card>
    </>
  );
}
