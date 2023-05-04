import { t, Trans } from '@lingui/macro';
import { Button, Heading, Loader, Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';
import HintIcon from '@/components/icons/hint';
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
    <div className="mt-[46px]">
      <Heading as={'h1'} className="flex items-center text-lg font-semibold capitalize">
        <Trans>Topics</Trans>
      </Heading>
      <p className="my-4">
        <Trans>
          You must have at least one topic in your project to publish and subscribe. Topics are
          categories or logs that hold messages and events in a logical order, allowing services and
          data sources to send and receive data between them with ease and accuracy.
        </Trans>
        <span className="ml-2">
          <HelpTooltip
            content={
              <p>
                <Trans>
                  {' '}
                  Messages and events are sent to and read from specific topics. Services that are{' '}
                  {''}
                  <span className="font-bold">producers, write</span> data to topics. Services that
                  are <span className="font-bold">consumers, read</span> data from topics. Topics
                  are multi-subscriber, which means that a topic can have zero, one, or multiple
                  consumers subscribing to that topic, with read access to the log.
                </Trans>
              </p>
            }
          >
            <button>
              <HintIcon />
            </button>
          </HelpTooltip>
        </span>
      </p>
      <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
        <div className="flex items-center gap-3"></div>
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
