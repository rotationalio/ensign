import { Heading } from '@rotational/beacon-core';

import AppLayout from '@/components/layout/AppLayout';
import OrganizationsTable from '@/features/organization/components/OrganizationTable';

import MemberDetails from './MemberDetails';

export default function MemberDetailsPage() {
  return (
    <>
      <AppLayout>
        <Heading as="h1" className="mb-4 text-lg font-semibold">
          User Profile
        </Heading>
        <MemberDetails />
        <OrganizationsTable />
      </AppLayout>
    </>
  );
}
