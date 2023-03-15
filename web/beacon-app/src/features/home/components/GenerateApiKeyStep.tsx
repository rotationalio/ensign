import { Button } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import React, { useEffect, useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';
import GenerateAPIKeyModal from '@/features/apiKeys/components/GenerateAPIKeyModal';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { useOrgStore } from '@/store';

export default function GenerateApiKeyStep() {
  const org = useOrgStore.getState() as any;
  const { projectID } = org;
  const { apiKeys } = useFetchApiKeys(projectID);
  const [isOpenAPIKeyDataModal, setIsOpenAPIKeyDataModal] = useState<boolean>(false);
  const [isOpenGenerateAPIKeyModal, setIsOpenGenerateAPIKeyModal] = useState<boolean>(false);
  const [key, setKey] = useState<any>(null);

  const alreadyHasKeys = apiKeys?.api_keys?.length > 0;

  const onOpenGenerateAPIKeyModal = () => {
    //if (alreadyHasKeys) return;
    setIsOpenGenerateAPIKeyModal(true);
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

  return (
    <>
      <CardListItem title="Step 2: Generate API Key">
        <div className="mt-5 flex flex-col gap-8 px-3 xl:flex-row">
          <ErrorBoundary
            fallback={
              <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
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
            <div className="flex flex-col justify-between sm:w-1/5">
              <Button
                className="h-[44px] w-[165px] text-sm"
                onClick={onOpenGenerateAPIKeyModal}
                disabled={alreadyHasKeys}
                data-testid="key"
              >
                Create API Key
              </Button>
              {alreadyHasKeys && (
                <div className="mx-auto  py-2">
                  <HeavyCheckMark className="h-12 w-12" />
                </div>
              )}
            </div>

            <ApiKeyModal open={isOpenAPIKeyDataModal} data={key} onClose={onCloseAPIKeyDataModal} />
            <GenerateAPIKeyModal
              open={isOpenGenerateAPIKeyModal}
              onClose={onCloseGenerateAPIKeyModal}
              onSetKey={setKey}
            />
          </ErrorBoundary>
        </div>
      </CardListItem>
    </>
  );
}
