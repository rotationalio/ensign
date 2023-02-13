import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';

import DashLayout from '@/components/layout/DashLayout';
import { useAuth } from '@/hooks/useAuth';

const PrivateRoute = () => {
  const { isAuthenticated } = useAuth();
  console.log('isAuthenticated', isAuthenticated);
  return isAuthenticated ? (
    <DashLayout>
      <Outlet />
    </DashLayout>
  ) : (
    <Navigate to="/" />
  );
};

export default PrivateRoute;
