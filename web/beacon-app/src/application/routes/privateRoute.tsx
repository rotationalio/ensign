import React, { Suspense } from 'react';
import { Navigate, Outlet } from 'react-router-dom';

import OvalLoader from '@/components/ui/OvalLoader';
import { useAuth } from '@/hooks/useAuth';

const DashLayout = React.lazy(() => import('@/components/layout/DashLayout'));

const PrivateRoute = () => {
  const { isAuthenticated } = useAuth();

  return isAuthenticated ? (
    <Suspense
      fallback={
        <div className="grid h-screen w-screen place-items-center">
          <OvalLoader width="50px" height="50px" />
        </div>
      }
    >
      <DashLayout>
        <Suspense
          fallback={
            <div className="flex items-center justify-center">
              <OvalLoader />
            </div>
          }
        >
          <Outlet />
        </Suspense>
      </DashLayout>
    </Suspense>
  ) : (
    <Navigate to="/" />
  );
};

export default PrivateRoute;
