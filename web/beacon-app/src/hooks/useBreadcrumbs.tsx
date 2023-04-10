import { useLocation, useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import BreadcrumbsIcon from '@/components/ui/Breadcrumbs/breadcrumbs-icon';

const useBreadcrumbs = () => {
  const { pathname } = useLocation();
  const navigateTo = useNavigate();

  const handleNavigateToHome = () => navigateTo(PATH_DASHBOARD.HOME);
  const splittedItems = pathname.split('/').filter(Boolean);
  const isClickable = splittedItems.length > 1;

  const items = splittedItems.map((item) =>
    item === 'app' ? (
      <button
        onClick={isClickable ? handleNavigateToHome : undefined}
        key={'home'}
        className={isClickable ? 'hover:underline' : undefined}
      >
        <BreadcrumbsIcon className="inline" /> Home
      </button>
    ) : (
      item
    )
  );

  return {
    items,
    separator: '/',
  };
};

export default useBreadcrumbs;
