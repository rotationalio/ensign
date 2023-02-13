import { AriaButton as Button, Toast } from '@rotational/beacon-core';
import { useState } from 'react';
import { Link } from 'react-router-dom';

import LightningBolt from '@/components/icons/lightning-bold-icon';
import Loader from '@/components/ui/Loader';
import { APP_ROUTE } from '@/constants';
import SuccessfulTenantCreationModal from '@/features/misc/components/onboarding/SuccessfulTenantCreationModal';
import { useCreateTenant } from '@/features/tenants/hooks/useCreateTenant';

export default function TenantQuickStart() {
  const [, setIsOpen] = useState(false);
  const tenant = useCreateTenant();

  const handleClose = () => setIsOpen(false);

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
      <div className="mx-auto max-w-2xl rounded-lg border border-solid border-primary-800">
        <div className="max-w-sm p-10 text-center">
          <h3 className="text-xl font-bold">Quick Start Tenant</h3>

          <LightningBolt className="mx-auto mt-5" />
          <p className="mt-12">
            We&#39;ll set up a tenant for you based on the closest region with availability x (by
            your IP address) and label it as your <span className="font-bold">Development</span>{' '}
            environment.
          </p>
          <p className="mt-6">You can change settings later and upgrade at any time.</p>
          <Link to={APP_ROUTE.TENANTS}>
            <Button
              isLoading={isFetchingTenant}
              color="secondary"
              variant="secondary"
              size="large"
              className="mt-28 w-48"
            >
              Create
            </Button>
          </Link>
          {}
        </div>
      </div>
    </>
  );
}
