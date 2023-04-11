/* eslint-disable react-hooks/exhaustive-deps */
import { Breadcrumbs, Heading, Loader } from '@rotational/beacon-core';
import invariant from 'invariant';
import { lazy, Suspense, useCallback, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import AppLayout from '@/components/layout/AppLayout';
import BreadcrumbsIcon from '@/components/ui/Breadcrumbs/breadcrumbs-icon';

import { useFetchProject } from '../hooks/useFetchProject';

const ProjectDetail = lazy(() => import('../components/ProjectDetail'));
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

const ProjectDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams<{ id: string }>();
  const { id: projectID } = param;

  invariant(projectID, 'project id is required');

  const { project } = useFetchProject(projectID);
  // this below is added to fix the issue of navigating to the project detail page
  useEffect(() => {
    if (!param || !projectID) {
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [param, navigate, projectID]);

  const CustomBreadcrumbs = useCallback(() => {
    return (
      <Breadcrumbs separator="/" className="ml-4 hidden md:block">
        <Breadcrumbs.Item
          onClick={() => navigate(PATH_DASHBOARD.HOME)}
          className="capitalize hover:underline"
        >
          <BreadcrumbsIcon className="inline" /> Home
        </Breadcrumbs.Item>
        <Breadcrumbs.Item className="!cursor-default capitalize">Projects</Breadcrumbs.Item>
        {project?.name ? <Breadcrumbs.Item>{project.name}</Breadcrumbs.Item> : null}
      </Breadcrumbs>
    );
  }, [project?.name, project?.id]);

  return (
    <AppLayout Breadcrumbs={<CustomBreadcrumbs />}>
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        Project Details
      </Heading>
      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <ProjectDetail projectID={projectID} />
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
