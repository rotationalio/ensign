import Loader from '@/components/ui/Loader';
import useDashOnboarding from '@/hooks/useDashOnboarding';

import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import { Sidebar } from './Sidebar';

type DashLayoutProps = {
  children?: React.ReactNode;
};

const DashLayout: React.FC<DashLayoutProps> = ({ children }) => {
  const { wasProfileFetched, isFetchingProfile } = useDashOnboarding();

  return (
    <div className="flex flex-col md:pl-[250px]">
      {isFetchingProfile && <Loader />}
      {wasProfileFetched ? (
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
