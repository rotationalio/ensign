import { Card, Heading } from '@rotational/beacon-core';

import Button from '@/components/ui/Button';

export default function AccessDocumentationStep() {
  return (
    <>
      <Card contentClassName="w-full min-h-[200px] border-2 rounded-md p-4">
        <Card.Header>
          <Heading as="h3" className="px-2 font-bold">
            Step 3: View Documentation
          </Heading>
        </Card.Header>
        <Card.Body>
          <div className="mt-5 flex flex-col gap-8 px-3 xl:flex-row">
            <p className="w-full text-sm sm:w-4/5">
              Love seeing examples with real code? Prefer watching tutorial videos? Still learning
              the basics? We’ve got you covered!
            </p>
            <div className="sm:w-1/5">
              <a
                href="https://ensign.rotational.dev/getting-started/"
                target="_blank"
                rel="noopener noreferrer"
                data-testid="viewDocsLink"
              >
                <Button className="h-[44px] w-[165px] text-sm" data-testid="viewDocs">
                  View Docs
                </Button>
              </a>
            </div>
          </div>
        </Card.Body>
      </Card>
    </>
  );
}
