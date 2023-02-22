import OrganizationsTable from '@/features/organization/components/OrganizationTable';

import MemberDetails from './MemberDetails';

export default function MemberDetailsPage() {
  return (
    <>
      <MemberDetails />
      <div className="mt-10">
        <OrganizationsTable />
      </div>
    </>
  );
}
