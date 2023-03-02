import { Button, Toast } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import React, { useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';

import { CardListItem } from '@/components/common/CardListItem';
import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';
import { useCreateProjectAPIKey } from '@/features/apiKeys/hooks/useCreateApiKey';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { getRecentTenant } from '@/utils/formatData';
export default function GenerateApiKeyStep() {
  const { tenants } = useFetchTenants();
  const recentTenant = getRecentTenant(tenants);
  const { createProjectNewKey, key, wasKeyCreated, isCreatingKey, hasKeyFailed, error } =
    useCreateProjectAPIKey(recentTenant?.id);
  const [isOpen, setOpen] = useState(wasKeyCreated);
  const handleCreateKey = () => {
    console.log('handleCreateKey');
    createProjectNewKey(recentTenant?.id);
  };

  if (hasKeyFailed) {
    // TODO: create handle error abstraction
    const errorData = error?.response?.data;
    const errorMessage =
      errorData ||
      errorData?.error ||
      errorData?.message ||
      errorData?.error_description ||
      errorData?.error?.error;
    console.log('errorMessage', errorMessage);
    toast.error(errorMessage || 'Something went wrong');

    <Toast
      isOpen={hasKeyFailed}
      variant="danger"
      title="Something went wrong, please try again later."
      description={errorMessage}
    />;
  }

  const onClose = () => setOpen(false);

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
                onClick={handleCreateKey}
                isLoading={isCreatingKey}
                disabled={wasKeyCreated}
                data-testid="key"
              >
                Create API Key
              </Button>
              {wasKeyCreated && <HeavyCheckMark className="h-16 w-16" />}
            </div>

            <Toaster />
            <ApiKeyModal open={isOpen} data={key} onClose={onClose} />
          </ErrorBoundary>
        </div>
      </CardListItem>
    </>
  );
}
