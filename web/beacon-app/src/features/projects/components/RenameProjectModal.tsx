import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import { useUpdateProject } from '../hooks/useUpdateProject';
import { Project } from '../types/Project';
import RenameProjectModalForm from './RenameProjectModalForm';

type ChangeRoleModalProps = {
  open: boolean;
  project: Project | null;
  handleModalClose: () => void;
};

// eslint-disable-next-line unused-imports/no-unused-vars
function RenameProjectModal({ open, handleModalClose, project }: ChangeRoleModalProps) {
  const { updateProject, wasProjectCreated } = useUpdateProject();

  useEffect(() => {
    if (wasProjectCreated) {
      toast.success('Success! You have renamed your project.');
    }
  }, [wasProjectCreated]);

  const handleSubmit = (values: any) => {
    const payload = {
      projectID: project?.id || '',
      projectPayload: {
        name: values['new-name'],
      },
    };
    updateProject(payload);
    handleModalClose();
  };

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
