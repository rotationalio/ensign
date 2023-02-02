import { Link } from 'react-router-dom';
import greenCircle from '/src/assets/icons/green-circle.svg'
import { routes } from '@/application';
import { memo } from 'react';


function AccessDashboard() {
  return (
    <div>
      <img src={greenCircle} alt="" />
      <div className="mt-4 ml-5">
      <Link to={routes.dashboard}>
        <span className="text-primary underline">View/Edit</span>
      </Link>
      </div>
    </div>
  );
}

export default memo(AccessDashboard)