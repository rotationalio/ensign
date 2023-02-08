import { Route } from 'react-router-dom';

import { routes } from '@/application';
import MainLayout from '@/components/layout/MainLayout';

import WelcomePage from './WelcomePage';

export default function WelcomePageLayout() {
  return (
    <>
      <Route path={routes.welcome} element={<MainLayout />}>
        <Route path={routes.welcome} element={<WelcomePage />} />
      </Route>
    </>
  );
}
