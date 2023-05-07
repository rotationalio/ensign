import { t, Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useCreateProject } from '@/features/projects/hooks/useCreateProject';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import NewProjectForm from './NewProjectForm';

type NewProjectModalProps = {
  isOpened: boolean;
  onClose: () => void;
};

function NewProjectModal({ onClose, isOpened }: NewProjectModalProps) {
  const navigateTo = useNavigate();
  const { createNewProject, isCreatingProject, wasProjectCreated, hasProjectFailed, error, reset } =
    useCreateProject();
  const { tenants } = useFetchTenants();
  const tenantID = tenants?.tenants[0]?.id;
  const handleSubmit = async (values: any) => {
    const payload = {
      ...values,
      tenantID: tenantID,
    };
    await createNewProject(payload);
  };

  useEffect(() => {
    if (wasProjectCreated) {
      toast.success(t`Success! You have created a new project.`);
      onClose();
      reset();
      navigateTo(`${PATH_DASHBOARD.PROJECTS}`);
    }
  }, [wasProjectCreated, onClose, reset, navigateTo]);

  useEffect(() => {
    if (hasProjectFailed) {
      toast.error(
        (error as any)?.response?.data?.error ||
          t`Could not create project. Please try again or contact support, if the problem continues.`
      );
      reset();
    }
  }, [hasProjectFailed, error, reset]);

  return (
    <>
      <Modal
        open={isOpened}
        onClose={onClose}
        containerClassName="max-h-[90vh] overflow-scroll max-w-[80vw] lg:max-w-[40vw] no-scrollbar"
        title={t`Create Project`}
        data-testid="newProjectModal"
      >
        <>
          <p>
            <Trans>
              Name your project and provide a description to get started. Next, you’ll create topics
              and finally API keys.
            </Trans>
          </p>
          <NewProjectForm onSubmit={handleSubmit} isSubmitting={isCreatingProject} />
        </>
      </Modal>
    </>
  );
}

export default NewProjectModal;
