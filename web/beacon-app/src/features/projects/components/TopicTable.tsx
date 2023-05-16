import { t, Trans } from '@lingui/macro';
import { Button, Heading, Loader, Table, Toast } from '@rotational/beacon-core';
import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';

import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';
import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';
import { Topic } from '@/features/topics/types/topicService';
import { formatDate } from '@/utils/formatDate';

import { getNormalizedDataStorage, getTopics } from '../util';
import { NewTopicModal } from './NewTopicModal';

export const TopicTable = () => {
  const initialColumns = useMemo(
    () => [
      { Header: t`Topic Name`, accessor: 'topic_name' },
      {
        Header: t`Publishers`,
        accessor: (t: Topic) => {
          const publishers = t?.publishers;
          return publishers || '---';
        },
      },
      {
        Header: t`Subscribers`,
        accessor: (t: Topic) => {
          const subscribers = t?.subscribers;
          return subscribers || '---';
        },
      },
      {
        Header: t`Data Storage`,
        accessor: (t: Topic) => {
          const value = t?.data_storage?.value;
          const units = t?.data_storage?.units;
          return getNormalizedDataStorage(value, units);
        },
      },
      { Header: t`Status`, accessor: 'status' },
      {
        Header: t`Date Created`,
        accessor: (date: any) => {
          return formatDate(new Date(date?.created));
        },
      },
    ],
    []
  ) as any;

  const [openNewTopicModal, setOpenNewTopicModal] = useState(false);
  const handleOpenNewTopicModal = () => setOpenNewTopicModal(true);
  const handleCloseNewTopicModal = () => setOpenNewTopicModal(false);

  const param = useParams<{ id: string }>();
  const { id: projectID } = param;
  const projID = projectID || (projectID as string);

  const { topics, isFetchingTopics, hasTopicsFailed, error } = useFetchTopics(projID);

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
    <div className="mt-[46px]" data-cy="topicComp">
      <Heading as={'h1'} className="flex items-center text-lg font-semibold capitalize">
        <Trans>Topics</Trans>
      </Heading>
      <p className="my-4">
        <Trans>
          You must have at least one topic in your project to publish and subscribe. Topics are
          categories or logs that hold messages and events in a logical order, allowing services and
          data sources to send and receive data between them with ease and accuracy.
        </Trans>
        <span className="ml-2" data-cy="topicHint">
          <HelpTooltip data-cy="topicInfo">
            <p>
              <Trans>
                {' '}
                Messages and events are sent to and read from specific topics. Services that are{' '}
                {''}
                <span className="font-bold">producers, write</span> data to topics. Services that
                are <span className="font-bold">consumers, read</span> data from topics. Topics are
                multi-subscriber, which means that a topic can have zero, one, or multiple consumers
                subscribing to that topic, with read access to the log.
              </Trans>
            </p>
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
