import { Button } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import bolt from '@/assets/images/lightning-bolt.png';

export default function TenantQuickStart() {
  return (
    <div className="mx-auto max-w-2xl rounded-lg border border-solid border-primary-800">
      <div className="max-w-sm p-10 text-center">
        <h3 className="text-xl font-bold">Quick Start Tenant</h3>
        <img src={bolt} alt="" className="mx-auto mt-5" />
        <p className="mt-6">
          We&#39;ll set up a tenant for you based on the closest region with availability x (by your
          IP address) and label it as your <span className="font-bold">Development</span>{' '}
          environment.
        </p>
        <p className="mt-6">You can change settings later and upgrade at any time.</p>
        <Link to="/">
          <Button color="secondary" size="large" className="mt-32">
            Create
          </Button>
        </Link>
      </div>
    </div>
  );
}
