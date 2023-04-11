import { Heading } from '@rotational/beacon-core';
import { useState } from 'react';

import ArrowDownUp from '@/components/icons/arrow-down-up';
import FunnelSimple from '@/components/icons/funnel-simple';
import ThreeDots from '@/components/icons/three-dots';
import Union from '@/components/icons/union';
import AppLayout from '@/components/layout/AppLayout';
import Button from '@/components/ui/Button';
import { USER_PERMISSIONS } from '@/constants/rolesAndStatus';
import { useFetchMember } from '@/features/members/hooks/useFetchMember';
import { usePermissions } from '@/hooks/usePermissions';
import { useOrgStore } from '@/store';

import AddNewMemberModal from '../components/AddNewMember/AddNewMemberModal';
import TeamsTable from '../components/TeamsTable';
export function TeamsPage() {
  const orgDataState = useOrgStore.getState() as any;

  const { member } = useFetchMember(orgDataState?.user);
  const { hasPermission } = usePermissions();
  const [isModalOpened, setIsModalOpened] = useState(false);
  const onClose = () => setIsModalOpened(false);
  const onOpen = () => setIsModalOpened(true);
  console.log('member', member);

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
          <ul className="flex items-center gap-3">
            <li className="flex items-center justify-center">
              <button>
                <Union className="fill-[#1C1C1C]" />
              </button>
            </li>
            <li className="flex items-center justify-center">
              <button>
                <FunnelSimple />
              </button>
            </li>
            <li className="flex items-center justify-center">
              <button>
                <ArrowDownUp />
              </button>
            </li>
            <li className="flex items-center justify-center">
              <button>
                <ThreeDots />
              </button>
            </li>
          </ul>
          <div>
            {hasPermission(USER_PERMISSIONS.COLLABORATORS_ADD) && (
              <Button className="flex items-center gap-1 text-xs" size="small" onClick={onOpen}>
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
