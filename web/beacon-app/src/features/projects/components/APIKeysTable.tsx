import { Table, Toast } from '@rotational/beacon-core';

import { TableHeading } from '@/components/common/TableHeader';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { formatDate } from '@/utils/formatDate';

export const APIKeysTable = () => {
  const { apiKeys, isFetchingApiKeys, hasApiKeysFailed, error } = useFetchApiKeys();

  if (isFetchingApiKeys) {
    // TODO: add loading state
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasApiKeysFailed}
        variant="danger"
        title="Sorry we are having trouble fetching your topics, please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const getApiKeys = (apikeys: any) => {
    if (!apikeys) return [];
    return Object.keys(apiKeys).map((key) => {
      const { name, owner, permissions, modifiers, created } = apiKeys[key];
      return { name, owner, permissions, modifiers, created };
    }) as any;
  };

  return (
    <div>
      <TableHeading>API Keys</TableHeading>
      <Table
        className="w-full"
        columns={[
          { Header: 'Name', accessor: 'name' },
          { Header: 'Permissions', accessor: 'permissions' },
          { Header: 'Owner', accessor: 'owner' },
          { Header: 'Last Used', accessor: 'modifiers' },
          {
            Header: 'Date Created',
            accessor: (date: any) => {
              return formatDate(new Date(date.created));
            },
          },
        ]}
        data={getApiKeys(apiKeys)}
      />
    </div>
  );
};

export default APIKeysTable;
