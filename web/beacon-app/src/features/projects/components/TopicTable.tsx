import { Trans } from '@lingui/macro';
import { Button, Loader, Table, Toast } from '@rotational/beacon-core';
import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';

import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';

import { getInitialColumns, getTopics } from '../util';
import { NewTopicModal } from './NewTopicModal';
import TopicTableHeader from './TopicTableHeader';

export const TopicTable = () => {
  const initialColumns = useMemo(() => getInitialColumns(), []) as any;

  const [openNewTopicModal, setOpenNewTopicModal] = useState(false);
  const handleOpenNewTopicModal = () => setOpenNewTopicModal(true);
  const handleCloseNewTopicModal = () => setOpenNewTopicModal(false);

  const param = useParams<{ id: string }>();
  const { id: projectID } = param;
  const projID = projectID || (projectID as string);

  const { topics, isFetchingTopics, hasTopicsFailed, error } = useFetchTopics(projID);

  console.log('topics data', topics); // do not remove this line. it is used for debugging since the topic id is not available in the UI

  if (isFetchingTopics) {
    //
    // TODO: add loading state
    return <Loader />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasTopicsFailed}
        duration={3000}
        variant="danger"
        title="Sorry we are having trouble fetching your topics, please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  return (
    <div className="mt-[46px]" data-cy="topicComp">
      <TopicTableHeader />
      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <div className="flex items-center gap-3"></div>
        <Button
          variant="primary"
          size="small"
          className="!text-xs"
          onClick={handleOpenNewTopicModal}
          data-cy="addTopic"
        >
          <Trans>+ New Topic</Trans>
        </Button>
        <NewTopicModal open={openNewTopicModal} handleClose={handleCloseNewTopicModal} />
      </div>
      <div className="overflow-hidden text-sm">
        <Table
          trClassName="text-sm"
          columns={initialColumns}
          data={getTopics(topics)}
          data-cy="topicTable"
        />
      </div>
    </div>
  );
};

export default TopicTable;
