import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { statusRequest } from '../api/StatusApiService';
import { StatusQuery } from '../types/StatusService';

export function useFetchStatus(): StatusQuery {
  const query = useQuery([RQK.STATUS], statusRequest(axiosInstance), {
    onError: (error: any) => {
      Sentry.captureException(error);
    },
  });

  return {
    getStatus: query.refetch,
    hasStatusFailed: query.isError,
    isFetchingStatus: query.isLoading,
    status: query.data,
    wasStatusFetched: query.isSuccess,
    error: query.error,
  };
}

export default useFetchStatus;
