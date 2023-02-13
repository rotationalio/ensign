import { Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { TableHeading } from '@/components/common/TableHeader';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import type { APIKey } from '@/features/apiKeys/types/apiKeyService';

export const APIKeysTable = () => {
  const [items, setItems] = useState<Omit<APIKey, 'id'>>();
  const [, setIsOpen] = useState(false);
  const handleClose = () => setIsOpen(false);

  const { apiKeys, isFetchingApiKeys, hasApiKeysFailed, wasApiKeysFetched, error } =
    useFetchApiKeys();

  if (isFetchingApiKeys) {
    // TODO: add loading state
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasApiKeysFailed}
        onClose={handleClose}
        variant="danger"
        title="Sorry we are having trouble fetching your topics, please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  if (wasApiKeysFetched && apiKeys) {
    // format apiKeys to match table
    const fk = Object.keys(apiKeys).map((key) => {
      const { name, owner, permissions, modifiers, created } = apiKeys[key];
      return { name, owner, permissions, modifiers, created };
    }) as any;
    setItems(fk);
  }

  return (
    <div>
      <TableHeading>API Keys</TableHeading>
      <Table
        columns={[
          { Header: 'Name', accessor: 'name' },
          { Header: 'Permissions', accessor: 'permissions' },
          { Header: 'Owner', accessor: 'owner' },
          { Created: 'Last Used', accessor: 'modifiers' },
          { Created: 'Date Created', accessor: 'created' },
        ]}
        data={items}
      />
    </div>
  );
};

export default APIKeysTable;
