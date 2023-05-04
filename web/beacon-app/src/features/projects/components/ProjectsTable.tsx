import { t } from '@lingui/macro';
import { Table } from '@rotational/beacon-core';
// import { ErrorBoundary } from '@sentry/react';
import React, { useCallback, useMemo } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import { formatDate } from '@/utils/formatDate';

import { Project } from '../types/Project';

type ProjectTableProps = {
  projects: Project[];
};

const ProjectsTable: React.FC<ProjectTableProps> = ({ projects }) => {
  const navigate = useNavigate();
  const initialColumns = useMemo(
    () => [
      { Header: t`Project ID`, accessor: 'id' },
      { Header: t`Project Name`, accessor: 'name' },
      {
        Header: t`Description`,
        accessor: (p: Project) => {
          const description = p?.description;
          if (!description) {
            return '---';
          }
          // cut off description at 100 characters
          return description?.length > 100
            ? `${description?.slice(0, 100)}...`
            : description || '---';
        },
      },
      {
        Header: 'Status',
        accessor: (p: Project) => {
          const status = p?.status;
          return status || '---';
        },
      },
      {
        Header: t`Date Created`,
        accessor: (date: any) => {
          return formatDate(new Date(date?.created));
        },
      },
      // { Header: 'Actions', accessor: 'actions' },
    ],
    []
  ) as any;

  const initialState = { hiddenColumns: ['id'] };

  const handleRenameProjectClick = (projectId: string) => {
    console.log('clicked!', projectId);
  };

  const handleChangeOwnerClick = (_projectId: string) => {};

  const getProjects = useCallback((projects: Project[]) => {
    return (projects || []).map((project: Project) => ({
      ...project,
      actions: [
        { label: 'Rename project', onClick: () => handleRenameProjectClick(project?.id) },
        { label: 'Change owner', onClick: () => handleChangeOwnerClick(project?.id) },
      ],
    }));
  }, []);

  const handleRedirection = (row: any) => {
    if (!row?.values?.id) {
      // add toast here
      toast.error(
        t`Sorry, we are having trouble redirecting you to your project. Please try again.`
      );
    }
    navigate(`${PATH_DASHBOARD.PROJECTS}/${row?.values?.id}`);
  };

  return (
    <div className="mx-4">
      {/* <ErrorBoundary
        fallback={
          <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
            <p>
              Sorry we are having trouble listing your projects, please refresh the page and try
              again.
            </p>
          </div>
        }
      > */}
      <Table
        trClassName="text-sm"
        columns={initialColumns}
        initialState={initialState}
        data={getProjects(projects) || []}
        onRowClick={(row: any) => {
          handleRedirection(row);
        }}
      />
      {/* </ErrorBoundary> */}
    </div>
  );
};

export default ProjectsTable;
