import { Loader, Table, Toast } from '@rotational/beacon-core';

import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';

import { getMembers } from '../util';

function TeamsTable() {
  const { members, isFetchingMembers, hasMembersFailed, error } = useFetchMembers();
  console.log(members);
  if (isFetchingMembers) {
    return <Loader />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasMembersFailed}
        variant="danger"
        title="We were unable to fetch your organizations. Please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

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
          { Header: 'Joined Date', accessor: 'date_added' },
          {
            // TODO: Make actions viewable only to members with owner and admin permission
            Header: 'Actions',
            Cell: () => {
              return <div>&hellip;</div>;
            },
          },
        ]}
        data={getMembers(members)}
      />
    </div>
  );
}

export default TeamsTable;
