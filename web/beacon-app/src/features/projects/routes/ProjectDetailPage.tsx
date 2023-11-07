/* eslint-disable react-hooks/exhaustive-deps */

import invariant from 'invariant';
import { lazy, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import AppLayout from '@/components/layout/AppLayout';

import ProjectActive from '../components/ProjectActive';
import ProjectBreadcrumbs from '../components/ProjectBreadcrumbs';
import ProjectDetailInfo from '../components/ProjectDetail/DetailInfo';
import ProjectDetailHeader from '../components/ProjectDetail/ProjectDetailHeader';
import ProjectSetup from '../components/ProjectSetup';
import ProjectSettings from '../components/Settings';
import { useFetchProject } from '../hooks/useFetchProject';
import useProjectSetup from '../hooks/useProjectSetup';

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

  useEffect(() => {
    // when user switch to another organization and the current project is not found then redirect to projects page
    if (error?.response?.status === 401) {
      navigate(PATH_DASHBOARD.PROJECTS);
    }
  }, [error]);

  return (
    <AppLayout Breadcrumbs={<ProjectBreadcrumbs project={project} />}>
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-4">
        <ProjectDetailHeader data={project} />
        <ProjectSettings data={project} />
      </div>
      <div className="mx-6">
        <ProjectDetailInfo data={project} />
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
        <TopicTable />
        <APIKeysTable projectID={projectID} />
      </div>
    </AppLayout>
  );
};

export default ProjectDetailPage;
