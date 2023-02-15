import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import statsRequest from './quickViewApiService';
import type { QuickViewDTO, QuickViewQuery } from './quickViewService';

function useFetchQuickView(stats: QuickViewDTO): QuickViewQuery {
  const query = useQuery([RQK.QUICK_VIEW, stats.id], () => statsRequest(axiosInstance)(stats));
  return {
    getQuickView: query.refetch,
    hasQuickViewFailed: query.isError,
    isFetchingQuickView: query.isLoading,
    quickView: query.data,
    wasQuickViewFetched: query.isSuccess,
    error: query.error,
  };
}

export default useFetchQuickView;
