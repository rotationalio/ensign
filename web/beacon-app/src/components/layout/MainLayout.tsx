import { Container } from '@rotational/beacon-core';
import { Outlet } from 'react-router-dom';

import { LandingFooter } from '@/components/auth/LandingFooter';
import { LandingHeader } from '@/components/auth/LandingHeader';

const MainLayout = () => {
  return (
    <>
      <LandingHeader />
      <Container>
        <Outlet />
      </Container>
      <LandingFooter />
    </>
  );
};

export default MainLayout;
