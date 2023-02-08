import { Outlet } from 'react-router-dom';

import { MainStyle } from './DashLayout.styles';
import { Sidebar } from './Sidebar';
import Topbar from './Topbar';

const DashLayout = () => {
  return (
    <div className="relative flex">
      <Sidebar />
      <Topbar />
      <MainStyle>
        <Outlet />
      </MainStyle>
    </div>
  );
};

export default DashLayout;
