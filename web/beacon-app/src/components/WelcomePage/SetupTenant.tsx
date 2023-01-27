import { Button } from '@rotational/beacon-core';

import data from '/src/assets/images/hosted-data-icon.png';

function SetupTenant() {
  return (
    <section className="grid max-w-6xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <img src={data} alt="" className="mx-auto mt-6" />
      <div>
        <h2 className="mt-8 font-bold">
          Set up Your Tenant <span className="font-normal">(required)</span>
        </h2>
        <p className="mt-8">
          A tenant is a collection of settings. The tenant is your locus of control when setting up
          projects and topics.
        </p>
      </div>
      <Button className="mx-auto mt-28 h-14 w-44 rounded bg-[#E66809] text-center text-2xl font-bold text-white">
        Set Up
      </Button>
    </section>
  );
}

export default SetupTenant;
