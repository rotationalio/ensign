import AppLayout from '@/components/layout/AppLayout';
import OrganizationsTable from '@/features/organization/components/OrganizationTable';

import MemberDetails from '../components/MemberDetails';

export default function MemberPage() {
  return (
    <>
      <AppLayout>
        <MemberDetails />
        <OrganizationsTable />
      </AppLayout>
    </>
  );
}
