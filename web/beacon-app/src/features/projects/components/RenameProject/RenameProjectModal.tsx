import { t } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import { useUpdateProject } from '../../hooks/useUpdateProject';
import type { Project } from '../../types/Project';
import RenameProjectForm from './RenameProjectForm';

type ChangeRoleModalProps = {
  open: boolean;
  project: Project;
  handleModalClose: () => void;
};

// eslint-disable-next-line unused-imports/no-unused-vars
function RenameProjectModal({ open, handleModalClose, project }: ChangeRoleModalProps) {
  const { updateProject, wasProjectUpdated } = useUpdateProject();

  useEffect(() => {
    if (wasProjectUpdated) {
      toast.success(t`Success! You have renamed your project.`);
    }
  }, [wasProjectUpdated]);

  const handleSubmit = (values: any) => {
    const payload = {
      projectID: project?.id || '',
      projectPayload: {
        name: values['name'],
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
      data-testid="rename-project-modal"
      onClose={handleModalClose}
    >
      <>
        <RenameProjectForm handleSubmit={handleSubmit} project={project} />
      </>
    </Modal>
  );
}

export default RenameProjectModal;
