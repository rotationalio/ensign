import { AriaButton } from '@rotational/beacon-core';

import EventingIcon from '@/components/icons/eventing-icon';

function ViewTutorials() {
  return (
    <section className="mx-auto grid max-w-4xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <EventingIcon className="mx-auto mt-6" />
      <div>
        <h2 className="mt-10 font-bold">
          View Tutorials <span className="font-normal">(optional)</span>
        </h2>
        <p className="mt-6">
          From quickstarts to detailed examples for data engineers, data scientists, and app
          developers, we&apos;ve got you covered.
        </p>
      </div>
      <div className="mx-auto mt-32">
        <AriaButton color="secondary" size="large" className="w-32">
          <a href="https://ensign.rotational.dev/getting-started/" target="_blank" rel="noreferrer">
            View
          </a>
        </AriaButton>
      </div>
    </section>
  );
}

export default ViewTutorials;
