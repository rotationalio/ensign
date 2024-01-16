import React, { Suspense } from 'react';
import { Navigate, Outlet } from 'react-router-dom';

import OvalLoader from '@/components/ui/OvalLoader';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { isOnboardedMember } from '@/features/members/utils';
import { isSandboxAccount } from '@/features/sandbox/util/utils';
import useAccountType from '@/hooks/useAccountType';
import { useAuth } from '@/hooks/useAuth';
const DashLayout = React.lazy(() => import('@/components/layout/DashLayout'));
const OnboardingLayout = React.lazy(() => import('@/components/layout/OnboardingLayout'));
const SandboxLayout = React.lazy(() => import('@/components/layout/SandboxLayout'));

const PrivateRoute = () => {
  const { profile: userProfile } = useFetchProfile();
  const { isAuthenticated } = useAuth();
  const isOnboarded = isOnboardedMember(userProfile?.onboarding_status);

  const accountType = useAccountType();
  const isSandboxAcct = isSandboxAccount(accountType);

  let Layout = OnboardingLayout;

  if (isOnboarded && isSandboxAcct) {
    Layout = SandboxLayout;
  } else if (isOnboarded && !isSandboxAcct) {
    Layout = DashLayout;
  }

  return isAuthenticated ? (
    <Suspense
      fallback={
        <div className="grid h-screen w-screen place-items-center">
          <OvalLoader width="50px" height="50px" />
        </div>
      }
    >
      <Layout>
        <Suspense
          fallback={
            <div className="flex items-center justify-center">
              <OvalLoader />
            </div>
          }
        >
          <Outlet />
        </Suspense>
      </Layout>
    </Suspense>
  ) : (
    <Navigate to="/" />
  );
};

export default PrivateRoute;
