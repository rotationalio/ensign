import { Heading, Table, Toast } from '@rotational/beacon-core';

import { formatDate } from '@/utils/formatDate';

import { useFetchTenants } from '../hooks/useFetchTenants';

export default function TenantTable() {
  const { tenants, isFetchingTenants, hasTenantsFailed, error } = useFetchTenants();

  if (isFetchingTenants) {
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasTenantsFailed}
        variant="danger"
        title="We were unable to fetch your tenants. Please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }
  return (
    <>
      <div className="rounded-lg bg-[#F7F9FB] py-2">
        <Heading as={'h2'} className="ml-4 text-lg font-bold">
          Tenants
        </Heading>
      </div>
      <Table
        columns={[
          { Header: 'Tenant Name', accessor: 'name' },
          { Header: 'Environment Label', accessor: 'environment_type' },
          /* { Header: 'Cloud Provider', accessor: 'cloud'},
            { Header: 'Region', accessor: 'region'}, */
          {
            Header: 'Date Created',
            accessor: (date: any) => {
              return formatDate.format(new Date(date.created));
            },
          },
        ]}
        data={tenants?.tenants}
      />
    </>
  );
}
