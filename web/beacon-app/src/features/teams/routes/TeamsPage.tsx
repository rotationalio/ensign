import { Heading } from '@rotational/beacon-core';
import { useState } from 'react';

import Union from '@/components/icons/union';
import AppLayout from '@/components/layout/AppLayout';
import Button from '@/components/ui/Button';
import { USER_PERMISSIONS } from '@/constants/rolesAndStatus';
import { usePermissions } from '@/hooks/usePermissions';

import AddNewMemberModal from '../components/AddNewMember/AddNewMemberModal';
import TeamsTable from '../components/TeamsTable';
export function TeamsPage() {
  const { hasPermission } = usePermissions();
  const [isModalOpened, setIsModalOpened] = useState(false);
  const onClose = () => setIsModalOpened(false);
  const onOpen = () => setIsModalOpened(true);

  return (
    <AppLayout>
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        Team
      </Heading>
      <p className="my-3 text-sm">
        Add team members to collaborate on your projects. Team members have access to projects
        across the organization.
      </p>
      <div>
        <div className="flex justify-between rounded-lg bg-[#F7F9FB] px-3 py-2">
          <div className="flex items-center gap-3"></div>
          <div>
            {hasPermission(USER_PERMISSIONS.COLLABORATORS_ADD) && (
              <Button
                data-cy="add-team-member"
                className="flex items-center gap-1 text-xs"
                size="small"
                onClick={onOpen}
              >
                <Union className="fill-white" />
                Team Member
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
