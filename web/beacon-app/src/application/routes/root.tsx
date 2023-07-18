import React from 'react';
import { createBrowserRouter, createRoutesFromElements, Outlet, Route } from 'react-router-dom';

import { ErrorPage } from '@/components/Error/ErrorPage';
import MainLayout from '@/components/layout/MainLayout';
import {
  LoginPage,
  RegistrationPage,
  SuccessfulAccountCreation,
  VerifyPage,
} from '@/features/auth';
import { SetupTenantPage, WelcomePage } from '@/features/onboarding';
import { inviteTeamMemberLoader, InviteTeamMemberVerification } from '@/features/teams';
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
const { ProjectsPage } = lazyImport(() => import('@/features/projects'), 'ProjectsPage');
const { ProjectDetailPage } = lazyImport(() => import('@/features/projects'), 'ProjectDetailPage');
const MemberDetailsPage = React.lazy(
  () => import('@/features/members/components/MemeberDetailsPage')
);
const { OrganizationPage } = lazyImport(
  () => import('@/features/organization'),
  'OrganizationPage'
);
const { TopicDetailPage } = lazyImport(() => import('@/features/topics'), 'TopicDetailPage');

const { TeamsPage } = lazyImport(() => import('@/features/teams'), 'TeamsPage');

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Root />} errorElement={<ErrorPage />}>
      <Route path="app" element={<PrivateRoute />}>
        <Route index element={<Home />} />
        <Route path="dashboard" element={<Home />} />
        <Route path="projects">
          <Route index element={<ProjectsPage />} />
          <Route path=":project-setup" element={<>Project setup</>} />
          <Route path=":id" element={<ProjectDetailPage />} />
        </Route>
        <Route path="topics/:id" element={<TopicDetailPage />} />
        <Route path="organization" element={<OrganizationPage />} />
        <Route path="profile" element={<MemberDetailsPage />} />
        <Route path="team" element={<TeamsPage />} />
      </Route>
      <Route element={<PublicRoutes />}>
        <Route path="register" element={<RegistrationPage />} />
        <Route path="/" element={<LoginPage />} />
        <Route path="onboarding/getting-started" element={<WelcomePage />} />
        <Route path="onboarding/setup" element={<SetupTenantPage />} />
        <Route
          path="invite"
          loader={inviteTeamMemberLoader}
          element={<InviteTeamMemberVerification />}
        />
      </Route>
      <Route element={<MainLayout />}>
        <Route path="verify" element={<VerifyPage />} />
        <Route path="verify-account" element={<SuccessfulAccountCreation />} />
      </Route>
    </Route>
  )
);

export default router;
