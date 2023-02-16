import { Route } from 'react-router-dom';

import { ROUTES } from '@/application';
import MainLayout from '@/components/layout/MainLayout';

import OnboardingCompletePage from './OnboardingCompletePage';

export default function OnboardingCompletePageLayout() {
  return (
    <>
      <Route path={ROUTES.COMPLETE} element={<MainLayout />}>
        <Route path={ROUTES.COMPLETE} element={<OnboardingCompletePage />} />
      </Route>
    </>
  );
}
