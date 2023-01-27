import { Button } from '@rotational/beacon-core';

import bolt from '/src/assets/images/lightning-bolt.png';

export default function TenantQuickStart() {
  return (
    <div className="rounded-lg border border-solid border-primary-800">
      <h3 className="font-bold">Quick Start Tenant</h3>
      <img src={bolt} alt="" />
      <p>
        We&#39;ll set up a tenant for you based on the closest region with availability x (by your
        IP address) and label it as your <span>Development</span> environment.
      </p>
      <p>You can change settings later and upgrade at any time.</p>
      <Button>Create</Button>
    </div>
  );
}
