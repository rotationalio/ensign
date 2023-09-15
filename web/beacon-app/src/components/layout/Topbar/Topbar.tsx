import { memo, ReactNode, useId } from 'react';

import { ProfileCard } from '@/components/common/ProfileCard/ProfileCard';
import Breadcrumbs from '@/components/ui/Breadcrumbs';
import { useAuth } from '@/hooks/useAuth';
import useBreadcrumbs from '@/hooks/useBreadcrumbs';

import ScheduleOfficeHours from '../../ScheduleOfficeHours/ScheduleOfficeHours';
import MobileNav from '../MobileNav/MobileNav';
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

  return (
    <>
      <Header className="flex flex-col-reverse items-baseline justify-center gap-2 bg-[#1D65A6] py-2 md:ml-[250px] md:min-h-[60px] md:border-b md:bg-white">
        <div className="flex w-11/12 justify-between">
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
              <ScheduleOfficeHours />
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
        <MobileNav />
      </Header>
    </>
  );
}

export default memo(Topbar);
