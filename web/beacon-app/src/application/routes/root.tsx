import { createBrowserRouter, createRoutesFromElements, Outlet, Route } from 'react-router-dom';

import { ErrorPage } from '@/components/ErrorPage';
import DashLayout from '@/components/layout/DashLayout';
import MainLayout from '@/components/layout/MainLayout';
import { LoginPage, RegistrationPage, SuccessfulAccountCreation } from '@/features/auth';
import { Home } from '@/features/home';
import { SetupTenantPage, WelcomePage } from '@/features/onboarding';

const Root = () => {
  return (
    <div>
      <Outlet />
    </div>
  );
};

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Root />} errorElement={<ErrorPage />}>
      <Route path="app" element={<DashLayout />}>
        <Route index element={<Home />} />
        <Route path="projects" element={<>Projects</>} />
      </Route>
      <Route element={<MainLayout />}>
        <Route path="register" element={<RegistrationPage />} />
        <Route path="/" element={<LoginPage />} />
        <Route path="verify-account" element={<SuccessfulAccountCreation />} />
        <Route path="onboarding/getting-started" element={<WelcomePage />} />
        <Route path="onboarding/setup" element={<SetupTenantPage />} />
      </Route>
    </Route>
  )
);

export default router;
