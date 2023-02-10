import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { apiKeysRequest } from '../api/apiKeysApiService';
import { APIKeysQuery } from '../types/apiKeyService';
export function useFetchApiKeys(): APIKeysQuery {
  const query = useQuery([RQK.APIKEYS] as const, () => apiKeysRequest(axiosInstance), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    // TODO: Change stale time
    staleTime: 1000 * 60 * 15,
  });

  return {
    getApiKeys: query.refetch,
    hasApiKeysFailed: query.isError,
    isFetchingApiKeys: query.isLoading,
    apiKeys: query.data,
    wasApiKeysFetched: query.isSuccess,
    error: query.error,
  };
}
