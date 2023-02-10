import { Heading, Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { useRegister } from '@/features/auth';

import { useFetchOrg } from '../../hooks/useFetchOrgDetail';

export default function OrganizationsTable() {
  const [, setIsOpen] = useState(false);
  const handleClose = () => setIsOpen(false);

  const { user } = useRegister();

  const { org_id } = user;

  const { org, isFetchingOrg, hasOrgFailed, error } = useFetchOrg(org_id);

  if (isFetchingOrg) {
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasOrgFailed}
        onClose={handleClose}
        variant="danger"
        title="We were unable to fetch your organizations. Please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const { name, created } = org;

  return (
    <>
      <div className="rounded-lg bg-[#F7F9FB] py-2">
        <Heading as={'h2'} className="ml-4 text-lg font-bold">
          Organizations
        </Heading>
      </div>
      <Table
        columns={[
          { Header: 'Organization Name', accessor: 'name' },
          { Header: 'Organization Role', accessor: 'role' },
          { Header: 'Projects', accessor: 'projects' },
          { Header: 'Date Created', accessor: 'created' },
        ]}
        data={[
          {
            name: { name },
            role: 'Owner',
            projects: '1',
            created: { created },
          },
        ]}
      />
    </>
  );
}
