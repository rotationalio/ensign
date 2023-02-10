import { Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { useCreateTenant } from '../hooks/useCreateTenant';

export default function TenantTable() {
  const [, setIsOpen] = useState(false);
  const handleClose = () => setIsOpen(false);

  const [items, setItems] = useState();

  const { tenant, isFetchingTenant, hasTenantFailed, wasTenantFetched, error } = useCreateTenant();

  if (isFetchingTenant) {
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasTenantFailed}
        onClose={handleClose}
        variant="danger"
        title="We were unable to fetch your tenants. Please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  // TODO: Add cloud provider and region once added to Tenant API.
  if (wasTenantFetched && tenant) {
    const ft = Object.keys(tenant).map((t) => {
      const { name, env, created } = tenant[t];
      return { name, env, created };
    }) as any;
    setItems(ft);
  }

  return (
    <>
      <div className="rounded-lg bg-[#F7F9FB] py-2">
        <h2 className="ml-4 text-lg font-bold">Tenants</h2>
      </div>
      <Table
        columns={[
          { Header: 'Tenant Name', accessor: 'name' },
          { Header: 'Environment Label', accessor: 'env' },
          /* { Header: 'Cloud Provider', accessor: 'cloud'},
            { Header: 'Region', accessor: 'region'}, */
          { Header: 'Date Created', accessor: 'created' },
        ]}
        data={items}
      />
    </>
  );
}
