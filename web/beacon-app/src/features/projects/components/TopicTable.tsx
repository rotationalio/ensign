import { t, Trans } from '@lingui/macro';
import { Heading, Loader, Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import Button from '@/components/ui/Button/Button';
import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';
import { Topic } from '@/features/topics/types/topicService';
import { formatDate } from '@/utils/formatDate';

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
          <Trans>Topics</Trans>
        </Heading>
        <Button
          variant="primary"
          size="small"
          className="!text-xs"
          onClick={handleOpenNewTopicModal}
        >
          <Trans>+ New Topic</Trans>
        </Button>
        <NewTopicModal open={openNewTopicModal} handleClose={handleCloseNewTopicModal} />
      </div>
      <div className="overflow-hidden text-sm">
        <Table
          trClassName="text-sm"
          columns={[
            { Header: t`Topic Name`, accessor: 'name' },
            { Header: t`Status`, accessor: 'status' },
            { Header: t`Publishers`, accessor: 'publishers' },
            { Header: t`Subscribers`, accessor: 'subscribers' },
            { Header: t`Data Storage`, accessor: 'data' },
            {
              Header: t`Date Created`,
              accessor: (date: any) => {
                return formatDate(new Date(date?.created));
              },
            },
          ]}
          data={(topics.topics as Topic[]) || []}
        />
      </div>
    </div>
  );
};

export default TopicTable;
