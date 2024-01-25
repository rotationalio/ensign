import MobileFooter from './MobileFooter';
import SandboxSidebar from './Sidebar/SandboxSidebar';

type SandboxLayoutProps = {
  children?: React.ReactNode;
};

function SandboxLayout({ children }: SandboxLayoutProps) {
  return (
    <div className="flex flex-col md:pl-[250px]" data-testid="sandbox-layout">
      <SandboxSidebar className="hidden md:block" />
      <div>{children}</div>
      <MobileFooter />
    </div>
  );
}

export default SandboxLayout;
