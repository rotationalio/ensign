/* eslint-disable react-hooks/exhaustive-deps */
import * as Tooltip from '@radix-ui/react-tooltip';
import { Breadcrumbs, Heading, Loader } from '@rotational/beacon-core';
import invariant from 'invariant';
import { lazy, Suspense, useCallback, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import HintIcon from '@/components/icons/hint';
import SettingIcon from '@/components/icons/setting';
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

  const getNormalizedProjectName = () => {
    return project?.name.split('-').join(' ');
  };

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
      <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
        <Heading as="h1" className="flex items-center text-lg font-semibold">
          <span className="mr-1 capitalize">{getNormalizedProjectName()}</span>
          <Tooltip.Provider>
            <Tooltip.Root>
              <Tooltip.Trigger asChild>
                <button className="">
                  <HintIcon />
                </button>
              </Tooltip.Trigger>
              <Tooltip.Portal>
                <Tooltip.Content
                  className="w-full max-w-[550px] rounded-md bg-[#2F4858] p-4 text-sm text-white"
                  sideOffset={5}
                  align="start"
                >
                  <table className="table-auto border-separate border-spacing-y-2">
                    <tbody>
                      <tr>
                        <td className="w-[150px] font-semibold">Project Status:</td>
                        <td>Inactive</td>
                      </tr>
                      <tr>
                        <td className="font-semibold">Description:</td>
                        <td>
                          Experiment to move from batch to stream processing for online learning
                          models on internal services
                        </td>
                      </tr>
                      <tr>
                        <td className="font-semibold">Owner:</td>
                        <td>Stephanie Kirby</td>
                      </tr>
                      <tr>
                        <td className="font-semibold">Created:</td>
                        <td>2022-Nov-21 15:35:02 GMT</td>
                      </tr>
                    </tbody>
                  </table>
                </Tooltip.Content>
              </Tooltip.Portal>
            </Tooltip.Root>
          </Tooltip.Provider>
        </Heading>
        <SettingIcon />
      </div>
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
