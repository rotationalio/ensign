import { t } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import { useUpdateProject } from '../../hooks/useUpdateProject';
import type { Project } from '../../types/Project';
import { UpdateProjectDTO } from '../../types/updateProjectService';
import EditProjectForm from './EditProjectForm';

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
      toast.success(t`Success! You have edited your project.`);
    }
  }, [wasProjectUpdated]);

  const handleSubmit = (values: any) => {
    const projectPayload = {} as UpdateProjectDTO['projectPayload'];
    if (values['name'] !== project.name && values['name'] !== '') {
      projectPayload['name'] = values['name'];
    }
    if (values['description'] !== project.description && values['description'] !== '') {
      projectPayload['description'] = values['description'];
    }
    const payload = {
      projectID: project?.id,
      projectPayload,
    };
    updateProject(payload);
    handleModalClose();
  };

  return (
    <Modal
      open={open}
      title="Edit Project"
      data-testid="edit-project-modal"
      data-cy="edit-project"
      onClose={handleModalClose}
    >
      <>
        <EditProjectForm handleSubmit={handleSubmit} project={project} />
      </>
    </Modal>
  );
}

export default RenameProjectModal;
