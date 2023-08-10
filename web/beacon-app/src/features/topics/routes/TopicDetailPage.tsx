import { Heading, Loader } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import AppLayout from '@/components/layout/AppLayout';
import DetailTooltip from '@/components/ui/Tooltip/DetailTooltip';
import { useFetchProject } from '@/features/projects/hooks/useFetchProject';

import AdvancedTopicPolicy from '../components/AdvancedTopicPolicy';
import EventDetailTable from '../components/EventDetailTable';
import TopicQuery from '../components/TopicQuery';
import TopicQuickView from '../components/TopicQuickView';
import TopicsBreadcrumbs from '../components/TopicsBreadcrumbs';
import TopicSettings from '../components/TopicSettings';
import TopicStateTag from '../components/TopicStateTag';
import EventDetailTableHeader from '../components/TopicUsage';
import { useFetchTopic } from '../hooks/useFetchTopic';
import { getFormattedTopicData } from '../utils';

const TopicDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams();
  const { id: topicID } = param as { id: string };
  const { topic, error, isFetchingTopic, wasTopicFetched } = useFetchTopic(topicID);
  const { project } = useFetchProject(topic?.project_id as string);

  console.log('[] topic', topic);

  // if user switch to another organization and topic is not found then
  // we need to redirect the user to the projects page
  useEffect(() => {
    if (error && error.response.status === 401) {
      navigate(PATH_DASHBOARD.PROJECTS);
    }
  }, [error, navigate]);

  return (
    <AppLayout
      Breadcrumbs={
        <TopicsBreadcrumbs
          data={{
            project: project,
            topic: topic,
          }}
        />
      }
    >
      {isFetchingTopic && <Loader />}
      {topic && wasTopicFetched && (
        <>
          <TopicQuickView topicID={topicID} />
          <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
            <Heading as="h1" className="flex items-center text-2xl font-semibold">
              <span className="mr-2" data-cy="topic-name">
                {topic?.topic_name}
              </span>
              <DetailTooltip data={getFormattedTopicData(topic)} />
              <span className="ml-4 mb-1.5">
                <TopicStateTag status={topic?.status} />
              </span>
            </Heading>
            <TopicSettings />
          </div>
          <EventDetailTableHeader />
          <EventDetailTable />
          <TopicQuery data={topic ?? []} />
          <AdvancedTopicPolicy />
        </>
      )}
    </AppLayout>
  );
};

export default TopicDetailPage;
