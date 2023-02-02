import { Link } from 'react-router-dom';
import checkmark from '/src/assets/images/checkmark.png';

import data from '/src/assets/images/hosted-data-icon.png';

import { routes } from '@/application';


export default function AccessDashboard() {
  return (
    <div className="grid max-w-6xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <img src={checkmark} alt="" />
      <Link to={routes.dashboard}>
        View/Edit
      </Link>
    </div>
  );
}
