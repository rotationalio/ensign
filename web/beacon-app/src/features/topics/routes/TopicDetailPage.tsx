import { Heading } from '@rotational/beacon-core';
import invariant from 'invariant';
import { useParams } from 'react-router-dom';

import AppLayout from '@/components/layout/AppLayout';
import DetailTooltip from '@/components/ui/Tooltip/DetailTooltip';
import { formatDate } from '@/utils/formatDate';

import AdvancedTopicPolicy from '../components/AdvancedTopicPolicy';
import TopicQuery from '../components/TopicQuery';
import TopicsBreadcrumbs from '../components/TopicsBreadcrumbs';
import TopicSettings from '../components/TopicSettings';
import { useFetchTopic } from '../hooks/useFetchTopic';
import { t } from '@lingui/macro';

const TopicDetailPage = () => {
  const param = useParams();
  const { id: topicID } = param;
  const { topic } = useFetchTopic(topicID as string);
  const topicData = [
    {
      label: t`Topic ID`,
      value: topic?.id,
    },
    {
      label: t`Status`,
      value: topic?.status,
    },
    {
      label: t`Created`,
      value: formatDate(new Date(topic?.created as string)),
    },
    {
      label: t`Modified`,
      value: formatDate(new Date(topic?.modified as string)),
    },
  ];

  invariant(topicID, 'topic id is required');
  return (
    <AppLayout Breadcrumbs={<TopicsBreadcrumbs topic={topic} />}>
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
        <Heading as="h1" className="flex items-center text-lg font-semibold">
          <span className="mr-2" data-cy="topic-name">
            {topic?.topic_name}
          </span>
          <DetailTooltip data={topicData} />
        </Heading>
        <TopicSettings />
      </div>

      <TopicQuery />
      <AdvancedTopicPolicy />
    </AppLayout>
  );
};

export default TopicDetailPage;
