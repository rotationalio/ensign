import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import statsRequest from './quickViewApiService';

function useFetchTenantQuickView(tenantID: string) {
  const query = useQuery([RQK.QUICK_VIEW, tenantID], () => statsRequest(axiosInstance)(tenantID), {
    enabled: !!tenantID,
    onError: (error) => {
      Sentry.captureException(error);
    },
  });
  return {
    getQuickView: query.refetch,
    hasQuickViewFailed: query.isError,
    isFetchingQuickView: query.isLoading,
    quickView: query.data,
    wasQuickViewFetched: query.isSuccess,
    error: query.error,
  };
}

export default useFetchTenantQuickView;
