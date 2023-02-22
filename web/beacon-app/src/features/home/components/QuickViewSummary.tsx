import { Trans } from '@lingui/macro';
import { Loader } from '@rotational/beacon-core';
import { Suspense } from 'react';

import { QuickView } from '@/components/common/QuickView';
import { SentryErrorBoundary } from '@/components/Error';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { useFetchQuickView } from '@/hooks/useFetchQuickView';
function QuickViewSummary() {
  const { tenants: t, getTenants } = useFetchTenants();

  if (!t) {
    getTenants();
  }

  const params = {
    key: 'tenant' as const,
    id: t?.tenants[0]?.id,
  };

  const { quickView, getQuickView } = useFetchQuickView(params);

  if (!quickView) {
    getQuickView();
  }

  return (
    <Suspense fallback={<Loader />}>
      <SentryErrorBoundary
        fallback={
          <div>
            <Trans>Something went wrong</Trans>
          </div>
        }
      >
        <QuickView data={quickView} />
      </SentryErrorBoundary>
    </Suspense>
  );
}

export default QuickViewSummary;
