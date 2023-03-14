import { memo } from 'react';
import { useLocation } from 'react-router-dom';

import Breadcrumbs from '@/components/ui/Breadcrumbs';
import BreadcrumbsIcon from '@/components/ui/Breadcrumbs/breadcrumbs-icon';

import MobileNav from '../MobileNav/MobileNav';
import { Header } from './Topbar.styles';

function Topbar() {
  const { pathname } = useLocation();
  const items = pathname.split('/').filter(Boolean);

  return (
    <Header className="flex flex-col-reverse items-baseline justify-center gap-2 bg-[#1D65A6] py-2 md:ml-[250px] md:min-h-[60px] md:border-b md:bg-white">
      <Breadcrumbs separator="/" className="ml-4 hidden md:block">
        {items.map((item) => (
          <Breadcrumbs.Item key={item} className="capitalize">
            {item === 'app' ? (
              <>
                <BreadcrumbsIcon className="inline" /> Home
              </>
            ) : (
              item
            )}
          </Breadcrumbs.Item>
        ))}
      </Breadcrumbs>
      <MobileNav />
    </Header>
  );
}

export default memo(Topbar);
