import { Button } from '@rotational/beacon-core';

import eventing from '/src/assets/images/eventing-icon.png';

function ViewTutorials() {
  return (
    <section className="grid max-w-6xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <img src={eventing} alt="" className="mx-auto mt-6" />
      <div>
        <h2 className="mt-10 font-bold">
          View Tutorials <span className="font-normal">(optional)</span>
        </h2>
        <p className="mt-6">
          From quickstarts to detailed examples for data engineers, data scientists, and app developers, we've got you covered.
        </p>
      </div>
      <Button 
        color="secondary"
        size="large"
        className="mx-auto"
        >
          <a href="https://ensign.rotational.dev/" target="_blank">View</a>
      </Button>
    </section>
  );
}

export default ViewTutorials;
