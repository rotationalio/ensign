import { Heading, Table, Toast } from '@rotational/beacon-core';

import Button from '@/components/ui/Button';
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
    <div className="text-sm">
      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <Heading as={'h1'} className="text-lg font-semibold">
          API Keys
        </Heading>
        <Button variant="primary" size="small" className="!text-xs">
          + Add new Key
        </Button>
      </div>
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
