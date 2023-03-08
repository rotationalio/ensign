import { Container } from '@rotational/beacon-core';
import React, { ReactNode } from 'react';

type PageProps = {
  children: ReactNode;
};

function AppLayout({ children }: PageProps) {
  return (
    <Container max={696} centered className="my-10 mt-8 px-4 xl:px-28">
      {children}
    </Container>
  );
}

export default AppLayout;
