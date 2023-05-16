import { Trans } from '@lingui/macro';
import { Button, Heading, mergeClassnames } from '@rotational/beacon-core';
import { useState } from 'react';

import Union from '@/components/icons/union';
import AppLayout from '@/components/layout/AppLayout';
import { USER_PERMISSIONS } from '@/constants/rolesAndStatus';
import { usePermissions } from '@/hooks/usePermissions';

import AddNewMemberModal from '../components/AddNewMember/AddNewMemberModal';
import TeamsTable from '../components/TeamsTable';
export function TeamsPage() {
  const { hasPermission } = usePermissions();
  const [isModalOpened, setIsModalOpened] = useState(false);
  const onClose = () => setIsModalOpened(false);
  const onOpen = () => setIsModalOpened(true);

  const hasPermissions = hasPermission(USER_PERMISSIONS.COLLABORATORS_ADD);

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
        <div
          className={mergeClassnames(
            'flex justify-between rounded-lg px-3 py-2',
            hasPermissions ? 'bg-[#F7F9FB]' : 'bg-neutral-white'
          )}
        >
          <div className="flex items-center gap-3"></div>
          <div>
            {hasPermissions && (
              <Button
                data-cy="add-team-member"
                className="flex items-center gap-1"
                size="medium"
                variant="primary"
                onClick={onOpen}
              >
                <Union className="fill-white" />
                <Trans>Team Member</Trans>
              </Button>
            )}
          </div>
        </div>
        <AddNewMemberModal isOpened={isModalOpened} onClose={onClose} />
        <TeamsTable />
      </div>
    </AppLayout>
  );
}
