import { Button } from '@rotational/beacon-core';

import bolt from '/src/assets/images/lightning-bolt.png';

export default function TenantQuickStart() {
  return (
    <div className="mx-auto max-w-2xl rounded-lg border border-solid border-primary-800">
      <div className="max-w-md p-10 text-center">
        <h3 className="text-xl font-bold">Quick Start Tenant</h3>
        <img src={bolt} alt="" className="mx-auto mt-5" />
        <p className="mt-6">
          We&#39;ll set up a tenant for you based on the closest region with availability x (by your
          IP address) and label it as your <span className="font-bold">Development</span>{' '}
          environment.
        </p>
        <p className="mt-6">You can change settings later and upgrade at any time.</p>
        <Button className="mt-20 h-14 w-44 rounded bg-danger-500 text-center text-xl font-bold text-white">Create</Button>
      </div>
    </div>
  );
}
