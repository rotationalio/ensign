import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';

function ProjectsTable() {
  const initialColumns = [
    { Header: 'Project ID', accessor: 'id' },
    { Header: 'Project Name', accessor: 'name' },
    { Header: 'status', accessor: 'Status' },
    {
      Header: 'Active Topics',
    },
    {
      Header: 'Data Storage',
    },
    {
      Header: 'Owner',
    },
    {
      Header: 'Date Created',
    },
    {
      Header: 'Actions',
    },
  ];

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
        <Table trClassName="text-sm" columns={initialColumns} data={[]} />
      </ErrorBoundary>
    </div>
  );
}

export default ProjectsTable;
