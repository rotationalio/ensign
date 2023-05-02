import { t, Trans } from '@lingui/macro';
import { Heading, Table, Toast } from '@rotational/beacon-core';
import { useEffect, useState } from 'react';

import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import ConfirmedIndicatorIcon from '@/components/icons/confirmedIndicatorIcon';
import PendingIndicatorIcon from '@/components/icons/pendingIndicatorIcon';
import RevokedIndicatorIcon from '@/components/icons/revokedIndicatorIcon';
import UnusedIndicatorIcon from '@/components/icons/unusedIndicatorIcon';
import Button from '@/components/ui/Button';
import { APIKEY_STATUS } from '@/constants/rolesAndStatus';
import GenerateAPIKeyModal from '@/features/apiKeys/components/GenerateAPIKeyModal';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { APIKeyStatus } from '@/features/apiKeys/types/apiKeyService';
import { formatDate } from '@/utils/formatDate';
import { capitalize } from '@/utils/strings';

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

  const statusIconMap = {
    [APIKEY_STATUS.ACTIVE]: <ConfirmedIndicatorIcon />,
    [APIKEY_STATUS.INACTIVE]: <PendingIndicatorIcon />,
    [APIKEY_STATUS.REVOKED]: <RevokedIndicatorIcon />,
    [APIKEY_STATUS.UNUSED]: <UnusedIndicatorIcon />,
  };

  return (
    <div className="mt-[46px]  border-y-neutral-600">
      <Heading as={'h1'} className="flex items-center text-lg font-semibold capitalize">
        <Trans>API Keys</Trans>
      </Heading>
      <p className="my-4">
        <Trans>
          API keys enable you to securely connect your data sources to Ensign. Generate at least one
          API key for your project. You can customize permissions.
        </Trans>
      </p>

      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <div className="flex items-center gap-3"></div>
        <div>
          <Button
            variant="primary"
            size="small"
            className="!text-xs"
            onClick={onOpenGenerateAPIKeyModal}
          >
            + Add New Key
          </Button>
        </div>
      </div>
      <Table
        trClassName="text-sm"
        className="w-full"
        columns={[
          { Header: t`Key Name`, accessor: 'name' },
          { Header: t`Permissions`, accessor: 'permissions' },
          {
            Header: t`Status`,
            accessor: (key: { status: APIKeyStatus }) => {
              return (
                <div className="flex items-center">
                  {statusIconMap[key.status]}
                  <span className="ml-1">{capitalize(key.status)}</span>
                </div>
              );
            },
          },
          {
            Header: t`Last Used`,
            accessor: (date: any) => {
              return formatDate(new Date(date?.last_activity));
            },
          },
          {
            Header: t`Date Created`,
            accessor: (date: any) => {
              return formatDate(new Date(date?.created));
            },
          },
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
