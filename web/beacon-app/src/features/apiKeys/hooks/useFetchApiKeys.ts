import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { apiKeysRequest } from '../api/apiKeysApiService';
import { APIKeysQuery } from '../types/apiKeyService';

export function useFetchApiKeys(projectID: string): APIKeysQuery {
  const query = useQuery(
    [RQK.API_KEYS, projectID],
    () => apiKeysRequest(axiosInstance)(projectID),
    {
      enabled: !!projectID,
      retry: 0,
    }
  );

  return {
    getApiKeys: query.refetch,
    hasApiKeysFailed: query.isError,
    isFetchingApiKeys: query.isLoading,
    apiKeys: query.data,
    wasApiKeysFetched: query.isSuccess,
    error: query.error,
  };
}
