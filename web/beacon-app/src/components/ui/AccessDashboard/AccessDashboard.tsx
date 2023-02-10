import { memo } from 'react';
import { Link } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';

function AccessDashboard() {
  return (
    <div>
      <HeavyCheckMark />
      <div className="mt-4 ml-2">
        <Link to={PATH_DASHBOARD.root}>
          <span className="text-primary underline">View/Edit</span>
        </Link>
      </div>
    </div>
  );
}

export default memo(AccessDashboard);
