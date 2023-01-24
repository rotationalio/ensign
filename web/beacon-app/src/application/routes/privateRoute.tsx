import React from 'react';
import { Navigate } from 'react-router-dom';

interface Props {
  component: React.ComponentType;
  path?: string;
}

export const PrivateRoute: React.FC<Props> = ({ component: RouteComponent }) => {
  const isAuthenticated = true;
  const hasRequiredRole = true;

  if (isAuthenticated && hasRequiredRole) {
    return <RouteComponent />;
  }

  if (isAuthenticated && !hasRequiredRole) {
    return <Navigate to="/unauthorized" />;
  }

  return <Navigate to="/" />;
};
