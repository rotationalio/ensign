import OrganizationsTable from '@/features/organizations/components/OrganizationTable/OrganizationTable';

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
