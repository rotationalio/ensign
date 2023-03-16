import { Button } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import React, { useEffect, useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';
import GenerateAPIKeyModal from '@/features/apiKeys/components/GenerateAPIKeyModal';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { useFetchTenantProjects } from '@/features/projects/hooks/useFetchTenantProjects';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
export default function GenerateApiKeyStep() {
  const { tenants } = useFetchTenants();
  const { projects } = useFetchTenantProjects(tenants?.tenants[0]?.id);
  const { apiKeys } = useFetchApiKeys(projects?.tenant_projects[0]?.id);
  const [isOpenAPIKeyDataModal, setIsOpenAPIKeyDataModal] = useState<boolean>(false);
  const [isOpenGenerateAPIKeyModal, setIsOpenGenerateAPIKeyModal] = useState<boolean>(false);
  const [key, setKey] = useState<any>(null);
  const [hasAlreadyGeneratedKey, setHasAlreadyGeneratedKey] = useState<boolean>(false);

  useEffect(() => {
    if (apiKeys?.api_keys?.length > 0) {
      setHasAlreadyGeneratedKey(true);
    }
  }, [apiKeys?.api_keys.length]);

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
            <div className="lg:w-1/55 mr-8 grid w-full place-items-center gap-3 md:w-1/5">
              <Button
                className="h-[44px] w-[165px] text-sm"
                onClick={onOpenGenerateAPIKeyModal}
                disabled={hasAlreadyGeneratedKey}
                data-testid="key"
              >
                Create API Key
              </Button>
              {apiKeys?.api_keys?.length > 0 && (
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
