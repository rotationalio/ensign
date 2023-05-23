import { t, Trans } from '@lingui/macro';
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
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const { name, created, owner, projects } = org;

  return (
    <>
      <Suspense fallback={<Loader />}>
        <SentryErrorBoundary
          fallback={
            <div>
              <Trans>Sorry, We were unable to fetch your workspaces. Please try again later.</Trans>
            </div>
          }
        >
          <div className="rounded-lg bg-[#F7F9FB] py-2">
            <Heading as={'h2'} className="ml-4 text-lg font-bold">
              <Trans>Workspaces</Trans>
            </Heading>
          </div>
          <div className="overflow-hidden text-sm" data-testid="orgTable">
            <Table
              trClassName="text-sm"
              columns={[
                { Header: t`Workspace Name`, accessor: 'name' },
                { Header: t`Workspace Owner`, accessor: 'role' },
                { Header: t`Projects`, accessor: 'projects' },
                {
                  Header: t`Date Created`,
                  accessor: (date: any) => {
                    return formatDate(new Date(date.created));
                  },
                },
              ]}
              data={[
                {
                  name: name,
                  role: owner,
                  projects: projects,
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
