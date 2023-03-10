import { AriaButton as Button, Heading } from '@rotational/beacon-core';

import WorldIcon from '@/components/icons/world-icon';

export default function CustomizeTenant() {
  return (
    <div className="mx-auto max-w-sm rounded-lg border border-solid border-primary-800">
      <div className="p-10 text-center">
        <Heading as="h3">Customize Tenant</Heading>

        <WorldIcon className="mx-auto mt-5" />
        <div className="ml-8 text-left">
          <p className="mt-6">
            You can customize your tenant settings, including regions (single and multi-region) and
            cloud providers.
          </p>
          <p className="mt-6">
            Customizing your tenant requires a paid plan. Upgrade now to set up a tenant specific to
            your development and modeling needs.
          </p>
        </div>
        <Button color="secondary" size="large" className="mt-16">
          Coming Soon
        </Button>
      </div>
    </div>
  );
}
