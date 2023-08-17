import { t, Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';

interface DeleteTopicModalProps {
  close: () => void;
  isOpen: boolean;
}

const DeleteTopicModal = ({ close, isOpen }: DeleteTopicModalProps) => {
  return (
    <Modal
      title={t`Delete Topic`}
      open={isOpen}
      onClose={close}
      containerClassName="max-w-md"
      data-cy="delete-topic-modal"
    >
      <>
        <p className="pb-4">
          <Trans>
            Please contact us at <span className="font-bold">support@rotational.io</span> to delete
            your topic. Please include your name, email, and topic name in your request to delete
            the topic. We promise there are real humans on the other end who will be ready to help.
            We're working on an automated process to delete topics and appreciate your patience.
          </Trans>
        </p>
        <p className="pb-4">
          <Trans>
            Deleting the topic will <span className="font-bold">permanently</span> destroy all data
            in the topic.
          </Trans>
        </p>
      </>
    </Modal>
  );
};

export default DeleteTopicModal;
