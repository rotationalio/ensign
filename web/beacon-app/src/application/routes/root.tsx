import { createBrowserRouter, createRoutesFromElements, Outlet, Route } from 'react-router-dom';

import ErrorPage from '@/components/ErrorPage';
import MainLayout from '@/components/layout/MainLayout';
import { Registration, SuccessfulAccountCreation } from '@/features/auth';
// import routers from features
// should we import all routes files in features folder automatically with a glob pattern?

const Root = () => {
  return (
    <div>
      <Outlet />
    </div>
  );
};

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<Root />} errorElement={<ErrorPage />}>
      <Route path="/auth" element={<MainLayout />}>
        <Route path="register" element={<Registration />} />
        <Route path="verify-account" element={<SuccessfulAccountCreation />} />
      </Route>
    </Route>
  )
);

export default router;
