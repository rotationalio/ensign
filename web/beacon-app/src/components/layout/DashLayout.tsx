import Loader from '@/components/ui/Loader';
import useDashOnboarding from '@/hooks/useDashOnboarding';

import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import { Sidebar } from './Sidebar';

type DashLayoutProps = {
  children?: React.ReactNode;
};

const DashLayout: React.FC<DashLayoutProps> = ({ children }) => {
  const { isMemberLoading, wasMemberFetched } = useDashOnboarding();

  return (
    <div className="flex flex-col md:pl-[250px]">
      {isMemberLoading && <Loader />}
      {wasMemberFetched ? (
        <>
          <Sidebar className="hidden md:block" />
          <MainStyle>{children}</MainStyle>
          <MobileFooter />
        </>
      ) : null}
    </div>
  );
};

export default DashLayout;
