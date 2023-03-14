import { Container } from '@rotational/beacon-core';
import React, { ReactNode } from 'react';

import Topbar from './Topbar';

type PageProps = {
  children: ReactNode;
  Breadcrumbs?: ReactNode;
};

// eslint-disable-next-line unused-imports/no-unused-vars
function AppLayout({ children, Breadcrumbs }: PageProps) {
  return (
    <>
      <Topbar Breadcrumbs={Breadcrumbs} />
      <Container max={696} centered className="my-10 mt-8 px-4 xl:px-28">
        {children}
      </Container>
    </>
  );
}

export default AppLayout;
