/* eslint-disable react-hooks/exhaustive-deps */
import { Breadcrumbs, Heading, Loader } from '@rotational/beacon-core';
import invariant from 'invariant';
import { lazy, Suspense, useCallback, useEffect } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import AppLayout from '@/components/layout/AppLayout';
import BreadcrumbsIcon from '@/components/ui/Breadcrumbs/breadcrumbs-icon';

import { useFetchProject } from '../hooks/useFetchProject';
import { formatProjectData } from '../util';
const ProjectDetail = lazy(() => import('../components/ProjectDetail'));
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

const ProjectDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams<{ id: string }>();

  invariant(param.id, 'id is required');

  const { project } = useFetchProject(param.id);
  const getFormattedProjectData = formatProjectData(project);
  console.log(project);
  useEffect(() => {
    if (!param || param.id === 'undefined' || param.id === 'null') {
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [param, navigate]);

  const CustomBreadcrumbs = useCallback(() => {
    return (
      <Breadcrumbs separator="/" className="ml-4 hidden md:block">
        <Breadcrumbs.Item className="capitalize">
          <Link to={PATH_DASHBOARD.HOME} className="hover:underline">
            <BreadcrumbsIcon className="inline" /> Home
          </Link>
        </Breadcrumbs.Item>
        <Breadcrumbs.Item className="!cursor-default capitalize">Projects</Breadcrumbs.Item>
        {project?.name ? <Breadcrumbs.Item>{project?.name}</Breadcrumbs.Item> : null}
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
        <ProjectDetail project={getFormattedProjectData} />
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
        <APIKeysTable projectID={param.id} />
      </Suspense>
    </AppLayout>
  );
};

export default ProjectDetailPage;
