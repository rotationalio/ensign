import { Table } from '@rotational/beacon-core';
import React from 'react';

function TeamsTable() {
  return (
    <div className="mx-4">
      <Table
        trClassName="text-sm"
        columns={[
          { Header: 'Name', accessor: 'name' },
          { Header: 'Email Address', accessor: 'email' },
          { Header: 'Role', accessor: 'role' },
          { Header: 'Status', accessor: 'status' },
          { Header: 'Last Activity', accessor: 'last_activity' },
          { Header: 'Joined Date', accessor: 'joined_date' },
          {
            Header: 'Actions',
          },
        ]}
        data={[]}
      />
    </div>
  );
}

export default TeamsTable;
