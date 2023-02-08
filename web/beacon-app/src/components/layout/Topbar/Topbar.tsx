import { memo } from 'react';
import { useLocation } from 'react-router-dom';

import Breadcrumbs from '@/components/ui/Breadcrumbs';
import BreadcrumbsIcon from '@/components/ui/Breadcrumbs/breadcrumbs-icon';

import { Header } from './Topbar.styles';

function Topbar() {
  const { pathname } = useLocation();
  const items = pathname.split('/').filter(Boolean);

  return (
    <Header>
      <Breadcrumbs separator="/" className="ml-4">
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
    </Header>
  );
}

export default memo(Topbar);
