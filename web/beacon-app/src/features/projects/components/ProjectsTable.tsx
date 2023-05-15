import { t } from '@lingui/macro';
import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import React, { useCallback, useMemo, useState } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import { formatDate } from '@/utils/formatDate';

import { RenameProjectModal } from '../components/RenameProject';
import { Project } from '../types/Project';

type ProjectTableProps = {
  projects: Project[];
  isLoading?: boolean;
};

const ProjectsTable: React.FC<ProjectTableProps> = ({ projects, isLoading = false }) => {
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
          return description?.length > 50
            ? `${description?.slice(0, 50)}...`
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
        Header: t`Active Topics`,
        accessor: 'active_topics',
      },
      {
        Header: t`Data Storage`,
        accessor: (p: Project) => {
          const value = p?.data_storage?.value;
          const units = p?.data_storage?.units;
          return String(value) + units || '0GB';
        },
      },
      {
        Header: t`Owner`,
        accessor: (p: Project) => {
          const name = p?.owner?.name;
          const picture = p?.owner?.picture;
          return (
            <div className="flex gap-1.5">
              <img src={picture} alt="" className="h-6 w-6 rounded-2xl" />
              <div className="mt-0.5">{name}</div>
            </div>
          );
        },
      },
      {
        Header: t`Date Created`,
        accessor: (date: any) => {
          return formatDate(new Date(date?.created));
        },
      },
      { Header: 'Actions', accessor: 'actions' },
    ],
    []
  ) as any;

  const initialState = { hiddenColumns: ['id'] };

  const [openRenameProjectModal, setOpenRenameProjectModal] = useState<{
    open: boolean;
    project: Project;
  }>({
    open: false,
    project: {} as Project,
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

  const handleRedirection = (row: any) => {
    if (!row?.values?.id) {
      toast.error(
        t`Sorry, we are having trouble redirecting you to your project. Please try again.`
      );
    }
    navigate(`${PATH_DASHBOARD.PROJECTS}/${row?.values?.id}`);
  };

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
        <RenameProjectModal {...openRenameProjectModal} handleModalClose={handleModalClose} />
        <Table
          trClassName="text-sm hover:bg-gray-100"
          columns={initialColumns}
          initialState={initialState}
          data={getProjects(projects) || []}
          onRowClick={(row: any) => {
            handleRedirection(row);
          }}
          isLoading={isLoading}
          data-cy="projectTable"
        />
      </ErrorBoundary>
    </div>
  );
};

export default ProjectsTable;
