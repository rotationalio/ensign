import { createBrowserRouter, createRoutesFromElements, Outlet, Route } from 'react-router-dom';

import { ErrorPage } from '@/components/ErrorPage';
import MainLayout from '@/components/layout/MainLayout';
import { LoginPage, Registration, SuccessfulAccountCreation } from '@/features/auth';

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
      <Route element={<MainLayout />}>
        <Route path="register" element={<Registration />} />
        <Route path="/" element={<LoginPage />} />
        <Route path="verify-account" element={<SuccessfulAccountCreation />} />
      </Route>
    </Route>
  )
);

export default router;
