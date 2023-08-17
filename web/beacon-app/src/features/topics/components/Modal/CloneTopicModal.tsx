import { t, Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';

interface CloneTopicModalProps {
  close: () => void;
  isOpen: boolean;
}

const CloneTopicModal = ({ close, isOpen }: CloneTopicModalProps) => {
  return (
    <Modal
      title={t`Clone Topic`}
      open={isOpen}
      onClose={close}
      containerClassName="max-w-md"
      data-cy="clone-topic-modal"
    >
      <>
        <p className="pb-4">
          <Trans>
            Please contact us at <span className="font-bold">support@rotational.io</span> to clone
            your topic. Please include your name, email, and topic name in your request to clone
            your topic. We promise there are real humans on the other end who will be ready to help.
            We're working on an automated process to clone topics and appreciate your patience.
          </Trans>
        </p>
        <p>
          <Trans>
            You can clone your topic to the current project, another existing project, or a new
            project. Note that topics must have unique names in each project.
          </Trans>
        </p>
      </>
    </Modal>
  );
};

export default CloneTopicModal;
