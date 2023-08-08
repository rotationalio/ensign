/* eslint-disable react-hooks/exhaustive-deps */
import { Heading, Loader } from '@rotational/beacon-core';
import invariant from 'invariant';
import { lazy, Suspense, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import AppLayout from '@/components/layout/AppLayout';
import DetailTooltip from '@/components/ui/Tooltip/DetailTooltip';

import ProjectActive from '../components/ProjectActive';
import ProjectBreadcrumbs from '../components/ProjectBreadcrumbs';
import ProjectSetup from '../components/ProjectSetup';
import ProjectSettings from '../components/Settings';
import { useFetchProject } from '../hooks/useFetchProject';
import useProjectSetup from '../hooks/useProjectSetup';
import { getFormattedProjectData } from '../util';
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));
import useProjectActive from '../hooks/useProjectActive';
const ProjectDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams<{ id: string }>();
  const { id: projectID } = param;
  const { isActive, setIsActive } = useProjectActive(projectID as string);
  invariant(projectID, 'project id is required');
  const { hasProject, hasTopics, hasApiKeys, warningMessage, hasAlreadySetup } =
    useProjectSetup(projectID);
  const { project, error } = useFetchProject(projectID);

  const getNormalizedProjectName = () => {
    return project?.name?.split('-').join(' ');
  };

  useEffect(() => {
    // when user switch to another organization and the current project is not found then redirect to projects page
    if (error?.response?.status === 401) {
      navigate(PATH_DASHBOARD.PROJECTS);
    }
  }, [error]);

  return (
    <AppLayout Breadcrumbs={<ProjectBreadcrumbs project={project} />}>
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
        <Heading as="h1" className="flex items-center text-lg font-semibold">
          <span className="mr-2" data-cy="project-name">
            {getNormalizedProjectName()}
          </span>
          <DetailTooltip data={getFormattedProjectData(project)} />
        </Heading>
        <ProjectSettings data={project} />
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
      {!isActive && hasAlreadySetup && (
        <ProjectActive onActive={setIsActive} projectID={projectID} />
      )}

      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <TopicTable />
      </Suspense>
      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <APIKeysTable projectID={projectID} />
      </Suspense>
    </AppLayout>
  );
};

export default ProjectDetailPage;
