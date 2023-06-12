import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useState } from 'react';

import { USER_PERMISSIONS } from '@/constants/rolesAndStatus';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import { usePermissions } from '@/hooks/usePermissions';
import { formatDate } from '@/utils/formatDate';

import { Member } from '../types/member';
import { getMembers } from '../util';
import ChangeRoleModal from './ChangeRoleModal';
import DeleteMemberModal from './DeleteMember/DeleteMemberModal';

interface Props {
  isLoading?: boolean;
}

function TeamsTable({ isLoading }: Props) {
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

  const initialColumns: any = [
    { Header: 'Name', accessor: 'name' },
    {
      Header: 'Status',
      accessor: 'status',
    },
    { Header: 'Email Address', accessor: 'email' },
    { Header: 'Role', accessor: 'role' },
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

  if (hasPermission(USER_PERMISSIONS.COLLABORATORS_EDIT || USER_PERMISSIONS.COLLABORATORS_REMOVE)) {
    initialColumns.push(actionsColumn);
  }

  return (
    <div className="mx-4" data-testid="teamTable">
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
          isLoading={isLoading}
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
