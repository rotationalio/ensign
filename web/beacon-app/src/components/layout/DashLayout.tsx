import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import { Sidebar } from './Sidebar';
import Topbar from './Topbar';

type DashLayoutProps = {
  children?: React.ReactNode;
};
const DashLayout: React.FC<DashLayoutProps> = ({ children }) => {
  return (
    <div className="flex md:pl-[250px]">
      <Sidebar className="hidden md:block" />
      <Topbar />
      <MainStyle>{children}</MainStyle>
      <MobileFooter />
    </div>
  );
};

export default DashLayout;
