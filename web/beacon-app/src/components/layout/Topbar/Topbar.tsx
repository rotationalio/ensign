import { Trans } from '@lingui/macro';
import { memo, ReactNode, useId, useState } from 'react';

import { ProfileCard } from '@/components/common/ProfileCard/ProfileCard';
import { MenuDropdownMenu } from '@/components/MenuDropdown/MenuDropdown';
import { useDropdownMenu } from '@/components/MenuDropdown/useDropdownMenu';
import Breadcrumbs from '@/components/ui/Breadcrumbs';
import { useFetchOrganizations } from '@/features/organization/hooks/useFetchOrganizations';
import { useAuth } from '@/hooks/useAuth';
import useBreadcrumbs from '@/hooks/useBreadcrumbs';
import { useOrgStore } from '@/store';

import ScheduleOfficeHours from '../../ScheduleOfficeHours/ScheduleOfficeHours';
import MobileNav from '../MobileNav/MobileNav';
import ProfileAvatar from '../ProfileAvatar/ProfileAvatar';
type TopBarProps = {
  Breadcrumbs?: ReactNode;
  isOnboarded?: boolean;
  profileData?: any;
};

function Topbar({ Breadcrumbs: CustomBreadcrumbs, isOnboarded, profileData }: TopBarProps) {
  const { items, separator } = useBreadcrumbs();
  const id = useId();
  const { logout } = useAuth();

  const Logout = () => {
    logout();
    window.location.href = '/';
  };

  const [isOpen, setIsOpen] = useState(false);

  const onOpenChange = () => {
    setIsOpen(!isOpen);
  };

  const appState = useOrgStore((state: any) => state) as any;
  const { organizations } = useFetchOrganizations();

  const { menuItems: dropdownItems } = useDropdownMenu({
    organizationsList: organizations?.organizations,
    currentOrg: appState?.orgID,
  });

  return (
    <div className="flex flex-col items-baseline justify-center gap-2 overflow-auto bg-[#1D65A6] py-4 md:border-b md:bg-white">
      <div className="flex w-[98%] justify-between py-2 xl:w-[92.5%]">
        {isOnboarded ? (
          <>
            {CustomBreadcrumbs ? (
              CustomBreadcrumbs
            ) : (
              <Breadcrumbs separator={separator} className="ml-4 hidden md:block">
                {items.map((item) => (
                  <Breadcrumbs.Item key={item + id} className="capitalize">
                    {item}
                  </Breadcrumbs.Item>
                ))}
              </Breadcrumbs>
            )}
            <div className="flex space-x-4">
              <ScheduleOfficeHours />
              <MenuDropdownMenu
                items={dropdownItems}
                trigger={<ProfileAvatar name={profileData?.name} />}
                onOpenChange={onOpenChange}
                isOpen={isOpen}
                data-cy="menu-dropdown"
              />
            </div>
          </>
        ) : (
          <>
            <span></span>
            <div className="flex items-center justify-end">
              <ProfileCard
                picture={profileData?.picture}
                owner_name={profileData?.email}
                cardSize="medium"
              />
              <button
                onClick={Logout}
                className="ml-4 pb-1 font-bold text-primary"
                data-cy="log-out-bttn"
              >
                <Trans>Log Out</Trans>
              </button>
            </div>
          </>
        )}
      </div>
      {!isOnboarded && <MobileNav />}
    </div>
  );
}

export default memo(Topbar);
