import { Heading } from '@rotational/beacon-core';

import AppLayout from '@/components/layout/AppLayout';
import TenantTable from '@/features/tenants/components/TenantTable';

import OrganizationDetails from '../components/OrganizationDetails';

export default function OrganizationPage() {
  return (
    <AppLayout>
      <div className="-ml-4">
        <Heading as="h1" className="mb-4 text-lg font-semibold">
          Organization Details
        </Heading>
      </div>
      <OrganizationDetails />
      <TenantTable />
    </AppLayout>
  );
}
