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
  const { updateProject, wasProjectUpdated } = useUpdateProject();

  useEffect(() => {
    if (wasProjectUpdated) {
      toast.success(t`Success! Your project's owner has been updated.`);
    }
  }, [wasProjectUpdated]);

  const handleSubmit = (values: any) => {
    const payload = {
      projectID: project?.id || '',
      projectPayload: {
        owner: {
          id: values?.new_owner,
        },
      },
    };
    updateProject(payload);
    handleModalClose();
  };

  return (
    <Modal
      open={open}
      containerClassName="w-[25vw]"
      title={t`Change Owner`}
      data-testid="prj-change-owner-modal"
      data-cy="change-proj-owner"
      onClose={handleModalClose}
    >
      <>
        <ChangeOwnerForm handleSubmit={handleSubmit} initialValues={project} />
      </>
    </Modal>
  );
}

export default ChangeOwnerModal;
