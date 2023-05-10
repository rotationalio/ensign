import { Modal } from '@rotational/beacon-core';

import { Project } from '../types/Project';
import RenameProjectModalForm from './RenameProjectModalForm';

type ChangeRoleModalProps = {
  open: boolean;
  project: Project | null;
  handleModalClose: () => void;
};

// eslint-disable-next-line unused-imports/no-unused-vars
function RenameProjectModal({ open, handleModalClose, project }: ChangeRoleModalProps) {
  const handleSubmit = () => {};

  return (
    <Modal
      open={open}
      title="Rename Project"
      containerClassName="overflow-scroll  max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
      data-testid="keyCreated"
      onClose={handleModalClose}
    >
      <>
        <RenameProjectModalForm handleSubmit={handleSubmit} project={project} />
      </>
    </Modal>
  );
}

export default RenameProjectModal;
