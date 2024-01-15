import { Container } from '@rotational/beacon-core';
import React, { ReactNode, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { isOnboardedMember } from '@/features/members/utils';

import SandboxBanner from './SanboxBanner/SandboxBanner';
import Topbar from './Topbar';

type PageProps = {
  children: ReactNode;
  Breadcrumbs?: ReactNode;
};

function AppLayout({ children, Breadcrumbs }: PageProps) {
  const navigate = useNavigate();
  const { profile: loaderData } = useFetchProfile();
  const isOnboarded = isOnboardedMember(loaderData?.onboarding_status);

  // if onboarded redirect to onboarded route
  useEffect(() => {
    if (!isOnboarded) {
      navigate(PATH_DASHBOARD.ONBOARDING);
    }
  }, [isOnboarded, navigate]);

  return (
    <>
      <Topbar Breadcrumbs={Breadcrumbs} isOnboarded={isOnboarded} profileData={loaderData} />
      {/* TODO: Display SandboxBanner only to user's with the sandbox account type. */}
      <SandboxBanner />
      <Container max={696} centered className="my-10 mt-8 px-4 xl:px-28">
        {children}
      </Container>
    </>
  );
}

export default AppLayout;
