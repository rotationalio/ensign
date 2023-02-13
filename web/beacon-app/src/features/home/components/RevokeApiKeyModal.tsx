import { AriaButton as Button, Checkbox, Modal } from '@rotational/beacon-core';
import { useState } from 'react';
import styled from 'styled-components';

import { Close } from '@/components/icons/close';

export type RevokeApiKeyModalProps = {
  open: boolean;
};

export default function RevokeApiKeyModal({ open }: RevokeApiKeyModalProps) {
  const [isOpen, setIsOpen] = useState(open);

  const closeModal = () => setIsOpen(false);

  return (
    <>
      <Modal open={isOpen} title="Revoke API Key" size="medium">
        <div className="gap-3">
          <Close onClick={closeModal} className="absolute top-6 right-8"></Close>
          <p className="my-3">
            Revoking the API key will result in producers and consumers connected to the topic to
            permanently lose access to the topic. To maintain access to the topic, generate a new
            API key and update your publishers and subscribers.
          </p>
          <p className="mb-2 mt-5">Check the box to revoke the API key.</p>

          <CheckboxFieldset>
            <Checkbox>
              I understand that revoking the API key will cause publishers and subscribers to lose
              access to the topic and may impact performance.
            </Checkbox>
          </CheckboxFieldset>
          <div className="my-5 text-center">
            <Button size="large" className="bg-[#DB3B00]">
              Take the reins
            </Button>
          </div>
        </div>
      </Modal>
    </>
  );
}

const CheckboxFieldset = styled.fieldset`
  label svg {
    min-width: 23px;
  }
`;
