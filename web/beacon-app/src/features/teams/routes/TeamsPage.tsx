import { Trans } from '@lingui/macro';
import { Button, Heading, mergeClassnames } from '@rotational/beacon-core';
import { useState } from 'react';

import RefreshIcon from '@/components/icons/refresh';
import Union from '@/components/icons/union';
import AppLayout from '@/components/layout/AppLayout';
import { USER_PERMISSIONS } from '@/constants/rolesAndStatus';
import { useFetchMembers } from '@/features/members/hooks/useFetchMembers';
import { usePermissions } from '@/hooks/usePermissions';

import AddNewMemberModal from '../components/AddNewMember/AddNewMemberModal';
import TeamsTable from '../components/TeamsTable';
export function TeamsPage() {
  const { hasPermission } = usePermissions();
  const [isModalOpened, setIsModalOpened] = useState(false);
  const onClose = () => setIsModalOpened(false);
  const onOpen = () => setIsModalOpened(true);

  const hasPermissions = hasPermission(USER_PERMISSIONS.COLLABORATORS_ADD);

  const { getMembers, isFetchingMembers } = useFetchMembers();

  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const refreshHandler = () => {
    setIsRefreshing(true);
    setTimeout(() => {
      getMembers();
      setIsRefreshing(false);
    }, 500);
  };

  return (
    <AppLayout>
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        <Trans>Team</Trans>
      </Heading>
      <p className="my-3 text-sm">
        <Trans>
          Add team members to collaborate on your projects. Team members have access to projects
          across the organization.
        </Trans>
      </p>
      <div>
        <div className="flex justify-between rounded-lg bg-[#F7F9FB] px-3 py-2">
          <div className="mt-3 ml-2">
            <button disabled={isFetchingMembers} onClick={refreshHandler}>
              <div className={mergeClassnames(isRefreshing ? 'animate-spin-slow' : '')}>
                <RefreshIcon />
              </div>
            </button>
          </div>
          <div className="flex items-center gap-3"></div>
          <div>
            <Button
              data-cy="add-team-member"
              className="flex items-center gap-1"
              size="medium"
              variant="primary"
              disabled={!hasPermissions}
              onClick={onOpen}
            >
              <Union className="fill-white" />
              <Trans>Team Member</Trans>
            </Button>
          </div>
        </div>
        <AddNewMemberModal isOpened={isModalOpened} onClose={onClose} />
        <TeamsTable isLoading={isRefreshing || isFetchingMembers} />
      </div>
    </AppLayout>
  );
}
