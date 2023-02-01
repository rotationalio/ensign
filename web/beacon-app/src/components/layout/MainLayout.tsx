import { Container } from '@rotational/beacon-core';
import { FC } from 'react';

import { LandingFooter } from '@/components/auth/LandingFooter';
import { LandingHeader } from '@/components/auth/LandingHeader';

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: FC<MainLayoutProps> = (props) => {
  return (
    <Container>
      <LandingHeader />
      {props.children}
      <LandingFooter />
    </Container>
  );
};

export default MainLayout;
