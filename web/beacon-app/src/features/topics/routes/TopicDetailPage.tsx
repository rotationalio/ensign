import { Trans } from '@lingui/macro';
import { Heading, Loader } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import AppLayout from '@/components/layout/AppLayout';
import DetailTooltip from '@/components/ui/Tooltip/DetailTooltip';
import { useFetchProject } from '@/features/projects/hooks/useFetchProject';

import AdvancedTopicPolicy from '../components/AdvancedTopicPolicy';
import EventDetailTable from '../components/EventDetailTable';
import EventDetailTableHeader from '../components/EventDetailTableHeader';
import TopicQuery from '../components/TopicQuery';
// import TopicQuickView from '../components/TopicQuickView';
import TopicsBreadcrumbs from '../components/TopicsBreadcrumbs';
import TopicSettings from '../components/TopicSettings';
import TopicStateTag from '../components/TopicStateTag';
import { useFetchTopic } from '../hooks/useFetchTopic';
import { getFormattedTopicData } from '../utils';

const TopicDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams();
  const { id: topicID } = param as { id: string };
  const { topic, error, isFetchingTopic, wasTopicFetched } = useFetchTopic(topicID);
  const { project } = useFetchProject(topic?.project_id as string);

  // if user switch to another organization and topic is not found then
  // we need to redirect the user to the projects page
  useEffect(() => {
    if (error && error?.response?.status === 401) {
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
          {/* <TopicQuickView topicID={topicID} /> */}
          <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-4">
            <Heading as="h1" className="flex items-center gap-3 text-2xl font-semibold">
              <Trans>Topic:</Trans>
              <span className="mr-2" data-cy="topic-name">
                {topic?.topic_name}
                &nbsp;
                <DetailTooltip data={getFormattedTopicData(topic)} />
              </span>
              <TopicStateTag status={topic?.status} />
            </Heading>
            <TopicSettings />
          </div>
          <div className="mx-6">
            <EventDetailTableHeader />
            <EventDetailTable />
            <TopicQuery data={topic ?? []} />
            <AdvancedTopicPolicy />
          </div>
        </>
      )}
    </AppLayout>
  );
};

export default TopicDetailPage;
