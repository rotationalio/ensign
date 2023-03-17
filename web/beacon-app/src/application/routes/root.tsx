import React from 'react';
import {
  createBrowserRouter,
  createRoutesFromElements,
  Navigate,
  Outlet,
  Route,
} from 'react-router-dom';

import { ErrorPage } from '@/components/Error/ErrorPage';
import { LoginPage, RegistrationPage, VerifyPage } from '@/features/auth';
import { SetupTenantPage, WelcomePage } from '@/features/onboarding';
import { lazyImport } from '@/utils/lazy-import';

import PrivateRoute from './privateRoute';
import PublicRoutes from './PublicRoutes';
const Root = () => {
  return (
    <div>
      <Outlet />
    </div>
  );
};

const { Home } = lazyImport(() => import('@/features/home'), 'Home');
const { ProjectDetailPage } = lazyImport(() => import('@/features/projects'), 'ProjectDetailPage');
const MemberDetailsPage = React.lazy(
  () => import('@/features/members/components/MemeberDetailsPage')
);
const { OrganizationPage } = lazyImport(
  () => import('@/features/organization'),
  'OrganizationPage'
);

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Root />} errorElement={<ErrorPage />}>
      <Route path="app" element={<PrivateRoute />}>
        <Route index element={<Home />} />
        <Route path="dashboard" element={<Home />} />
        <Route path="projects" element={<>Projects List</>} />
        <Route path="projects/:id" element={<ProjectDetailPage />} />
        <Route path="organization" element={<OrganizationPage />} />
        <Route path="profile" element={<MemberDetailsPage />} />
        <Route path="*" element={<Navigate to="/app" />} />
      </Route>
      <Route element={<PublicRoutes />}>
        <Route path="register" element={<RegistrationPage />} />
        <Route path="/" element={<LoginPage />} />
        <Route path="onboarding/getting-started" element={<WelcomePage />} />
        <Route path="onboarding/setup" element={<SetupTenantPage />} />
      </Route>
      <Route path="verify" element={<VerifyPage />} />
    </Route>
  )
);

export default router;
