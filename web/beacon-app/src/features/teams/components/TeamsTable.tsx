import { Loader, Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import ConfirmedIndicator from '@/components/icons/confirmedIndicator';
import PendingIndicator from '@/components/icons/pendingIndicator';
import { MEMBER_STATUS } from '@/constants/rolesAndStatus';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import { formatDate } from '@/utils/formatDate';

import { Member, MemberStatus } from '../types/member';
import { getMembers } from '../util';
import ChangeRoleModal from './ChangeRoleModal';

function TeamsTable() {
  const { members, isFetchingMembers, hasMembersFailed, error } = useFetchMembers();
  const [openChangeRoleModal, setOpenChangeRoleModal] = useState<{
    opened: boolean;
    member?: Member;
  }>({
    opened: false,
    member: undefined,
  });

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

  const handleOpenChangeRoleModal = (member: Member) =>
    setOpenChangeRoleModal({ member, opened: true });

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
            accessor: (m: { status: MemberStatus }) => {
              return (
                <div className="flex items-center">
                  {m.status === MEMBER_STATUS.CONFIRMED && <ConfirmedIndicator />}
                  {m.status === MEMBER_STATUS.PENDING && <PendingIndicator />}
                  <span className="pl-1">{m.status}</span>
                </div>
              );
            },
          },
          {
            Header: 'Last Activity',
            accessor: (date: any) => {
              return formatDate(new Date(date?.last_activity));
            },
          },
          {
            Header: 'Joined Date',
            accessor: (date: any) => {
              return formatDate(new Date(date?.date_added));
            },
          },
          {
            Header: 'Actions',
            accessor: 'actions',
          },
        ]}
        data={getMembers(members, { handleOpenChangeRoleModal })}
      />
      <ChangeRoleModal
        openChangeRoleModal={openChangeRoleModal}
        setOpenChangeRoleModal={setOpenChangeRoleModal}
      />
    </div>
  );
}

export default TeamsTable;
