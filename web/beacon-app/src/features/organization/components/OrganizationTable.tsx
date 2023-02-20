import { Heading, Table, Toast } from '@rotational/beacon-core';
import { useEffect, useState } from 'react';

import { queryCache } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { useFetchOrg } from '../hooks/useFetchOrgDetail';

export default function OrganizationsTable() {
  const [, setIsOpen] = useState(false);
  const handleClose = () => setIsOpen(false);

  const [orgID, setOrgID] = useState<any>();

  const orgs = queryCache.find(RQK.ORG_DETAIL) as any;

  useEffect(() => {
    if (orgs) {
      setOrgID(orgs[0].id as string);
    }
  }, [orgs]);

  const { org, isFetchingOrg, hasOrgFailed, error } = useFetchOrg(orgID);

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
