import { Container } from '@rotational/beacon-core';
import React, { ReactNode } from 'react';

import useUserLoader from '@/features/members/loaders/userLoader';
import { isOnboardedMember } from '@/features/members/utils';

import Topbar from './Topbar';

type PageProps = {
  children: ReactNode;
  Breadcrumbs?: ReactNode;
};

function AppLayout({ children, Breadcrumbs }: PageProps) {
  const { member: loaderData } = useUserLoader();
  const isOnboarded = isOnboardedMember(loaderData?.onboarding_status);

  return (
    <>
      {isOnboarded && <Topbar Breadcrumbs={Breadcrumbs} />}
      <Container max={696} centered className="my-10 mt-8 px-4 xl:px-28">
        {children}
      </Container>
    </>
  );
}

export default AppLayout;
