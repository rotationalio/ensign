import { Button, Modal } from '@rotational/beacon-core';

import { Close as CloseIcon } from '@/components/icons/close';

interface DeleteOrgModalProps {
  close: () => void;
  isOpen: boolean;
}

export default function DeleteOrgModal({ close, isOpen }: DeleteOrgModalProps) {
  return (
    <div>
      <Modal title="Delete Organization" open={isOpen} containerClassName="relative max-w-md">
        <>
          <Button
            onClick={close}
            variant="ghost"
            className="absolute top-4 right-4 border-none py-2"
          >
            <CloseIcon />
          </Button>

          <p className="pb-4">
            Please contact us at <span className="font-bold">support@rotational.io</span> to delete
            your organization. Please include your name, email, and Org ID in your request to delete
            the organization. We are working on an automated process to delete organizations and
            appreciate your patience.
          </p>
          <p className="pb-4">
            Please note that deleting the organization will{' '}
            <span className="font-bold">permanently</span> delete the Tenant, Project, Users, Topics
            and all data associated with the organization.
          </p>
        </>
      </Modal>
    </div>
  );
}
