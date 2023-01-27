import { Button } from '@rotational/beacon-core';

import world from '/src/assets/images/world-icon.png';

export default function CustomizeTenant() {
  return (
    <div className="mx-auto max-w-2xl rounded-lg border border-solid border-primary-800">
      <div className="max-w-md p-10 text-center">
        <h3 className="text-xl font-bold">Customize Tenant</h3>
        <img src={world} alt="" className="mx-auto mt-5" />
        <p className="mt-6">
          You can customize your tenant settings, including regions (single and multi-region) and
          cloud providers.
        </p>
        <p className="mt-6">
          Customizing your tenant requires a paid plan. Upgrade now to set up a tenant specific to
          your development and modeling needs.
        </p>
        <Button className="mx-auto mt-10 h-14 w-40 rounded bg-[#E66809] text-center text-xl font-bold text-white">
          Coming Soon
        </Button>
      </div>
    </div>
  );
}
