import { t, Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';
import { useParams } from 'react-router-dom';

import { useCreateTopic } from '../hooks/useCreateTopic';
import NewTopicModalForm from './NewTopicModalForm';

export const NewTopicModal = ({
  open,
  handleClose,
}: {
  open: boolean;
  handleClose: () => void;
}) => {
  const { createTopic, wasTopicCreated, isCreatingTopic, hasTopicFailed, error, reset } =
    useCreateTopic();

  const param = useParams<{ id: string }>();
  const { id: projectID } = param;
  const projID = projectID || (projectID as string);

  const handleSubmitTopicForm = async (values: any) => {
    const payload = {
      ...values,
      projectID: projID,
    };
    await createTopic(payload);
  };

  useEffect(() => {
    if (wasTopicCreated) {
      toast.success(t`Success! Your have created a new topic.`);
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
        title={
          <h1>
            <Trans>New Topic</Trans>
          </h1>
        }
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
