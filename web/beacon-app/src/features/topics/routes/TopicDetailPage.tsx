import { Heading } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import AppLayout from '@/components/layout/AppLayout';
import DetailTooltip from '@/components/ui/Tooltip/DetailTooltip';

import AdvancedTopicPolicy from '../components/AdvancedTopicPolicy';
import EventDetailTable from '../components/EventDetailTable';
import TopicQuery from '../components/TopicQuery';
import TopicQuickView from '../components/TopicQuickView';
import TopicsBreadcrumbs from '../components/TopicsBreadcrumbs';
import TopicSettings from '../components/TopicSettings';
import { useFetchTopic } from '../hooks/useFetchTopic';
import { getFormattedTopicData } from '../utils';
const TopicDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams();
  const { id: topicID } = param as { id: string };
  const { topic, error } = useFetchTopic(topicID);

  // if user switch to another organization and topic is not found then
  // we need to redirect the user to the projects page
  useEffect(() => {
    if (error && error.response.status === 401) {
      navigate(PATH_DASHBOARD.PROJECTS);
    }
  }, [error, navigate]);

  return (
    <AppLayout Breadcrumbs={<TopicsBreadcrumbs topic={topic} />}>
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
        <Heading as="h1" className="flex items-center text-lg font-semibold">
          <span className="mr-2" data-cy="topic-name">
            {topic?.topic_name}
          </span>
          <DetailTooltip data={getFormattedTopicData(topic)} />
        </Heading>
        <TopicSettings />
      </div>
      <TopicQuickView topicID={topicID} />
      <EventDetailTable />
      <TopicQuery />
      <AdvancedTopicPolicy />
    </AppLayout>
  );
};

export default TopicDetailPage;
