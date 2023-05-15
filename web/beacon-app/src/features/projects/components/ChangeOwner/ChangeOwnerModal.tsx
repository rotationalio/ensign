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
      title={t`Change Owner`}
      containerClassName="overflow-scroll  max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
      data-testid="prj-change-owner-modal"
      onClose={handleModalClose}
    >
      <>
        <ChangeOwnerForm handleSubmit={handleSubmit} initialValues={project} />
      </>
    </Modal>
  );
}

export default ChangeOwnerModal;
