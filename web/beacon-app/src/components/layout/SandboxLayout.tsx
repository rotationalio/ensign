import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import SandboxSidebar from './Sidebar/SandboxSidebar';

type SandboxLayoutProps = {
  children?: React.ReactNode;
};

function SandboxLayout({ children }: SandboxLayoutProps) {
  return (
    <div className="flex flex-col md:pl-[250px]">
      <SandboxSidebar className="hidden md:block" />
      <MainStyle>{children}</MainStyle>
      <MobileFooter />
    </div>
  );
}

export default SandboxLayout;
