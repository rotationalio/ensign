import { DefaultToastOptions, Toaster } from 'react-hot-toast';

import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import { Sidebar } from './Sidebar';
type DashLayoutProps = {
  children?: React.ReactNode;
};

const TOAST_DURATION = 5 * 1000;

const toasterOptions: DefaultToastOptions = {
  duration: TOAST_DURATION,
};

const DashLayout: React.FC<DashLayoutProps> = ({ children }) => {
  return (
    <div className="flex flex-col md:pl-[250px]">
      <Sidebar className="hidden md:block" />
      <MainStyle>{children}</MainStyle>
      <MobileFooter />
      <Toaster toastOptions={toasterOptions} />
    </div>
  );
};

export default DashLayout;
