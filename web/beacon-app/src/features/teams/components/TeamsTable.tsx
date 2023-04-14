import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useState } from 'react';

import ConfirmedIndicatorIcon from '@/components/icons/confirmedIndicatorIcon';
import PendingIndicatorIcon from '@/components/icons/pendingIndicatorIcon';
import { MEMBER_STATUS, USER_PERMISSIONS } from '@/constants/rolesAndStatus';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import { usePermissions } from '@/hooks/usePermissions';
import { formatDate } from '@/utils/formatDate';

import { Member, MemberStatus } from '../types/member';
import { getMembers } from '../util';
import ChangeRoleModal from './ChangeRoleModal';
import DeleteMemberModal from './DeleteMember/DeleteMemberModal';

function TeamsTable() {
  const { members } = useFetchMembers();
  const { hasPermission } = usePermissions();

  const [openChangeRoleModal, setOpenChangeRoleModal] = useState<{
    opened: boolean;
    member?: Member;
  }>({
    opened: false,
    member: undefined,
  });

  const [openDeleteMemberModal, setOpenDeleteMemberModal] = useState<{
    opened: boolean;
    member?: Member;
  }>({
    opened: false,
    member: undefined,
  });

  const handleOpenChangeRoleModal = (member: Member) =>
    setOpenChangeRoleModal({ member, opened: true });

  const handleOpenDeleteMemberModal = (member: Member) =>
    setOpenDeleteMemberModal({ member, opened: true });

  const handleOncloseDeleteMemberModal = () => setOpenDeleteMemberModal({ opened: false });

  const memberStatusIconMap = {
    [MEMBER_STATUS.CONFIRMED]: <ConfirmedIndicatorIcon />,
    [MEMBER_STATUS.PENDING]: <PendingIndicatorIcon />,
  };

  const initialColumns = [
    { Header: 'Name', accessor: 'name' },
    { Header: 'Email Address', accessor: 'email' },
    { Header: 'Role', accessor: 'role' },
    {
      Header: 'Status',
      accessor: (m: { status: MemberStatus }) => {
        return (
          <div className="flex items-center">
            {memberStatusIconMap[m.status]}
            <span className="ml-1">{m.status}</span>
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
  ];

  const actionsColumn = { Header: 'Actions', accessor: 'actions' };

  {
    hasPermission(USER_PERMISSIONS.COLLABORATORS_EDIT || USER_PERMISSIONS.COLLABORATORS_REMOVE) &&
      initialColumns.push(actionsColumn);
  }

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
        <Table
          trClassName="text-sm"
          columns={initialColumns}
          data={getMembers(members, {
            handleOpenChangeRoleModal,
            handleOpenDeleteMemberModal,
          })}
        />
        <ChangeRoleModal
          openChangeRoleModal={openChangeRoleModal}
          setOpenChangeRoleModal={setOpenChangeRoleModal}
        />
        <DeleteMemberModal
          onOpen={openDeleteMemberModal}
          onClose={handleOncloseDeleteMemberModal}
        />
      </ErrorBoundary>
    </div>
  );
}

export default TeamsTable;
