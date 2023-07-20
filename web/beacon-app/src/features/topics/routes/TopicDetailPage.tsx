import { Heading } from '@rotational/beacon-core';
import invariant from 'invariant';
import { useParams } from 'react-router-dom';

import AppLayout from '@/components/layout/AppLayout';

import AdvancedTopicPolicy from '../components/AdvancedTopicPolicy';
import TopicQuery from '../components/TopicQuery';
import TopicQuickView from '../components/TopicQuickView';
import TopicsBreadcrumbs from '../components/TopicsBreadcrumbs';
import TopicSettings from '../components/TopicSettings';
import { useFetchTopic } from '../hooks/useFetchTopic';
const TopicDetailPage = () => {
  const param = useParams();
  const { id: topicID } = param;
  const { topic } = useFetchTopic(topicID as string);

  invariant(topicID, 'topic id is required');
  return (
    <AppLayout Breadcrumbs={<TopicsBreadcrumbs topic={topic} />}>
      <TopicQuickView topicID={topicID} />
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
        <Heading as="h1" className="flex items-center text-lg font-semibold">
          <span className="mr-2" data-cy="topic-name">
            {topic?.topic_name}
          </span>
        </Heading>
        <TopicSettings />
      </div>

      <TopicQuery />
      <AdvancedTopicPolicy />
    </AppLayout>
  );
};

export default TopicDetailPage;
