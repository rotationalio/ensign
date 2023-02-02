import { routes } from '@/application';
import { Button } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

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
          Your tenant is your team's control panel for all projects and topics. Specify preferences around encryption, privacy, and locality (e.g. for GDPR, CCPA, etc).
        </p>
      </div>
      <Link to={routes.setup}>
        <Button
          color="secondary"
          size="large" 
          className="mx-auto">
            Set Up
        </Button>
      </Link>
    </section>
  );
}

export default SetupTenant;
