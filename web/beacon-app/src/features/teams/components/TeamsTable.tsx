import { Loader, Table, Toast } from '@rotational/beacon-core';

import ConfirmedIndicator from '@/components/icons/confirmedIndicator';
import PendingIndicator from '@/components/icons/pendingIndicator';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import { formatDate } from '@/utils/formatDate';

import { getMembers } from '../util';

function TeamsTable() {
  const { members, isFetchingMembers, hasMembersFailed, error } = useFetchMembers();

  if (isFetchingMembers) {
    return <Loader />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasMembersFailed}
        variant="danger"
        title="We were unable to fetch your members. Please try again later."
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
          {
            Header: 'Status',
            accessor: (m: any) => {
              return (
                <div className="flex items-center">
                  {m.status === 'Confirmed' && <ConfirmedIndicator />}
                  {m.status === 'Pending' && <PendingIndicator />}
                  <span className="pl-1">{m.status}</span>
                </div>
              );
            },
          },
          {
            Header: 'Last Activity',
            accessor: (date: any) => {
              return formatDate(new Date(date.last_activity));
            },
          },
          {
            Header: 'Joined Date',
            accessor: (date: any) => {
              return formatDate(new Date(date.date_added));
            },
          },
          {
            // TODO: Make actions viewable only to members with owner and admin permission
            Header: 'Actions',
            Cell: () => {
              return (
                <button type="button" className="mb-2 text-2xl font-bold">
                  &hellip;
                </button>
              );
            },
          },
        ]}
        data={getMembers(members)}
      />
    </div>
  );
}

export default TeamsTable;
