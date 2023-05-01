import { Button, Heading, Loader, Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';
import { Topic } from '@/features/topics/types/topicService';

import { NewTopicModal } from './NewTopicModal';

export const TopicTable = () => {
  const [openNewTopicModal, setOpenNewTopicModal] = useState(false);
  const handleOpenNewTopicModal = () => setOpenNewTopicModal(true);
  const handleCloseNewTopicModal = () => setOpenNewTopicModal(false);

  const { getTopics, topics, isFetchingTopics, hasTopicsFailed, error } = useFetchTopics();

  if (!topics) {
    getTopics();
  }

  if (isFetchingTopics) {
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
    <div className="my-5">
      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <Heading as="h1" className="text-lg font-semibold">
          Topics
        </Heading>
        <Button
          variant="primary"
          size="small"
          className="!text-xs"
          onClick={handleOpenNewTopicModal}
        >
          + New Topic
        </Button>
        <NewTopicModal open={openNewTopicModal} handleClose={handleCloseNewTopicModal} />
      </div>
      <div className="overflow-hidden text-sm">
        <Table
          trClassName="text-sm"
          columns={[
            { Header: 'Topics ID', accessor: 'id' },
            { Header: 'Name', accessor: 'name' },
          ]}
          data={(topics.topics as Topic[]) || []}
        />
      </div>
    </div>
  );
};

export default TopicTable;
