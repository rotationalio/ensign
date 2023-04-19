import { Modal } from '@rotational/beacon-core';

import { Close } from '@/components/icons/close';
import Button from '@/components/ui/Button/Button';

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
    <div className="relative">
      <Modal
        open={open}
        title="Rename Project"
        containerClassName="overflow-scroll  max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="keyCreated"
      >
        <>
          <RenameProjectModalForm handleSubmit={handleSubmit} project={project} />
          <Button
            onClick={handleModalClose}
            variant="ghost"
            className="absolute top-2 right-2 min-h-fit min-w-fit py-2"
          >
            <Close className="text-primary-900" />
          </Button>
        </>
      </Modal>
    </div>
  );
}

export default RenameProjectModal;
