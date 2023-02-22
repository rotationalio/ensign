import { t, Trans } from '@lingui/macro';
import { Button, Toast } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import React, { useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';

import { CardListItem } from '@/components/common/CardListItem';
import { ApiKeyModal } from '@/components/common/Modal/ApiKeyModal';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';
import { useCreateAPIKey } from '@/features/apiKeys/hooks/useCreateApiKey';

export default function GenerateApiKeyStep() {
  const { createNewKey, key, wasKeyCreated, isCreatingKey, hasKeyFailed, error } =
    useCreateAPIKey();
  const [isOpen, setOpen] = useState(wasKeyCreated);
  const handleCreateKey = () => {
    createNewKey();
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

    toast.error(errorMessage || t`Something went wrong`);

    <Toast
      isOpen={hasKeyFailed}
      variant="danger"
      title={t`Something went wrong, please try again later.`}
      description={errorMessage}
    />;
  }

  const onClose = () => setOpen(false);

  return (
    <>
      <CardListItem title={t`Step 2: Generate API Key`}>
        <div className="mt-5 flex flex-col gap-8 px-3 md:flex-row">
          <ErrorBoundary
            fallback={
              <div className="item-center my-auto flex w-full justify-center text-center font-bold text-danger-500">
                <p>
                  <Trans>
                    Sorry we are having trouble creating your API key, please try again.
                  </Trans>
                </p>
              </div>
            }
          >
            <p className="w-full sm:w-4/5">
              <Trans>
                API keys enable you to securely connect your data sources to Ensign. Each key
                consists of two parts - a ClientID and a ClientSecret. You’ll need both to establish
                a client connection, create Ensign topics, publishers, and subscribers. Keep your
                API keys private -- if you misplace your keys, you can revoke them and generate new
                ones.
              </Trans>
            </p>
            <div className="sm:w-1/5">
              <Button
                className="h-[44px] w-[165px] text-sm"
                onClick={handleCreateKey}
                isLoading={isCreatingKey}
                disabled={wasKeyCreated}
              >
                <Trans>Create API Key</Trans>
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
