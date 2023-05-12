import { t } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import { useUpdateProject } from '../../hooks/useUpdateProject';
import type { Project } from '../../types/Project';
import ChangeOwnerForm from './ChangeOwnerForm';

type ChangeRoleModalProps = {
  open: boolean;
  project: Project;
  handleModalClose: () => void;
};

// eslint-disable-next-line unused-imports/no-unused-vars
function ChangeOwnerModal({ open, handleModalClose, project }: ChangeRoleModalProps) {
  const { updateProject, wasProjectCreated } = useUpdateProject();

  useEffect(() => {
    if (wasProjectCreated) {
      toast.success(t`Success! You have renamed your project.`);
    }
  }, [wasProjectCreated]);

  const handleSubmit = (values: any) => {
    console.log(values);
    const payload = {
      ownerID: values['new_owner'],
      projectID: project?.id || '',
    } as any;
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
        <ChangeOwnerForm handleSubmit={handleSubmit} initialValues={project} />
      </>
    </Modal>
  );
}

export default ChangeOwnerModal;
