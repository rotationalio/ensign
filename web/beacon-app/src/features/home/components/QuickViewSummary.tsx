import { Loader } from '@rotational/beacon-core';
import { Suspense } from 'react';

import { QuickView } from '@/components/common/QuickView';
import { SentryErrorBoundary } from '@/components/Error';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { useFetchTenantQuickView } from '@/hooks/useFetchQuickView';

import { getDefaultHomeStats, getHomeStatsHeaders } from '../util';
function QuickViewSummary() {
  const { tenants } = useFetchTenants();

  const { quickView } = useFetchTenantQuickView(tenants?.tenants[0]?.id);

  const getQuickViewData = () => {
    if (!quickView) return getDefaultHomeStats();
    return quickView;
  };

  return (
    <Suspense fallback={<Loader />}>
      <SentryErrorBoundary fallback={<div>Something went wrong</div>}>
        <QuickView data={getQuickViewData()} headers={getHomeStatsHeaders()} />
      </SentryErrorBoundary>
    </Suspense>
  );
}

export default QuickViewSummary;
