/* eslint-disable react-hooks/exhaustive-deps */
import { Heading, Loader } from '@rotational/beacon-core';
import invariant from 'invariant';
import { lazy, Suspense, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import AppLayout from '@/components/layout/AppLayout';

import ProjectBreadcrumbs from '../components/ProjectBreadcrumbs';
import ProjectDetailTooltip from '../components/ProjectDetailTooltip';
import ProjectSetup from '../components/ProjectSetup';
import ProjectSettings from '../components/Settings';
import { useFetchProject } from '../hooks/useFetchProject';
import useProjectSetup from '../hooks/useProjectSetup';
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

const ProjectDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams<{ id: string }>();
  const { id: projectID } = param;

  invariant(projectID, 'project id is required');
  const { hasProject, hasTopics, hasApiKeys, warningMessage, hasAlreadySetup } =
    useProjectSetup(projectID);
  const { project } = useFetchProject(projectID);

  const getNormalizedProjectName = () => {
    return project?.name.split('-').join(' ');
  };

  useEffect(() => {
    if (!param || !projectID) {
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [param, navigate, projectID]);

  return (
    <AppLayout Breadcrumbs={<ProjectBreadcrumbs project={project} />}>
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
        <Heading as="h1" className="flex items-center text-lg font-semibold">
          <span className="mr-2 capitalize">{getNormalizedProjectName()}</span>
          <ProjectDetailTooltip data={project} />
        </Heading>
        <ProjectSettings />
      </div>
      {!hasAlreadySetup && (
        <ProjectSetup
          warningMessage={warningMessage}
          config={{
            isProjectCreated: hasProject,
            isAPIKeyCreated: hasApiKeys,
            isTopicCreated: hasTopics,
          }}
        />
      )}
      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <APIKeysTable projectID={projectID} />
      </Suspense>
      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <TopicTable />
      </Suspense>
    </AppLayout>
  );
};

export default ProjectDetailPage;
