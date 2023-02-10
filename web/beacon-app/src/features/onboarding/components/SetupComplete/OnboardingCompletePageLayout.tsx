import { Route } from 'react-router-dom';

import { routes } from '@/application';
import MainLayout from '@/components/layout/MainLayout';

import OnboardingCompletePage from './OnboardingCompletePage';

export default function OnboardingCompletePageLayout() {
  return (
    <>
      <Route path={routes.complete} element={<MainLayout />}>
        <Route path={routes.complete} element={<OnboardingCompletePage />} />
      </Route>
    </>
  );
}
