import { Container } from '@rotational/beacon-core';
import { Outlet } from 'react-router-dom';

import { LandingFooter } from '@/components/auth/LandingFooter';
import { LandingHeader } from '@/components/auth/LandingHeader';
import { WelcomePage } from '../../features/onboarding/components/Welcome';

function MainLayout () {
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
