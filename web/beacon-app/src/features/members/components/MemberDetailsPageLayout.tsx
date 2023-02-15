import { Route } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import DashLayout from '@/components/layout/DashLayout';

import MemberDetailsPage from './MemeberDetailsPage';

export default function MemberDetailsPageLayout() {
  return (
    <>
      <Route path={PATH_DASHBOARD.profile} element={<DashLayout />}>
        <Route path={PATH_DASHBOARD.profile} element={<MemberDetailsPage />} />
      </Route>
    </>
  );
}
