// import { Trans } from '@lingui/macro';
import { Container } from '@rotational/beacon-core';
import React, { ReactNode, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { /* EXTERNAL_LINKS, */ PATH_DASHBOARD } from '@/application';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { isOnboardedMember } from '@/features/members/utils';

// import Alert from '../common/Alert/Alert';
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
      {/* TODO: Move to a SandboxAlert component and add a check for the user's account type
      and only display if they are a sandbox user. */}
      {/*  <Alert>
        <div className="flex h-auto w-full flex-col items-center justify-center gap-x-4 bg-[#192E5B] py-6 text-center font-bold text-white lg:flex-row">
          <p>
            <Trans>
              You are using the Ensign Sandbox. Ready to deploy your models to production?
            </Trans>
          </p>
          <div className="mt-4 rounded-md border border-white bg-[#316D3C] px-4 py-1 lg:mt-0">
            <a href={EXTERNAL_LINKS.ENSIGN_PRICING} target="_blank" rel="noreferrer">
              <Trans>Upgrade</Trans>
            </a>
          </div>
        </div>
      </Alert> */}
      <Container max={696} centered className="my-10 mt-8 px-4 xl:px-28">
        {children}
      </Container>
    </>
  );
}

export default AppLayout;
