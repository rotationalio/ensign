import { Button } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import React, { useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';
import { Toast } from '@/components/ui/Toast';
import { useCreateProjectAPIKey } from '@/features/apiKeys/hooks/useCreateApiKey';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { APIKeyDTO } from '@/features/apiKeys/types/createApiKeyService';
import { useFetchPermissions } from '@/hooks/useFetchPermissions';
import { useOrgStore } from '@/store';

import GenerateAPIKeyModal from './GenerateAPIKeyModal';

export default function GenerateApiKeyStep() {
  const org = useOrgStore.getState() as any;
  const { projectID } = org;

  const { permissions } = useFetchPermissions();

  const { createProjectNewKey, key, wasKeyCreated, isCreatingKey, hasKeyFailed, error } =
    useCreateProjectAPIKey();
  const { apiKeys } = useFetchApiKeys(projectID);
  const [openAPIKeyDataModal, setOpenAPIKeyDataModal] = useState(false);

  const [openGenerateAPIKeyModal, setOpenGenerateAPIKeyModal] = useState<boolean>(false);
  // eslint-disable-next-line unused-imports/no-unused-vars
  const handleCreateKey = ({ name, permissions }: any) => {
    const payload = {
      projectID,
      name,
      permissions,
    } satisfies APIKeyDTO;

    createProjectNewKey(payload);
    if (wasKeyCreated) {
      setOpenAPIKeyDataModal(true);
    }

    // TODO: create handle error abstraction
  };

  const alreadyHasKeys = apiKeys?.api_keys?.length > 0;

  const onOpenGenerateAPIKeyModal = () => {
    setOpenGenerateAPIKeyModal(true);
  };
  const onCloseGenerateAPIKeyModal = () => {
    if (wasKeyCreated) {
      setOpenAPIKeyDataModal(true);
    }
    setOpenGenerateAPIKeyModal(false);
  };

  if (hasKeyFailed || error) {
    // TODO: create handle error abstraction
    // const errorData = error?.response?.data;
    // const errorMessage =
    //   errorData ||
    //   errorData?.error ||
    //   errorData?.message ||
    //   errorData?.error_description ||
    //   errorData?.error?.error;
    // console.log('errorMessage', errorMessage);

    <Toast
      isOpen={hasKeyFailed}
      variant="danger"
      description={(error as any)?.response?.data?.error || 'Something went wrong'}
    />;
  }

  const onCloseAPIKeyDataModal = () => setOpenAPIKeyDataModal(false);

  return (
    <>
      <CardListItem title="Step 2: Generate API Key">
        <div className="mt-5 flex flex-col gap-8 px-3 xl:flex-row">
          <ErrorBoundary
            fallback={
              <div className="item-center my-auto flex w-full justify-center text-center font-bold text-danger-500">
                <p>Sorry we are having trouble creating your API key, please try again.</p>
              </div>
            }
          >
            <p className="w-full text-sm sm:w-4/5">
              API keys enable you to securely connect your data sources to Ensign. Each key consists
              of two parts - a ClientID and a ClientSecret. You’ll need both to establish a client
              connection, create Ensign topics, publishers, and subscribers. Keep your API keys
              private -- if you misplace your keys, you can revoke them and generate new ones.
            </p>
            <div className="sm:w-1/5">
              <Button
                className="h-[44px] w-[165px] text-sm"
                onClick={onOpenGenerateAPIKeyModal}
                isLoading={isCreatingKey}
                disabled={alreadyHasKeys}
                data-testid="key"
              >
                Create API Key
              </Button>
              {alreadyHasKeys && <HeavyCheckMark className="h-16 w-16" />}
            </div>

            <ApiKeyModal open={openAPIKeyDataModal} data={key} onClose={onCloseAPIKeyDataModal} />
            <GenerateAPIKeyModal
              data={permissions}
              onSuccessfulCreate={wasKeyCreated}
              isLoading={isCreatingKey}
              open={openGenerateAPIKeyModal}
              onCloseModal={onCloseGenerateAPIKeyModal}
              onCreateNewKey={handleCreateKey}
              setOpenGenerateAPIKeyModal={setOpenGenerateAPIKeyModal}
            />
          </ErrorBoundary>
        </div>
      </CardListItem>
    </>
  );
}
