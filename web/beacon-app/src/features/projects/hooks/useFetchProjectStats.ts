import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import projectStatsApiRequest from '../api/projectStatsApiService';

function useFetchProjectStats(tenantID: string) {
  const query = useQuery(
    [RQK.PROJECT_QUICK_VIEW, tenantID],
    () => projectStatsApiRequest(axiosInstance)(tenantID),
    {
      retry: 0,
      enabled: !!tenantID,
    }
  );

  return {
    getProjectQuickView: query.refetch,
    hasProjectQuickViewFailed: query.isError,
    isFetchingProjectQuickView: query.isLoading,
    projectQuickView: query.data,
    wasProjectQuickViewFetched: query.isSuccess,
    error: query.error,
  };
}

export default useFetchProjectStats;
