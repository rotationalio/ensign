import { t, Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import { useCreateTopic } from '../hooks/useCreateTopic';
import { useFetchTenantProjects } from '../hooks/useFetchTenantProjects';
import NewTopicModalForm from './NewTopicModalForm';

export const NewTopicModal = ({
  open,
  handleClose,
}: {
  open: boolean;
  handleClose: () => void;
}) => {
  const { tenants } = useFetchTenants();
  const { projects } = useFetchTenantProjects(tenants?.tenant[0]?.id);
  const projectID = projects?.project[0]?.id;

  const handleSubmitTopicForm = async (values: any) => {
    const payload = {
      ...values,
      projectID: projectID,
    };
    await createTopic(payload);
  };

  const { createTopic, wasTopicCreated, isCreatingTopic, hasTopicFailed, error, reset } =
    useCreateTopic();

  useEffect(() => {
    if (wasTopicCreated) {
      toast.success(t`Success! You have created a new topic.`);
      handleClose();
      reset();
    }
  }, [wasTopicCreated, handleClose, reset]);

  useEffect(() => {
    if (hasTopicFailed) {
      toast.error(
        (error as any)?.response?.data?.error ||
          t`Could not create topic. Please try again or contact support, if the problem continues.`
      );
      reset();
    }
  }, [hasTopicFailed, error, reset]);

  return (
    <>
      <Modal
        open={open}
        title={t`Name Topic`}
        containerClassName="max-h-[90vh] overflow-scroll max-w-[80vw] lg:max-w-[40vw] no-scrollbar"
        onClose={handleClose}
        data-testid="topicModal"
      >
        <>
          <p className="text-sm">
            <Trans>
              Each topic has a name that is unique across the tenant. Topic names are a combination
              of letters, numbers, underscores, or dashes. Topic names cannot have spaces or begin
              with an underscore or dash. Topic names are case insensitive.
            </Trans>
          </p>
          <p className="mt-2 text-sm">
            <Trans>Example topic name:</Trans> Fuzzy_Topic_Name-001
          </p>
          <NewTopicModalForm onSubmit={handleSubmitTopicForm} isSubmitting={isCreatingTopic} />
        </>
      </Modal>
    </>
  );
};
