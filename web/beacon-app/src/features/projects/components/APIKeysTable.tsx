import { Heading, Table, Toast } from '@rotational/beacon-core';
import { useEffect, useState } from 'react';

import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import Button from '@/components/ui/Button';
import GenerateAPIKeyModal from '@/features/apiKeys/components/GenerateAPIKeyModal';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';

import { getApiKeys } from '../util';
interface APIKeysTableProps {
  projectID: string;
}

export const APIKeysTable = ({ projectID }: APIKeysTableProps) => {
  const { apiKeys, isFetchingApiKeys, hasApiKeysFailed, error } = useFetchApiKeys(projectID);
  const [isOpenAPIKeyDataModal, setIsOpenAPIKeyDataModal] = useState<boolean>(false);
  const [isOpenGenerateAPIKeyModal, setIsOpenGenerateAPIKeyModal] = useState<boolean>(false);
  const [key, setKey] = useState<any>(null);
  const onOpenGenerateAPIKeyModal = () => {
    setIsOpenGenerateAPIKeyModal(true);
  };

  const onCloseGenerateAPIKeyModal = () => {
    setIsOpenGenerateAPIKeyModal(false);
  };

  const onCloseAPIKeyDataModal = () => setIsOpenAPIKeyDataModal(false);

  useEffect(() => {
    if (key) {
      setIsOpenAPIKeyDataModal(true);
    }
  }, [key]);
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

  return (
    <div className="text-sm">
      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <Heading as={'h1'} className="text-lg font-semibold">
          API Keys
        </Heading>
        <Button
          variant="primary"
          size="small"
          className="!text-xs"
          onClick={onOpenGenerateAPIKeyModal}
        >
          + Add New Key
        </Button>
      </div>
      <Table
        className="w-full"
        columns={[
          { Header: 'Name', accessor: 'name' },
          { Header: 'Client ID', accessor: 'client_id' },
        ]}
        data={getApiKeys(apiKeys)}
      />
      <ApiKeyModal open={isOpenAPIKeyDataModal} data={key} onClose={onCloseAPIKeyDataModal} />
      <GenerateAPIKeyModal
        open={isOpenGenerateAPIKeyModal}
        onClose={onCloseGenerateAPIKeyModal}
        onSetKey={setKey}
      />
    </div>
  );
};

export default APIKeysTable;
