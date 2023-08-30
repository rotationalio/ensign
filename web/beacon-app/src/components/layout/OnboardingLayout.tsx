import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import { OnBoardingSidebar } from './Sidebar';

type OnboardingLayoutProps = {
  children?: React.ReactNode;
};

const OnboardingLayout: React.FC<OnboardingLayoutProps> = ({ children }) => {
  return (
    <div className="flex flex-col md:pl-[250px]">
      <OnBoardingSidebar className="hidden md:block" />
      <MainStyle>{children}</MainStyle>
      <MobileFooter />
    </div>
  );
};

export default OnboardingLayout;
