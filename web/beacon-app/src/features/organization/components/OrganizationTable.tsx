import { Heading, Loader, Table, Toast } from '@rotational/beacon-core';
import { Suspense } from 'react';

import { SentryErrorBoundary } from '@/components/Error';
import { useOrgStore } from '@/store';
import { formatDate } from '@/utils/formatDate';

import { useFetchOrg } from '../hooks/useFetchOrgDetail';

export default function OrganizationsTable() {
  const orgDataState = useOrgStore.getState() as any;
  const { org, isFetchingOrg, hasOrgFailed, error } = useFetchOrg(orgDataState.org);

  if (isFetchingOrg) {
    return <Loader />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasOrgFailed}
        variant="danger"
        title="We were unable to fetch your organizations. Please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const { id, name, created, owner } = org;

  return (
    <>
      <Suspense fallback={<Loader />}>
        <SentryErrorBoundary
          fallback={
            <div>Sorry, We were unable to fetch your organizations. Please try again later.</div>
          }
        >
          <div className="rounded-lg bg-[#F7F9FB] py-2">
            <Heading as={'h2'} className="ml-4 text-lg font-bold">
              Organizations
            </Heading>
          </div>
          <div className="overflow-hidden text-sm">
            <Table
              columns={[
                { Header: 'Organization ID', accessor: 'id' },
                { Header: 'Organization Name', accessor: 'name' },
                { Header: 'Organization Owner', accessor: 'role' },
                { Header: 'Projects', accessor: 'projects' },
                {
                  Header: 'Date Created',
                  accessor: (date: any) => {
                    return formatDate(new Date(date.created));
                  },
                },
              ]}
              data={[
                {
                  id: id,
                  name: name,
                  role: owner,
                  created: created,
                },
              ]}
            />
          </div>
        </SentryErrorBoundary>
      </Suspense>
    </>
  );
}
