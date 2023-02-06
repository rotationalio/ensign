import { Outlet } from 'react-router-dom';
import { Container } from '@rotational/beacon-core';

import { LandingFooter } from '@/components/auth/LandingFooter';
import { LandingHeader } from '@/components/auth/LandingHeader';

const MainLayout = () => {
  return (
    <>
      <LandingHeader />
      <Container className="place-items-center">
        <Outlet />
      </Container>
      <LandingFooter />
    </>
  );
};

export default MainLayout;
