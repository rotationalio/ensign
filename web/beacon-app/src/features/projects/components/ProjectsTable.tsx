import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useCallback, useState } from 'react';

import { formatDate } from '@/utils/formatDate';

import { Project } from '../types/Project';
import RenameProjectModal from './RenameProjectModal';

type ProjectTableProps = {
  projects: Project[];
};

const initialColumns = [
  { Header: 'Project ID', accessor: 'id' },
  { Header: 'Project Name', accessor: 'name' },
  {
    Header: 'Status',
    accessor: () => {
      return <p className="text-center">-</p>;
    },
  },
  {
    Header: 'Active Topics',
    accessor: () => {
      return <p className="text-center">-</p>;
    },
  },
  {
    Header: 'Data Storage',
    accessor: () => {
      return <p className="text-center">-</p>;
    },
  },
  {
    Header: 'Owner',
    accessor: () => {
      return <p className="text-center">-</p>;
    },
  },
  {
    Header: 'Date Created',
    accessor: (date: any) => {
      return formatDate(new Date(date?.created));
    },
  },
  { Header: 'Actions', accessor: 'actions' },
];

function ProjectsTable({ projects }: ProjectTableProps) {
  const [openRenameProjectModal, setOpenRenameProjectModal] = useState<{
    open: boolean;
    project: Project | null;
  }>({
    open: false,
    project: null,
  });

  const handleRenameProjectClick = (project: Project) => {
    setOpenRenameProjectModal({ open: true, project });
  };

  const handleChangeOwnerClick = (_projectId: string) => {};

  const getProjects = useCallback((projects: Project[]) => {
    return (projects || []).map((project: Project) => ({
      ...project,
      actions: [
        { label: 'Rename project', onClick: () => handleRenameProjectClick(project) },
        { label: 'Change owner', onClick: () => handleChangeOwnerClick(project?.id) },
      ],
    }));
  }, []);

  const handleModalClose = () =>
    setOpenRenameProjectModal({ ...openRenameProjectModal, open: false });

  return (
    <div className="mx-4">
      <ErrorBoundary
        fallback={
          <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
            <p>
              Sorry we are having trouble listing your members, please refresh the page and try
              again.
            </p>
          </div>
        }
      >
        <Table trClassName="text-sm" columns={initialColumns} data={getProjects(projects) || []} />
        <RenameProjectModal {...openRenameProjectModal} handleModalClose={handleModalClose} />
      </ErrorBoundary>
    </div>
  );
}

export default ProjectsTable;
