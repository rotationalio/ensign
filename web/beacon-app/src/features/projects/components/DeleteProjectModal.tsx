import { Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
interface DeleteProjectModalProps {
  close: () => void;
  isOpen: boolean;
}

export default function DeleteProjectModal({ close, isOpen }: DeleteProjectModalProps) {
  return (
    <Modal
      title="Delete Project"
      open={isOpen}
      onClose={close}
      containerClassName="max-w-md"
      data-testid="delete-prj-modal"
    >
      <>
        <p className="pb-4">
          <Trans>
            Please contact us at <span className="font-bold">support@rotational.io</span>
            to delete your project. Please include your name, email, and project name in your
            request to delete the project. We promise there are real humans on the other end who
            will be ready to help. We’re working on an automated process to delete and appreciate
            your patience.
          </Trans>
        </p>
        <p className="pb-4">
          <Trans>
            Please note that deleting the project will{' '}
            <span className="font-bold">permanently</span>
            delete the project, API keys, topics, and all data associated with the project.
          </Trans>
        </p>
      </>
    </Modal>
  );
}
