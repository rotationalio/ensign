import { Container } from '@rotational/beacon-core';
import { Outlet } from 'react-router-dom';

import { LandingFooter } from '@/components/auth/LandingFooter';
import { LandingHeader } from '@/components/auth/LandingHeader';
import { WelcomePage } from '../ui/WelcomePage';

const MainLayout = () => {
  return (
    <Container>
      <LandingHeader />
      <Outlet />
      <WelcomePage />
      <LandingFooter />
    </Container>
  );
};

export default MainLayout;
