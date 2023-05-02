import { t } from '@lingui/macro';
import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useCallback } from 'react';

import { formatDate } from '@/utils/formatDate';

import { Project } from '../types/Project';

type ProjectTableProps = {
  projects: Project[];
};

const initialColumns = [
  { Header: t`Project Name`, accessor: 'name' },
  {
    Header: t`Description`,
    accessor: (p: Project) => {
      const description = p?.description;
      if (!description) {
        return '---';
      }
      // cut off description at 100 characters
      return description?.length > 100 ? `${description?.slice(0, 100)}...` : description || '---';
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
];

function ProjectsTable({ projects }: ProjectTableProps) {
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

  return (
    <div className="mx-4">
      <ErrorBoundary
        fallback={
          <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
            <p>
              Sorry we are having trouble listing your projects, please refresh the page and try
              again.
            </p>
          </div>
        }
      >
        <Table trClassName="text-sm" columns={initialColumns} data={getProjects(projects) || []} />
      </ErrorBoundary>
    </div>
  );
}

export default ProjectsTable;
