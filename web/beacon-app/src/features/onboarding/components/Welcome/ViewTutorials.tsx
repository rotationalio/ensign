import { AriaButton } from '@rotational/beacon-core';

import eventing from '/src/assets/images/eventing-icon.png';

function ViewTutorials() {
  return (
    <section className="grid grid-cols-3 mx-auto max-w-4xl rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <img src={eventing} alt="" className="mx-auto mt-6" />
      <div>
        <h2 className="mt-10 font-bold">
          View Tutorials <span className="font-normal">(optional)</span>
        </h2>
        <p className="mt-6">
          From quickstarts to detailed examples for data engineers, data scientists, and app developers, we've got you covered.
        </p>
      </div>
      <div className="mx-auto mt-32">
      <AriaButton 
        color="secondary"
        size="large"
        >
          <a href="https://ensign.rotational.dev/getting-started/" target="_blank">View</a>
      </AriaButton>
      </div>
    </section>
  );
}

export default ViewTutorials;
