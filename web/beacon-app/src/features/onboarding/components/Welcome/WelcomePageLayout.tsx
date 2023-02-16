import { Route } from 'react-router-dom';

import { ROUTES } from '@/application';
import MainLayout from '@/components/layout/MainLayout';

import WelcomePage from './WelcomePage';

export default function WelcomePageLayout() {
  return (
    <>
      <Route path={ROUTES.WELCOME} element={<MainLayout />}>
        <Route path={ROUTES.WELCOME} element={<WelcomePage />} />
      </Route>
    </>
  );
}
