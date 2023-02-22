import { t, Trans } from '@lingui/macro';
import { Heading, Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { useFetchTenants } from '../hooks/useFetchTenants';

export default function TenantTable() {
  const [, setIsOpen] = useState(false);
  const handleClose = () => setIsOpen(false);

  const { getTenants, tenants, isFetchingTenants, hasTenantsFailed, error } = useFetchTenants();

  if (!tenants) {
    getTenants();
  }

  if (isFetchingTenants) {
    return (
      <div>
        <Trans>Loading...</Trans>
      </div>
    );
  }

  if (error) {
    return (
      <Toast
        isOpen={hasTenantsFailed}
        onClose={handleClose}
        variant="danger"
        title={t`We were unable to fetch your tenants. Please try again later.`}
        description={(error as any)?.response?.data?.error}
      />
    );
  }
  return (
    <>
      <div className="rounded-lg bg-[#F7F9FB] py-2">
        <Heading as={'h2'} className="ml-4 text-lg font-bold">
          <Trans>Tenants</Trans>
        </Heading>
      </div>
      <Table
        columns={[
          { Header: t`Tenant Name`, accessor: 'name' },
          { Header: t`Environment Label`, accessor: 'env' },
          /* { Header: 'Cloud Provider', accessor: 'cloud'},
            { Header: 'Region', accessor: 'region'}, */
          { Header: t`Date Created`, accessor: 'created' },
        ]}
        data={tenants?.tenants}
      />
    </>
  );
}
