import { AriaButton, Card, Heading } from '@rotational/beacon-core';

import eventingIcon from '/src/components/icons/eventingIcon.svg';

function ViewTutorials() {
  return (
    <Card className="mx-auto max-w-5xl py-10 text-xl">
      <div className="grid grid-cols-3 md:grid">
        <img src={eventingIcon} alt="" className="mx-auto mt-6" />
        <div className="text-center">
          <Heading as="h2">
            View Tutorials <span className="font-normal">(optional)</span>
          </Heading>
          <p className="mt-6">
            From quickstarts to detailed examples for data engineers, data scientists, and app
            developers, we&apos;ve got you covered.
          </p>
        </div>
        <div className="mx-auto mt-8 md:mt-16">
          <AriaButton color="secondary" size="large" className="w-32">
            <a
              href="https://ensign.rotational.dev/getting-started/"
              target="_blank"
              rel="noreferrer"
            >
              View
            </a>
          </AriaButton>
        </div>
      </div>
    </Card>
  );
}

export default ViewTutorials;
