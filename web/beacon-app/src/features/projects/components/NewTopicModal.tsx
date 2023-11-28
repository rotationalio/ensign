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
        containerClassName="w-[25vw]"
        title={t`Create New Data Topic`}
        onClose={handleClose}
        data-testid="topicModal"
      >
        <>
          <p className="my-4">
            <Trans>
              A topic is a <span className="font-semibold">labeled</span> stream of information
              related to your use case. So you have to name your topic to create it. Topic names
              are:
            </Trans>
          </p>
          <ul className="mb-4 ml-5 list-outside list-disc">
            <li>
              <Trans>Unique and case insensitive</Trans>
            </li>
            <li>
              <Trans>A combination of letters, numbers, underscores, or dashes</Trans>
            </li>
            <li>
              <Trans>Cannot contain spaces or begin with a number, underscore, or dash</Trans>
            </li>
          </ul>
          <p className="mb-4 mt-2">
            <Trans>
              We recommend a naming convention that contains a data descriptor prefix and data type
              suffix. Examples:
            </Trans>
          </p>
          <ul className="mb-4 ml-5 list-outside list-disc">
            <li>instances-json</li>
            <li>hotels-avro</li>
            <li>flights-parquet</li>
            <li>weather-xml</li>
            <li>earthquake-csv</li>
          </ul>
          <NewTopicModalForm onSubmit={handleSubmitTopicForm} isSubmitting={isCreatingTopic} />
        </>
      </Modal>
    </>
  );
};
