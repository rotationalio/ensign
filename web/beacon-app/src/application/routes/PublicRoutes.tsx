import React from 'react';
import { Navigate } from 'react-router-dom';

import MainLayout from '@/components/layout/MainLayout';
import { useAuth } from '@/hooks/useAuth';

function PublicRoutes() {
  const { isAuthenticated } = useAuth();

  if (isAuthenticated) return <Navigate to="/app" />;

  return <MainLayout />;
}

export default PublicRoutes;
