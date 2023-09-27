import { Loader } from '@rotational/beacon-core';
import { Suspense, useEffect, useState } from 'react';

import { QuickView } from '@/components/common/QuickView';
import { SentryErrorBoundary } from '@/components/Error';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { useFetchTenantQuickView } from '@/hooks/useFetchQuickView';

import { getDefaultHomeStats, getHomeStatsHeaders } from '../util';
function QuickViewSummary() {
  const [quickViewData, setQuickViewData] = useState<any>(getDefaultHomeStats()); // by default we will show empty values
  const { tenants } = useFetchTenants();

  // console.log('[] tenants', tenants);

  const { quickView, error } = useFetchTenantQuickView(tenants?.tenants[0]?.id);

  useEffect(() => {
    if (quickView) {
      setQuickViewData(quickView);
    }
  }, [quickView]);

  useEffect(() => {
    if (error) {
      setQuickViewData(getDefaultHomeStats());
    }
  }, [error]);

  return (
    <Suspense fallback={<Loader />}>
      <SentryErrorBoundary fallback={<div>Something went wrong</div>}>
        <QuickView data={quickViewData} headers={getHomeStatsHeaders()} />
      </SentryErrorBoundary>
    </Suspense>
  );
}

export default QuickViewSummary;
