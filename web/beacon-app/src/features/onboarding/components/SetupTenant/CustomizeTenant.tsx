import { AriaButton as Button } from '@rotational/beacon-core';

import WorldIcon from '@/components/icons/world-icon';

export default function CustomizeTenant() {
  return (
    <div className="mx-auto max-w-2xl rounded-lg border border-solid border-primary-800">
      <div className="max-w-sm p-10 text-center">
        <h3 className="text-xl font-bold">Customize Tenant</h3>

        <WorldIcon className="mx-auto mt-5" />
        <p className="mt-6">
          You can customize your tenant settings, including regions (single and multi-region) and
          cloud providers.
        </p>
        <p className="mt-6">
          Customizing your tenant requires a paid plan. Upgrade now to set up a tenant specific to
          your development and modeling needs.
        </p>
        <Button color="secondary" size="large" className="mt-16">
          Coming Soon
        </Button>
      </div>
    </div>
  );
}
