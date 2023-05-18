import { t, Trans } from '@lingui/macro';
import { Button, Heading, Table, Toast } from '@rotational/beacon-core';
import { useEffect, useState } from 'react';

import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';
import GenerateAPIKeyModal from '@/features/apiKeys/components/GenerateAPIKeyModal';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { formatDate } from '@/utils/formatDate';

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
    // console.log('onOpenGenerateAPIKeyModal', isOpenGenerateAPIKeyModal);
    // console.log('isOpenAPIKeyDataModal', isOpenAPIKeyDataModal);
    setIsOpenGenerateAPIKeyModal(true);

    // console.log('onOpenGenerateAPIKeyModal', isOpenGenerateAPIKeyModal);
  };

  const onCloseGenerateAPIKeyModal = () => {
    setIsOpenGenerateAPIKeyModal(false);
  };

  const onCloseAPIKeyDataModal = () => {
    setIsOpenAPIKeyDataModal(false);
  };

  useEffect(() => {
    if (key) {
      setIsOpenAPIKeyDataModal(true);
    }
  }, [key]);

  if (isFetchingApiKeys) {
    // TODO: add loading state
    return <div>Loading...</div>;
  }

  if (hasApiKeysFailed) {
    <Toast
      isOpen={hasApiKeysFailed}
      variant="danger"
      description={(error as any)?.response?.data?.error}
    />;
  }

  const initialColumns: any = [
    { Header: t`Key Name`, accessor: 'name' },
    {
      Header: t`Status`,
      accessor: 'status',
    },
    { Header: t`Permissions`, accessor: 'permissions' },
    {
      Header: t`Last Used`,
      accessor: (date: any) => {
        return formatDate(new Date(date?.last_used));
      },
    },
    {
      Header: t`Date Created`,
      accessor: (date: any) => {
        return formatDate(new Date(date?.created));
      },
    },
  ];

  return (
    <div className="mt-[46px]  border-y-neutral-600" data-cy="keyComp">
      <Heading as={'h1'} className="flex items-center text-lg font-semibold capitalize">
        <Trans>API Keys</Trans>
      </Heading>
      <div className="flex space-x-1">
        <p className="my-4">
          <Trans>
            API keys enable you to securely connect your data sources to Ensign. Generate at least
            one API key for your project. You can customize permissions.
          </Trans>
          <span className="ml-2" data-cy="keyHint">
            <HelpTooltip data-cy="keyInfo">
              <p>
                <Trans>
                  Each key consists of two parts - a ClientID and a ClientSecret. You'll need both
                  to establish a client connection, create Ensign topics, publishers, and
                  subscribers. Keep your API keys private -- if you misplace your keys, you can
                  revoke them and generate new ones.
                </Trans>
              </p>
            </HelpTooltip>
          </span>
        </p>
      </div>
      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <div className="flex items-center gap-3"></div>
        <div>
          <Button
            variant="primary"
            size="small"
            className="px-5 !text-xs"
            onClick={onOpenGenerateAPIKeyModal}
            data-cy="addKey"
          >
            + New Key
          </Button>
        </div>
      </div>
      <Table
        trClassName="text-sm"
        columns={initialColumns}
        data={getApiKeys(apiKeys)}
        data-cy="keyTable"
      />
      <ApiKeyModal open={isOpenAPIKeyDataModal} data={key} onClose={onCloseAPIKeyDataModal} />
      <GenerateAPIKeyModal
        open={isOpenGenerateAPIKeyModal}
        onClose={onCloseGenerateAPIKeyModal}
        onSetKey={setKey}
        projectId={projectID}
      />
    </div>
  );
};

export default APIKeysTable;
