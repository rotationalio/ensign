import { Button, Modal } from '@rotational/beacon-core';

import { Close as CloseIcon } from '@/components/icons/close';
interface CancelAcctModalProps {
  close: () => void;
  isOpen: boolean;
}

export default function CancelAcctModal({ close, isOpen }: CancelAcctModalProps) {
  return (
    <Modal
      title="Cancel Account"
      open={isOpen}
      containerClassName="max-w-md"
      data-testid="cancelAcctModal"
    >
      <>
        <Button
          onClick={close}
          variant="ghost"
          className="bg-transparent absolute -right-10 top-5 border-none border-none p-2 p-2"
        >
          <CloseIcon />
        </Button>
        <p className="pb-4">
          Please contact us at <span className="font-bold">support@rotational.io</span> to cancel
          your account. Please include your name, email, and Org ID in your request to cancel your
          account. We are working on an automated process to cancel accounts and appreciate your
          patience.
        </p>
        <p className="pb-4">
          You are the <span className="font-bold">Owner</span> of this account. If you cancel your
          account, your Organization, Tenant, Project, and Topic and all associated data will be{' '}
          <span className="font-bold">permanently</span> deleted.
        </p>
      </>
    </Modal>
  );
}
