import { Table, Toast } from '@rotational/beacon-core';

import { TableHeading } from '@/components/common/TableHeader';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';

interface APIKeysTableProps {
  projectID: string;
}

export const APIKeysTable = ({ projectID }: APIKeysTableProps) => {
  const { apiKeys, isFetchingApiKeys, hasApiKeysFailed, error } = useFetchApiKeys(projectID);

  if (isFetchingApiKeys) {
    // TODO: add loading state
    return <div>Loading...</div>;
  }

  if (error) {
    <Toast
      isOpen={hasApiKeysFailed}
      variant="danger"
      title="Sorry we are having trouble fetching your API Keys, please try again later."
      description={(error as any)?.response?.data?.error}
    />;
  }

  const getApiKeys = (apikeys: any) => {
    if (!apikeys?.api_keys || apikeys?.api_keys.length === 0) return [];
    return Object.keys(apiKeys?.api_keys).map((key) => {
      const { id, name, client_id } = apiKeys.api_keys[key];
      return { id, name, client_id };
    }) as any;
  };

  return (
    <div>
      <TableHeading>API Keys</TableHeading>
      <Table
        className="w-full"
        columns={[
          { Header: 'ID', accessor: 'id' },
          { Header: 'Name', accessor: 'name' },
          { Header: 'Client ID', accessor: 'client_id' },
        ]}
        data={getApiKeys(apiKeys)}
      />
    </div>
  );
};

export default APIKeysTable;
