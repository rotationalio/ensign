import { memo, ReactNode, useId, useState } from 'react';

// import { EXTERNAL_LINKS } from '@/application';
// import Alert from '@/components/common/Alert/Alert';
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
import { Header } from './Topbar.styles';
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
    <>
      <Header className="flex flex-col items-baseline justify-center gap-2 bg-[#1D65A6] md:ml-[250px] md:min-h-[60px] md:border-b md:bg-white">
        {/* TODO: Display alert banner when user is on the sandbox. */}
        {/* <Alert>
          <div className="flex h-auto w-full flex-col items-center justify-center gap-x-4 bg-[#192E5B] py-6 text-center font-bold text-white lg:flex-row">
            <p>You are using the Ensign Sandbox. Ready to deploy your models to production?</p>
            <div className="mt-4 rounded-md border border-white bg-[#316D3C] py-1 px-4 lg:mt-0">
              <a href={EXTERNAL_LINKS.ENSIGN_PRICING} target="_blank" rel="noreferrer">
                Upgrade
              </a>
            </div>
          </div>
        </Alert> */}
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
              <div className="flex h-20  items-center justify-end">
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
                  Log Out
                </button>
              </div>
            </>
          )}
        </div>
        {!isOnboarded && <MobileNav />}
      </Header>
    </>
  );
}

export default memo(Topbar);
