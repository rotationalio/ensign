import { memo, ReactNode, useId } from 'react';

import Breadcrumbs from '@/components/ui/Breadcrumbs';
import useBreadcrumbs from '@/hooks/useBreadcrumbs';

import ScheduleOfficeHours from '../../ScheduleOfficeHours/ScheduleOfficeHours';
import MobileNav from '../MobileNav/MobileNav';
import { Header } from './Topbar.styles';

type TopBarProps = {
  Breadcrumbs?: ReactNode;
};

function Topbar({ Breadcrumbs: CustomBreadcrumbs }: TopBarProps) {
  const { items, separator } = useBreadcrumbs();
  const id = useId();

  return (
    <>
      <Header className="flex flex-col-reverse items-baseline justify-center gap-2 bg-[#1D65A6] py-2 md:ml-[250px] md:min-h-[60px] md:border-b md:bg-white">
        <div className="flex w-11/12 justify-between">
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
        </div>
        <MobileNav />
      </Header>
    </>
  );
}

export default memo(Topbar);
