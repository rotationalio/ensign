import { AriaButton as Button, Heading, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import LightningBolt from '@/components/icons/lightning-bold-icon';
import Loader from '@/components/ui/Loader';
import SuccessfulTenantCreationModal from '@/features/misc/components/onboarding/SuccessfulTenantCreationModal';
import { useCreateTenant } from '@/features/tenants/hooks/useCreateTenant';

export default function TenantQuickStart() {
  const [, setIsOpen] = useState(false);
  const tenant = useCreateTenant();

  const handleClose = () => setIsOpen(false);

  const handleTenantCreation = () => {
    tenant.createTenant();
  };

  const { isFetchingTenant, hasTenantFailed, wasTenantFetched, error } = tenant;

  if (isFetchingTenant) {
    return <Loader />;
  }

  if (wasTenantFetched) {
    return <SuccessfulTenantCreationModal open={true} />;
  }

  return (
    <>
      <Toast
        isOpen={hasTenantFailed}
        onClose={handleClose}
        variant="danger"
        title="Unable to create Quick Start Tenant. Please try again."
        description={(error as any)?.response?.data?.error}
      />
      <div className="mx-auto max-w-sm rounded-lg border border-solid border-primary-800">
        <div className="p-10 text-center">
          <Heading as="h3" className="pt-1">
            Quick Start Tenant
          </Heading>

          <LightningBolt className="mx-auto mt-5" />
          <div className="ml-6 text-left">
            <p className="mt-11">
              We&#39;ll set up a tenant for you based on the closest region with availability x (by
              your IP address) and label it as your <span className="font-bold">Development</span>{' '}
              environment.
            </p>
            <p className="mt-6">You can change settings later and upgrade at any time.</p>
          </div>
          <Button
            isLoading={isFetchingTenant}
            onClick={handleTenantCreation}
            variant="secondary"
            size="large"
            className="mt-8 w-48 lg:mt-[106px]"
          >
            Create
          </Button>
          {}
        </div>
      </div>
    </>
  );
}
