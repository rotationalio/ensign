import { Container } from '@rotational/beacon-core';
import React, { ReactNode, useEffect } from 'react';
import { useLoaderData, useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { isOnboardedMember } from '@/features/members/utils';

import Topbar from './Topbar';

type PageProps = {
  children: ReactNode;
  Breadcrumbs?: ReactNode;
};

function AppLayout({ children, Breadcrumbs }: PageProps) {
  const navigate = useNavigate();
  const loaderData = useLoaderData() as any;
  const isOnboarded = isOnboardedMember(loaderData?.member?.status);

  useEffect(() => {
    if (isOnboarded) {
      navigate(PATH_DASHBOARD.ONBOARDING);
    }
  }, [isOnboarded, navigate]);

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
