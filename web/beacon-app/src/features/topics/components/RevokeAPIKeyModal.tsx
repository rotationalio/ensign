import { t, Trans } from '@lingui/macro';
import { Button, Checkbox, Modal } from '@rotational/beacon-core';
import { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import styled from 'styled-components';

import { useDeleteAPIKey } from '@/features/apiKeys/hooks/useDeleteApiKey';
import { APIKey } from '@/features/apiKeys/types/apiKeyService';
type RevokeAPIKeyModalProps = {
  onOpen: {
    opened: boolean;
    key?: APIKey;
  };
  onClose: () => void;
};

// TODO: This component needs to be refactored

const RevokeAPIKeyModal = ({ onOpen, onClose }: RevokeAPIKeyModalProps) => {
  const { opened, key } = onOpen;
  const { deleteApiKey, wasKeyDeleted, hasKeyDeletedFailed, reset, isDeletingKey, error } =
    useDeleteAPIKey(key?.id || '');

  const [isChecked, setIsChecked] = useState(false);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const handleCheckboxChange = () => {
    setIsChecked(!isChecked);
  };

  const deleteAPIKeyHandler = () => {
    deleteApiKey();
  };

  useEffect(() => {
    if (hasKeyDeletedFailed) {
      reset();
      handleCheckboxChange();
      onClose();

      toast.error(
        error?.response?.data?.error ||
          t`Sorry, but we had a problem while trying to revoke the API key. Please try again. If the issue persists, contact our support team for assistance.`
      );
    }
  }, [hasKeyDeletedFailed, onClose, reset, handleCheckboxChange, error]);

  useEffect(() => {
    if (wasKeyDeleted) {
      reset();
      handleCheckboxChange();
      onClose();
      toast.success(t`API Key was successfully deleted.`);
      onClose();
    }
  }, [wasKeyDeleted, onClose, reset, handleCheckboxChange]);

  return (
    <Modal title={t`Revoke API Key`} open={opened} onClose={onClose} containerClassName="max-w-md">
      <>
        <p className="pb-4">
          <Trans>
            Revoking the API key will result in producers and consumers connected to the topic to{' '}
            <span className="font-bold">permanently</span> lose access to the topic. To maintain
            access to the topic, generate a new API key and update your publishers and subscribers.
          </Trans>
        </p>
        <p className="pb-4">
          <Trans>Check the box to revoke the API key.</Trans>
        </p>
        <div className="pb-6">
          <span className="font-bold">Key Name:</span> {key?.name}
        </div>
        <CheckboxFieldset onClick={handleCheckboxChange} className="pb-8">
          <Checkbox>
            <Trans>
              I understand that revoking the API key will cause publishers and subscribers to lose
              access to the topic and may impact performance.
            </Trans>
          </Checkbox>
        </CheckboxFieldset>
        <div className="mx-auto w-[150px] pb-4">
          <Button
            variant="secondary"
            disabled={!isChecked}
            onClick={deleteAPIKeyHandler}
            isLoading={isDeletingKey}
          >
            <Trans>Revoke API Key</Trans>
          </Button>
        </div>
        <div className="mx-auto w-[150px] pb-4">
          <Button variant="ghost" onClick={onClose} className="w-[130px] bg-[#000000B2] text-white">
            <Trans>Cancel</Trans>
          </Button>
        </div>
      </>
    </Modal>
  );
};

const CheckboxFieldset = styled.fieldset`
  label svg {
    min-width: 23px;
  }
`;

export default RevokeAPIKeyModal;
