import React, { Suspense } from 'react';
import { Navigate, Outlet, useLoaderData } from 'react-router-dom';

import OvalLoader from '@/components/ui/OvalLoader';
import { isOnboardedMember } from '@/features/members/utils';
import { useAuth } from '@/hooks/useAuth';
const DashLayout = React.lazy(() => import('@/components/layout/DashLayout'));
const OnboardingLayout = React.lazy(() => import('@/components/layout/OnboardingLayout'));

const PrivateRoute = () => {
  const { isAuthenticated } = useAuth();
  const loaderData = useLoaderData() as any;
  // console.log('[] loaderData', loaderData?.member);
  const isOnboarded = isOnboardedMember(loaderData?.member?.status);

  const Layout = !isOnboarded ? DashLayout : OnboardingLayout;

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
