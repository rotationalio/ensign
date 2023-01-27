import { Button } from '@rotational/beacon-core';

import world from '/src/assets/images/world-icon.png';

export default function CustomizeTenant() {
  return (
    <div className="rounded-lg border border-solid border-primary-800">
      <h3 className="font-bold">Customize Tenant</h3>
      <img src={world} alt="" />
      <p>
        You can customize your tenant settings, including regions (single and multi-region) and
        cloud providers.
      </p>
      <p>
        Customizing your tenant requires a paid plan. Upgrade now to set up a tenant specific to
        your development and modeling needs.
      </p>
      <Button>Coming Soon</Button>
    </div>
  );
}
