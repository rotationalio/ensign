import { Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';

import NewTopicModalForm from './NewTopicModalForm';

export const NewTopicModal = ({
  open,
  handleClose,
}: {
  open: boolean;
  handleClose: () => void;
}) => {
  const handleSubmitTopicForm = () => {};

  return (
    <>
      <Modal
        open={open}
        title={
          <h1>
            <Trans>New Topic</Trans>
          </h1>
        }
        containerClassName="max-h-[90vh] overflow-scroll max-w-[80vw] lg:max-w-[40vw] no-scrollbar"
        onClose={handleClose}
        data-testid="keyModal"
      >
        <>
          <p className="text-sm">
            <Trans>
              Each topic has a name that is unique across the tenant. Topic names are a combination
              of letters, numbers, underscores, or dashes. Topic names cannot have spaces or begin
              with an underscore or dash. Topics names are case insensitive.
            </Trans>
          </p>
          <p className="mt-2 text-sm">
            <Trans>Example topic name:</Trans> Fuzzy_Topic_Name-001
          </p>
          <NewTopicModalForm handleSubmit={handleSubmitTopicForm} />
        </>
      </Modal>
    </>
  );
};
