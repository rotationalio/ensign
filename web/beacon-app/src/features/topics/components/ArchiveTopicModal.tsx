import { Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';

interface ArchiveTopicModalProps {
  close: () => void;
  isOpen: boolean;
}

const ArchiveTopicModal = ({ close, isOpen }: ArchiveTopicModalProps) => {
  return (
    <Modal title="Archive Topic" open={isOpen} onClose={close} containerClassName="max-w-md">
      <>
        <p className="pb-4">
          <Trans>
            Please contact us at <span className="font-bold">support@rotational.io</span> to archive
            your project. Please include your name, email, and topic name in your request to archive
            the topic. We promise there are real humans on the other end who will be ready to help.
            We're working on an automated process to archive topics and appreciate your patience.
          </Trans>
        </p>
        <p className="pb-4">
          <Trans>
            Archiving the topic means no more additional data will be written to the topic. The
            topic is “frozen”.
          </Trans>
        </p>
      </>
    </Modal>
  );
};

export default ArchiveTopicModal;
