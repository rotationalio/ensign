import { MainStyle } from './DashLayout.styles';
import { Sidebar } from './Sidebar';
import Topbar from './Topbar';

type DashLayoutProps = {
  children?: React.ReactNode;
};
const DashLayout: React.FC<DashLayoutProps> = ({ children }) => {
  return (
    <div className="relative flex">
      <Sidebar />
      <Topbar />
      <MainStyle>{children}</MainStyle>
    </div>
  );
};

export default DashLayout;
